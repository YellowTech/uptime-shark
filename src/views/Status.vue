<template>
  <h1 v-if="this.$store.state.error">Error Occurred</h1>
  <div class="heading-secondary">Services</div>
  <h1 v-if="!this.$store.state.loaded">Loading</h1>
  <div v-else>
    <div v-if="!this.$store.state.error" class="service-list bubble">
      <ServiceEntry v-for="item in monitorList" :key="item.name" :serviceName="item.name" :serviceStatus="item.status"></ServiceEntry>
    </div>

  </div>

</template>

<script setup lang="ts">
  import ServiceEntry from '@/components/ServiceEntry.vue';
  import { computed } from 'vue'
  import store from '../store'

  // const props = defineProps([''])

  const monitorList = computed(() => {
    return store.state.monitors
  })
</script>