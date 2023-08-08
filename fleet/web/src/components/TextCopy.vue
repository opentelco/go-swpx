<script setup lang="ts">

import { ref } from 'vue'
import { useQuasar, copyToClipboard } from 'quasar'

const $q = useQuasar()


const props = defineProps<{
  text: string
}>()

const copy = (text: string) => {
  copyToClipboard(text).then(() => {
    $q.notify({
      type: 'positive',
      message: text + ' copied to clipboard',
      icon: 'announcement'
    })
  }).catch(() => {
    $q.notify({
      type: 'negative',
      message: text + ' copied to clipboard',
      icon: 'announcement'
    })
  })
}

let hover = ref(false)

</script>

<template>
  <div>
    <div @mouseover="hover = true" @mouseleave="hover = false" @mousedown="copy(props.text)" class="hoverEffect cursor-pointer" :data-replace="props.text">
      <span>{{ props.text }}</span>

      <!-- <q-tooltip anchor="bottom middle" self="top middle">
        <q-btn v-if="hover" dense flat round icon="content_copy" size="xs" /> Click to copy to clipboard
      </q-tooltip> -->

    </div>

  </div>
</template>

<style>
.hoverEffect {
  overflow: hidden;
  position: relative;
  display: inline-block;
  line-height: 1.5;
}

.hoverEffect::before,
.hoverEffect::after {
 content: '';
  position: absolute;
  width: 100%;
  left: 0;
}
.hoverEffect::before {
  background-color: #276c85;

  height: 2px;
  bottom: 0;
  transform-origin: 100% 50%;
  transform: scaleX(0);
  transition: transform .3s cubic-bezier(0.76, 0, 0.24, 1);
}
.hoverEffect::after {
  content: attr(data-replace);
  height: 100%;
  top: 0;
  transform-origin: 100% 50%;
  transform: translate3d(200%, 0, 0);
  transition: transform .3s cubic-bezier(0.76, 0, 0.24, 1);
  color: #276c85;
}

.hoverEffect:hover::before {
  content: "\E84E";
  transform-origin: 0% 50%;
  transform: scaleX(1);
}
.hoverEffect:hover::after {
  transform: translate3d(0, 0, 0);
}

.hoverEffect span {
  display: inline-block;
  transition: transform .3s cubic-bezier(0.76, 0, 0.24, 1);
}

.hoverEffect:hover span {
  transform: translate3d(-200%, 0, 0);
}
</style>
