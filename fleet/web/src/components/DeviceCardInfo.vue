<script setup lang="ts" >

import { Device } from '../gql/graphql'

const props = defineProps<{
  device: Device
}>()

</script>

<template>
  <q-card class="shadow-2 card">
    {{ props.device.hostname || 'Unknown' }}

  </q-card>

</template>

<style lang="sass">
.card
  padding: 10px
  width: 100%
  border-radius: 15px
</style>
