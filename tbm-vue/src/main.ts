/**
 * main.ts
 *
 * Bootstraps Vuetify and other plugins then mounts the App`
 */

// Plugins
import { registerPlugins } from '@/plugins'

// Components
import App from './App.vue'

// Store
import { store, key } from "@/store";

// Composables
import { createApp } from 'vue'

const app = createApp(App)
app.use(store, key);
registerPlugins(app)

app.mount('#app')
