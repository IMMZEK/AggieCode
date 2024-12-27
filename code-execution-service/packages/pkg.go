package pkg

import (
	"code-execution-service/packages/lang"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"
)

var (
	ErrInvalidRequest       = errors.New("The request parameters provided are invalid.")
	ErrLanguageNotSupported = errors.New("The specified programming language is not supported.")
	ErrMethodNotSupported   = errors.New("The specified method is not supported for execution.")
)

type ExecutionRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
	Method   string `json:"method"`
}

type ExecutionResponse struct {
	Output        string `json:"output"`
	Error         string `json:"error,omitempty"`
	StatusMessage string `json:"status_message,omitempty"`
}

type ExecutionService struct {
	containers  map[string]string
	RateLimiter *RateLimiter
	Sanitizer   *Sanitizer
}

type RateLimiter struct {
	limit int
}

type Sanitizer struct {
	maxCodeLength int
}

type SanitizationError struct {
	Message string `json:"message"`
	Details string `json:"details"`
}

func (e *SanitizationError) Error() string {
	return e.Message
}

func (s *Sanitizer) SanitizeCode(code, language string) error {
	if len(code) > s.maxCodeLength {
		return &SanitizationError{
			Message: "Code length exceeds maximum limit",
			Details: "Max length allowed is " + string(rune(s.maxCodeLength)),
		}
	}

	systemPatterns := []string{
		`(?i)(subprocess|exec\.|shell|eval|child_process)`,
		`(?i)(io/ioutil|os\.Open|os\.Create|os\.Remove)`,
		`(?i)(net\.Listen|net\.Dial|http\.|urllib|axios)`,
	}
	if matched, err := matchPatterns(systemPatterns, code); err != nil || matched {
		return &SanitizationError{
			Message: "Prohibited system-level access detected",
			Details: "Code contains restricted system operations",
		}
	}

	var restrictedPatterns []string
	switch language {
	case "python":
		if strings.Contains(code, "import") || strings.Contains(code, "from") {
			restrictedPatterns = []string{
				`^import\s+(?!math|random|datetime|json|re|string|collections|itertools|functools|typing).*$`,
				`^from\s+(?!math|random|datetime|json|re|string|collections|itertools|functools|typing)\s+import.*$`,
			}
		}
		restrictedPatterns = append(restrictedPatterns, []string{
			`__import__`, `globals|locals|vars`, `getattr|setattr|delattr`,
			`pip|setuptools|pkg_resources`,
		}...)
	case "go":
		safePackages := []string{
			"fmt", "strings", "strconv", "math", "time", "encoding/json", "errors",
			"sort", "regexp",
		}

		if strings.Contains(code, "import") {
			lines := strings.Split(code, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "import") {
					importMatch := regexp.MustCompile(`^import\s+"([^"]+)"`).FindStringSubmatch(line)
					if importMatch != nil {
						pkg := importMatch[1]
						isSafe := false
						for _, safePkg := range safePackages {
							if pkg == safePkg {
								isSafe = true
								break
							}
						}
						if !isSafe {
							return &SanitizationError{
								Message: "Prohibited go code pattern detected",
								Details: "Unauthorized import: " + pkg,
							}
						}
					}
				}
			}
		}

		restrictedPatterns = []string{
			`unsafe\.`, `reflect\.`, `plugin\.`, `go/ast`,
			`syscall\.`, `debug\.`, `runtime\.`, `os\.Exit`, `panic\(`,
		}
	case "js":
		if strings.Contains(code, "require") || strings.Contains(code, "import") {
			restrictedPatterns = []string{
				`require\(.*\)`, `import\s+.*\s+from`, `import\s*{.*}`,
			}
		}
		restrictedPatterns = append(restrictedPatterns, []string{
			`process`, `global`, `Buffer`, `__proto__`, `prototype`,
			`fs`, `child_process`, `eval`, `Function`, `process\.env`}...)
	default:
		return errors.New("unsupported language: " + language)
	}

	if len(restrictedPatterns) > 0 {
		if matched, err := matchPatterns(restrictedPatterns, code); err != nil || matched {
			return &SanitizationError{
				Message: "Prohibited " + language + " code pattern detected",
				Details: "Unauthorized module or operation",
			}
		}
	}

	return nil
}

func matchPatterns(patterns []string, text string) (bool, error) {
	for _, pattern := range patterns {
		match, err := regexp.MatchString(pattern, text)
		if err != nil {
			return false, err
		}
		if match {
			return true, nil
		}
	}
	return false, nil
}

func NewExecutionService() *ExecutionService {
	return &ExecutionService{
		containers: map[string]string{
			"cpp":    "cpp-executor",
			"java":   "java-executor",
			"js":     "js-executor",
			"python": "python-executor",
		},
		RateLimiter: NewRateLimiter(),
		Sanitizer:   NewSanitizer(1000),
	}
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		limit: 100,
	}
}

func NewSanitizer(maxSize int) *Sanitizer {
	return &Sanitizer{
		maxCodeLength: maxSize,
	}
}

func (s *ExecutionService) validateRequest(req *ExecutionRequest) error {
	if req.Language == "" || req.Code == "" || req.Method == "" {
		return ErrInvalidRequest
	}

	if _, supported := s.containers[req.Language]; !supported {
		return ErrLanguageNotSupported
	}

	if req.Method != "run" && req.Method != "test" {
		return ErrMethodNotSupported
	}

	return nil
}

func (s *ExecutionService) HandleExecute(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var req ExecutionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, ErrInvalidRequest.Error(), http.StatusBadRequest)
		return
	}

	// Validate the request parameters
	if err := s.validateRequest(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Sanitize the code
	if err := s.Sanitizer.SanitizeCode(req.Code, req.Language); err != nil {
		response := ExecutionResponse{
			Error:         err.Error(),
			StatusMessage: "Code Sanitization Error",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Execute the code
	output, err := ExecuteCode(req.Language, req.Code, req.Method)
	response := ExecutionResponse{
		Output:        output,
		StatusMessage: "Accepted",
	}

	if err != nil {
		response.Error = err.Error()
		response.StatusMessage = "Runtime Error"
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func ExecuteCode(language, code, method string) (string, error) {
	containerName, ok := map[string]string{
		"cpp":    "cpp-executor",
		"java":   "java-executor",
		"js":     "js-executor",
		"python": "python-executor",
		"go":     "go-executor",
	}[language]

	if !ok {
		return "", ErrLanguageNotSupported
	}

	switch language {
	case "cpp":
		return lang.ExecuteCppCode(containerName, code)
	case "java":
		return lang.ExecuteJavaCode(containerName, code)
	case "js":
		return lang.ExecuteJsCode(containerName, code)
	case "python":
		return lang.ExecutePythonCode(containerName, code)
	case "go":
		return lang.ExecuteGoCode(containerName, code)
	default:
		return "", ErrLanguageNotSupported
	}
}
