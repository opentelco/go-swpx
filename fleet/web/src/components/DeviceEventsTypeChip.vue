<script setup lang="ts">
import { DeviceEventType } from '../gql/graphql'

const props = defineProps<{
  type: DeviceEventType
}>()
</script>

<template>
  <q-chip class="shadow-1 subtitle">
    <q-avatar color="primary">
      <q-icon color="white" name="fa-solid fa-network-wired" v-if="props.type == DeviceEventType.Device" />
      <q-icon color="white" name="fa-solid fa-file-code" v-else />
    </q-avatar>
    <span class="">
      {{ props.type }}
    </span>
  </q-chip>
</template>

<style lang="sass">
.type-card
  width: 100%
  max-width: 160px
  max-height: 70px
  font-size: 0.8em
  padding: 0px
  // round corners
  border-radius: 10px
  // inverted cool shadow

.subtitle
  font-size: 0.9em
  text-transform: uppercase



</style>
