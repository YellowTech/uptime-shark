<template>
  <div class="heading-secondary">Login</div>
  <div>
    <div class="row u-mbsmall">
      <div class="input-group">
        <input type="password" v-model="password" placeholder="Password" v-on:keyup.enter="sendLogin"/>
      </div>
      <div>
        <a class="btn positive" href="#" @click="sendLogin()">Login</a>
      </div>
      <div class="heading-tertiary">{{ message }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref } from 'vue'
  import store  from '../store'


  const password = ref<string>("")
  const message = ref<string>("")

  function sendLogin() {
    message.value = "Loading"
    fetch(store.state.apiDomain + "/api/auth/login", {
      method: 'POST',
      headers: {
        'Accept': 'application/json, text/plain, */*',
        'Content-Type': 'application/json',
      },
      credentials: "include",
      body: JSON.stringify({password: password.value})
    }).then(res => res.json())
      .then(data => {
        if (data.error) {
          message.value = "Wrong Password"
        } else {
          message.value = "Success"
        }
        store.dispatch('fetchData')
        store.dispatch('checkLogin')
        });
  }
</script>