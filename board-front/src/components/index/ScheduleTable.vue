<script setup>
import axios from "axios";
import {ref, onMounted, onUnmounted, watch} from "vue";

// 当前时间
const currentTime = ref(new Date());

// 更新当前时间
const updateCurrentTime = () => {
  currentTime.value = new Date();
};

// 设置定时器
onMounted(() => {
  // 每秒更新当前时间
  const intervalId = setInterval(updateCurrentTime, 5000);

  // 清除定时器
  onUnmounted(() => {
    clearInterval(intervalId);
  });
});

// 定义状态
const all_schedule = ref('');
const today_schedule = ref('');
const start_schedule = ref([]);
const end_schedule = ref([]);

// 获取星期几
const daysOfWeek = ["周日", "周一", "周二", "周三", "周四", "周五", "周六"];
const dayName = ref(daysOfWeek[new Date().getDay()]);
function updateSchedule() {
  // 获取数据
  axios.get('/api/schedule/get-schedule').then(response => {
    all_schedule.value = JSON.parse(response.data.data.Res.Schedule);
    const this_day = new Date().getDay() || 7;
    today_schedule.value = all_schedule.value[this_day.toString()];
    start_schedule.value = all_schedule.value[8];
    end_schedule.value = all_schedule.value[9];
  });
}


// 判断当前时间是否在某个时间段内
function isWithinTimeRange(startTime, endTime) {
  const nowTime = currentTime.value.getHours() * 60 + currentTime.value.getMinutes();
  const [startHour, startMinute] = startTime.split(":").map(Number);
  const [endHour, endMinute] = endTime.split(":").map(Number);
  const startMinutes = startHour * 60 + startMinute;
  const endMinutes = endHour * 60 + endMinute;
  return nowTime >= startMinutes && nowTime <= endMinutes;
}

function isLargerTimeRange(startTime, endTime) {
  const nowTime = currentTime.value.getHours() * 60 + currentTime.value.getMinutes();
  const [startHour, startMinute] = startTime.split(":").map(Number);
  const [endHour, endMinute] = endTime.split(":").map(Number);
  const startMinutes = startHour * 60 + startMinute;
  const endMinutes = endHour * 60 + endMinute;
  return nowTime >= startMinutes && nowTime >= endMinutes;
}

function getStyle(index) {
  const isInTimeRange = isWithinTimeRange(start_schedule.value[index], end_schedule.value[index]);
  const isLargerRange = isLargerTimeRange(start_schedule.value[index], end_schedule.value[index]);
  if (isLargerRange) {
    return {
      'background-color': '#3daee920',
      'font-weight': 450
    };
  } else if (isInTimeRange) {
    return {
      'background-color': '#93cee97f',
      'font-weight': 800
    };
  }
  return {};
}

function getTextStyle(index) {
  const isInTimeRange = isWithinTimeRange(start_schedule.value[index], end_schedule.value[index]);
  const isLargerRange = isLargerTimeRange(start_schedule.value[index], end_schedule.value[index]);
  if (isLargerRange) {
    return {
      'font-weight': 450
    };
  } else if (isInTimeRange) {
    return {
      'font-weight': 800
    };
  }
  return {};
}
updateSchedule()
</script>

<template>
  <div class="flex mr-2" @click="updateSchedule">
    <!-- 显示当前星期几 -->
    <div
      class="h-16 w-16 rounded-2 backdrop-blur-lg mt-2 ml-2 flex items-center justify-center flex-grow"
      style="background-color: #3daee920">
      <a class="text-white text-3xl font-450">{{ dayName }}</a>
    </div>

    <!-- 显示今天的课程，并判断时间区间 -->
    <div
      class="h-16 w-16 rounded-2 backdrop-blur-lg mt-2 ml-2 flex items-center justify-center flex-grow"
      v-for="(item, index) in today_schedule" :key="index"
      :style="getStyle(index)">
      <a class="text-white text-3xl"
         :style="getTextStyle(index)">
        {{ item }}
      </a>
    </div>
  </div>
</template>

<style scoped>
</style>
