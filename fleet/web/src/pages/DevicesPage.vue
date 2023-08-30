<template>
  <div class="q-pa-md q-gutter-sm">
    <q-breadcrumbs align="left" class="q-ml-md">
      <q-breadcrumbs-el icon="home" to="/" />
      <q-breadcrumbs-el label="Devices" icon="fa-solid fa-network-wired" to="/devices" />
    </q-breadcrumbs>
  </div>
  <div class="q-pa-md">
    <DeviceList />
  </div>
</template>

<script setup lang="ts">
import DeviceList from '../components/DeviceList.vue'



</script>
