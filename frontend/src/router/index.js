import { createRouter, createWebHistory } from 'vue-router'
import StartPage from '../components/StartPage.vue'
import ChatBox from '../components/ChatBox.vue'
import AuthPage from '../components/AuthPage.vue'

const routes = [
  { path: '/messenger/', component: StartPage },
  { path: '/messenger/chat', component: ChatBox },
  { path: '/messenger/auth', component: AuthPage }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
