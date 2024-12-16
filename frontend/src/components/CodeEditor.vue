<!-- CodeEditor.vue -->
<template>
  <div class="code-editor-container">
    <!-- Toolbar -->
    <div :class="['editor-toolbar', isDarkMode ? 'dark' : 'light']">
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
          <option v-for="size in fontSizes" :key="size" :value="size">{{ size }}</option>
        </select>
      </div>

      <!-- Language Selector -->
      <div class="toolbar-item">
        <label for="language-select">Language:</label>
        <select id="language-select" v-model="language" class="styled-select">
          <option v-for="lang in languages" :key="lang.value" :value="lang.value">{{ lang.label }}</option>
        </select>
      </div>

      <!-- Run Button -->
      <button class="run-button" @click="runCode">Run</button>
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
import { ref, watch } from "vue";
import { VueMonacoEditor } from "@guolao/vue-monaco-editor";

export default {
  components: {
    VueMonacoEditor,
  },
  props: {
    isDarkMode: {
      type: Boolean,
      required: true,
    },
  },
  setup(props) {
    const code = ref(`"""
Welcome to AggieCode Collaborative IDE
This is Python code!
"""

def hello_aggie():
    print("Hello, AggieCode!")

hello_aggie()`);

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

<style>
.code-editor-container {
  display: flex;
  flex-direction: column;
  height: 100%;
}

/* Toolbar */
.editor-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.5rem 1rem;
  border-bottom: 1px solid var(--border-color);
  transition: background-color 0.3s, color 0.3s;
}

/* The background and text colors for the toolbar now depend on the parent's theme classes */
.editor-toolbar.light {
  background-color: #f9f9f9;
  color: #333;
}

.editor-toolbar.dark {
  background-color: #2a2a2a;
  color: #f5f5f5;
}

/* Toolbar Items */
.toolbar-item {
  display: flex;
  align-items: center;
  margin-right: 1.5rem;
}

.toolbar-item label {
  margin-right: 0.5rem;
  font-size: 0.9rem;
  font-weight: bold;
}

/* Styled Dropdowns inherit global vars */
.styled-select {
  padding: 0.4rem 0.6rem;
  font-size: 0.9rem;
  border: 1px solid var(--border-color);
  border-radius: 5px;
  background-color: var(--dropdown-bg);
  color: var(--dropdown-color);
  outline: none;
  transition: background-color 0.3s, color 0.3s, border-color 0.3s;
}

.styled-select:hover {
  border-color: #888;
}

.styled-select:focus {
  border-color: #00bcd4;
}

/* Run Button */
.run-button {
  padding: 0.5rem 1rem;
  background-color: #00bcd4;
  color: #fff;
  border: none;
  border-radius: 5px;
  font-size: 0.9rem;
  font-weight: bold;
  cursor: pointer;
  transition: background-color 0.3s, transform 0.2s;
}

.run-button:hover {
  background-color: #0097a7;
}

.run-button:active {
  transform: scale(0.95);
}

/* Monaco Editor */
.monaco-editor-wrapper {
  flex: 1;
  width: 100%;
}
</style>
