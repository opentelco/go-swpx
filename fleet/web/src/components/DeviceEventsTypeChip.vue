<script setup lang="ts">
import { DeviceEventType } from '../gql/graphql'

const props = defineProps<{
  type: DeviceEventType
}>()
</script>

<template>
  <q-avatar color="primary" size="30px">

    <q-icon color="white" name="fa-solid fa-network-wired" v-if="props.type == DeviceEventType.Device" />

    <q-icon color="white" name="fa-solid fa-file-code" v-else />
    <q-tooltip>
      type <q-icon name="arrow_right_alt" /> {{ props.type }}
    </q-tooltip>
  </q-avatar>
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
