<script setup lang="ts">
import { EventConnection } from '../gql/graphql'

import DeviceEventOutcomeChip from './DeviceEventOutcomeChip.vue'
import DeviceEventsTypeChip from './DeviceEventsTypeChip.vue';
import DeviceEventActionChip from './DeviceEventActionChip.vue';
import { computed, ref, toRef } from 'vue'
import { cpuUsage } from 'process';

let search = ref('')
let filter = ref(null)

const props = defineProps<{
  events: EventConnection
  limit: number
  offset: number
}>()

const propsLimit = toRef(props, 'limit');
const propsOffset = toRef(props, 'offset');
const propsEvents = toRef(props, 'events')

propsLimit.value = 20

const columns = [
  { name: 'timestamp', label: 'Timestamp', field: 'createdAt', align: 'left' },
  { name: 'message', label: 'Message', field: 'message', align: 'left' },
  { name: 'type', label: 'Type', field: 'type', align: 'left' },
  { name: 'outcome', label: 'Outcome', field: 'outcome', align: 'left' },
  { name: 'action', label: 'Action', field: 'action', align: 'left' },
]



const options = [
  {
    label: 'Google',
    value: 'goog',
    icon: 'mail'
  },
  {
    label: 'Facebook',
    value: 'fb',
    icon: 'bluetooth'
  },
  {
    label: 'Twitter',
    value: 'twt',
    icon: 'map'
  },
  {
    label: 'Apple',
    value: 'app',
    icon: 'golf_course'
  },
  {
    label: 'Oracle',
    value: 'ora',
    disable: true,
    icon: 'casino'
  }
]


function calculatePagination(limit: number, offset: number, totalCount: number | null | undefined): { page: number, rowsPerPage: number, rowsNumber: number } {
  const rowsNumber = totalCount;
  const rowsPerPage = limit;

  const page = limit > 0 ? Math.floor(offset / limit) + 1 : 1;
  return { page, rowsPerPage, rowsNumber };
}

// computed pagination with the help of the calculatePagination fuction
const pagination = computed(() => {
  return calculatePagination(propsLimit.value, propsOffset.value, propsEvents.value.pageInfo.total)
});

const reFetch = ((tp) => {
  const { page, rowsPerPage, sortBy, descending } = tp.pagination
  const f = tp.filter

  localPagination.value.page = page
  localPagination.value.rowsPerPage = rowsPerPage

  propsLimit.value = rowsPerPage
  propsOffset.value = localPagination.value.page*rowsPerPage


  search.value = f
  console.log(search)


})



const localPagination = ref({
  sortBy: 'desc',
  descending: false,
  page: 1,
  rowsPerPage: 10,
  rowsNumber: 50
})



</script>

<template>
  {{ pagination }}
  <br>
  {{ localPagination }}
  <div class="q-pa-md">
    <div class="row justify-center">
      <div class="">
        <!-- <Timeline size="md"/> -->
      </div>
    </div>
    <div class="row">
      <div class="col col-12">
        <q-table title="Events" :rows="props.events.events" :columns="columns" row-key="id"
          class="q-pa-lg rounded shadow-3" :filter="search" :pagination="localPagination" @request="reFetch">

          <template v-slot:top-right>
            <q-input debounce="300" v-model="search" outlined rounded input-class="text-right"
              class="q-ml-md q-pr-lg search-input">
              <template v-slot:append>
                <q-icon v-if="search === ''" name="search" />
                <q-icon v-else name="clear" class="cursor-pointer" @click="search = ''" />
              </template>
            </q-input>


            <q-select color="orange" multiple filled v-model="filter" :options="options" stack-label>

              <template v-if="filter" v-slot:append>
                <q-icon name="cancel" @click.stop.prevent="filter = null" class="cursor-pointer" />
              </template>

              <template v-slot:selected-item="scope">
                <q-chip removable dense square color="white" text-color="primary" class=""
                  @remove="scope.removeAtIndex(scope.index)" :tabindex="scope.tabindex">
                  <q-avatar color="primary" text-color="white" :icon="scope.opt.icon" />
                  {{ scope.opt.label }}
                </q-chip>
              </template>
            </q-select>




          </template>

          <template v-slot:body="p">

            <q-tr :props="props">

              <q-td key="timestamp" :props="p">
                {{ p.row.createdAt }}
              </q-td>

              <q-td key="message" :props="p">
                {{ p.row.message }}
              </q-td>

              <q-td key="type" :props="p">
                <DeviceEventsTypeChip :type="p.row.type" class="q-mr-sm" />

              </q-td>

              <q-td key="outcome" :props="p">
                <DeviceEventOutcomeChip :outcome="p.row.outcome" class="q-mr-sm" />
              </q-td>

              <q-td key="action" :props="p">
                <DeviceEventActionChip :action="p.row.action" class="q-mr-sm" />
              </q-td>

            </q-tr>
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
