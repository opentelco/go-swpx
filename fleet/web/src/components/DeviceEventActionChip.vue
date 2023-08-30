<script setup lang="ts">
import { DeviceEventAction } from 'src/gql/graphql';

let defaults = {
  size: 'sm'
}

let props = defineProps<{
  action: DeviceEventAction

}>()




// switch icon depending on the DeviceEventAction
// return icon and color matching the action (ActionIcon)
const actionIcon = (action: DeviceEventAction) => {
  switch (action) {
    case DeviceEventAction.CollectConfig:
      return { icon: 'fa-solid fa-code-merge', bg: 'secondary', fg: 'white'}
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
    <q-avatar size="30px" :color="actionIcon(props.action).bg">
      <q-icon :color="actionIcon(props.action).fg" :name="actionIcon(props.action).icon" />
      <q-tooltip >
      action <q-icon name="arrow_right_alt"/> {{ props.action }}
    </q-tooltip>
    </q-avatar>

</template>

<style lang="sass">
</style>
