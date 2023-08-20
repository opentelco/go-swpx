<template>
  <q-page class="q-pa-md ">
    <div v-if="loading" class="text-h6">Loading...</div>

    <q-banner inline-actions class="text-white bg-red" v-if="error">
      Could not fetch device: {{ error }}
    </q-banner>
    <div v-if="result" class="q-pa-md example-row-variable-width">
      <div class="row">
        <div class="text-h4 q-pa-md col-12">
          {{ result.device.hostname }}
        </div>
        <div class="row">
          <DeviceStatusBadge :status="result.device.status" />
          <DeviceStateBadge :state="result.device.state" />
        </div>
        <q-card class="shadow-2 card">
          <q-tabs v-model="tab" dense class="text-grey" active-color="primary" indicator-color="primary" align="justify">
            <q-route-tab label="info" name="info" icon="mail" :to="tabRouteTo('')" />
            <q-route-tab label="events" name="events" icon="event" :to="tabRouteTo('events')" />
            <q-route-tab label="configuration" name="configuration" icon="info" :to="tabRouteTo('configuration')" />
            <q-route-tab label="changes" name="changes" icon="info" :to="tabRouteTo('changes')" />
            <q-route-tab label="stanza" name="stanza" icon="info" :to="tabRouteTo('stanza')" />
          </q-tabs>
          <q-separator />
          <q-tab-panels v-model="tab" animated>
            <q-tab-panel name="info">
              <div class="row">
                <div class="col-xl-3 col-lg-4 col-md-6">
                  <div class="text-h6">Info</div>
                  <DeviceCardInfo :device="result.device" />
                </div>
                <div class="col-xl-3 col-lg-4 col-md-6">
                  <div class="text-h6">Status</div>

                </div>
              </div>

            </q-tab-panel>

            <q-tab-panel name="events">
              DevicePageLimit: {{  eventLimit  }}
              <DeviceEventsTable :events="result.device.events" :limit="eventLimit" :offset="eventOffset"/>
            </q-tab-panel>

            <q-tab-panel name="configuration">
              <DeviceConfigHistory :configs="result.device.configurations"/>
            </q-tab-panel>

            <q-tab-panel name="changes">
              <div class="text-h6">changes</div>
              Lorem ipsum dolor sit amet consectetur adipisicing elit.
            </q-tab-panel>

            <q-tab-panel name="stanza">
              <div class="text-h6">stanza</div>
              Lorem ipsum dolor sit amet consectetur adipisicing elit.
            </q-tab-panel>

          </q-tab-panels>
        </q-card>

      </div>
    </div>

  </q-page>
</template>

<script setup lang="ts">

import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { Device } from '../gql/graphql'
import { useRoute } from 'vue-router'
import { ref } from 'vue'
import DeviceCardInfo from '../components/DeviceCardInfo.vue'
import DeviceStatusBadge from '../components/DeviceStatusBadge.vue'
import DeviceStateBadge from '../components/DeviceStateBadge.vue'
import DeviceConfigHistory from '../components/DeviceConfigHistory.vue'
import DeviceEventsTable from '../components/DeviceEventsTable.vue'

const route = useRoute()

const tabRouteTo = (tab: string) => {
  return '/devices/' + route.params.id + '/' + tab
}

// remap the response so it matches the response
type response = {
  device: Device
}

let tab = ref('info')

let getDevice = gql`query Device (
  $deviceId: ID!
  $stanzaLimit: Int
  $stanzaOffset: Int
  $configLimit: Int
  $configOffset: Int
  $changesLimit: Int
  $changesOffset: Int
  $eventLimit: Int
  $eventOffset: Int,
  ) {
    device(id: $deviceId) {
        id
        hostname
        domain
        managementIp
        serialNumber
        model
        version
        networkRegion
        pollerResourcePlugin
        pollerProviderPlugin
        state
        status
        lastSeen
        createdAt
        updatedAt
        lastReboot
        stanzas(limit: $stanzaLimit, offset: $stanzaOffset) {
            stanzas {
                id
                name
                description
                template
                content
                revert_template
                revert_content
                device_type
                updatedAt
                createdAt
                appliedAt
            }
            pageInfo {
                limit
                offset
                total
                count
            }
        }
        configurations(limit: $configLimit, offset: $configOffset ) {
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
        changes(limit: $changesLimit, offset: $changesOffset) {
            changes {
                id
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
        events(limit: $eventLimit, offset: $eventOffset) {
            events {
                id
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
        schedules {
            interval
            type
            lastRun
            active
            failedCount
        }
      }
}`


let eventLimit = ref(10)
let eventOffset = ref(0)

let { result, loading, error } = useQuery<response>(getDevice, {
  deviceId: route.params.id,
  stanzaLimit: 10,
  stanzaOffset: 0,
  configLimit: 10,
  configOffset: 0,
  changesLimit: 10,
  changesOffset: 0,
  eventLimit: eventLimit,
  eventOffset: eventOffset
})

</script>

<style lang="sass">
.card
  padding: 10px
  width: 100%
  border-radius: 15px


.example-row-variable-width1
  .row
    background: rgba(#aa0, .1)
  .row > div
    padding: 10px 15px
    background: rgba(#999,.15)
    border: 1px solid rgba(#999,.2)
  .row + .row
    margin-top: 1rem
</style>
