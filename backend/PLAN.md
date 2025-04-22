# Backend Implementation Plan - Code Execution Service (Go)

This document outlines the steps to implement the Go-based Code Execution Service (CES) for AggieCode, focusing on secure, multi-language execution using Docker sandboxing.

**Phase 1: Basic Go Service Setup**

*   [ ] **Initialize/Verify Go Module:** Ensure `code-execution-service/go.mod` is correctly set up.
*   [ ] **Setup HTTP Server:** Implement a basic HTTP server in `code-execution-service/main.go` using `net/http` (or choose and integrate Gin/Fiber if preferred).
*   [ ] **Define API Endpoint (`/api/execute`):**
    *   [ ] Create handler function for `POST /api/execute`.
    *   [ ] Define request JSON structure (e.g., `{ "language": "python", "code": "print('hello')", "stdin": "" }`). Use structs for unmarshalling.
    *   [ ] Define response JSON structure (e.g., `{ "stdout": "hello\n", "stderr": "", "error": null, "execution_time_ms": 50 }`). Use structs for marshalling.
    *   [ ] Add basic request validation (language supported, code present).
*   [ ] **Logging:** Integrate a structured logging library (e.g., `log/slog` introduced in Go 1.21, or `zerolog`, `zap`).
*   [ ] **Configuration:** Manage configuration (e.g., server port, Docker settings, resource limits) via environment variables or config files.

**Phase 2: Docker Integration & Basic Execution**

*   [ ] **Docker Client:** Add the Docker Go SDK (`docker/docker/client`) as a dependency.
*   [ ] **Container Execution Logic:**
    *   [ ] Implement a function/service to handle code execution for a given language.
    *   [ ] **Temporary Files:** Create a secure temporary directory on the host for each execution request.
    *   [ ] **Write Code:** Write the user's code to the appropriate file within the temporary directory (e.g., `main.py`, `main.js`, `main.cpp`).
    *   [ ] **Select Image:** Determine the correct Docker image based on the requested language (using the existing Dockerfiles in `executors/`).
    *   [x] **Define the command to execute the code (e.g., `["python", "main.py"]`, `["node", "main.js"]`). Handle compilation steps for C++, Java, Go within the container command or entrypoint script.
        *   [x] Attach to the container to stream `stdout` and `stderr`.
        *   [x] Wait for the container to finish.
    *   [x] **Capture Output:** Read `stdout` and `stderr` from the container logs/streams.
    *   [x] **Cleanup:** Ensure the container is removed (`--rm` equivalent) and the temporary directory is deleted after execution.

**Phase 3: Security & Resource Limits**

*   [x] **Time Limits:** Implement execution time limits using `context.WithTimeout` around the Docker container execution logic (e.g., `ContainerWait`).
*   [x] **Memory Limits:** Configure memory limits when creating the container using the Docker SDK (e.g., `HostConfig.Memory`).
*   [x] **CPU Limits:** Configure CPU limits (e.g., `HostConfig.NanoCPUs`).
*   [x] **Network Isolation:** Disable networking for execution containers (`HostConfig.NetworkMode="none"`).
*   [x] **Prevent Fork Bombs/Resource Exhaustion:**
    *   [x] Set process limits (`ulimit`) within the Docker containers via `HostConfig.Ulimits`.
    *   [x] Consider limiting the number of concurrent executions the service handles.
*   [x] **Read-Only Filesystem:** Explore running containers with a read-only root filesystem (`HostConfig.ReadonlyRootfs=true`), mounting only the necessary code directory as writable temporarily if needed, or read-only after writing.

**Phase 4: Language-Specific Handling**

*   [x] **Compilation:** Refine the commands/entrypoints in the language Dockerfiles (`executors/`) to handle compilation steps for C++, Java, and Go before execution. Capture compilation errors separately in `stderr`.
*   [x] **Input (stdin):** Implement passing the `stdin` from the API request to the container's standard input stream.
*   [x] **Error Handling:** Differentiate between compilation errors, runtime errors, timeouts, and resource limit exceeded errors in the API response.

**Phase 5: Testing & Refinement**

*   [x] **Unit Tests:** Write unit tests for helper functions, request/response handling.
*   [ ] **Integration Tests:** Write integration tests that use the Docker SDK to test the container execution logic for each language, including error cases and resource limits.
*   [x] **API Tests:** Test the `/api/execute` endpoint thoroughly.
*   [ ] **Dockerfile Review:** Ensure language Dockerfiles are minimal, secure, and efficient.
*   [x] **Concurrency:** Test concurrent execution requests.

**Phase 6: Deployment**

*   [ ] **Containerize CES:** Create a Dockerfile for the Go CES application itself.
*   [ ] **Deployment Strategy:** Plan deployment (e.g., Cloud Run, Kubernetes, simple VM). Ensure Docker (or the chosen container runtime) is available in the deployment environment.

**Future Enhancements (Post-MVP)**

*   [ ] **gVisor Integration:** Explore replacing/augmenting Docker with gVisor for enhanced sandboxing security.
*   [ ] **Queueing System:** Implement a job queue (e.g., Redis-based) if concurrent load becomes high.
*   [ ] **Caching:** Cache compilation results for languages like C++/Java/Go if the same code is executed repeatedly (though likely low ROI for typical student lab usage).
*   [ ] **More Languages:** Add support for other languages by creating corresponding Dockerfiles and execution logic.
*   **Monitoring:** Integrate with Prometheus/Grafana for metrics (execution times, error rates, resource usage).
