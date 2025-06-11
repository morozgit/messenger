import { createRouter, createWebHistory } from 'vue-router'
import StartPage from '../components/StartPage.vue'
import ChatBox from '../components/ChatBox.vue'
import AuthPage from '../components/AuthPage.vue'

const routes = [
  { path: '/', component: StartPage },
  { path: '/chat', component: ChatBox },
  { path: '/auth', component: AuthPage }
]

const router = createRouter({
  history: createWebHistory('/messenger/'),
  routes,
})

export default router
