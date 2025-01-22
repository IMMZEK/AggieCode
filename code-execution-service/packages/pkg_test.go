package pkg

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateRequest(t *testing.T) {
	service := NewExecutionService()

	tests := []struct {
		name    string
		request ExecutionRequest
		wantErr error
	}{
		{
			name: "valid python request",
			request: ExecutionRequest{
				Language: "python",
				Code:     "print('hello')",
				Method:   "run",
			},
			wantErr: nil,
		},
		{
			name: "empty language",
			request: ExecutionRequest{
				Language: "",
				Code:     "print('hello')",
				Method:   "run",
			},
			wantErr: ErrInvalidRequest,
		},
		{
			name: "unsupported language",
			request: ExecutionRequest{
				Language: "rust",
				Code:     "println!('hello')",
				Method:   "run",
			},
			wantErr: ErrLanguageNotSupported,
		},
		{
			name: "invalid method",
			request: ExecutionRequest{
				Language: "python",
				Code:     "print('hello')",
				Method:   "invalid",
			},
			wantErr: ErrMethodNotSupported,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateRequest(&tt.request)
			if err != tt.wantErr {
				t.Errorf("validateRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSanitizeCode(t *testing.T) {
	sanitizer := NewSanitizer(1000)

	tests := []struct {
		name     string
		code     string
		language string
		wantErr  bool
	}{
		{
			name:     "valid python code",
			code:     "print('hello')",
			language: "python",
			wantErr:  false,
		},
		{
			name:     "python code with system access",
			code:     "import os\nos.system('rm -rf /')",
			language: "python",
			wantErr:  true,
		},
		{
			name:     "valid go code",
			code:     "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"hello\")\n}",
			language: "go",
			wantErr:  false,
		},
		{
			name:     "go code with unsafe imports",
			code:     "package main\n\nimport \"os\"\n\nfunc main() {\n\tos.Exit(1)\n}",
			language: "go",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sanitizer.SanitizeCode(tt.code, tt.language)
			if (err != nil) != tt.wantErr {
				t.Errorf("SanitizeCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleExecute(t *testing.T) {
	service := NewExecutionService()

	tests := []struct {
		name           string
		request        ExecutionRequest
		wantStatus     int
		wantErrInBody  bool
		wantOutputNull bool
	}{
		{
			name: "valid python request",
			request: ExecutionRequest{
				Language: "python",
				Code:     "print('hello')",
				Method:   "run",
			},
			wantStatus:     http.StatusOK,
			wantErrInBody:  false,
			wantOutputNull: false,
		},
		{
			name: "invalid language",
			request: ExecutionRequest{
				Language: "invalid",
				Code:     "print('hello')",
				Method:   "run",
			},
			wantStatus:     http.StatusOK,
			wantErrInBody:  true,
			wantOutputNull: true,
		},
		{
			name: "dangerous code",
			request: ExecutionRequest{
				Language: "python",
				Code:     "import os; os.system('rm -rf /')",
				Method:   "run",
			},
			wantStatus:     http.StatusOK,
			wantErrInBody:  true,
			wantOutputNull: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, _ := json.Marshal(tt.request)
			req := httptest.NewRequest(http.MethodPost, "/execute", bytes.NewBuffer(reqBody))
			w := httptest.NewRecorder()

			service.HandleExecute(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("HandleExecute() status = %v, want %v", w.Code, tt.wantStatus)
			}

			var response ExecutionResponse
			json.NewDecoder(w.Body).Decode(&response)

			if tt.wantErrInBody && response.Error == "" {
				t.Error("HandleExecute() expected error in response body, got none")
			}

			if !tt.wantErrInBody && response.Error != "" {
				t.Errorf("HandleExecute() unexpected error in response body: %v", response.Error)
			}

			if tt.wantOutputNull && response.Output != "" {
				t.Error("HandleExecute() expected null output, got output")
			}
		})
	}
}

func TestRateLimiter(t *testing.T) {
	limiter := NewRateLimiter(2, 1) // 2 requests per minute, burst of 1
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	limitedHandler := limiter.Limit(handler)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "127.0.0.1"

	// First request should succeed
	w1 := httptest.NewRecorder()
	limitedHandler.ServeHTTP(w1, req)
	if w1.Code != http.StatusOK {
		t.Errorf("First request: got status %v, want %v", w1.Code, http.StatusOK)
	}

	// Second request should be rate limited
	w2 := httptest.NewRecorder()
	limitedHandler.ServeHTTP(w2, req)
	if w2.Code != http.StatusTooManyRequests {
		t.Errorf("Second request: got status %v, want %v", w2.Code, http.StatusTooManyRequests)
	}
}
