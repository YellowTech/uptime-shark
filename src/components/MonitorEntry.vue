<template>
  <div class="monitor-entry">
    <span class="monitor-pills-box">
      <div class="monitor-pills-pill tooltip tooltip-fade" :class="[!log.Failed?'positive':'negative', 'monitor-status']" :data-tooltip="log.Message + ' - ' + timeConverter(log.Time)" v-for="log in props.monitorEntry.logs" :key="log.Time"></div>
    </span>

    <span :class="[props.monitorEntry.status?'positive':'negative', 'monitor-status']">{{uptimePercent * 100}}%</span>
    <span class="monitor-name">{{ props.monitorEntry.name }}</span>

    <!-- <p>{{props.monitorEntry.id}}</p>
    <p>{{ props.monitorEntry.status }}</p>
    <p>{{ props.monitorEntry.url }}</p>
    <p v-for="log in props.monitorEntry.logs" :key="log.Time">{{log.Message}} - {{timeConverter(log.Time)}}</p>
    <p>{{props.monitorEntry.statusMessage}}</p> -->
  </div>
</template>

<script setup lang="ts">
  import { computed, defineProps } from 'vue'
  import type { Monitor } from '@/store'

  const uptimePercent = computed(() => {
    return Number(props.monitorEntry.logs.filter(entry => !entry.Failed).length / props.monitorEntry.logs.length).toFixed(2)
  })

  const props = defineProps<{
    monitorEntry: Monitor, 
  }>()

  function timeConverter(UNIX_timestamp:number){
    var a = new Date(UNIX_timestamp * 1000);
    var months = ['Jan','Feb','Mar','Apr','May','Jun','Jul','Aug','Sep','Oct','Nov','Dec'];
    var year = a.getFullYear();
    var month = months[a.getMonth()];
    var date = a.getDate();
    var hour = a.getHours() < 10 ? '0' + a.getHours() : a.getHours();
    var min = a.getMinutes() < 10 ? '0' + a.getMinutes() : a.getMinutes();
    var sec = a.getSeconds() < 10 ? '0' + a.getSeconds() : a.getSeconds();
    var time = date + ' ' + month + ' ' + year + ' ' + hour + ':' + min + ':' + sec ;
    return time;
  }
</script>
