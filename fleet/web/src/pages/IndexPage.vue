<template>
  <q-page class="row items-center justify-evenly">
    <q-input
      rounded
      outlined
      v-model="text"
      input-class="text-right"
      class="q-ml-md"
      style="width: 400px"
      hint="what do you want to do today?"
    >
      <template v-slot:append>
        <q-icon v-if="text === ''" name="search" />
        <q-icon v-else name="clear" class="cursor-pointer" @click="text = ''" />
      </template>
    </q-input>
  </q-page>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue';

export default defineComponent({
  setup() {
    return {
      text: ref(''),
    };
  },
  name: 'IndexPage',
  components: {},
});
</script>
