<script setup lang="ts">
import { DeviceEventOutcome } from 'src/gql/graphql';

let props = defineProps<{
  outcome: DeviceEventOutcome
}>()



// switch icon depending on the DeviceEventOutcome
// return icon and color matching the outcome (OutcomeIcon)
const outcomeIcon = (outcome: DeviceEventOutcome) => {
  switch (outcome) {
    case DeviceEventOutcome.Success:
      return { icon: 'check', bg: 'green', fg: 'white' }
    case DeviceEventOutcome.Failure:
      return { icon: 'warning', bg: 'red', fg: 'white' }
    default:
      return { icon: 'help', bg: 'yellow', fg: 'black' }
  }

}


</script>

<template>
  <q-avatar :color="outcomeIcon(props.outcome).bg" size="30px">
    <q-icon :color="outcomeIcon(props.outcome).fg" :name="outcomeIcon(props.outcome).icon" />
    <q-tooltip>
      outcome <q-icon name="arrow_right_alt" /> {{ props.outcome }}
    </q-tooltip>
  </q-avatar>
</template>

<style lang="sass">
</style>
