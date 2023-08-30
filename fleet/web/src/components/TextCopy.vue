<script setup lang="ts">

import { ref, useSlots } from 'vue'
import { useQuasar, copyToClipboard } from 'quasar'

const $q = useQuasar()

let copied = ref(false)
let slots = useSlots()

const props = defineProps<{
  text: string
  toCopy?: string
  noIcon?: boolean
}>()


// copy toCopy if defined else text
const whatText = () => {
  slots.default?.()[0]
  if (props.toCopy) {
    return props.toCopy
  } else {
    return props.text
  }
}

const copy = (text: string) => {
  copyToClipboard(text).then(() => {
    copied.value = true
    setTimeout(() => {
      copied.value = false
      hover.value = false
    }, 1000)

    // $q.notify({
    //   type: 'positive',
    //   message: text + ' copied to clipboard',
    //   icon: 'announcement'
    // })
  }).catch(() => {
    // $q.notify({
    //   type: 'negative',
    //   message: text + ' copied to clipboard',
    //   icon: 'announcement'
    // })
  })
}

let hover = ref(false)

</script>

<template>
  <div @mouseover="hover = true" @mouseleave="hover = false" @mousedown="copy(whatText())"
    class="hoverEffect cursor-pointer copyDiv" :data-replace="props.text">
    <div class="text-copy">
      <span>
        {{ props.text }}
      </span>
      <span class="icon-copy">
      <q-icon v-show="hover && !noIcon && !copied" name="content_copy" class="q-ml-sm" />
      <q-icon name="check" color="green" class="q-ml-sm" v-show="!noIcon && copied" />
    </span>
    </div>
  </div>
</template>

<style>
.icon-copy {
  width: 20px + 5px;
  height: auto;
  float: right;

}

.floating {
  position: absolute;
  background-color: white;
  border-radius: 5px;
  padding: 5px;
  box-shadow: 0px 0px 5px 0px rgba(0, 0, 0, 0.75);
  animation: fadein 0.5s;
  z-index: 1000;
  float: right

}

.copyDiv {
  width: fit-content;
  overflow: auto;

}

.hoverEffect:hover {
  animation: shake 0.5s;
}

@keyframes shake {
  0% {
    transform: translateX(0)
  }

  25% {
    transform: translateX(3px);
  }


  50% {
    transform: translateX(-3px);
  }

  100% {
    transform: translateX(0px);
  }
}
</style>
