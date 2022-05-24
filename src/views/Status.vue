<template>
  <div class="heading-secondary">Monitors</div>
  <h1 v-if="!$store.state.loaded">Loading</h1>
  <div v-else>
    <div v-if="!$store.state.error" class="monitor-list bubble">
      <MonitorEntry v-for="item in monitorList" :key="item.name" :monitorEntry="item"></MonitorEntry>
    </div>
  </div>
</template>

<script setup lang="ts">
  import MonitorEntry from '@/components/MonitorEntry.vue';
  import { computed } from 'vue'
  import store from '../store'

  const monitorList = computed(() => {
    return store.state.monitors.filter(monitor => monitor.enabled)
  })
</script>