<template>
  <Login v-if="!$store.state.authenticated" />
  <MonitorEdit v-else />
</template>

<script setup lang="ts">
  import MonitorEdit from '@/components/MonitorEdit.vue';
  import Login from '@/components/Login.vue';

</script>