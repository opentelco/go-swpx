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
      return '~assets/vendor_cisco.svg'
    case 'juniper':
      return '~assets/vendor_juniper.svg'
    case 'arista':
      return '~assets/vendor_arista.svg'
    case 'mikrotik':
      return '~assets/vendor_mikrotik.svg'
    case 'ubiquiti':
      return '~assets/vendor_ubiquiti.svg'
    case 'huawei':
      return '~assets/vendor_huawei.svg'
    case 'dell':
      return '~assets/vendor_dell.svg'
    case 'hp':
      return '~assets/vendor_hp.svg'
    case 'fortinet':
      return '~assets/vendor_fortinet.svg'
    case 'palo alto':
      return '~assets/vendor_paloalto.svg'
    default:
      return 'img/vendor/vendor_unknown.svg'
  }
}
</script>

<template>
  <q-img class="rounded-borders" :src="vendorImage(props.vendor)" />
</template>
