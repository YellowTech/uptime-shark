<template>
  <h1 v-if="this.$store.state.error">Oops! An Error Occurred...</h1>
  <div v-else>
    <MonitorEdit></MonitorEdit>
  </div>
</template>

<script setup lang="ts">
  import MonitorEdit from '@/components/MonitorEdit.vue';

</script>