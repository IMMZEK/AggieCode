# AggieCode

# AggieCode: Real-time Collaborative IDE for Students

<!-- [![Build Status](https://github.com/IMMZEK/AggieCode/actions/workflows/build.yml/badge.svg)](https://github.com/IMMZEK/AggieCode/actions/workflows/build.yml) -->
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)

<p align="center">
  <img src="docs/logo.jpg" alt="AggieCode Logo" width="300">
</p>

## Overview

AggieCode is a web-based, real-time collaborative Integrated Development Environment (IDE) designed to enhance first-year engineering students' learning experience engaged in team-based lab environments. It empowers students to collaboratively write, execute, and debug code while sharing outputs in real time. AggieCode aims to provide an intuitive interface, inspired by Visual Studio Code and powered by the Monaco Editor, it also provides intelligent code autocompletion powered by Machine Learning.AggieCode is a web-based, real-time collaborative Integrated Development Environment (IDE) specifically designed to enhance the learning experience of first-year engineering students in team-based lab environments. It empowers students to collaborate on writing, executing, and debugging code while sharing outputs in real time. AggieCode aims to provide an intuitive interface, inspired by Visual Studio Code, and is powered by the Monaco Editor. Additionally, it offers intelligent code autocompletion powered by Machine Learning.

## Key Features

*   **Real-time Collaborative Coding:**
    *   Shared file editing with live cursor synchronization.
    *   Simultaneous code editing by multiple users.
*   **Multi-Language Support:**
    *   Code execution for Python, C++, and more.
    *   Shared output window visible to all collaborators.
*   **User-Friendly Interface:**
    *   Inspired by VS Code, providing a familiar and intuitive coding experience.
    *   Responsive design for seamless use on various devices.
*   **Secure Authentication:**
    *   **Firebase Authentication** for secure user login and session management.
*   **Efficient Backend:**
    *   Node.js backend for the main application logic.
    *   **Go-based microservice for secure code execution.**
    *   **(Planned) Go-based microservice for ML-powered autocompletion.**
    *   **(Planned) Go-based server for real-time collaboration using Y.js.**
    *   Optimized for fast response times and minimal resource usage.
*   **ML-Powered Code Autocompletion:**
    *   **Intelligent code suggestions** powered by a **CodeBERT** model from Hugging Face.
    *   **Context-aware autocompletion** to improve coding efficiency.
    *   **Support for Python** (initially) with plans to expand to other languages.
*   **Demo-Ready:**
    *   Easily deployable for personal showcases, classroom demonstrations, and local development.
    *   **Free Tier Friendly:** Designed to work within the free tier limits of cloud providers like DigitalOcean, Linode, or Vultr, as well as the option to utilize the free tier of Google Cloud Platform (GCP) for initial development and testing.
* **Theme Toggling**:
    * Light and Dark theme are available to provide the user a more personalized experience.

## Technology Stack

### Frontend

*   **Framework:** [Vue.js](https://vuejs.org/)
*   **UI Library:** [Vuetify](https://vuetifyjs.com/) (Material Design components)
*   **Code Editor:** [Monaco Editor](https://microsoft.github.io/monaco-editor/) (via [monaco-vue](https://www.npmjs.com/package/@guolao/vue-monaco-editor))
*   **State Management:** [Vuex](https://vuex.vuejs.org/) (or [Pinia](https://pinia.vuejs.org/))
*   **Styling:** [Tailwind CSS](https://tailwindcss.com/) (Optional, for rapid UI adjustments)
*   **Icons**: [Font Awesome](https://fontawesome.com/)
*   **Authentication:** [Firebase Authentication](https://firebase.google.com/docs/auth)
*   **Real-time Collaboration:** [Y.js](https://github.com/yjs/yjs)

### Backend

*   **Main Backend:** [Node.js](https://nodejs.org/) with [Express.js](https://expressjs.com/)
*   **Code Execution Service:** **Go** with `net/http` or a framework like `Gin` or `Fiber`, **Docker** or **gVisor** or **WebAssembly** for sandboxing.
*   **(Planned) ML Autocompletion Service:** Go, `net/http` (or framework), `gomlx` or `tensorflow/tensorflow/go`, potentially TensorFlow Serving or TorchServe.
*   **(Planned) Real-Time Collaboration (Y.js) Server:** Go, `net/http` (or framework), `go-yjs` (or alternative), `gorilla/websocket`.
*   **Authentication:** [Firebase Admin SDK](https://firebase.google.com/docs/admin/setup) (for Node.js and Go)
*   **Real-Time Collaboration (Partially Implemented):**
    *   [Y.js](https://github.com/yjs/yjs) (CRDT-based real-time state management) - *Planned*
    *   [y-websocket](https://github.com/yjs/y-websocket) (for document synchronization) - *Planned*
    *   [Socket.IO](https://socket.io/) (for auxiliary real-time events) - *Planned*
*   **Database:**
    *   [PostgreSQL](https://www.postgresql.org/) (for persistent storage) - *Planned*
    *   [Redis](https://redis.io/) (for in-memory caching) - *Planned*

### Infrastructure

*   **Development:** Local development with the option to use the free tier of cloud providers (for Demo Purposes as of now).
*   **Deployment (Initial Stages):**
    *   **Free Tier Friendly:** Designed to utilize the free tiers of GCP (as of now).
    *   **Google Cloud Platform (GCP):**  Utilizes GCP's free tier for initial development and testing (utilizing Cloud Functions or Cloud Run for the backend and Cloud Storage).

### Monitoring and Testing

*   **Error Tracking:** [Sentry](https://sentry.io/) - *Planned*
*   **Performance Monitoring:** [Prometheus](https://prometheus.io/) + [Grafana](https://grafana.com/) - *Planned*
*   **Frontend Testing:**
    *   [Cypress](https://www.cypress.io/) (for end-to-end testing) - *Planned*
    *   [Jest](https://jestjs.io/) (for unit testing) - *Planned*
*   **Backend Testing:**
    *   [Mocha](https://mochajs.org/) or [Jest](https://jestjs.io/) with [supertest](https://github.com/visionmedia/supertest) (for API testing) - *Planned*

## Project Roadmap

### Phase 1: Initial Setup (Completed)

*   ✅ Set up Vue.js project with Vuetify and Vuex/Pinia.
*   ✅ Integrate Monaco Editor using `monaco-vue`.
*   ✅ Implement basic routing (login, IDE workspace).

### Phase 2: Authentication with Firebase (In Progress)

*   Implement user authentication using Firebase Authentication.
*   Store access tokens in the frontend (Vuex/Pinia and localStorage).

### Phase 3: Basic ML Autocompletion (In Progress)

*   ✅ Select a suitable pre-trained open-source model (**CodeBERT**).
*   Create a backend API endpoint (`/api/autocomplete`) to handle autocompletion requests.
*   Load the CodeBERT model and tokenizer using the Hugging Face `transformers` library.
*   Implement input preprocessing and inference logic in the backend.
*   Integrate the autocompletion API with the `CodeEditor.vue` component.
*   Implement debouncing on the frontend.
*   Display suggestions in the Monaco Editor.

### Phase 4: Enhanced Autocompletion

*   Fine-tune the CodeBERT model (if necessary) on a dataset of Python code.
*   Optimize model inference for latency (experiment with `generate()` parameters).
*   Implement basic caching on the backend.
*   Improve context handling and UI/UX for displaying suggestions.

### Phase 5: Collaboration (Partially Implemented)

*   Integrate Y.js and `y-websocket` for real-time document synchronization.
*   Implement Socket.IO for real-time event broadcasting.
*   Develop collaborative editing features:
    *   Cursor highlighting.
    *   File sharing.

### Phase 6: Code Execution (In Progress)

*   **Develop a Go-based microservice for secure code execution.**
*   **Use Docker, gVisor, or WebAssembly for sandboxing.**
*   Implement resource limits (CPU, memory, time).
*   Support multiple languages (Python, C++, Java, etc.).
*   Integrate with the frontend to display output in a shared window.

### Phase 7: Advanced Autocompletion and Language Expansion

*   Add support for multiple languages (e.g., C++, Java, JavaScript).
*   Implement more advanced features:
    *   Type inference.
    *   Snippet completion.

### Phase 8: Session Management

*   Implement session-based user management.
*   Enable joining specific lab sessions via session codes.

### Phase 9: Finalize UI & Testing

*   Refine/Re-design UI using Vuetify.
*   Implement comprehensive testing (real-time collaboration, edge cases, security).
*   Integrate Sentry for error tracking.

### Phase 10: Scaling and Refinement

*   Migrate to a dedicated model serving framework (if needed).
*   Scale infrastructure based on projected usage during the semester.
*   Set up monitoring and alerting.
*   (Reach) Establish a continuous training/fine-tuning pipeline.

## Architecture

### Frontend Components

*   **`FileExplorer.vue`:** Manages the file tree, allowing users to create, open, rename, and delete files.
*   **`CodeEditor.vue`:** The core Monaco Editor component for real-time collaborative coding, now with ML-powered autocompletion using CodeBERT.
*   **`OutputWindow.vue`:** Displays the output of code executed via Judge0.
*   **`CollaborationPanel.vue`:** Shows a list of online collaborators and their cursor positions.
*   **`Navbar.vue`:** Provides navigation, theme switching, and other global actions.
*   **`OAuth2Callback.vue`:** Replaced with Firebase for Authentification.
*   **`services/auth.js`:** Manages Firebase authentication related logic such as token refresh and generating authorization URL.
*   **`store/modules/auth.js` (Vuex):** Manages the Firebase authentication state.

### Backend Architecture

*   **Main Backend (Node.js):**
    *   Handles general API requests, user sessions (with Firebase), and potentially serves as a proxy to other services.
    *   API Endpoints:
        *   `POST /api/token`:  (Removed since we are using Firebase)
        *   `POST /api/refresh`: (Removed since we are using Firebase)
        *   `GET /api/config` (Optional): Provides public configuration information to the frontend (e.g., Google Client ID).
        *   `POST /files/content`: Saves file content to the server (currently stored as text blobs in PostgreSQL for the demo). - *Planned*
*   **Code Execution Service (Go):**
    *   Provides secure code execution in a sandboxed environment.
    *   API Endpoint:
        *   `POST /api/execute`: Accepts code, language, and optional input, executes the code, and returns the output.
*   **ML Autocompletion Service (Go) - *Planned*:**
    *   Provides intelligent code suggestions using a pre-trained CodeBERT model.
    *   API Endpoint:
        *   `POST /api/autocomplete`: Accepts code, cursor position, and language, and returns a list of autocompletion suggestions.
*   **Real-Time Collaboration Server (Go or Node.js) - *Partially Implemented*:**
    *   Handles real-time document synchronization using Y.js and potentially other real-time events using Socket.IO.
    *   May be implemented in Node.js (using `y-websocket`) or in Go (using `go-yjs` or an alternative technology if `go-yjs` is not viable).
*   **WebSocket Events (Planned):**
    *   `cursor-update`: Broadcasts real-time cursor position changes to all collaborators in the same session.
    *   `file-update`: Synchronizes file content changes across all clients.
    *   `output-update`: Shares the output of code execution with all collaborators.

## Acknowledgements

*   [Vue.js](https://vuejs.org/)
*   [Vuetify](https://vuetifyjs.com/)
*   [Monaco Editor](https://microsoft.github.io/monaco-editor/)
*   [Y.js](https://github.com/yjs/yjs)
*   [Judge0 API](https://api.judge0.com/)
*   [Google Cloud Platform](https://cloud.google.com/)
*   [Hugging Face Transformers](https://huggingface.co/docs/transformers/index)
*   [CodeBERT](https://huggingface.co/microsoft/codebert-base)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.