<template>
  <div class="chat-app">
    <!-- User sidebar -->
    <aside class="user-sidebar">
      <div class="sidebar-header">
        <h2>Chat App</h2>
        <div class="user-badge" v-if="currentUser">
          <span class="user-avatar">👤</span>
          <span class="username">{{ currentUser }}</span>
        </div>
      </div>
      
      <div class="user-list">
        <h3 class="section-title">Active Users</h3>
        <div 
          v-for="(user, index) in users" 
          :key="index"
          class="user-item"
          :class="{
            'active': user.username === selectedUser,
            'current-user': user.username === currentUser
          }"
          @click="user.username !== currentUser && selectUser(user.username)"
        >
          <span class="user-avatar">{{ user.username[0].toUpperCase() }}</span>
          <span class="user-name">{{ user.username }}</span>
          <span v-if="user.username === currentUser" class="you-badge">you</span>
        </div>
      </div>
      
      <button 
        v-if="currentUser" 
        @click="logout" 
        class="logout-btn"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
          <polyline points="16 17 21 12 16 7"></polyline>
          <line x1="21" y1="12" x2="9" y2="12"></line>
        </svg>
        Logout
      </button>
    </aside>

    <!-- Main chat area -->
    <main class="chat-area">
      <div class="chat-header">
        <template v-if="!currentUser">
          <div class="status-message">
            <span class="status-icon">👤</span>
            <span>Please log in to start chatting</span>
          </div>
        </template>
        <template v-else-if="!selectedUser">
          <div class="status-message">
            <span class="status-icon">💬</span>
            <span>Select a user to start chatting</span>
          </div>
        </template>
        <template v-else>
          <div class="chat-partner">
            <div class="partner-avatar">{{ selectedUser[0].toUpperCase() }}</div>
            <div class="partner-info">
              <div class="partner-name">{{ selectedUser }}</div>
              <div class="chat-status">Online</div>
            </div>
          </div>
        </template>
      </div>

      <div class="messages-container">
        <div 
          v-for="(msg, index) in messages" 
          :key="index"
          class="message"
          :class="{
            'outgoing': msg.author === currentUser,
            'incoming': msg.author !== currentUser
          }"
        >
          <div class="message-avatar" v-if="msg.author !== currentUser">
            {{ msg.author[0].toUpperCase() }}
          </div>
          <div class="message-content">
            <div class="message-sender" v-if="msg.author !== currentUser">
              {{ msg.author }}
            </div>
            <div class="message-text">{{ msg.content }}</div>
            <div class="message-time">{{ formatTime(msg.timestamp) }}</div>
          </div>
        </div>
      </div>

      <form 
        v-if="selectedUser" 
        @submit.prevent="sendMessage"
        class="message-input-area"
      >
        <input 
          v-model="content" 
          placeholder="Type your message..." 
          required
          class="message-input"
        />
        <button type="submit" class="send-button">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="22" y1="2" x2="11" y2="13"></line>
            <polygon points="22 2 15 22 11 13 2 9 22 2"></polygon>
          </svg>
        </button>
      </form>
    </main>
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
  messages.value = []
}

const sendMessage = () => {
  if (!currentUser.value || !selectedUser.value || !content.value || !socket) return;

  const message = {
    author: currentUser.value,
    recipient: selectedUser.value,
    content: content.value,
    timestamp: new Date().toISOString(),
    isLocal: true
  };

  if (message.recipient !== currentUser.value) {
    messages.value.push(message);
  }

  socket.send(JSON.stringify(message));
  content.value = '';
};

const formatTime = (timestamp) => {
  if (!timestamp) return ''
  const date = new Date(timestamp)
  return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

const logout = () => {
  localStorage.removeItem('username')
  currentUser.value = ''
  selectedUser.value = ''
  router.push('/')
}

onMounted(() => {
  const savedUsername = localStorage.getItem('username')
  if (savedUsername) {
    currentUser.value = savedUsername
  }

  fetchUsers()

  const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
  const host = window.location.hostname
  const port = 8080

  socket = new WebSocket(`${protocol}://${host}:${port}/messenger/api/ws`)

  socket.addEventListener('open', () => {
    console.log('WebSocket connected')
  })

  socket.addEventListener('error', (error) => {
    console.error('WebSocket error:', error)
  })

  socket.addEventListener('close', () => {
    console.log('WebSocket disconnected')
  })

  socket.addEventListener('message', event => {
  const msg = JSON.parse(event.data);
  
  if (msg.author === currentUser.value && !msg.isLocal) {
    return;
  }
  
  if (!msg.timestamp) {
    msg.timestamp = new Date().toISOString();
  }
  
  messages.value.push(msg);
  })
})

onBeforeUnmount(() => {
  if (socket) socket.close()
})
</script>

<style scoped>
/* Base styles */
:root {
  --primary-color: #4361ee;
  --primary-light: #e0e7ff;
  --secondary-color: #3f37c9;
  --text-color: #2b2d42;
  --text-light: #8d99ae;
  --bg-color: #f8f9fa;
  --card-bg: #ffffff;
  --border-color: #e9ecef;
  --outgoing-bg: #4361ee;
  --outgoing-text: #ffffff;
  --incoming-bg: #f1f3f5;
  --incoming-text: #2b2d42;
  --error-color: #ef233c;
  --success-color: #2ec4b6;
}

* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

body {
  font-family: 'Segoe UI', system-ui, -apple-system, sans-serif;
  color: var(--text-color);
  background-color: var(--bg-color);
}

/* Chat app layout */
.chat-app {
  display: flex;
  height: 90vh;
  width: 90%; /* Занимает 90% ширины родителя */
  min-width: 1200px; /* Минимальная ширина */
  max-width: 1800px; /* Максимальная ширина */
  margin: 0 auto;
  background-color: var(--card-bg);
  box-shadow: 0 0 20px rgba(255, 0, 255, 0.2);
  border-radius: 12px;
  overflow: hidden;
}

/* User sidebar */
.user-sidebar {
  width: 350px;
  background-color: var(--card-bg);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  padding: 1.5rem 1rem;
}

.sidebar-header {
  padding-bottom: 1.5rem;
  border-bottom: 1px solid var(--border-color);
  margin-bottom: 1.5rem;
}

.sidebar-header h2 {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--primary-color);
  margin-bottom: 1rem;
}

.user-badge {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.5rem;
  background-color: var(--primary-light);
  border-radius: 8px;
}

.user-avatar {
  width: 32px;
  height: 32px;
  background-color: var(--primary-color);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 0.875rem;
}

.username {
  font-weight: 600;
  font-size: 0.875rem;
}

.user-list {
  flex: 1;
  overflow-y: auto;
  margin-bottom: 1rem;
}

.section-title {
  font-size: 0.75rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--text-light);
  margin-bottom: 1rem;
}

.user-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-bottom: 0.25rem;
}

.user-item:hover {
  background-color: var(--primary-light);
}

.user-item.active {
  background-color: var(--primary-light);
  font-weight: 600;
}

.user-item.current-user {
  background-color: #0d4d8d;
  cursor: default;
}

.user-name {
  flex: 1;
  font-size: 0.875rem;
}

.you-badge {
  font-size: 0.75rem;
  background-color: var(--primary-light);
  color: var(--primary-color);
  padding: 0.25rem 0.5rem;
  border-radius: 12px;
}

.logout-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  padding: 0.75rem;
  background-color: var(--error-color);
  color: white;
  border: none;
  border-radius: 8px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.logout-btn:hover {
  background-color: #d90429;
}

/* Chat area */
.chat-area {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.chat-header {
  padding: 1.5rem;
  border-bottom: 1px solid var(--border-color);
}

.status-message {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: var(--text-light);
}

.status-icon {
  font-size: 1.25rem;
}

.chat-partner {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.partner-avatar {
  width: 40px;
  height: 40px;
  background-color: var(--primary-color);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 1rem;
}

.partner-info {
  display: flex;
  flex-direction: column;
}

.partner-name {
  font-weight: 600;
}

.chat-status {
  font-size: 0.75rem;
  color: var(--success-color);
}

/* Messages container */
.messages-container {
  flex: 1;
  padding: 1.5rem;
  overflow-y: auto;
  background-color: #2e3031;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.message {
  display: flex;
  gap: 0.75rem;
  max-width: 70%;
}

.message.incoming {
  align-self: flex-start;
}

.message.outgoing {
  align-self: flex-end;
  flex-direction: row-reverse;
}

.message-avatar {
  width: 32px;
  height: 32px;
  background-color: var(--primary-color);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 0.875rem;
  flex-shrink: 0;
}

.message-content {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.message-sender {
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--text-light);
}

.message-text {
  padding: 0.75rem 1rem;
  border-radius: 12px;
  font-size: 0.875rem;
  line-height: 1.4;
}

.incoming .message-text {
  background-color: var(--incoming-bg);
  color: var(--incoming-text);
  border-top-left-radius: 4px;
}

.outgoing .message-text {
  background-color: var(--outgoing-bg);
  color: var(--outgoing-text);
  border-top-right-radius: 4px;
}

.message-time {
  font-size: 0.625rem;
  color: var(--text-light);
  text-align: right;
}

/* Message input area */
.message-input-area {
  padding: 1rem 1.5rem;
  border-top: 1px solid var(--border-color);
  display: flex;
  gap: 0.75rem;
  background-color: var(--card-bg);
}

.message-input {
  flex: 1;
  padding: 0.75rem 1rem;
  border: 1px solid var(--border-color);
  border-radius: 24px;
  font-size: 0.875rem;
  transition: all 0.2s ease;
}

.message-input:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 2px var(--primary-light);
}

.send-button {
  width: 48px;
  height: 48px;
  background-color: var(--primary-color);
  color: white;
  border: none;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
}

.send-button:hover {
  background-color: var(--secondary-color);
}

/* Scrollbar styling */
::-webkit-scrollbar {
  width: 6px;
}

::-webkit-scrollbar-track {
  background: var(--border-color);
  border-radius: 3px;
}

::-webkit-scrollbar-thumb {
  background: var(--text-light);
  border-radius: 3px;
}

::-webkit-scrollbar-thumb:hover {
  background: var(--primary-color);
}
</style>