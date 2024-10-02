<template>
  <v-app-bar :elevation="0" rounded class="border-b">
    <v-app-bar-nav-icon v-if="is_mob"></v-app-bar-nav-icon>
    <v-app-bar-title>
      <h2>高三3班</h2>
    </v-app-bar-title>
  </v-app-bar>
</template>

<script setup>
import { useDisplay } from 'vuetify'
import { computed } from 'vue'
const emit = defineEmits(['update:is_show'])
const { mobile } = useDisplay()
const is_mob = computed(() => {
  return mobile.value
})

import {useRouter} from "vue-router";
const router = useRouter()
</script>
