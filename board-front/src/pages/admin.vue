<template>
  <v-app>
    <v-navigation-drawer v-model="is_show" :location="mob ? 'top' : 'left'">
      <v-list density="compact" nav>
        <v-list-item v-bind:key='tab.to' v-for="tab in tabs" color="primary" :prepend-icon="tab.icon" :title="tab.title"
                     :to="'/' + tab.to"></v-list-item>
      </v-list>
    </v-navigation-drawer>
    <v-main>
      <router-view v-slot="{ Component }">
        <v-scroll-x-transition mode="out-in">
          <component :is="Component"/>
        </v-scroll-x-transition>
      </router-view>
    </v-main>
  </v-app>
</template>

<script setup>
import {ref, computed} from 'vue'
import {useDisplay} from 'vuetify'
import {useRoute, useRouter} from "vue-router";

const router = useRouter()
if (window.location.pathname === '/setting')
  router.push("/admin/schedule-setting")

const is_show = ref(false)
const {mobile} = useDisplay()
if (!mobile.value)
  is_show.value = true
const mob = computed(() => {
  return mobile.value
})

const tabs = [
  {
    title: "班级设置",
    icon: "mdi-segment",
    to: "admin/schedule-setting"
  },
  {
    title: "时间管理",
    icon: "mdi-heart",
    to: "admin/count-setting"
  },
  {
    title: "公告管理",
    icon: "mdi-school",
    to: "admin/notice-setting"
  },
  {
    title: "基础设置",
    icon: "mdi-bed",
    to: "admin/basic-setting"
  },
]

document.title = "后台管理";
</script>

<style>
</style>
