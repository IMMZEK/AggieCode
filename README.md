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
    *   **OAuth 2.0 with Google as the identity provider** for secure user login and session management.
    *   **Authorization Code Grant flow** for enhanced security.
    *   **Backend server** to handle sensitive operations like token exchange and refresh.
*   **Efficient Backend:**
    *   Lightweight Node.js backend for code execution, synchronization, and handling OAuth 2.0 flows.
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
*   **OAuth 2.0 Handling:** [axios](https://axios-http.com/) for making requests to the backend for token management.

### Backend

*   **Runtime:** [Node.js](https://nodejs.org/)
*   **Framework:** [Express.js](https://expressjs.com/)
*   **Authentication:**
    *   **OAuth 2.0** using the **Authorization Code Grant** flow.
    *   **Google as the identity provider.**
    *   Secure handling of **Client Secret** and token exchange on the backend.
*   **ML Model Integration:**
    *   [Hugging Face `Inference API`](https://huggingface.co/docs/api-inference/index) to load and run the **CodeBERT** model for autocompletion.
    *   API endpoint (`/api/autocomplete`) to handle autocompletion requests.
*   **Real-Time Collaboration (Partially Implemented):**
    *   [Y.js](https://github.com/yjs/yjs) (CRDT-based real-time state management) - *Planned*
    *   [y-websocket](https://github.com/yjs/y-websocket) (for document synchronization) - *Planned*
    *   [Socket.IO](https://socket.io/) (for auxiliary real-time events) - *Planned*
*   **Code Execution:**
    *   [Judge0 API](https://api.judge0.com/) (for server-side Python and C++ execution) - *Planned*
    *   [Pyodide](https://pyodide.org/en/stable/) (Optional, for browser-based Python execution) - *Planned*
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

### Phase 2: Authentication with OAuth 2.0 (In Progress)

*   Deciding Between:
    *   A - Implement OAuth 2.0 Authorization Code Grant flow.
    *   B - Use Google (Firebase) as the identity provider.
*   Create a backend server (Express.js) to handle token exchange and refresh.
*   ✅ Securely store the Client Secret on the backend.
*   Store access tokens and refresh tokens (if applicable) in the frontend (Vuex/Pinia and localStorage).
*   ✅ Implement token refresh logic.

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

### Phase 6: Code Execution

*   Integrate Judge0 API for server-side code execution.
*   Design a shared output window.
*   Explore Pyodide for client-side Python execution.

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
*   **`OAuth2Callback.vue`:** Handles the redirect from Google after authorization, exchanges the code for tokens, and stores them.
*   **`services/oauth2.js`:** Manages the OAuth 2.0 flow, including generating the authorization URL, exchanging the code for tokens, and refreshing tokens (by interacting with the backend).
*   **`store/modules/oauth2.js` (Vuex):**  Manages the OAuth 2.0 state (access token, refresh token, expiration) in the Vuex store.

### Backend Architecture

*   **API Endpoints:**
    *   `POST /api/token`: Handles the exchange of the authorization code for an access token and refresh token (securely using the Client Secret).
    *   `POST /api/refresh`: Refreshes the access token using the refresh token.
    *   `GET /api/config` (Optional): Provides public configuration information to the frontend (e.g., Google Client ID).
    *   `POST /api/autocomplete`: Handles autocompletion requests, loads the CodeBERT model, performs inference, and returns suggestions.
    *   `POST /files/content`: Saves file content to the server (currently stored as text blobs in PostgreSQL for the demo). - *Planned*
    *   `POST /execute`: Sends code to Judge0 for execution and returns the results. - *Planned*
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