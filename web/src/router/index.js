import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/Login.vue'
import AIChatView from '../views/AIChat.vue'

const routes = [
  {
    path: '/',
    name: 'ai-chat',
    component: AIChatView,
    meta: { requiresAuth: true }
  },
  {
    path: '/login',
    name: 'login',
    component: LoginView
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

// 路由守卫
router.beforeEach((to, from, next) => {
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
  const token = localStorage.getItem('token')
  
  if (requiresAuth && !token) {
    next('/login')
  } else {
    next()
  }
})

export default router