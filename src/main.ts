import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'

import "@fontsource/lato/300.css"
import "@fontsource/lato/400.css"
// import "./scss/custom.scss"
// // import "bootstrap/dist/css/bootstrap.min.css"
// import "bootstrap"
import "./scss/main.scss"

createApp(App).use(store).use(router).mount('#app')
