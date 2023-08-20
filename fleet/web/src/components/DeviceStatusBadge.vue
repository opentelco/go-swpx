<script setup lang="ts">
import { DeviceStatus } from '../gql/graphql'
import { ref } from 'vue'

let hover = ref(false)
const props = defineProps<{
  status: DeviceStatus
}>()



// colorFromStatus is a function that get the prop status and switch cases all the possible values
// and return a bg color and fg color matching the status
// (the bg color is the color of the chip and the fg color is the color of the icon)
const colorFromStatus = (status: DeviceStatus) => {
  switch (status) {
    case DeviceStatus.Reachable:
      return { bg: 'primary', fg: 'white' }
    case DeviceStatus.Unreachable:
      return { bg: 'negative', fg: 'white' }
    default:
      return { bg: 'grey', fg: 'white' }
  }
}




// statusIcon is a function that get the prop status and switch cases all the possible values
// and return an icon matching the status
const statusIcon = (status: DeviceStatus) => {
  switch (status) {
    case DeviceStatus.Reachable:
      return 'fa-solid fa-heart-pulse'
    case DeviceStatus.Unreachable:
      return 'remove_circle'
    default:
      return 'help'
  }
}

// statusTooltip is a function that get the prop status and switch cases all the possible values
// and return a tooltip matching the status
const statusTooltip = (status: DeviceStatus) => {
  let pre = 'Status: '
  switch (status) {
    case DeviceStatus.Reachable:
      return pre + ' Device is reachable'
    case DeviceStatus.Unreachable:
      return pre + 'Fleet could not reach the device'
    case DeviceStatus.New:
      return pre + 'Device is new, no checks have been run yet'

    default:
      return pre + 'Device status is unknown'
  }
}



</script>

<template>
  <div>
    <q-chip size="sm" outline square :color="colorFromStatus(props.status).bg" :text-color="colorFromStatus(props.status).fg" :icon="statusIcon(props.status)">
      {{ props.status }}
      <q-tooltip class="shadow-4">
      {{ statusTooltip(props.status) }}
    </q-tooltip>
    </q-chip>


  </div>

</template>
