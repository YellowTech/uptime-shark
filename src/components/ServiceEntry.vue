<template>
  <div class="service-entry">
    <h1>{{ props.monitorEntry.name }}</h1>
    <p>{{ props.monitorEntry.status }}</p>
    <p>{{ props.monitorEntry.url }}</p>
    <p v-for="log in props.monitorEntry.logs" :key="log.Time">{{log.Message}}-{{timeConverter(log.Time)}}</p>
  </div>
</template>

<!-- <script lang="ts">
import { Options, Vue } from 'vue-class-component';

@Options({
  props: {
    name: String,
    status: Boolean
  }
})
export default class ServiceEntry extends Vue {
  name!: string
  status!: boolean
}
</script> -->

<script setup lang="ts">
  import { defineProps } from 'vue'
  import type { Monitor } from '@/store'

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
