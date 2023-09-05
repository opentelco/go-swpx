<script setup lang="ts">
import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { ListDeviceChangesResponse } from '../gql/graphql'

import { getRelativeTimestamp } from '../utils/time.ts'

import DeviceEventOutcomeChip from './DeviceEventOutcomeChip.vue'
import DeviceEventsTypeChip from './DeviceEventsTypeChip.vue';
import DeviceEventActionChip from './DeviceEventActionChip.vue';
import { ref, toRefs, onMounted } from 'vue'
import TextCopy from './TextCopy.vue';


let getDeviceEvents = gql`query DeviceEvents (
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




const props = defineProps<{
  deviceId: string,
}>()

let { deviceId } = toRefs(props)
let eventTableRef = ref()




const columns = [
  { name: 'timestamp', label: 'Timestamp', field: 'createdAt', align: 'left' },
  { name: 'message', label: 'Message', field: 'message', align: 'left' },
  { name: 'tags', label: '', field: 'tags', align: 'center' },
]


let limit = ref(10)
let offset = ref(0)


type response = {
  deviceChanges: ListDeviceChangesResponse
}


let { onResult, result, loading, error, refetch } = useQuery<response>(getDeviceEvents, {
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
let reFetch = ((props) => {
  const { page, rowsPerPage, sortBy, descending } = props.pagination


  if (page !== pagination.value.page) {
    offset.value = (page - 1) * rowsPerPage
    refetch()
  }

  if (rowsPerPage !== limit.value) {
    limit.value = rowsPerPage
    offset.value = (page - 1) * rowsPerPage
    refetch()
  }
  pagination.value.page = page
  pagination.value.rowsPerPage = rowsPerPage
  pagination.value.sortBy = sortBy
  pagination.value.descending = descending
})


onMounted(() => {
  eventTableRef.value.requestServerInteraction()
})
</script>

<template>
  {{ error }}
  <div class="q-pa-md">
    <div class="row justify-center">
      <div class="">
      </div>
    </div>
    <div class="row">
      <div class="col col-12">
        <q-table grid :loading="loading" title="Events" :rows="result?.deviceChanges.changes" :columns="columns"
          row-key="id" binary-state-sort class="q-pa-lg rounded shadow-3" v-model:pagination="pagination"
          @request="reFetch" ref="eventTableRef">

          <template v-slot:item="p">
            <q-list :props="p" class="rounded-borders col-12" bordered>
              <q-item clickable v-ripple class="q-pa-md ">
                <q-item-section avatar top>
                  <q-avatar>
                    <q-icon name="fa-solid fa-pen-to-square" />
                  </q-avatar>
                </q-item-section>

                <q-item-section top class="col-2 gt-sm">
                  <q-item-label class="q-mt-sm">
                    {{ p.row.type }}
                  </q-item-label>
                </q-item-section>
                <q-item-section top>
                  <q-item-label lines="1">
                    <span class="text-weight-medium q-mr-md"> <q-badge>{{ p.row.field }}</q-badge></span>
                  </q-item-label>

                  <q-item-label caption lines="1">
                    <q-icon name="fa-solid fa-arrow-left"/>

                    {{ p.row.oldValue }}


                  </q-item-label>
                  <q-item-label lines="1" class="q-mt-xs text-body2 text-weight-bold text-primary text-uppercase">
                    <q-icon name="fa-solid fa-arrow-right"/> {{ p.row.newValue }}
                  </q-item-label>
                </q-item-section>

                <q-item-section top side>
                  <div class="text-grey-8 q-pa-xs">
                    {{ getRelativeTimestamp(p.row.createdAt) }}
                    <q-tooltip>
                      {{ p.row.createdAt }}
                    </q-tooltip>
                  </div>
                  <div class="text-grey-8 q-gutter-xs">

                  </div>
                </q-item-section>
              </q-item>

            </q-list>

          </template>
          <template v-slot:no-data="{ message }">
            <div class="full-width row flex-center q-mt-lg ">
              <div class="col-12 col-md-auto text-center q-mt-lg">
                <q-img src="img/undraw_empty_re_opql.svg" style="max-width: 400px" />
                <h5 class="title">No Changes</h5>
                <p class="text-body1">
                  If anything has changes on a device it will show up here. <br>
                  Reason:
                  {{ message }}
                </p>
              </div>

            </div>
          </template>
        </q-table>
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
