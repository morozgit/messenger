<template>
  <div class="auth-container">
    <h1>{{ isSignUp ? 'Sign up' : 'Sign in' }}</h1>

    <form @submit.prevent="submit">
      <input v-model="username" placeholder="Username" required />
      <input v-model="password" placeholder="Password" type="password" required />

      <input
        v-if="isSignUp"
        v-model="confirmPassword"
        placeholder="Confirm Password"
        type="password"
        required
      />

      <button type="submit">{{ isSignUp ? 'Create account' : 'Log in' }}</button>
    </form>

    <p class="switch" @click="isSignUp = !isSignUp">
      {{ isSignUp ? 'Already have an account? Sign in' : 'No account? Sign up' }}
    </p>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

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
  console.log('Submitting', {
  username: username.value,
  password: password.value,
  confirmPassword: confirmPassword.value,
  isSignUp: isSignUp.value,
})

  try {
    if (isSignUp.value) {
      await axios.post('http://localhost:8080/add_users', {
        username: username.value,
        password: password.value,
      })
    } else {
      await axios.post('http://localhost:8080/login', {
        username: username.value,
        password: password.value,
      }, {
        headers: { 'Content-Type': 'application/json' }
      })
    }

    localStorage.setItem('username', username.value)
    router.push('/chat')
  } catch (err) {
    console.error('Error:', err)
    alert(err.response?.data?.error || err.message || 'Unknown error')
  }
}
</script>

<style scoped>
.auth-container {
  max-width: 400px;
  margin: 100px auto;
  text-align: center;
}

input {
  width: 100%;
  padding: 0.5rem;
  margin-top: 1rem;
}

button {
  margin-top: 1rem;
  padding: 0.5rem 1rem;
}

.switch {
  margin-top: 1rem;
  color: #f9f9f9;
  text-decoration: underline;
  cursor: pointer;
}
</style>
