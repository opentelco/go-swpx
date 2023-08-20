<script setup lang="ts">
import { DeviceEventAction } from 'src/gql/graphql';

let defaults = {
  size: 'sm'
}

let props = defineProps<{
  action: DeviceEventAction
  size: string

}>()


type ActionIcon = {
  icon: string,
  fg: string
  bg: string
}

// switch icon depending on the DeviceEventAction
// return icon and color matching the action (ActionIcon)
const actionIcon = (action: DeviceEventAction) => {
  switch (action) {
    case DeviceEventAction.CollectConfig:
      return { icon: 'fa-solid fa-file-code', bg: 'primary', fg: 'white'}
    case DeviceEventAction.CollectDevice:
      return { icon: 'query_stats', bg: 'secondary', fg: 'white' }
    case DeviceEventAction.Create:
      return { icon: 'fa-regular fa-file', bg: 'blue', fg: 'white' }
    case DeviceEventAction.Update:
      return { icon: 'edit', bg: 'yellow', fg: 'white' }
    default:
      return { icon: 'help', bg: 'yellow', fg: 'black' }
  }

}

</script>

<template>
  <q-chip class="shadow-1 subtitle">
    <q-avatar :color="actionIcon(props.action).bg">
      <q-icon :color="actionIcon(props.action).fg" :name="actionIcon(props.action).icon" />
    </q-avatar>
      <span class="">
        {{ props.action }}
      </span>
  </q-chip>
</template>

<style lang="sass">
</style>
