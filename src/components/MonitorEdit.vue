<template>
  <div class="heading-secondary">Edit Monitor</div>
  <h1 v-if="!$store.state.loaded">Loading</h1>
  <div v-else>
    <div class="monitor-edit-list">
      <a href="#" class="monitor-edit-list-new u-mtsmall" @click="resetEdit()">+ New Monitor</a>
      <a href="#" class="monitor-edit-list-entry u-mtsmall" v-for="item in monitorList" :key="item.name" @click="chooseMonitor(item)">{{item.name}}
      </a>
    </div>

    <div class="u-mtsmall" v-if="chosen">
      <div class="heading-tertiary" v-if="monitorEdit.id == ''">New Monitor: {{monitorEdit.name}}</div>
      <div class="heading-tertiary" v-else>Edit Monitor: {{monitorEdit.name}}</div>
      <div class="row u-mbsmall" v-if="monitorEdit.mode == 'http'">
        <div class="col-1-of-2 input-group">
          <label for="name">Name</label>
          <input id="name" v-model="monitorEdit.name" placeholder="Monitor Name" />
        </div>
        <div class="col-1-of-2 input-group">
          <label for="interval">Interval (check every {{humanTime(monitorEdit.interval)}})</label>
          <input id="interval" v-model="monitorEdit.interval" type="number" min="5" />
        </div>
        <div class="col-1-of-2 input-group">
          <label for="status">Enabled</label>
          <input id="status" v-model="monitorEdit.enabled" type="checkbox"/>
        </div>
        <div class="col-1-of-2 input-group">
          <label for="inverted">Inverted (reachable is bad)</label>
          <input id="inverted" v-model="monitorEdit.inverted" type="checkbox" />
        </div>
        <div class="col-1-of-2 input-group">
          <label for="url">URL</label>
          <input id="url" v-model="monitorEdit.url" placeholder="https://example.com" />
        </div>
      </div>
      <a class="btn positive" href="#" @click="sendEdit()">Submit</a>
      <span class="heading-tertiary">{{ message }}</span>{{ errorMessage }}
      <span class="float-right" v-if="monitorEdit.id != ''">
        <a class="btn negative" href="#" v-if="!confirmation" @click="confirmation = true">Remove Monitor "{{monitorEdit.name}}"</a>
        <a class="btn negative" href="#" v-if="confirmation" @click="sendRemove()">Are you sure?</a>
      </span>
    </div>
    <div v-else>
      <div class="heading-tertiary u-mtsmall">{{ message }}</div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { computed, reactive, ref } from 'vue'
  import store, { Monitor, MonitorEdit } from '../store'

  // const props = defineProps([''])

  const monitorList = computed(() => {
    return store.state.monitors
  })

  const chosen = ref<boolean>(false)
  const confirmation = ref<boolean>(false)
  const monitorEdit = reactive<MonitorEdit>(
      {
        id: "",
        name: "New Monitor",
        interval: 60,
        enabled: true,
        inverted: false,
        mode: "http",
        url: "http://example.com",
      }
    )

  const message = ref<string>("")
  const errorMessage = ref<string>("")

  function chooseMonitor(monitor: Monitor) {
    // if clicked same twice
    if(monitorEdit.id == monitor.id) {
      resetEdit();
      chosen.value = false;
    } else {
      monitorEdit.id = monitor.id
      monitorEdit.name = monitor.name
      monitorEdit.interval = monitor.interval
      monitorEdit.enabled = monitor.enabled
      monitorEdit.inverted = monitor.inverted
      monitorEdit.mode = monitor.mode
      monitorEdit.url = monitor.url
      chosen.value = true
      confirmation.value = false
      message.value = ""
      errorMessage.value = ""
    }
  }

  function resetEdit() {
    monitorEdit.id = ""
    monitorEdit.name = ""
    monitorEdit.interval = 300
    monitorEdit.enabled = true
    monitorEdit.inverted = false
    monitorEdit.mode = "http"
    monitorEdit.url = ""
    chosen.value = true
    confirmation.value = false
    message.value = ""
    errorMessage.value = ""
  }

  function sendEdit() {
    message.value = "Processing"
    fetch(store.state.apiDomain + "/api/edit", {
      method: 'POST',
      headers: {
        'Accept': 'application/json, text/plain, */*',
        'Content-Type': 'application/json'
      },
      credentials: "include",
      body: JSON.stringify(monitorEdit)
    }).then(res => res.json())
      .then(res => {
        if (!res.error) {
          resetEdit()
          chosen.value = false
          message.value = "Success"
        } else {
          message.value = "Failure"
          errorMessage.value = "  " + res.error
        }
        console.log(res)
        store.dispatch('fetchData')
        });
  }

  function sendRemove() {
    fetch(store.state.apiDomain + "/api/remove", {
      method: 'POST',
      headers: {
        'Accept': 'application/json, text/plain, */*',
        'Content-Type': 'application/json'
      },
      credentials: "include",
      body: JSON.stringify(monitorEdit)
    }).then(res => res.json())
      .then(res => {
        if (!res.error) {
          resetEdit()
          chosen.value = false
          message.value = "Deletion Successful"
        } else {
          message.value = "Failure"
          errorMessage.value = "  " + res.error
        }

        console.log(res)
        store.dispatch('fetchData')
        });
  }

  var reduceInterval: number

  function humanTime(seconds: number) {
    let time = 1
    let ret = ""
    if (seconds >= 1000000000000000000) {
      clearInterval(reduceInterval)
      reduceInterval = setInterval(function () {
          if(monitorEdit.interval > 1000){
            monitorEdit.interval = Math.round(monitorEdit.interval * 0.85)
          } else {
            clearInterval(reduceInterval);
          }
      }, 20);
      return "this happens if you don't listen..."
    } else if (seconds >= 10000000000000000) {
      return "please stop"
    } else if (seconds >= 100000000000000) {
      return "are you serious?"
    } else if (seconds >= 1000000000000) {
      return "would you really need that?"
    } else if (seconds >= 10000000000) {
      return "why would you do this?"
    } else if (seconds >= 157680000) {
      return "probably never"
    } else if (seconds >= 31536000) {
      time = seconds/31536000
      ret = " year"
    } else if (seconds >= 604800) {
      time = seconds/604800
      ret = " week"
    } else if (seconds >= 86400) {
      time = seconds/86400
      ret = " day"
    } else if (seconds >= 3600) {
      time = seconds/3600 
      ret = " hour"
    } else if (seconds >= 60) {
      time = seconds/60 
      ret = " minute"
    } else {
      time = seconds 
      ret = " second"
    }
    time = Math.round (time*100) / 100
    return time + ret + (time===1?"":"s")
  }
</script>