<script setup lang="ts">
import { Configuration, ListConfigurationsResponse } from '../gql/graphql'
import { toRefs, ref, computed } from 'vue'
import { getRelativeTimestamp } from '../utils/time'
import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import TextCopy from './TextCopy.vue';

let getDeviceConfigurations = gql`query DeviceConfigurations (
  $deviceId: ID!
  $limit: Int
  $offset: Int
  ) {
    configurations( params: { deviceId: $deviceId, limit: $limit, offset: $offset } ) {
        configurations {
            id
            configuration
            changes
            checksum
            createdAt
        }
        pageInfo {
            limit
            offset
            total
            count
        }
      }
}`




const props = defineProps<{
  deviceId: string,
}>()

let { deviceId } = toRefs(props)

let limit = ref(3)
let offset = ref(0)



type response = {
  configurations: ListConfigurationsResponse
}


let { onResult, result, loading, error, refetch } = useQuery<response>(getDeviceConfigurations, {
  deviceId: deviceId.value,
  limit: limit,
  offset: offset
})

const pagination = ref({
  sortBy: 'desc',
  descending: false,
  page: 1,
  rowsPerPage: limit.value,
  rowsNumber: 99
})


onResult(queryResult => {
  if (queryResult.data !== undefined) {
    pagination.value.rowsNumber = queryResult.data.configurations.pageInfo.total
  }
})

// computed pagination with the help of the calculatePagination fuction
let reFetch = ((page) => {
  offset.value = (page - 1) * pagination.value.rowsPerPage
  refetch()
})



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



const maxPage = computed(() => {
  return Math.ceil(pagination.value.rowsNumber / pagination.value.rowsPerPage)
})

</script>

<template>
  <div class="q-pa-md">
    <div class="row justify-center">
      <div class="text-h2 q-mt-lg">
        Configuration history
      </div>
    </div>
    <div class="row justify-center">
      <div class="col-12 q-pa-lg flex flex-center">
        <q-pagination v-model="pagination.page" :max=maxPage input @update:model-value="reFetch" v-if="maxPage>1"/>
      </div>
    </div>

    <div class="row justify-center q-mb-xl">
      <div class="col-12" v-if="result?.configurations.configurations.length > 0">
        <q-timeline :layout="layout" color="secondary">
          <q-timeline-entry v-for="cfg, i in result?.configurations.configurations" :key="cfg.id" :side="side(i)"
            :icon="hasChanged(cfg) ? 'fa-solid fa-code-merge' : 'fa-regular fa-file-code'"
            :color="hasChanged(cfg) ? 'blue' : 'primary'" v-if="result?.configurations.configurations.length > 0">
            <template v-slot:title>
              {{ hasChanged(cfg) ? 'Configuration changed' : 'Configration baseline' }}
            </template>
            <template v-slot:subtitle>
              <div>
                <TextCopy :text=getRelativeTimestamp(cfg.createdAt) :toCopy="cfg.createdAt" noIcon
                  :class="side(i) == 'right' ? 'float-right' : ''" />
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

              <div style="width: 100%;">

                <div  v-if="hasChanged(cfg)">

                  <q-scroll-area style="height: 400px;">
                    <pre class="q-m-lg">{{ cfg.changes }}</pre>
                  </q-scroll-area>
                </div>

                <div class="row q-mt-lg text-center config-history-item"
                  :class="[side(i) === 'left' ? 'float-right' : 'float-left']" v-else>
                  <div class="col q-mt-lg">
                    <q-img src="img/undraw_light_the_fire_gt58.svg" style="max-width: 200px" />
                    <h5 class="title">Configuration baseline</h5>
                    <p class="text-body1">
                      The first configuration we collected for the devce.<br>

                    </p>
                  </div>
                </div>
              </div>
            </div>
          </q-timeline-entry>
        </q-timeline>
      </div>
      <div class="col-12 col-md-auto text-center q-mt-lg" v-else>
        <q-img src="img/undraw_empty_re_opql.svg" style="max-width: 400px" />
        <h5 class="title">No stored configurations</h5>
        <p class="text-body1">
          Get a configuration for the device and it will show up here.
        </p>
        <q-btn color="primary" icon="fa-solid fa-file-code" class="">
          <span class="q-pl-sm">Collect config</span>
        </q-btn>

      </div>
    </div>
  </div>
</template>


<style lang="sass">
.config-history-item
  max-width:600px
  height: 400px
.title
  margin-bottom: 3px

  </style>
