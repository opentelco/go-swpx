<script setup lang="ts">
import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { DeviceChange, DeviceEvent, ListDeviceChangesResponse, ListDeviceEventsResponse } from '../gql/graphql'

import { ref, toRefs, onMounted, computed } from 'vue'



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

      deviceEvents( params: { deviceId: $deviceId, limit: $limit, offset: $offset } ) {
      events {
            id
            deviceId
            type
            message
            action
            outcome
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

function formatDateToYYYYMMDD(date: Date): string {
  return date.toISOString().slice(0, 10);
}

// Timeline is a computed value that returns a map of TimelineObject  ordered by createdAt
const timeline = computed(() => {
  let timeline: Map<string, TimelineObject[]> = new Map()
  let changes = result.value?.deviceChanges.changes
  let events = result.value?.deviceEvents.events




  if (changes !== undefined && changes !== null) {
    changes.forEach((change) => {
      let ts = formatDateToYYYYMMDD(new Date(change.createdAt))

      if (timeline.has(ts)) {
        let items = timeline.get(ts)
        items?.push({ change: change })
        timeline.set(ts, items!)
      } else {
        timeline.set(ts, [{ change: change }])
      }
    })
  }

  if (events !== undefined && events !== null) {
    events.forEach((event) => {
      let ts = formatDateToYYYYMMDD(new Date(event.createdAt))

      if (timeline.has(ts)) {
        let items = timeline.get(ts)
        items?.push({ event: event })
        timeline.set(ts, items!)
      } else {
        timeline.set(ts, [{ event: event }])
      }
    })
  }

  let ordredMap = new Map([...timeline.entries()].sort((a, b) => {
    return new Date(a[0]).getTime() - new Date(b[0]).getTime()
  }))
  return Array.from(ordredMap)
})

// function to convert timeline to a array of key and values




const props = defineProps<{
  deviceId: string,
}>()

let { deviceId } = toRefs(props)

let limit = ref(10)
let offset = ref(0)



type response = {
  deviceChanges: ListDeviceChangesResponse
  deviceEvents: ListDeviceEventsResponse
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
        <q-pagination v-model="pagination.page" :max=maxPage input @update:model-value="reFetch" />
      </div>
    </div>
    <div class="row">
      <div class="col col-12">
        {{ error }}
        <q-timeline layout="loose" color="secondary" class="" v-if="timeline.length > 0" :loading=loading>
          <q-timeline-entry :side="isEven(i)" icon="edit" v-for="(item, key) in timeline" :key="i" class="q-mb-xl">
            <template v-slot:title>
              <span class="text-weight-medium q-mr-md">
                {{ item[0] }}
              </span>
            </template>
            <template v-slot:subtitle>
              <span class="text-weight-medium q-mr-md">
                {{ item[1].length }} happenings
              </span>
            </template>

            <div class="text-body1">

              <q-card style="max-width:700px" class="q-pa-lg q-ma-lg">
                <div class="row" v-for="(h, i) in item[1]" :key="i">

                  <div v-if="h.event" class="col col-8 q-pa-md ">
                    <q-card class="my-card shadow-2 q-pa-sm"  bordered>
                      <q-item>
                        <q-item-section avatar>
                          <q-avatar>
                            <q-icon name="fa-solid fa-envelope-open-text"/>
                          </q-avatar>
                        </q-item-section>

                        <q-item-section>
                          <q-item-label>{{ h.event?.type }}</q-item-label>
                          <q-item-label caption>
                            {{ h.event?.action }}
                          </q-item-label>
                        </q-item-section>
                      </q-item>

                      <q-separator />
                      <div class="text-body1">
                        {{ h.event?.message }}
                      </div>
                    </q-card>
                  </div>

                  <div v-if="h.change" class="col col-8 q-pa-md ">
                    <q-card class="my-card shadow-2 q-pa-sm"  bordered>
                      <q-item>
                        <q-item-section avatar>
                          <q-avatar>
                            <q-icon name="fa-solid fa-pen-to-square"/>
                          </q-avatar>
                        </q-item-section>

                        <q-item-section>
                          <q-item-label>Field <q-badge align="middle">{{ h.change?.field }}</q-badge> changed</q-item-label>
                          <q-item-label caption>
                            {{ h.change?.newValue }} <q-icon name="fa-solid fa-arrow-right"/> {{ h.change?.oldValue }}
                          </q-item-label>
                        </q-item-section>
                      </q-item>

                      <q-separator />
                      <div class="text-body1">
                        {{ h.event?.message }}
                      </div>
                    </q-card>
                  </div>

                </div>
                <q-separator class="q-mb-md q-mt-md" />
                <div class="row">
                  <div class="col text-subtitle1">
                  </div>
                  <div class="col text-body1">
                  </div>
                </div>
                <div class="row">
                  <div class="col text-subtitle1">
                  </div>
                  <div class="col text-body1">
                  </div>
                </div>
                <div class="row q-mt-md">
                  <div class="col-12 text-caption">
                    <div class="text-right">
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
            <h5 class="title">Quite times</h5>
            <p class="text-body1">
              If anything happens on the device props it will show up here. <br>
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
