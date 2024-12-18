import { createApp } from 'vue';
import App from './App.vue';
import { install as VueMonacoEditorPlugin } from '@guolao/vue-monaco-editor';

/* import the fontawesome core */
import { library } from '@fortawesome/fontawesome-svg-core'

/* import font awesome icon component */
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

/* import icons */
import { faUserSecret } from '@fortawesome/free-solid-svg-icons'

/* add icons to the library */
library.add(faUserSecret)

const app = createApp(App);

app.use(VueMonacoEditorPlugin, {
  paths: {
    vs: 'https://cdn.jsdelivr.net/npm/monaco-editor@0.43.0/min/vs',
  },
});

app.component('font-awesome-icon', FontAwesomeIcon)

app.mount('#app');