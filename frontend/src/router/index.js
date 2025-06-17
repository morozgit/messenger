import { createRouter, createWebHistory } from 'vue-router'
import ChatPage from '../components/ChatPage.vue'
import AuthPage from '../components/AuthPage.vue'

const routes = [
  { path: '/', component: AuthPage },
  { path: '/chat', component: ChatPage },
]

const router = createRouter({
  history: createWebHistory('/messenger/'),
  routes,
})

export default router
