/**
 * plugins/vuetify.ts
 *
 * Framework documentation: https://vuetifyjs.com`
 */

// Styles
import '@mdi/font/css/materialdesignicons.css'
import 'vuetify/styles'

// Composables
import { createVuetify } from 'vuetify'

// https://vuetifyjs.com/en/introduction/why-vuetify/#feature-guides
export default createVuetify({
  theme: {
    defaultTheme: "default",
    themes: {
      default: {
        dark: false,
        colors: {
          primary: "rgb(45, 97, 255)",
          background: "rgb(210, 216, 222)",
          oddRow: "rgb(226, 229, 231)",
          surface: "rgb(235, 235, 235)",
          appBar: "rgb(0, 33, 72)"
        }
      }
    }
  }
})
