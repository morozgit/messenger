<template>
  <div class="container">
    <h1>Welcome to the Chat App</h1>

    <div class="button-row">
      <button @click="goTo('signin')">Sign in</button>
      <button @click="goTo('signup')">Sign up</button>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

const goTo = (mode) => {
  router.push({ path: '/auth', query: { mode } })
}
</script>

<style scoped>
.container {
  padding: 2rem;
  text-align: center;
}

.button-row {
  margin-top: 1.5rem;
  display: flex;
  justify-content: center;
  gap: 1rem;
}

button {
  padding: 0.5rem 1rem;
  font-size: 1rem;
  cursor: pointer;
}
</style>
