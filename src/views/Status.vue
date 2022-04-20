<template>
  <h1 v-if="this.$store.state.error">Oops! An Error Occurred...</h1>
  <div v-else>
    <div class="heading-secondary">Monitors</div>
    <h1 v-if="!this.$store.state.loaded">Loading</h1>
    <div v-else>
      <div v-if="!this.$store.state.error" class="monitor-list bubble">
        <MonitorEntry v-for="item in monitorList" :key="item.name" :monitorEntry="item"></MonitorEntry>
      </div>
    </div>
  </div>

</template>

<script setup lang="ts">
  import MonitorEntry from '@/components/MonitorEntry.vue';
  import { computed } from 'vue'
  import store from '../store'

  // const props = defineProps([''])

  const monitorList = computed(() => {
    return store.state.monitors
  })
</script>