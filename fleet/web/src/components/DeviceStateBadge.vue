<script setup lang="ts">
import { DeviceState } from '../gql/graphql'
import { ref } from 'vue'
const props = defineProps<{
  state: DeviceState
}>()

let hover = ref(false)


// colorFromState is a function that get the prop state and switch cases all the possible values
// and return a bg color and fg color matching the state
// (the bg color is the color of the chip and the fg color is the color of the icon)
const colorFromState = (state: DeviceState) => {
  switch (state) {
    case DeviceState.Active:
      return { bg: 'positive', fg: 'white' }
    case DeviceState.Inactive:
      return { bg: 'negative', fg: 'white' }
    case DeviceState.Rouge:
      return { bg: 'warning', fg: 'black' }
    case DeviceState.New:
      return { bg: 'blue', fg: 'white' }
    default:
      return { bg: 'grey', fg: 'white' }
  }
}


// stateIcon is a function that get the prop state and switch cases all the possible values
// and return an icon matching the state
const stateIcon = (state: DeviceState) => {
  switch (state) {
    case DeviceState.Active:
      return 'check_circle'
    case DeviceState.Inactive:
      return 'remove_circle'
    case DeviceState.Rouge:
      return 'warning'
    case DeviceState.New:
      return 'add_circle'
    default:
      return 'help'
  }
}

// stateTooltip is a function that get the prop state and switch cases all the possible values
// and return a tooltip matching the state
const stateTooltip = (state: DeviceState) => {
  let pre = 'State: '
  switch (state) {
    case DeviceState.Active:
      return pre + ' Polling is enabled'
    case DeviceState.Inactive:
      return pre + 'Polling is disabled'
    case DeviceState.Rouge:
      return pre + 'Device is rouge'
    case DeviceState.New:
      return pre + 'Device is new'
    default:
      return pre + 'Device is unknown'
  }
}

</script>

<template>
  <div>
    <q-chip size="sm" outline square :color="colorFromState(props.state).bg" :text-color="colorFromState(props.state).fg"
      :icon="stateIcon(props.state)">
      {{ props.state }}
      <q-tooltip class="shadow-4">
        {{ stateTooltip(props.state) }}
      </q-tooltip>
    </q-chip>

  </div>
</template>
