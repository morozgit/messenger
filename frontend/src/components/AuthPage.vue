<template>
  <div class="auth-wrapper">
    <div class="auth-card">
      <div class="logo">🔺 Messenger</div>
      <h2>{{ isSignUp ? 'Create Account' : 'Welcome Back' }}</h2>

      <input type="text" placeholder="Username" v-model="username" required />
      <input type="password" placeholder="Password" v-model="password" required />

      <input
        v-if="isSignUp"
        type="password"
        placeholder="Confirm Password"
        v-model="confirmPassword"
        required
      />

      <button class="login-btn" @click.prevent="submit">
        {{ isSignUp ? 'Sign Up' : 'Log In' }}
      </button>

      <div class="divider">OR</div>

      <div class="bottom-links">
        <span @click="isSignUp = !isSignUp" style="cursor: pointer;">
          {{ isSignUp ? 'Already have an account? Sign in' : 'No account? Sign up' }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'

const router = useRouter()
const route = useRoute()

const isSignUp = ref(false)
const username = ref('')
const password = ref('')
const confirmPassword = ref('')

onMounted(() => {
  isSignUp.value = route.query.mode === 'signup'
})

const submit = async () => {
  if (!username.value || !password.value) return

  if (isSignUp.value && password.value !== confirmPassword.value) {
    alert('Passwords do not match')
    return
  }

  try {
    if (isSignUp.value) {
      await axios.post('/messenger/api/add_users', {
        username: username.value,
        password: password.value,
      })
    } else {
      await axios.post('/messenger/api/login', {
        username: username.value,
        password: password.value,
      }, {
        headers: { 'Content-Type': 'application/json' }
      })
    }

    localStorage.setItem('username', username.value)
    router.push('/chat')
  } catch (err) {
    console.error(err)
    alert(err.response?.data?.error || err.message || 'Unknown error')
  }
}
</script>

<style scoped>
html, body {
  margin: 0;
  padding: 0;
  height: 100%;
}


.auth-card {
  background-color: rgba(0, 0, 0, 0.65);
  padding: 2rem;
  border-radius: 20px;
  width: 100%;
  max-width: 360px;
  box-shadow: 0 0 20px rgba(255, 0, 255, 0.2);
  text-align: center;
  backdrop-filter: blur(10px);
  color: white;
}

.logo {
  font-weight: bold;
  font-size: 1.4rem;
  margin-bottom: 1rem;
}

input[type="text"],
input[type="password"] {
  width: 93%;
  padding: 0.7rem;
  margin-bottom: 1rem;
  border: none;
  border-radius: 8px;
  background-color: #2a2a2a;
  color: white;
}

.remember {
  display: flex;
  align-items: center;
  justify-content: start;
  font-size: 0.9rem;
  margin-bottom: 1rem;
}

.remember input {
  margin-right: 0.5rem;
}

.login-btn {
  width: 100%;
  padding: 0.7rem;
  background-color: #ff4081;
  border: none;
  border-radius: 8px;
  color: white;
  font-weight: 600;
  cursor: pointer;
  margin-bottom: 1rem;
}

.divider {
  display: flex;
  align-items: center;
  text-align: center;
  color: #aaa;
  margin: 1rem 0;
}

.divider::before,
.divider::after {
  content: '';
  flex: 1;
  border-bottom: 1px solid #444;
}

.divider:not(:empty)::before {
  margin-right: 0.5em;
}

.divider:not(:empty)::after {
  margin-left: 0.5em;
}

.metamask-btn {
  background-color: white;
  color: black;
  padding: 0.7rem;
  width: 100%;
  border-radius: 8px;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  cursor: pointer;
}

.bottom-links {
  font-size: 0.85rem;
  margin-top: 1rem;
}

.bottom-links a {
  color: #00bfff;
  text-decoration: none;
  margin: 0 0.3rem;
  cursor: pointer;
}
</style>

