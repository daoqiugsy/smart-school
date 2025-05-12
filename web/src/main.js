import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import api from './services/api' // 导入 axios 实例

const app = createApp(App)

app.config.globalProperties.$axios = api // 将 axios 实例挂载到全局属性

app.use(router).mount('#app')
