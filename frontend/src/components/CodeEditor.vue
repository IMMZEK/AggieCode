<template>
  <div class="code-editor-container">
    <!-- Toolbar -->
    <div
      :class="[
        'editor-toolbar',
        isDarkMode ? 'dark' : 'light',
        isDarkMode ? '' : 'light-toolbar',
      ]"
    >
      <!-- Theme Selector -->
      <div class="toolbar-item">
        <label for="theme-select">Theme:</label>
        <select id="theme-select" v-model="theme" class="styled-select">
          <option value="vs">Light</option>
          <option value="vs-dark">Dark</option>
        </select>
      </div>

      <!-- Font Size Selector -->
      <div class="toolbar-item">
        <label for="font-size-select">Font Size:</label>
        <select id="font-size-select" v-model="fontSize" class="styled-select">
          <option v-for="size in fontSizes" :key="size" :value="size">{{
            size
          }}</option>
        </select>
      </div>

      <!-- Language Selector -->
      <div class="toolbar-item">
        <label for="language-select">Language:</label>
        <select id="language-select" v-model="language" class="styled-select">
          <option
            v-for="lang in languages"
            :key="lang.value"
            :value="lang.value"
            >{{ lang.label }}</option
          >
        </select>
      </div>

      <!-- Run Button -->
      <div class="toolbar-item run-button-container">
        <button class="run-button" @click="runCode">
          <i class="fas fa-play"></i> Run
        </button>
      </div>
    </div>

    <!-- Monaco Editor -->
    <div class="monaco-editor-wrapper">
      <vue-monaco-editor
        v-model:value="code"
        :language="language"
        :theme="theme"
        :options="editorOptions"
        @mount="onEditorMounted"
      />
    </div>
  </div>
</template>

<script>
import { ref, watch, onMounted } from "vue";
import { VueMonacoEditor } from "@guolao/vue-monaco-editor";
import { faPlay } from "@fortawesome/free-solid-svg-icons";
import { library } from "@fortawesome/fontawesome-svg-core";

library.add(faPlay);

export default {
  components: {
    VueMonacoEditor,
  },
  props: {
    isDarkMode: {
      type: Boolean,
      required: true,
    },
    activeFile: {
      type: String,
      default: null,
    },
  },
  setup(props) {
    const code = ref(``);

    const languages = ref([
      { value: "python", label: "Python" },
      { value: "javascript", label: "JavaScript" },
      { value: "java", label: "Java" },
      { value: "cpp", label: "C++" },
      { value: "html", label: "HTML" },
      { value: "css", label: "CSS" },
      { value: "json", label: "JSON" },
    ]);

    const theme = ref(props.isDarkMode ? "vs-dark" : "vs");
    watch(() => props.isDarkMode, (newValue) => {
      theme.value = newValue ? "vs-dark" : "vs";
    });

    const fontSizes = ref([10, 12, 14, 16, 18, 20, 24, 28]);
    const fontSize = ref(14);

    const language = ref("python");

    const editorOptions = ref({
      automaticLayout: true,
      fontSize: fontSize.value,
      minimap: { enabled: true },
      scrollBeyondLastLine: false,
    });

    watch(fontSize, (newFontSize) => {
      editorOptions.value = { ...editorOptions.value, fontSize: newFontSize };
    });

    const onEditorMounted = (editor) => {
      console.log("Monaco Editor is ready:", editor);
    };

    const runCode = () => {
      console.log("Code is being executed...");
      console.log(code.value);
    };

    onMounted(() => {
  if (props.activeFile) {
    // Update the default code based on the active file
    switch (props.activeFile) {
      case "main.js":
        code.value = `// Welcome to AggieCode Collaborative IDE\n// This is JavaScript code!\n\nconsole.log("Hello, AggieCode!");`;
        language.value = "javascript";
        break;
      case "App.vue":
      case "CodeEditor.vue":
        code.value = `<!-- ${props.activeFile} -->\n<!-- Your Vue.js code here -->`;
        language.value = "html";
        break;
      case "index.html":
        code.value = `<!DOCTYPE html>\n<html lang="en">\n<head>\n    <meta charset="UTF-8">\n    <meta name="viewport" content="width=device-width, initial-scale=1.0">\n    <title>AggieCode - Collaborative IDE</title>\n</head>\n<body>\n    <div id="app"></div>\n    <script type="module" src="/src/main.js"><\/script>\n</body>\n</html>`;
        language.value = "html";
        break;
      case "style.css":
        code.value = `/* ${props.activeFile} */\n/* Your CSS code here */`;
        language.value = "css";
        break;
      default:
        code.value = `// Welcome to AggieCode Collaborative IDE\n// This is a new file!`;
        // TODO: Default to Python with example code.
        language.value = "javascript"; // Default to JavaScript
        break;
    }
  }
});

watch(() => props.activeFile, (newFile) => {
  if (newFile) {
    // Update the code and language based on the active file
    switch (newFile) {
      case "main.js":
        code.value = `// Welcome to AggieCode Collaborative IDE\n// This is JavaScript code!\n\nconsole.log("Hello, AggieCode!");`;
        language.value = "javascript";
        break;
      case "App.vue":
      case "CodeEditor.vue":
        code.value = `<!-- ${newFile} -->\n<!-- Your Vue.js code here -->`;
        language.value = "html";
        break;
      case "index.html":
        code.value = `<!DOCTYPE html>\n<html lang="en">\n<head>\n    <meta charset="UTF-8">\n    <meta name="viewport" content="width=device-width, initial-scale=1.0">\n    <title>AggieCode - Collaborative IDE</title>\n</head>\n<body>\n    <div id="app"></div>\n    <script type="module" src="/src/main.js"><\/script>\n</body>\n</html>`;
        language.value = "html";
        break;
      case "style.css":
        code.value = `/* ${newFile} */\n/* Your CSS code here */`;
        language.value = "css";
        break;
      default:
        code.value = `// Welcome to AggieCode Collaborative IDE\n// This is a new file!`;
        language.value = "javascript"; // Default to JavaScript
        break;
    }
  }
});

    return {
      code,
      theme,
      fontSize,
      fontSizes,
      language,
      languages,
      editorOptions,
      onEditorMounted,
      runCode,
    };
  },
};
</script>

<style scoped>
.code-editor-container {
  display: flex;
  flex-direction: column;
  height: 100%;
}

/* Toolbar */
.editor-toolbar {
  display: flex;
  align-items: center;
  justify-content: flex-start; /* Align items to the start */
  padding: 0.5rem 1rem;
  border-bottom: 1px solid var(--light-border);
  transition: background-color 0.3s, color 0.3s;
  flex-wrap: wrap; /* Allow items to wrap */
}

.editor-toolbar.light-toolbar {
  border-bottom: 1px solid var(--light-border);
}

.dark .editor-toolbar {
  border-bottom-color: var(--dark-border);
}

/* Toolbar Items */
.toolbar-item {
  display: flex;
  align-items: center;
  margin-right: 1.5rem;
  margin-bottom: 0.5rem; /* margin-bottom for wrapping */
}

.toolbar-item label {
  margin-right: 0.5rem;
  font-size: 0.9rem;
  font-weight: bold;
  color: var(--light-text);
}

.dark .toolbar-item label {
  color: var(--dark-text);
}

/* Styled Dropdowns */
.styled-select {
  padding: 0.4rem 0.6rem;
  font-size: 0.9rem;
  border: 1px solid var(--light-border);
  border-radius: 5px;
  background-color: var(--light-button);
  color: var(--light-text);
  outline: none;
  transition: background-color 0.3s, color 0.3s, border-color 0.3s;
}

.dark .styled-select {
  background-color: var(--dark-button);
  color: var(--dark-text);
  border-color: var(--dark-border);
}

.styled-select:hover {
  border-color: #888;
}

.styled-select:focus {
  border-color: var(--accent-color);
}

/* Run Button */
.run-button {
  padding: 0.4rem 0.8rem;
  background-color: var(--accent-color); /* Aggie Maroon */
  color: #fff;
  border: none;
  border-radius: 5px;
  font-size: 0.9rem;
  font-weight: bold;
  cursor: pointer;
  transition: background-color 0.3s, transform 0.2s;
  display: flex; /* Use flexbox */
  align-items: center; /* Align items vertically */
  gap: 0.5rem; /* Space between icon and text */
}

.run-button-container {
  margin-left: auto;
}

.dark .run-button {
  background-color: var(--dark-button);
}

.run-button:hover {
  background-color: #5a0000; /* Darker shade of Aggie Maroon */
}

.dark .run-button:hover {
  background-color: #505050;
}

.run-button:active {
  transform: scale(0.95);
}

/* Run Button Icon */
.run-button i {
  font-size: 1rem;
}

/* Monaco Editor */
.monaco-editor-wrapper {
  flex: 1;
  width: 100%;
  border: none; /* Remove any default borders */
}
</style>