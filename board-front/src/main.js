/**
 * main.js
 *
 * Bootstraps Vuetify and other plugins then mounts the App`
 */

// Plugins
import {registerPlugins} from '@/plugins'
import 'virtual:uno.css'
import {pl, zhHans} from 'vuetify/locale'
// Components
import App from './App.vue'
const vuetify = createVuetify({
  locale: {
    locale: 'zhHans',
    current: 'zhHans',
    messages: { zhHans, pl },
  },
})

// Composables
import {createApp} from 'vue'
import {createVuetify} from "vuetify";

const app = createApp(App)
registerPlugins(app)
app.use(vuetify)
app.mount('#app')

