<template>
  <div class="container">
    <div class="header row">
      <div class="col-1-of-2">
        <div class="page-title">
          <router-link to="/">
            <svg class="svganim" width="100" height="100" version="1.1" viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
              <g fill="none">
                <path d="m11.6 61.22a40 40 0 0 1 17.03-45.03 40 40 0 0 1 47.98 3.962 40 40 0 0 1 9.415 47.22" stop-color="#000000" pathLength="100"/>
                <path d="m10.43 60.17c10.77 0.2854 21.55 0.5708 30.3 3.278s15.48 7.836 21.47 10.72 11.23 3.523 16.48 4.161" pathLength="100"/>
                <path d="m11.91 60.4c10.46-9.06 20.91-18.12 30.16-21.71s17.27-1.725 25.29 0.1457c-4.788 6.39-9.576 12.78-8.469 17.16 1.107 4.383 8.107 6.758 13.86 8.177 5.756 1.42 10.27 1.884 14.78 2.348" pathLength="100"/>
                <path d="m14.87 70.02c6.739-0.197 13.48-0.394 20.27 2.093 6.792 2.487 13.64 7.658 18.93 10.55s9.036 3.513 12.78 4.132" pathLength="100"/>
                <path d="m21.96 79.43c4.237 0.6038 8.474 1.208 11.75 2.557 3.274 1.349 5.589 3.444 7.766 4.842 2.177 1.397 4.213 2.096 6.251 2.795" pathLength="100"/>
              </g>
            </svg>
          </router-link>
          <router-link to="/" class="heading-primary">Uptime Shark</router-link>
        </div>
      </div>
      <div class="col-1-of-2 nav">
          <router-link :to="{ name: 'Status'}">Status</router-link>
          <router-link :to="{ name: 'Edit'}">Edit</router-link>
      </div>
    </div>

    <h1 v-if="this.$store.state.error">Oops! An Error Occurred... <br> {{this.$store.state.errorMessage}}</h1>
    <router-view v-else/>
  </div>
</template>

<script setup lang="ts">
  import { onBeforeMount } from 'vue'
  import store from './store'
  onBeforeMount(() => {
    store.dispatch('fetchData')
    store.dispatch('checkLogin')
    
    setInterval(function() {
      store.dispatch('fetchData')
    }, 5000);

    setInterval(function() {
      if (store.state.authenticated)
        store.dispatch('checkLogin')
    }, 10000);

    // You can clear a periodic function by uncommenting:
    // clearInterval(intervalId);    
  })
</script>
