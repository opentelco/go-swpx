<script setup lang="ts">
import { ref } from 'vue'
import { DeviceEventOutcome } from 'src/gql/graphql';

let props = defineProps<{
  outcome: DeviceEventOutcome
  size: string | 'md'
}>()


type OutcomeIcon = {
  icon: string,
  fg: string
  bg: string
}

// switch icon depending on the DeviceEventOutcome
// return icon and color matching the outcome (OutcomeIcon)
const outcomeIcon = (outcome: DeviceEventOutcome) => {
  switch (outcome) {
    case DeviceEventOutcome.Success:
      return { icon: 'check', bg: 'green', fg: 'white' }
    case DeviceEventOutcome.Failure:
      return { icon: 'close', bg: 'red', fg: 'white' }
    default:
      return { icon: 'help', bg: 'yellow', fg: 'black' }
  }

}


</script>

<template>
  <q-chip class="shadow-1 subtitle">
    <q-avatar :color="outcomeIcon(props.outcome).bg">
      <q-icon :color="outcomeIcon(props.outcome).fg" :name="outcomeIcon(props.outcome).icon" />
    </q-avatar>
      <span class="">
        {{ props.outcome }}
      </span>
  </q-chip>
</template>

<style lang="sass">
</style>
