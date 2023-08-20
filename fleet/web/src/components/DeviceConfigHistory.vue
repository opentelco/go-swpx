<script setup lang="ts">
import { Configuration,ConfigurationConnection } from '../gql/graphql'
import { ref } from 'vue'
import TextCopy from './TextCopy.vue';

const props = defineProps<{
  configs: ConfigurationConnection
}>()

const layout = 'loose'

// if int is een return left otherwise return right
const side = (int: number) => {
  if (int % 2 == 0) {
    return 'left'
  } else {
    return 'right'
  }
}

const hasChanged = (cfg: Configuration) => {
  if (cfg.changes) {
    return true
  } else {
    return false
  }
}

// format date to relative date, '3 days and 3 hours ago' etc , if older than 1 month return absolute date
const formatDate = (date: string) => {
  let d = new Date(date)
  let now = new Date()
  let diff = now.getTime() - d.getTime()
  let days = Math.floor(diff / (1000 * 60 * 60 * 24))
  let hours = Math.floor(diff / (1000 * 60 * 60))
  let minutes = Math.floor(diff / (1000 * 60))
  let seconds = Math.floor(diff / (1000))

  if (days > 30) {
    return d.toLocaleDateString()
  } else if (days > 0) {
    return days + ' days ago'
  } else if (hours > 0) {
    return hours + ' hours ago'
  } else if (minutes > 0) {
    return minutes + ' minutes ago'
  } else if (seconds > 0) {
    return seconds + ' seconds ago'
  } else {
    return 'just now'
  }
}

</script>

<template>
  <div>
    <q-timeline :layout="layout" color="secondary">
      <q-timeline-entry heading>
        Configurations timeline
      </q-timeline-entry>

      <q-timeline-entry v-for="cfg,i in props.configs.configurations" :key="cfg.id" :side="side(i)" :icon="hasChanged(cfg) ? 'fa-solid fa-code-merge' : 'fa-regular fa-file-code'" :color="hasChanged(cfg) ? 'blue' : 'primary' ">
        <template v-slot:title>
          {{ hasChanged(cfg) ? 'Configuration changed' : 'Configration baseline'}}
        </template>
        <template v-slot:subtitle>
          <div >
          <TextCopy :text=formatDate(cfg.createdAt) :toCopy="cfg.createdAt" noIcon :class="side(i) == 'right' ? 'float-right' : ''"/>
        </div>

        </template>
        <div>
          <q-btn name="show" icon="fa-solid fa-eye" flat dense round class="q-mr-md">
            <q-tooltip>Show configuration</q-tooltip>
          </q-btn>
          <q-btn name="download" icon="fa-solid fa-download" flat dense round class="q-mr-md">
            <q-tooltip>Download configuration</q-tooltip>
          </q-btn>
          <q-btn name="hash" icon="fa-solid fa-fingerprint" flat dense round class="q-mr-md">
            <q-tooltip>Configuration checksum ({{ cfg.checksum }})</q-tooltip>
          </q-btn>
          <q-btn name="compare" icon="fa-solid fa-code-merge" flat dense round class="q-mr-md" disabled>
            <q-tooltip>Compare configuration</q-tooltip>
          </q-btn>

          <q-scroll-area
          v-if="hasChanged(cfg)"
        style="height: 300px; width:100%"
        >
          <pre class="q-m-lg">{{cfg.changes}}</pre>
          </q-scroll-area>
          <div v-else>Baseline Configuration fetched from device</div>
        </div>
      </q-timeline-entry>
    </q-timeline>
  </div>
</template>
