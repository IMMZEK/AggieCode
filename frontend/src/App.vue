<template>
  <div :class="['app', isDarkMode ? 'dark' : 'light']">
    <!-- Header Section -->
    <header class="app-header">
      <button @click="toggleDarkMode" class="theme-toggle">
        <span v-if="isDarkMode">🌙</span>
        <span v-else>☀️</span>
      </button>
      <h1>AggieCode IDE</h1>
      <div class="login-avatar"></div>
    </header>

    <!-- Main Layout -->
    <div class="app-layout">
      <!-- Left Sidebar (File Management) -->
      <aside class="sidebar">
        <h2>Files</h2>
        <ul>
          <li>main.js</li>
          <li>App.vue</li>
          <li>CodeEditor.vue</li>
        </ul>
      </aside>

      <!-- Center: IDE Editor -->
      <main class="editor-container">
        <CodeEditor :isDarkMode="isDarkMode" />
      </main>

      <!-- Right Sidebar -->
      <aside class="right-sidebar">
        <!-- Terminal Section -->
        <section class="terminal">
          <header class="terminal-header">
            <h2>Output</h2>

          </header>
          <div class="terminal-output">
            <p>// Output will appear here...</p>
          </div>
        </section>

        <!-- Chat Section -->
        <section class="chat">
          <header class="chat-header">
            <h2>Chat</h2>
          </header>
          <div class="chat-box">
            <ul>
              <li><strong>User1:</strong> Hello!</li>
              <li><strong>User2:</strong> Hi, let's start coding!</li>
            </ul>
            <div class="chat-input-container">
              <input type="text" placeholder="Type a message..." />
              <button>Send</button>
            </div>
          </div>
        </section>
      </aside>
    </div>

    <!-- Footer Section -->
    <footer class="app-footer">
      <p>Built with ❤️ by the AggieCode Team</p>
    </footer>
  </div>
</template>

<script>
import { ref } from "vue";
import CodeEditor from "./components/CodeEditor.vue";

export default {
  components: {
    CodeEditor,
  },
  setup() {
    const isDarkMode = ref(false);

    const toggleDarkMode = () => {
      isDarkMode.value = !isDarkMode.value;
    };

    return {
      isDarkMode,
      toggleDarkMode,
    };
  },
};
</script>

<style>
/* Shared Styles */
body, html {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
  font-family: 'Inter', 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  height: 100%;
}

.app {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

/* Light Mode Styles */
.light {
  background: linear-gradient(to bottom right, #f3f4f6, #e9ecef);
  color: #333;
}

.light .app-header {
  background-color: #800000;
  color: #ffffff;
}

.light .sidebar {
  background-color: #f9f9f9;
  color: #333;
  border-right: 1px solid #ddd;
}

.light .editor-container {
  background-color: #ffffff;
  color: #000;
}

.light .right-sidebar {
  background-color: #f9f9f9;
  color: #333;
}

.light .terminal {
  background-color: #1e1e1e;
  color: #ffffff;
}

.light .chat {
  background-color: #ffffff;
  color: #333;
}

.light .app-footer {
  background-color: #800000;
  color: #ffffff;
}

/* Dark Mode Styles */
.dark {
  background: linear-gradient(to bottom right, #1e1e1e, #121212);
  color: #f5f5f5;
}

.dark .app-header {
  background-color: #3d3d3d;
  color: #f5f5f5;
}

.dark .sidebar {
  background-color: #252525;
  color: #f5f5f5;
  border-right: 1px solid #444;
}

.dark .editor-container {
  background-color: #2a2a2a;
  color: #f5f5f5;
}

.dark .right-sidebar {
  background-color: #252525;
  color: #f5f5f5;
}

.dark .terminal {
  background-color: #000;
  color: #0f0;
}

.dark .chat {
  background-color: #2a2a2a;
  color: #f5f5f5;
}

.dark .app-footer {
  background-color: #3d3d3d;
  color: #f5f5f5;
}

/* Header */
.app-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 1rem;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.app-header h1 {
  margin: 0;
  font-size: 1.4rem;
  letter-spacing: 0.5px;
}

.login-avatar {
  height: 40px;
  width: 40px;
  background-color: #ffffff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.8rem;
  cursor: pointer;
}

.theme-toggle {
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
}

/* Layout */
.app-layout {
  display: flex;
  flex: 1;
}

/* Sidebar (File Management) */
.sidebar {
  width: 200px;
  padding: 1rem;
}

/* Editor */
.editor-container {
  flex: 1;
  margin: 1rem;
  border-radius: 12px;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

/* Right Sidebar */
.right-sidebar {
  width: 300px;
  display: flex;
  flex-direction: column;
}

/* Terminal Section */
.terminal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem 1rem;
}

.run-button {
  padding: 0.5rem 1rem;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

/* Chat Section */
.chat-header {
  padding: 0.5rem 1rem;
}

.chat-box {
  flex: 1;
  padding: 1rem;
}

.chat-input-container {
  display: flex;
  padding: 0.5rem;
}

.chat-input-container input {
  flex: 1;
  margin-right: 0.5rem;
}

/* Footer */
.app-footer {
  text-align: center;
  padding: 0.75rem 0;
}
</style>