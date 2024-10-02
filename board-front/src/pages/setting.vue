<template>
  <v-app>
    <v-app-bar :elevation="0" rounded class="border-b">
      <v-app-bar-title class="flex">
        <h2>{{ className }}</h2>
      </v-app-bar-title>
      <v-btn @click="unBind">解绑</v-btn>
    </v-app-bar>
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
import AppBar from "@/components/AppBar.vue";
import axios from "axios";
function unBind() {
  axios.get('/api/device/un-bind')
  location.reload()
}
const className = ref(''); // 存储格式化的班级名称

function testBind() {
  axios.get('/api/device/get-class').then(res => {
    const grade = res.data.data.grade; // 获取年级，例如 1 代表高一
    const classNum = res.data.data.class; // 获取班级号，例如 1 代表 1 班

    // 将年级数字转换为对应的中文
    const gradeMap = {
      1: '一',
      2: '二',
      3: '三'
    };

    // 格式化班级名称为 "高一1班"
    className.value = `高${gradeMap[grade]}${classNum}班`;
  });
}


onMounted(() => {
  testBind()
})
const router = useRouter()
if (window.location.pathname === '/setting' || window.location.pathname === '/setting/')
  router.push("/setting/schedule-setting")

const is_show = ref(false)
const {mobile} = useDisplay()
if (!mobile.value)
  is_show.value = true
const mob = computed(() => {
  return mobile.value
})

const tabs = [
  {
    title: "课程设置",
    icon: "mdi-segment",
    to: "setting/schedule-setting"
  },
  {
    title: "时间管理",
    icon: "mdi-heart",
    to: "setting/count-setting"
  },
  {
    title: "公告管理",
    icon: "mdi-school",
    to: "setting/notice-setting"
  },
  {
    title: "基础设置",
    icon: "mdi-bed",
    to: "setting/basic-setting"
  },
]

document.title = "后台管理";
</script>

<style>
</style>
