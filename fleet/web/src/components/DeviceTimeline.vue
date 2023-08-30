<script setup lang="ts">
import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { DeviceChange, DeviceEvent, ListDeviceChangesResponse } from '../gql/graphql'

import { getRelativeTimestamp } from '../utils/time'


import { ref, toRefs, onMounted, computed } from 'vue'
import TextCopy from './TextCopy.vue';


let getDeviceChanges = gql`query DeviceChanges (
  $deviceId: ID!
  $limit: Int
  $offset: Int
  ) {
    deviceChanges( params: { deviceId: $deviceId, limit: $limit, offset: $offset } ) {
      changes {
          id
          deviceId
          field
          oldValue
          newValue
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


type TimelineObject = {
  change?: DeviceChange
  event?: DeviceEvent
}

// Timeline is a array of TimelineObject
let timeline = ref<TimelineObject[]>([])


const props = defineProps<{
  deviceId: string,
}>()

let { deviceId } = toRefs(props)

let limit = ref(10)
let offset = ref(0)



type response = {
  deviceChanges: ListDeviceChangesResponse
}


let { onResult, result, loading, error, refetch } = useQuery<response>(getDeviceChanges, {
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
    pagination.value.rowsNumber = queryResult.data.deviceChanges.pageInfo.total
  }
})

// computed pagination with the help of the calculatePagination fuction
let reFetch = ((page) => {
    offset.value = (page - 1) * pagination.value.rowsPerPage
    refetch()
  })

// isEven returns right if even left is odd
const isEven = (num: number) => {
  if (num % 2 == 0) {
    return 'right'
  } else {
    return 'left'
  }
}

// computed maxPAge value from pagination page and rowsPerPage etc
const maxPage = computed(() => {
  return Math.ceil(pagination.value.rowsNumber / pagination.value.rowsPerPage)
})

onMounted(() => {
  // eventTableRef.value.requestServerInteraction()
})
</script>

<template>
  {{ error }}
  <div class="q-pa-md">
    <div class="row justify-center">
      <div class="text-h2 q-mt-lg">
        Timeline
      </div>
    </div>
    <div class="row justify-center">
      <div class="q-pa-lg flex flex-center">
        <q-pagination v-model="pagination.page" :max=maxPage input @update:model-value="reFetch"/>
      </div>
    </div>
    <div class="row">
      <div class="col col-12">
        {{ error }}
        <q-timeline layout="loose" color="secondary" class="" v-if="result?.deviceChanges.changes?.length > 0">
          <q-timeline-entry :side="isEven(i)" icon="edit" v-for="p, i in result?.deviceChanges.changes" :key="p.id"
            class="q-mb-xl">
            <template v-slot:title>
              <span class="text-weight-medium q-mr-md"> Field <q-chip square color="primary" text-color="white" dense
                  class="text-uppercase">{{ p.field }}</q-chip>
                changed</span>
            </template>
            <template v-slot:subtitle>
              {{ getRelativeTimestamp(p.createdAt) }}
              <q-tooltip>{{ p.createdAt }}</q-tooltip>
            </template>

            <div class="text-body1">

              <q-card style="max-width:700px" class="q-pa-lg q-ma-lg">
                <div class="row">
                  <div class="col">
                    The field changed from <q-badge color="grey">{{ p.oldValue }}</q-badge> <q-icon
                      name="arrow_forward" />
                    <q-badge color="green">{{ p.newValue }}</q-badge>
                  </div>
                </div>
                <q-separator class="q-mb-md q-mt-md" />
                <div class="row">
                  <div class="col text-subtitle1">
                    <q-icon name="commit" />
                    Before Change
                  </div>
                  <div class="col text-body1">
                    <TextCopy :text="p.oldValue" />
                  </div>
                </div>
                <div class="row">
                  <div class="col text-subtitle1">
                    <q-icon name="merge" />
                    After Change
                  </div>
                  <div class="col text-body1">
                    <TextCopy :text="p.newValue" />
                  </div>
                </div>
                <div class="row q-mt-md">
                  <div class="col-12 text-caption">
                    <div class="text-right">
                      Change was done <span class="text-bold">{{ p.createdAt }} </span>
                    </div>
                  </div>
                </div>
              </q-card>
            </div>
          </q-timeline-entry>
        </q-timeline>




        <div class="full-width row flex-center q-mt-lg " v-else>
          <div class="col-12 col-md-auto text-center q-mt-lg">
            <q-img src="img/undraw_empty_re_opql.svg" style="max-width: 400px" />
            <h5 class="title">No changes</h5>
            <p class="text-body1">
              If anything changes on the device props it will show up here. <br>
            </p>
          </div>

        </div>

      </div>



    </div>

  </div>
</template>

<style lang="sass">

.rounded
  border-radius: 10px

.search-input
  width: 350px


</style>
