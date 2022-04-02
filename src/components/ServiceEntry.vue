<template>
  <div class="service-entry">
    <h1>{{ props.serviceName }}</h1>
    <p>{{ props.serviceStatus }}</p>
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

  const props = defineProps(['serviceName', 'serviceStatus'])

</script>
