
<template>
  <div>
    <q-banner inline-actions class="text-white bg-red" v-if="error">
      Could not fetch devices: {{ error }}
    </q-banner>
    <div class="q-pa-md">
      <q-table :grid="grid" flat bordered title="Devices" v-if="result" :rows="result?.devices?.devices" :columns="columns"
        row-key="name" :filter="filter" :loading="loading">
        <template v-slot:top-right>

          <q-input debounce="300" v-model="filter" outlined rounded input-class="text-right"
            class="q-ml-md q-pr-lg search-input">
            <template v-slot:append>
              <q-icon v-if="filter === ''" name="search" />
              <q-icon v-else name="clear" class="cursor-pointer" @click="filter = ''" />
            </template>
          </q-input>

          <q-btn round color="primary" @click="grid = !grid">
            <q-icon name="table_rows" v-if="grid" />
            <q-icon name="dashboard" v-else />
          </q-btn>

        </template>


        <template v-slot:item="props">
          <div class="q-pa-xs col-xs-12 col-sm-6 col-md-4">
            <q-card class="my-card" bordered>
              <q-card-section horizontal>
                <q-card-section class="q-pt-xs">

                  <div class="q-pt-sm">
                    <DeviceStateBadge :state="props.row.state" class="" />
                    <DeviceStatusBadge :status="props.row.status" />
                  </div>

                  <div class="text-h5 q-mt-sm q-mb-xs">{{ props.row.hostname }}</div>
                  <q-list>
                    <q-item>
                      <q-item-section avatar>
                        <q-icon color="primary" name="public" />
                      </q-item-section>

                      <q-item-section>
                        <q-item-label>Management IP</q-item-label>
                        <q-item-label caption>
                          <TextCopy :text="props.row.managementIp" />
                        </q-item-label>
                      </q-item-section>

                    </q-item>

                    <q-item>
                      <q-item-section avatar>
                        <q-icon color="primary" name="hub" />
                      </q-item-section>
                      <q-item-section>
                        <q-item-label>Network Region</q-item-label>
                        <q-item-label caption>
                          <TextCopy :text="props.row.networkRegion" />
                        </q-item-label>
                      </q-item-section>

                    </q-item>


                  </q-list>

                  <div class="text-caption text-grey">
                  </div>
                </q-card-section>

                <q-card-section class="col-5 flex flex-center">
                  <DeviceVendorImage :vendor="props.row.model" />
                </q-card-section>
              </q-card-section>

              <q-separator />

              <q-card-actions >
                <q-btn flat icon="fa-solid fa-file-code" class="q-mr-md">
                  Collect Config
                </q-btn>

                <q-btn flat icon="query_stats">
                  Discover
                </q-btn>
                <q-space />
                <q-btn flat icon="fa-solid fa-square-up-right">
                  Open
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
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { ListDeviceResponse } from '../gql/graphql'
import DeviceStateBadge from './DeviceStateBadge.vue'
import DeviceStatusBadge from './DeviceStatusBadge.vue'
import TextCopy from './TextCopy.vue'
import DeviceVendorImage from './DeviceVendorImage.vue'

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

let searchInput = ref('')
let filter = ref('')


const onInput = (event) => {
  // Handle any immediate input logic here if needed
  console.log('Immediate input handling:', event.target.value);
};

const { result, loading, error } = useQuery<ListDeviceResponse>(gql`
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
  search: searchInput
})

</script>

<style lang="sass" scoped>
.search-input
  width: 300px

.my-card
  width: 100%
  max-width: 500px
  border-bottom-color: rgb(82, 84, 196)

</style>
