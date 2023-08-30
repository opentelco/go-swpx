<template>
  <q-banner inline-actions class="text-white bg-red" v-if="error">
    Could not fetch device: {{ error }}
  </q-banner>
  <div class="q-pa-md ">
    <q-breadcrumbs align="left" class="q-ml-md">
      <q-breadcrumbs-el icon="home" to="/" />
      <q-breadcrumbs-el label="Devices" icon="fa-solid fa-network-wired" to="/devices" />
      <q-breadcrumbs-el :label="result?.device.hostname" icon="fa-solid fa-cube" :to="'/devices/' + route.params.id" />
      <q-breadcrumbs-el label="configuration" icon="fa-regular fa-file-code"
        :to="'/devices/' + route.params.id + '/configuration'" v-if="route.name == 'deviceConfiguration'" />
      <q-breadcrumbs-el label="events" icon="fa-solid fa-envelope-open-text"
        :to="'/devices/' + route.params.id + '/events'" v-if="route.name == 'deviceEvents'" />
      <q-breadcrumbs-el label="changes" icon="fa-solid fa-pen-to-square" :to="'/devices/' + route.params.id + '/changes'"
        v-if="route.name == 'deviceChanges'" />
        <q-breadcrumbs-el label="timeline" icon="fa-solid fa-pen-to-square" :to="'/devices/' + route.params.id + '/timeline'"
        v-if="route.name == 'deviceTimeline'" />
      <q-breadcrumbs-el label="stanzas" icon="data_object" :to="'/devices/' + route.params.id + '/stanza'"
        v-if="route.name == 'deviceStanza'" />
    </q-breadcrumbs>
  </div>
  <div v-if="result" class="q-px-lg q-pb-xl">

    <q-card class="q-pa-xl rounder-borders shadow-5" bordered style="max-width:1600px">

      <div class="row justify-between">
        <div class="col col-md-auto self-start q-pa-md">
          <q-item dense>
            <q-item-section avatar class="q-mr-md">
              <DeviceVendorImage vendor="huawei" style="size: 80px" />
            </q-item-section>
            <q-item-section style="min-width:300px">
              <q-item-label class="text-weight-medium text-h5">
                <TextCopy :text="result.device.hostname" class="q-mr-md" />
              </q-item-label>
              <q-item-label caption class>
                <TextCopy :text="result.device.managementIp" />
              </q-item-label>
            </q-item-section>
            <q-item-section side top>
              <div class="badge-header">
                <DeviceStatusBadge :status="result.device.status" />
                <DeviceStateBadge :state="result.device.state" />
              </div>
            </q-item-section>
          </q-item>

        </div>
        <div class="col col-2 q-pa-md">
          <q-item clickable>
            <q-item-section avatar>
              <q-icon color="primary" name="running_with_errors" />
            </q-item-section>
            <q-item-section>
              <q-item-label class="text-weight-medium ">
                Workflows
              </q-item-label>
              <q-item-label caption class>
                <span v-if="result?.device.lastSeen != ''">
                  X running
                </span>
                <span class="" v-else>
                  not set
                </span>
              </q-item-label>
            </q-item-section>
            <q-item-secction side>
              <q-btn round  icon="fa-solid fa-external-link" size="xs" color="primary" class="q-pa-sm"/>
            </q-item-secction>
          </q-item>
        </div>

        <div class="col col-2 q-pa-md">
          <q-item>
            <q-item-section avatar>
              <q-icon color="primary" name="fa-regular fa-clock" />
            </q-item-section>
            <q-item-section>
              <q-item-label class="text-weight-medium ">Last Seen</q-item-label>
              <q-item-label caption class>
                <span v-if="result?.device.lastSeen != ''">
                  {{ getRelativeTimestamp(result?.device.lastSeen) }}
                  <q-tooltip>
                    {{ result?.device.lastSeen }}
                  </q-tooltip>
                </span>
                <span class="" v-else>
                  not set
                </span>
              </q-item-label>
            </q-item-section>
          </q-item>
        </div>

      </div>
      <div class="row">
        <div class="col-12">

          <q-tabs v-model="tab" class="text-grey" active-color="primary" indicator-color="primary" align="justify">
            <q-route-tab label="device" name="spotlight" icon="fa-regular fa-lightbulb" :to="tabRouteTo('')" />
            <q-route-tab label="ports" name="ports" icon="fa-solid fa-ethernet" disabled/>


            <q-route-tab label="configuration" name="configuration" icon="fa-regular fa-file-code"
              :to="tabRouteTo('configuration')">
              <q-badge color="deep-orange" floating v-if="!hasConfigurations()">
                <q-icon name="warning" />
              </q-badge>
            </q-route-tab>
            <q-route-tab label="stanza" name="stanza" icon="data_object" :to="tabRouteTo('stanza')" />
            <q-route-tab label="timeline" name="timeline" icon="fa-solid fa-timeline" :to="tabRouteTo('timeline')"/>
            <q-route-tab label="changes" name="changes" icon="fa-solid fa-pen-to-square" :to="tabRouteTo('changes')">
              <q-badge color="deep-orange" floating v-if="!hasChanges()">
                <q-icon name="warning" />
              </q-badge>
            </q-route-tab>
            <q-route-tab label="events" name="events" icon="fa-solid fa-envelope-open-text" :to="tabRouteTo('events')" />

          </q-tabs>
          <q-separator />
          <q-tab-panels v-model="tab" animated  transition-prev="fade" transition-next="fade">
            <q-tab-panel name="spotlight">
              <div class="row q-pb-md">
                <div class="col col-4">
                  <q-list dense>
                    <q-item>
                      <q-item-section avatar>
                        <q-icon color="primary" name="fa-solid fa-fingerprint" />
                      </q-item-section>
                      <q-item-section>
                        <q-item-label class="text-weight-medium ">Fingerprint</q-item-label>
                        <q-item-label caption class>
                          <span v-if="result?.device.serialNumber != ''">
                            {{ result?.device.serialNumber }}
                          </span>
                          <span class="" v-else>
                            not set
                          </span>
                        </q-item-label>
                      </q-item-section>
                    </q-item>
                    <q-item class="q-mt-sm">
                      <q-item-section avatar><i class=""></i>
                        <q-icon color="primary" name="fa-regular fa-registered" />
                      </q-item-section>
                      <q-item-section>
                        <q-item-label class="text-weight-medium ">Model</q-item-label>
                        <q-item-label caption class>
                          <span v-if="result?.device.model != ''">
                            {{ result?.device.model }}
                          </span>
                          <span class="" v-else>
                            not set
                          </span>
                        </q-item-label>
                      </q-item-section>
                    </q-item>
                    <q-item class="q-mt-sm">
                      <q-item-section avatar><i class=""></i>
                        <q-icon color="primary" name="flash_on" />
                      </q-item-section>
                      <q-item-section>
                        <q-item-label class="text-weight-medium ">SWP Plugin (Resource / Provider) </q-item-label>
                        <q-item-label caption class>
                          <span v-if="result?.device.pollerResourcePlugin != ''">
                            {{ result?.device.pollerResourcePlugin }}
                            <q-tooltip> Resource plugin: {{ result?.device.pollerResourcePlugin }}</q-tooltip>
                          </span>
                          <span class="" v-else>
                            not set
                          </span>
                          /
                          <span v-if="result?.device.pollerProviderPlugin != ''">
                            {{ result?.device.pollerProviderPlugin }}
                            <q-tooltip> Provider plugin: {{ result?.device.pollerProviderPlugin }}</q-tooltip>
                          </span>
                          <span class="" v-else>
                            not set
                          </span>
                        </q-item-label>
                      </q-item-section>
                    </q-item>
                  </q-list>
                </div>
                <div class="col col-4">

                  <q-card class="rounder-borders" flat>
                    <q-card-section>
                      <q-card-section class="q-pt-xs col-8">
                        <q-list>
                          <q-item v-if="result?.device.lastSeen != ''">
                            <q-item-section avatar>
                              <q-icon color="primary" name="fa-regular fa-clock" />
                            </q-item-section>
                            <q-item-section>
                              <q-item-label class="text-weight-medium">Last Seen</q-item-label>
                              <q-item-label caption>
                                {{ result?.device.lastSeen }}
                              </q-item-label>
                            </q-item-section>
                          </q-item>

                          <q-item v-if="result?.device.lastReboot != ''">
                            <q-item-section avatar>
                              <q-icon color="primary" name="restart_alt" />
                            </q-item-section>
                            <q-item-section>
                              <q-item-label class="text-weight-medium">Last Reboot</q-item-label>
                              <q-item-label caption>
                                {{ result?.device.lastReboot }}
                              </q-item-label>
                            </q-item-section>
                          </q-item>


                          <q-item v-if="result?.device.model != ''">
                            <q-item-section avatar>
                              <q-icon color="primary" name="tag" />
                            </q-item-section>
                            <q-item-section>
                              <q-item-label class="text-weight-medium">Model</q-item-label>
                              <q-item-label caption>
                                <TextCopy :text="result?.device.model" />
                              </q-item-label>
                            </q-item-section>
                          </q-item>
                          <q-item v-if="result?.device.version != ''">
                            <q-item-section avatar>
                              <q-icon color="primary" name="class" />
                            </q-item-section>
                            <q-item-section>
                              <q-item-label class="text-weight-medium">Version</q-item-label>
                              <q-item-label caption lines="2">
                                <TextCopy :text="result?.device.version" />
                              </q-item-label>
                            </q-item-section>
                          </q-item>
                        </q-list>
                      </q-card-section>
                    </q-card-section>
                  </q-card>
                </div>
                <div class="col col-4 ">
                </div>
              </div>
              <div class="row">
                <div class="col col-4 q-pa-md">
                </div>
                <div class="col col-4 q-pa-md">
                  <q-card class="rounded-borders ">
                    <q-card-section>
                      <q-card-section class="q-pt-xs col-8">
                        <q-list>
                          <q-item>
                            <q-item-section avatar>
                              <q-icon color="primary" name="dns" />
                            </q-item-section>
                            <q-item-section>
                              <q-item-label class="text-weight-medium">Domain</q-item-label>
                              <q-item-label caption>
                                {{ result?.device.domain }}
                              </q-item-label>
                            </q-item-section>
                          </q-item>

                          <q-item>
                            <q-item-section avatar>
                              <q-icon color="primary" name="fa-solid fa-fingerprint" />
                            </q-item-section>
                            <q-item-section>
                              <q-item-label class="text-weight-medium">Serial Number</q-item-label>
                              <q-item-label caption>
                                {{ result?.device.serialNumber }}
                              </q-item-label>
                            </q-item-section>
                          </q-item>

                          <q-item>
                            <q-item-section avatar>
                              <q-icon color="primary" name="hub" />
                            </q-item-section>
                            <q-item-section>
                              <q-item-label class="text-weight-medium">Network Region</q-item-label>
                              <q-item-label caption>
                                <TextCopy :text="result?.device.networkRegion" />
                              </q-item-label>
                            </q-item-section>
                          </q-item>
                          <q-item v-if="result?.device.model != ''">
                            <q-item-section avatar>
                              <q-icon color="primary" name="tag" />
                            </q-item-section>
                            <q-item-section>
                              <q-item-label class="text-weight-medium">Model</q-item-label>
                              <q-item-label caption>
                                <TextCopy :text="result?.device.model" />
                              </q-item-label>
                            </q-item-section>
                          </q-item>
                          <q-item v-if="result?.device.version != ''">
                            <q-item-section avatar>
                              <q-icon color="primary" name="class" />
                            </q-item-section>
                            <q-item-section>
                              <q-item-label class="text-weight-medium">Version</q-item-label>
                              <q-item-label caption lines="2">
                                <TextCopy :text="result?.device.version" />
                              </q-item-label>
                            </q-item-section>
                          </q-item>
                        </q-list>
                      </q-card-section>
                    </q-card-section>
                  </q-card>

                </div>
              </div>

            </q-tab-panel>

            <q-tab-panel name="stanza">
              <div class="text-h6">stanza</div>
              Lorem ipsum dolor sit amet consectetur adipisicing elit.
            </q-tab-panel>

            <q-tab-panel name="configuration">
              <DeviceConfigTable :deviceId="route.params.id" />
            </q-tab-panel>

            <q-tab-panel name="timeline">
              <DeviceTimeline :deviceId="route.params.id" />
            </q-tab-panel>

            <q-tab-panel name="changes">
              Changes in a table
            </q-tab-panel>
            <q-tab-panel name="events">
              <DeviceEventsTable :deviceId="route.params.id" />
            </q-tab-panel>



          </q-tab-panels>
        </div>
      </div>
    </q-card>
  </div>
</template>

<script setup lang="ts">

import { useQuery } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { Device } from '../gql/graphql'
import { useRoute } from 'vue-router'
import { ref } from 'vue'
import { getRelativeTimestamp } from '../utils/time'

import DeviceStatusBadge from '../components/DeviceStatusBadge.vue'
import DeviceStateBadge from '../components/DeviceStateBadge.vue'
import DeviceConfigTable from '../components/DeviceConfigTable.vue'
import DeviceEventsTable from '../components/DeviceEventsTable.vue'
import DeviceTimeline from 'src/components/DeviceTimeline.vue'
import DeviceVendorImage from 'src/components/DeviceVendorImage.vue'
import TextCopy from 'src/components/TextCopy.vue'

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

// hasConfigurations checks if the device has configurations and returns true or false
const hasConfigurations = () => {
  if (result) {
    if (result.value?.device == null) {
      return false
    }
    if (result.value?.device.configurations == null) {
      return false
    }
    if (result.value?.device.configurations.pageInfo == null) {
      return false
    }
    if (result.value?.device.configurations.pageInfo.total == null) {
      return false
    }
    return result.value?.device.configurations.pageInfo.total > 0 ? true : false
  }
}

const hasChanges = () => {
  if (result) {
    if (result.value?.device == null) {
      return false
    }
    if (result.value?.device.changes == null) {
      return false
    }
    if (result.value?.device.changes.pageInfo == null) {
      return false
    }
    if (result.value?.device.changes.pageInfo.total == null) {
      return false
    }
    return result.value?.device.changes.pageInfo.total > 0 ? true : false
  }
}

</script>

<style lang="sass">
.badge-header
  display: flex

</style>
