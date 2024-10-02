<script setup>
import axios from 'axios';
import BaseData from "@/components/tools/BaseData.vue";
import WeatherWidget from "@/components/tools/WeatherWidget.vue";

let intervalId = null;
onMounted(() => {
  getCountDown()
})

onUnmounted(() => {
  if (intervalId) {
    clearInterval(intervalId);
  }
})

const noticeHTML = ref('')
const noticeSwitch = ref(false)
function getConfig() {
  axios.get("/api/setting/get-setting-value/?key=notice")
    .then((res) => {
      noticeHTML.value = res.data.data.Res.Value;
      //console.log(slogan_update_time.value);
    });
  axios.get("/api/setting/get-setting-value/?key=notice_switch")
    .then((res) => {
      noticeSwitch.value = res.data.data.Res.Value === 'true';
    });
  getCountDown()
}

onMounted(() => {
  // 初始调用
  getConfig();
  // 设置一个定时器来定期调用 getConfig
  setInterval(getConfig, 1000 * 60);
});

const count_down_list = ref('')

function getCountDown() {
  axios.get("/api/count-down/get-event").then((res) => {
    count_down_list.value = res.data.data.Res;
    count_down_list.value.forEach((item) => {
      item.left = getDayDiff(item.Date); // 使用 item.Date 计算并设置 left
    });

    // 过滤掉 left 小于 0 的元素，并对剩下的元素按 left 值从大到小排序
    count_down_list.value = count_down_list.value
      .filter(item => item.left >= 0)
      .sort((a, b) => b.left - a.left);

    //console.log(count_down_list.value)
  })
}

function getDayDiff(targetDate) {
  const now = new Date();
  const diff = new Date(targetDate).getTime() - now.getTime();
  return Math.floor(diff / (1000 * 3600 * 24))
}

</script>

<template>
  <div class="flex w-full" @click="getConfig">
    <div class="flex flex-col w-7/11 ml-2 mr-2">
      <div class=" w-full bg-transparent backdrop-blur-xl rounded-4 flex  flex-col pr-2 mb-2 pt-1 pb-1">
        <div
          class="h-32 rounded-2 backdrop-blur-xl ml-2 flex items-center justify-center flex-grow"
          v-for="index in count_down_list.length">
          <a class="text-white text-24 font-500">距离</a>
          <a class="text-white text-24 font-900">{{ count_down_list[index - 1].Event }}</a>
          <a class="text-white text-24 font-500">还有</a>
          <a class="text-white text-24 font-900">{{ count_down_list[index - 1].left }}</a>
          <a class="text-white text-24 font-500">天</a>
        </div>
      </div>
      <div class=" w-full bg-transparent backdrop-blur-xl rounded-4 flex pr-2 mb-2" v-if="noticeSwitch">
        <div class="w-16 rounded-2 backdrop-blur-xl mt-2 mb-2 ml-2 flex flex-col flex-grow p-2">
          <div class="flex justify-start mb-1">
            <a class="text-white font-700 w-full text-2xl" v-html="noticeHTML"/>
          </div>
        </div>
      </div>
    </div>
    <div class="flex flex-col w-4/11 mr-2">
      <BaseData :is-outside="true"/>
      <WeatherWidget/>
    </div>

  </div>


</template>

<style scoped>

</style>
