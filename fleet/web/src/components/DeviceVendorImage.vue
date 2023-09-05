<script setup lang="ts">

const props = defineProps<{
  vendor: string
}>()

// vendorImage is a function that get the prop vendor and switch cases all the possible values
// and return an image matching the vendor
const vendorImage = (vendor: string) => {

  if (vendor == undefined) {
    return 'img/vendor/vendor_unknown.svg'
  }

  switch (vendor.toLowerCase()) {
    case 'cisco':
      return '/img/vendor_logos/cisco.svg'
    case 'juniper':
      return '/img/vendor_logos/juniper.svg'
    case 'arista':
      return '/img/vendor_logos/arista.svg'
    case 'mikrotik':
      return '/img/vendor_logos/mikrotik.svg'
    case 'ubiquiti':
      return '/img/vendor_logos/ubiquiti.svg'
    case 'huawei':
      return '/img/vendor_logos/huawei.svg'
    case 'dell':
      return '/img/vendor_logos/dell.svg'
    case 'hp':
      return '/img/vendor_logos/hp.svg'
    case 'fortinet':
      return '~assets/vendor_logos/fortinet.svg'
    case 'palo alto':
      return '~assets/vendor_logos/alto.svg'
    default:
      return 'img/vendor/vendor_unknown.svg'
  }
}
</script>

<template>
  <q-img class="rounded-borders" :src="vendorImage(props.vendor)" fit="scale-down"/>
</template>

<style lang="sass">
</style>
