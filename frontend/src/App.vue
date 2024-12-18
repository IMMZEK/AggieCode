<template>
  <div :class="['app', isDarkMode ? 'dark' : 'light']">
    <!-- Header -->
    <header class="app-header">
      <div class="header-left">
        <button class="hamburger-menu" @click="toggleLeftSidebar">
          <span></span>
          <span></span>
          <span></span>
        </button>
        <h1>AggieCode IDE</h1>
      </div>
      <div class="header-right">
        <button @click="toggleDarkMode" class="theme-toggle">
          <span v-if="isDarkMode">🌙</span>
          <span v-else>☀️</span>
        </button>
        <div class="user-profile" @click="toggleProfileDropdown">
          <span class="user-initials"></span>
          <ul class="profile-dropdown" v-show="showProfileDropdown">
            <li>Profile</li>
            <li>Settings</li>
            <li @click="logout">Logout</li>
          </ul>
        </div>
      </div>
    </header>

    <!-- Main Layout -->
    <div class="app-layout">
      <!-- Left Sidebar (File Management) -->
      <aside class="sidebar" :class="{ collapsed: isLeftSidebarCollapsed }">
        <div class="sidebar-header">
          <h2 v-show="!isLeftSidebarCollapsed" class="sidebar-title">Files</h2>
          <button class="collapse-button" @click="toggleLeftSidebar">
            <i
              :class="
                isLeftSidebarCollapsed
                  ? 'fas fa-angle-right'
                  : 'fas fa-angle-left'
              "
            ></i>
          </button>
        </div>
        <ul v-show="!isLeftSidebarCollapsed" class="file-list">
          <li
            v-for="file in files"
            :key="file.name"
            :class="{ active: activeFile === file.name }"
            @click="selectFile(file.name)"
          >
            <i :class="getFileIcon(file.name)"></i>
            {{ file.name }}
          </li>
        </ul>
      </aside>

      <!-- Center: IDE Editor -->
      <main class="editor-container">
        <CodeEditor :isDarkMode="isDarkMode" :activeFile="activeFile" />
      </main>

      <!-- Right Sidebar -->
      <aside class="right-sidebar" :class="{ collapsed: isRightSidebarCollapsed }">
        <div class="sidebar-header">
          <h2 v-show="!isRightSidebarCollapsed" class="sidebar-title">
            {{ activeRightSidebarTabTitle }}
          </h2>
          <button class="collapse-button" @click="toggleRightSidebar">
            <i
              :class="
                isRightSidebarCollapsed
                  ? 'fas fa-angle-left'
                  : 'fas fa-angle-right'
              "
            ></i>
          </button>
        </div>
        <div class="tabs" v-show="!isRightSidebarCollapsed">
          <button
            v-for="tab in rightSidebarTabs"
            :key="tab.id"
            :class="{ active: activeRightSidebarTab === tab.id }"
            @click="activeRightSidebarTab = tab.id"
          >
            {{ tab.label }}
          </button>
        </div>
        <div class="tab-content" v-show="!isRightSidebarCollapsed">
          <section class="terminal" v-if="activeRightSidebarTab === 'output'">
            <div class="terminal-output">
              <p>// Output will appear here...</p>
            </div>
          </section>
          <section class="chat" v-if="activeRightSidebarTab === 'chat'">
            <div class="chat-box">
              <ul>
                <li><strong>User1:</strong> Hello!</li>
                <li><strong>User2:</strong> Hi, let's start coding!</li>
              </ul>
            </div>
            <div class="chat-input-container">
              <input type="text" placeholder="Type a message..." />
              <button>Send</button>
            </div>
          </section>
          <section class="debug" v-if="activeRightSidebarTab === 'debug'">
            <!-- Placeholder content for debug panel -->
            <p>Debugging controls and information will be here.</p>
          </section>
          <section
            class="extensions"
            v-if="activeRightSidebarTab === 'extensions'"
          >
            <!-- Placeholder content for extensions panel -->
            <p>Available extensions will be listed here.</p>
          </section>
        </div>
      </aside>
    </div>

    <!-- Footer -->
    <footer class="app-footer">
      <p>Built with ❤️ by the AggieCode Team</p>
    </footer>
  </div>
</template>

<script>
import { ref, computed } from "vue";
import CodeEditor from "./components/CodeEditor.vue";
import {
  faJs,
  faPython,
  faHtml5,
  faCss3Alt,
  faJava,
} from "@fortawesome/free-brands-svg-icons";
import { faFile, faBug } from "@fortawesome/free-solid-svg-icons";
import { library } from "@fortawesome/fontawesome-svg-core";

library.add(faJs, faPython, faHtml5, faCss3Alt, faJava, faFile, faBug);

export default {
  components: {
    CodeEditor,
  },
  setup() {
    const isDarkMode = ref(false);
    const isLeftSidebarCollapsed = ref(false);
    const isRightSidebarCollapsed = ref(false);
    const activeRightSidebarTab = ref("output");
    const showProfileDropdown = ref(false);
    const activeFile = ref("main.js");
    const files = ref([
      { name: "main.js", icon: "fab fa-js" },
      { name: "App.vue", icon: "fab fa-vuejs" },
      { name: "CodeEditor.vue", icon: "fab fa-vuejs" },
      { name: "index.html", icon: "fab fa-html5" },
      { name: "style.css", icon: "fab fa-css3-alt" },
    ]);

    const rightSidebarTabs = ref([
      { id: "output", label: "Output" },
      { id: "chat", label: "Chat" },
      { id: "debug", label: "Debug" },
      { id: "extensions", label: "Extensions" },
    ]);

    const activeRightSidebarTabTitle = computed(() => {
      const activeTab = rightSidebarTabs.value.find(
        (tab) => tab.id === activeRightSidebarTab.value
      );
      return activeTab ? activeTab.label : "";
    });

    const toggleDarkMode = () => {
      isDarkMode.value = !isDarkMode.value;
    };

    const toggleLeftSidebar = () => {
      isLeftSidebarCollapsed.value = !isLeftSidebarCollapsed.value;
    };

    const toggleRightSidebar = () => {
      isRightSidebarCollapsed.value = !isRightSidebarCollapsed.value;
    };

    const toggleProfileDropdown = () => {
      showProfileDropdown.value = !showProfileDropdown.value;
    };

    const selectFile = (filename) => {
      activeFile.value = filename;
    };

    const getFileIcon = (filename) => {
      const ext = filename.split(".").pop();
      switch (ext) {
        case "js":
          return "fab fa-js";
        case "vue":
          return "fab fa-vuejs";
        case "html":
          return "fab fa-html5";
        case "css":
          return "fab fa-css3-alt";
        case "py":
          return "fab fa-python";
        case "java":
          return "fab fa-java";
        default:
          return "fas fa-file";
      }
    };

    const logout = () => {
      // Placeholder for logout functionality
      console.log("Logout");
      showProfileDropdown.value = false;
    };

    return {
      isDarkMode,
      isLeftSidebarCollapsed,
      isRightSidebarCollapsed,
      activeRightSidebarTab,
      activeRightSidebarTabTitle,
      rightSidebarTabs,
      showProfileDropdown,
      files,
      activeFile,
      toggleDarkMode,
      toggleLeftSidebar,
      toggleRightSidebar,
      toggleProfileDropdown,
      selectFile,
      getFileIcon,
      logout,
    };
  },
};
</script>

<style>
/* Variables for Light Theme */
:root {
--light-bg: #f6f8fa;
--light-sidebar-bg: #e8eaed;
--light-text: #24292e;
--light-border: #d1d5da;
--accent-color: #800000; /* Aggie Maroon */
--light-hover: #f0f2f5;
--light-active: #e0e4e8;
--light-button: #fafbfc;
}

/* Variables for Dark Theme */
.dark {
--dark-bg: #181818;
--dark-sidebar-bg: #282828;
--dark-text: #f0f4f8;
--dark-border: #444;
--dark-hover: #404040;
--dark-active: #505050;
--dark-button: #383838;
}

/* Shared Styles */
body,
html {
margin: 0;
padding: 0;
box-sizing: border-box;
font-family: "Inter", "Segoe UI", Tahoma, Geneva, Verdana, sans-serif;
height: 100%;
overflow: hidden;
}

.app {
display: flex;
flex-direction: column;
min-height: 100vh;
}

/* Light Mode Styles */
.light {
background-color: var(--light-bg);
color: var(--light-text);
}

.light .app-header {
background-color: var(--accent-color);
color: #fff;
}

.light .sidebar {
background-color: var(--light-sidebar-bg);
border-right: 1px solid var(--light-border);
}

.light .sidebar-header {
border-bottom: 1px solid var(--light-border);
}

.light .file-list li {
color: var(--light-text);
}

.light .file-list li.active,
.light .file-list li:hover {
background-color: var(--light-hover);
}

.light .right-sidebar {
background-color: var(--light-sidebar-bg);
border-left: 1px solid var(--light-border);
}

.light .terminal,
.light .chat {
background-color: var(--light-bg);
color: var(--light-text);
border: 1px solid var(--light-border);
}

.light .app-footer {
background-color: var(--accent-color);
color: #fff;
}

.light .styled-select {
background-color: var(--light-button);
color: var(--light-text);
border: 1px solid var(--light-border);
}

/* Dark Mode Styles */
.dark {
background-color: var(--dark-bg);
color: var(--dark-text);
}

.dark .app-header {
background-color: var(--dark-sidebar-bg);
}

.dark .sidebar {
background-color: var(--dark-sidebar-bg);
border-right: 1px solid var(--dark-border);
}

.dark .sidebar-header {
border-bottom: 1px solid var(--dark-border);
}

.dark .file-list li {
color: var(--dark-text);
}

.dark .file-list li.active,
.dark .file-list li:hover {
background-color: var(--dark-hover);
}

.dark .right-sidebar {
background-color: var(--dark-sidebar-bg);
border-left: 1px solid var(--dark-border);
}

.dark .terminal,
.dark .chat {
background-color: var(--dark-bg);
color: var(--dark-text);
border: 1px solid var(--dark-border);
}

.dark .app-footer {
background-color: var(--dark-sidebar-bg);
}

.dark .styled-select {
background-color: var(--dark-button);
color: var(--dark-text);
border: 1px solid var(--dark-border);
}

/* Header */
.app-header {
display: flex;
justify-content: space-between;
align-items: center;
padding: 0.5rem 1rem;
box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
z-index: 10; /* Header should stay on top */
}

.header-left {
display: flex;
align-items: center;
}

.hamburger-menu {
background: none;
border: none;
cursor: pointer;
margin-right: 1rem;
display: flex;
flex-direction: column;
justify-content: space-around;
width: 20px;
height: 20px;
}

.hamburger-menu span {
display: block;
width: 100%;
height: 2px;
background-color: var(--light-text);
transition: transform 0.3s, opacity 0.3s;
}

.dark .hamburger-menu span {
background-color: var(--dark-text);
}

.header-right {
display: flex;
align-items: center;
}

.app-header h1 {
margin: 0;
font-size: 1.4rem;
letter-spacing: 0.5px;
}

.user-profile {
height: 36px;
width: 36px;
background-color: #fff;
border-radius: 50%;
cursor: pointer;
position: relative;
display: flex;
align-items: center;
justify-content: center;
}

.dark .user-profile {
background-color: var(--dark-button);
}

.user-initials {
font-size: 0.9rem;
color: var(--accent-color);
}

.dark .user-initials {
color: var(--dark-text);
}

.profile-dropdown {
position: absolute;
top: 40px;
right: 0;
background-color: var(--light-bg);
border: 1px solid var(--light-border);
border-radius: 4px;
list-style: none;
padding: 0.5rem 0;
z-index: 20;
min-width: 120px;
box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
}

.dark .profile-dropdown {
background-color: var(--dark-bg);
border-color: var(--dark-border);
}

.profile-dropdown li {
padding: 0.5rem 1rem;
cursor: pointer;
white-space: nowrap;
}

.profile-dropdown li:hover {
background-color: var(--light-hover);
}

.dark .profile-dropdown li:hover {
background-color: var(--dark-hover);
}

.theme-toggle {
background: none;
border: none;
font-size: 1.5rem;
cursor: pointer;
margin-right: 1rem;
}

/* Layout */
.app-layout {
display: flex;
flex: 1;
}

/* Sidebar (File Management) */
.sidebar {
width: 240px;
transition: width 0.3s ease;
user-select: none;
position: relative;
z-index: 5;
}

.sidebar.collapsed {
width: 48px; /* Width when collapsed */
}

.sidebar.collapsed .sidebar-header h2 {
display: none; /* Hide title when collapsed */
}

.sidebar-header {
display: flex;
justify-content: space-between;
align-items: center;
padding: 0.5rem;
}

.sidebar-title {
flex-grow: 1; /* Allow the title to take up space */
text-align: center; /* Center the title */
margin-right: auto; /* Push the title to the left */
}

.collapse-button {
background: none;
border: none;
cursor: pointer;
font-size: 1rem;
color: var(--light-text);
margin-left: auto; /* Push the button to the right */
}

.dark .collapse-button {
color: var(--dark-text);
}

.file-list {
list-style: none;
padding: 0;
}

.file-list li {
display: flex;
align-items: center;
padding: 0.5rem 1rem;
cursor: pointer;
border-left: 3px solid transparent;
}

.file-list li i {
margin-right: 0.5rem;
width: 1.2em;
text-align: center;
}

.file-list li.active {
border-left-color: var(--accent-color);
background-color: var(--light-active);
}

.file-list li:hover {
background-color: var(--light-hover);
}

.dark .file-list li.active,
.dark .file-list li:hover {
background-color: var(--dark-hover);
}

/* Editor */
.editor-container {
flex: 1;
display: flex;
flex-direction: column;
padding: 1rem;
overflow: hidden;
}

/* Right Sidebar */
.right-sidebar {
width: 300px;
transition: width 0.3s ease;
user-select: none;
position: relative;
z-index: 5;
}

.right-sidebar.collapsed {
width: 48px;
}

.right-sidebar.collapsed .tabs,
.right-sidebar.collapsed .tab-content {
display: none;
}

.tabs {
display: flex;
border-bottom: 1px solid var(--light-border);
}

.dark .tabs {
border-bottom-color: var(--dark-border);
}

.tabs button {
background: none;
border: none;
border-bottom: 2px solid transparent;
padding: 0.5rem 1rem;
cursor: pointer;
font-weight: 600;
color: var(--light-text);
}

.dark .tabs button {
color: var(--dark-text);
}

.tabs button.active {
border-bottom-color: var(--accent-color);
}

.tab-content {
padding: 1rem;
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
border: 1px solid var(--light-border);
border-radius: 4px;
margin-bottom: 1rem;
max-height: 200px;
overflow-y: auto;
}

.dark .chat-box {
border-color: var(--dark-border);
}

.chat-box ul {
list-style: none;
padding: 0;
margin: 0;
}

.chat-box li {
margin-bottom: 0.5rem;
}

.chat-input-container {
display: flex;
align-items: center;
}

.chat-input-container input {
flex: 1;
padding: 0.5rem;
border: 1px solid var(--light-border);
border-radius: 4px;
margin-right: 0.5rem;
}

.dark .chat-input-container input {
border-color: var(--dark-border);
background-color: var(--dark-bg);
color: var(--dark-text);
}

.chat-input-container button {
padding: 0.5rem 1rem;
border: none;
border-radius: 4px;
background-color: var(--accent-color);
color: white;
cursor: pointer;
}

/* Footer */
.app-footer {
text-align: center;
padding: 0.75rem 0;
}

/* Add Font Awesome styles for icons */
.fa-js,
.fa-vuejs,
.fa-html5,
.fa-css3-alt,
.fa-python,
.fa-java,
.fa-file {
color: var(--accent-color); 
}

.dark .fa-js,
.dark .fa-vuejs,
.dark .fa-html5,
.dark .fa-css3-alt,
.dark .fa-python,
.dark .fa-java,
.dark .fa-file {
color: var(--dark-text);
}

.fa-bug {
color: red; /* Debug icon  color*/
}
</style>