<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { Device, ListDeviceResponse } from '../gql/graphql'
import DeviceStateBadge from './DeviceStateBadge.vue'
import DeviceStatusBadge from './DeviceStatusBadge.vue'
import TextCopy from './TextCopy.vue'
import DeviceVendorImage from './DeviceVendorImage.vue'

let eventTableRef = ref()

let grid = ref(true)


const columns: Array<any> = [
  {
    name: 'hostname',
    field: 'hostname',
    required: true,
    label: 'Hostname',
    align: 'left',
    sortable: true
  },
  { name: 'managementIp', align: 'center', label: 'Management IP', field: 'managementIp', sortable: true },
]

let limit = ref(10)
let offset = ref(0)
let filter = ref('')

// remap the response so it matches the response
type response = {
  devices: ListDeviceResponse
}


const { onResult, result, loading, error, refetch } = useQuery<response>(gql`
  query Devices ($limit: Int!,  $search: String!){
    devices(params: { limit: $limit, search: $search } ) {
        pageInfo {
            limit
            offset
            total
            count
        }
        devices {
            id
            hostname
            managementIp
            serialNumber
            model
            networkRegion
            pollerResourcePlugin
            pollerProviderPlugin
            version
            state
            status
            lastSeen
            createdAt
            updatedAt
            lastReboot
        }
      }
    }
`, {
  limit: 10,
  search: filter
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
    pagination.value.rowsNumber = queryResult.data.devices.pageInfo.total
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

interface CardDevice extends Device {
  loading: boolean;
}

const data = computed(() => {
  if (result.value && result.value.devices) {
    return result.value.devices.devices.map((device: CardDevice) => {
      const x = Math.floor(Math.random() * 2)
      return {
        ...device,
        loading: x,
      }
    })
  }
  return []
})

</script>

<template>
  <q-banner inline-actions class="text-white bg-red" v-if="error">
    Could not fetch devices: {{ error }}
  </q-banner>

  <div class="q-pa-md">
    <div class="text-h4 q-pa-md">Devices</div>
    <q-card class="shadow-2 card">
      <q-table :grid="grid" flat bordered :rows="data" :columns="columns" row-key="name" :filter="filter"
        :loading="loading" @request="reFetch" ref="eventTableRef" v-model:pagination="pagination">
        <template v-slot:top-right>
          <q-input debounce="300" v-model="filter" outlined rounded input-class="text-right"
            class="q-ml-md q-pr-lg search-input">
            <template v-slot:append>
              <q-icon v-if="filter === ''" name="search" />
              <q-icon v-else name="clear" class="cursor-pointer" @click="filter = ''" />
            </template>
          </q-input>

          <q-btn round color="primary" @click="grid = !grid" size="sm">
            <q-icon name="table_rows" v-if="grid" />
            <q-icon name="dashboard" v-else />
          </q-btn>

        </template>


        <template v-slot:item="props">
          <div class="q-pa-lg col-sm-12 col-md-6 col-lg-5 col-xl-4">
            <q-card class="card shadow-3">
              <q-item>
                <q-item-section class="col-3 vendor-avatar">
                  <DeviceVendorImage vendor="huawei" />
                </q-item-section>

                <q-item-section>
                  <q-item-label>
                  </q-item-label>
                  <q-item-label class="text-weight-medium">
                    <TextCopy :text="props.row.hostname" />
                  </q-item-label>
                  <q-item-label caption>
                    <TextCopy :text="props.row.managementIp" />
                  </q-item-label>

                </q-item-section>
                <DeviceStateBadge :state="props.row.state" size="sm" />
                <DeviceStatusBadge :status="props.row.status" size="sm" />
              </q-item>


              <q-separator class="q-mt-sm q-mb-sm" />



              <q-card-section horizontal>
                <q-card-section class="q-pt-xs col-8">
                  <q-list>
                    <q-item v-if="props.row.lastSeen != ''">
                      <q-item-section avatar>
                        <q-icon color="primary" name="timer" />
                      </q-item-section>
                      <q-item-section>
                        <q-item-label class="text-weight-medium">Last Seen</q-item-label>
                        <q-item-label caption>
                          {{ props.row.lastSeen }}
                        </q-item-label>
                      </q-item-section>
                    </q-item>

                    <q-item v-if="props.row.lastReboot != ''">
                      <q-item-section avatar>
                        <q-icon color="primary" name="restart_alt" />
                      </q-item-section>
                      <q-item-section>
                        <q-item-label class="text-weight-medium">Last Reboot</q-item-label>
                        <q-item-label caption>
                          {{ props.row.lastReboot }}
                        </q-item-label>
                      </q-item-section>
                    </q-item>

                    <q-item v-if="props.row.networkRegion != ''">
                      <q-item-section avatar>
                        <q-icon color="primary" name="hub" />
                      </q-item-section>
                      <q-item-section>
                        <q-item-label class="text-weight-medium">Network Region</q-item-label>
                        <q-item-label caption>
                          <TextCopy :text="props.row.networkRegion" />
                        </q-item-label>
                      </q-item-section>
                    </q-item>
                    <q-item v-if="props.row.model != ''">
                      <q-item-section avatar>
                        <q-icon color="primary" name="tag" />
                      </q-item-section>
                      <q-item-section>
                        <q-item-label class="text-weight-medium">Model</q-item-label>
                        <q-item-label caption>
                          <TextCopy :text="props.row.model" />
                        </q-item-label>
                      </q-item-section>
                    </q-item>
                    <q-item v-if="props.row.version != ''">
                      <q-item-section avatar>
                        <q-icon color="primary" name="class" />
                      </q-item-section>
                      <q-item-section>
                        <q-item-label class="text-weight-medium">Version</q-item-label>
                        <q-item-label caption lines="2">
                          <TextCopy :text="props.row.version" />
                        </q-item-label>
                      </q-item-section>
                    </q-item>
                  </q-list>
                </q-card-section>

                <q-separator vertical v-if="props.row.loading" />
                <q-card-section class="col-4 items-center q-pa-xl" v-if="props.row.loading">
                  <q-item-label caption class="q-mb-sm">
                    Getting data from device...
                  </q-item-label>
                  <q-spinner-ios color="primary" size="50px" />
                </q-card-section>
              </q-card-section>



              <q-separator />
              <q-card-actions>
                <q-btn flat icon="fa-solid fa-file-code">
                  <q-tooltip class="shadow-4">Collect configuration from device</q-tooltip>
                </q-btn>
                <q-btn flat icon="query_stats">
                  <q-tooltip class="shadow-4">Collect information about device</q-tooltip>
                </q-btn>

                <q-space />
                <q-btn flat icon="fa-solid fa-square-up-right" :to="'/devices/' + props.row.id">
                  <q-tooltip class="shadow-4">Open device to see details</q-tooltip>
                </q-btn>
              </q-card-actions>
            </q-card>

          </div>
        </template>

        <template v-slot:no-data="{ message }">
          <div class="full-width row flex-center text-accent q-gutter-sm">
            <q-icon size="2em" name="sentiment_dissatisfied" />
            <span>
              Well this is sad... {{ message }}
            </span>
          </div>
        </template>

      </q-table>
    </q-card>
  </div>
</template>

<style lang="sass" scoped>
.search-input
  width: 350px

.card
  padding: 10px
  width: 100%
  border-radius: 15px
.vendor-avatar
  max-height: 80px
</style>
