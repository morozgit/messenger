<template>
  <div class="chat-container">
    <!-- Список пользователей -->
    <div class="users">
      <h3>Users</h3>
      <ul>
        <li
          v-for="(user, index) in users"
          :key="index"
          :class="{ active: user.username === selectedUser, you: user.username === currentUser }"
          @click="user.username !== currentUser && selectUser(user.username)"
        >
          {{ user.username }}
        </li>
      </ul>
      <button v-if="currentUser" @click="logout" class="logout-button">Logout</button>
    </div>

    <!-- Статус текущего чата -->
    <div class="status">
      <div v-if="!currentUser">👤 Please log in</div>
      <div v-else-if="!selectedUser">💬 Select someone to chat with</div>
      <div v-else>
        You are: <strong>{{ currentUser }}</strong><br />
        Chatting with: <strong>{{ selectedUser }}</strong>
      </div>
    </div>

    <!-- Область сообщений -->
    <div class="chat">
      <div class="messages">
        <div v-for="(msg, index) in messages" :key="index">
          <strong>{{ msg.author }}:</strong> {{ msg.content }}
        </div>
      </div>

      <form v-if="selectedUser" @submit.prevent="sendMessage">
        <input v-model="content" placeholder="Type message..." required />
        <button type="submit">Send</button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'

const router = useRouter()

const messages = ref([])
const users = ref([])
const content = ref('')
const selectedUser = ref('')
const currentUser = ref('')
let socket = null

const fetchUsers = async () => {
  try {
    const res = await axios.get('/messenger/api/users')
    users.value = res.data
  } catch (err) {
    console.error('Failed to load users:', err)
  }
}

const selectUser = (name) => {
  if (!currentUser.value) {
    currentUser.value = name
  }
  selectedUser.value = name
}

const sendMessage = () => {
  if (!currentUser.value || !selectedUser.value || !content.value || !socket) {
    console.warn('Missing fields')
    return
  }

  socket.send(JSON.stringify({
    author: currentUser.value,
    recipient: selectedUser.value,
    content: content.value,
  }))

  content.value = ''
}

const logout = () => {
  localStorage.removeItem('username')
  currentUser.value = ''
  selectedUser.value = ''
  router.push('/auth')
}

onMounted(() => {
  const savedUsername = localStorage.getItem('username')
  if (savedUsername) {
    currentUser.value = savedUsername
  }

  fetchUsers()

  socket = new WebSocket('ws://ubuntuserver:8080/ws')
  socket.addEventListener('open', () => {
    console.log('WebSocket connected')
  })
  socket.addEventListener('message', event => {
    const msg = JSON.parse(event.data)
    messages.value.push(msg)
  })
})

onBeforeUnmount(() => {
  if (socket) socket.close()
})
</script>

<style scoped>
.chat-container {
  display: flex;
  gap: 1rem;
}

.users {
  width: 220px;
  border-right: 1px solid #ccc;
  padding-right: 1rem;
}

.chat {
  flex: 1;
}

.status {
  margin-bottom: 1rem;
}

.messages {
  max-height: 200px;
  overflow-y: auto;
  border: 1px solid #ccc;
  padding: 10px;
  margin-bottom: 1em;
}

form {
  display: flex;
  gap: 8px;
  margin-top: 1em;
}

input {
  flex: 1;
}

.active {
  font-weight: bold;
  background-color: #eef;
}

.you {
  color: green;
}

ul li {
  cursor: pointer;
  margin-bottom: 4px;
}

.logout-button {
  margin-top: 1rem;
  background-color: #f44336;
  color: white;
  border: none;
  padding: 6px 12px;
  cursor: pointer;
  border-radius: 4px;
}
</style>
