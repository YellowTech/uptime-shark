<template>
  <div class="heading-secondary">Services</div>
  <div class="service-list">
    <ServiceEntry v-for="item in statusList" :key="item.name" :serviceName="item.name" :serviceStatus="item.status"></ServiceEntry>
  </div>

</template>

<script setup lang="ts">
  import ServiceEntry from '@/components/ServiceEntry.vue';
  import { computed } from 'vue'
  import store from '../store'
  import {Istate} from '../store/interfaces'

  // const props = defineProps([''])

  const statusList = computed(() => {
    return store.state.status as unknown as Istate
  })
</script>