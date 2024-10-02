<script setup>
import axios from 'axios';
import WeatherWidget from "@/components/tools/WeatherWidget.vue";
import {ref} from "vue";

const nowDate = ref(['11:40 PM', '第36周', '2024/9/16'])
const dayPast = ref(0)
const noticeHTML = ref('<p><br></p>')

function formatTime(now) {
  const hours = now.getHours();
  const minutes = now.getMinutes();
  const seconds = now.getSeconds();
  const ampm = hours >= 12 ? 'PM' : 'AM';
  const formattedHours = hours % 12 || 12;
  const formattedMinutes = minutes.toString().padStart(2, '0');
  const formattedSeconds = seconds.toString().padStart(2, '0');
  return `${formattedHours}:${formattedMinutes}:${formattedSeconds} ${ampm}`;
}

function getWeekName(now) {
  // 获取星期几
  const daysOfWeek = ["周日", "周一", "周二", "周三", "周四", "周五", "周六"];
  return daysOfWeek[now.getDay()]
}

function formatDate(now) {
  const year = now.getFullYear();
  const month = now.getMonth() + 1; // getMonth() is 0-based
  const day = now.getDate();
  return `${year}/${month}/${day}`;
}

function daysToWeeksAndDays(days) {
  const weeks = Math.floor(days / 7);  // 计算周数
  const remainingDays = days % 7;      // 计算剩余天数
  return `${weeks}周`; // 返回格式化字符串
}

function getDayProgressPercentage(now) {
  // 创建当天的 6:00 和 22:00 时间
  const startTime = new Date(now.getFullYear(), now.getMonth(), now.getDate(), 6, 0, 0);
  const endTime = new Date(now.getFullYear(), now.getMonth(), now.getDate(), 22, 0, 0);
  // 计算总的秒数
  const totalDuration = endTime - startTime;
  // 如果当前时间在 6 点之前，返回 0
  if (now < startTime) {
    return 0;
  }
  // 如果当前时间在 22 点之后，返回 100
  if (now > endTime) {
    return 100;
  }
  // 计算已经过去的秒数，并转换为百分比
  const passedTime = now - startTime;
  return ((passedTime / totalDuration) * 100).toFixed(3);
}

// 设置时间
function setTime() {
  const now = new Date(); // 获取当前时间
  nowDate.value[0] = formatTime(now)
  nowDate.value[1] = formatDate(now)
  nowDate.value[2] = className.value
  // 判断当前时间是否是凌晨2点
  const hours = now.getHours(); // 获取当前小时
  const minutes = now.getMinutes(); // 获取当前小时
  const seconds = now.getSeconds(); // 获取当前小时
  if (hours === 2 && minutes === 0 && seconds === 0) {
    location.reload()
  } // 判断是否为2点
  dayPast.value = getDayProgressPercentage(now);
}

let intervalId = null;
onMounted(() => {
  setTime()
  getCountDown()
  intervalId = setInterval(setTime, 1000);
})

onUnmounted(() => {
  if (intervalId) {
    clearInterval(intervalId);
  }
})

const slogan = ref({
  text: "",
  from: ""
})

function updateSlogan() {
  axios.get("https://international.v1.hitokoto.cn/?c=i").then((res) => {
    slogan.value.text = res.data.hitokoto
    if (res.data.from_who)
      slogan.value.from = '——' + res.data.from_who + '《' + res.data.from + '》'
    else slogan.value.from = '——' + '《' + res.data.from + '》'
  })
}

const slogan_update_time = ref(10);

function getConfig() {
  axios.get("/api/setting/get-setting-value/?key=slogan_time")
    .then((res) => {
      slogan_update_time.value = res.data.data.Res.Value;
      //console.log(slogan_update_time.value);
    });
  axios.get("/api/setting/get-setting-value/?key=notice")
    .then((res) => {
      noticeHTML.value = res.data.data.Res.Value;
      //console.log(slogan_update_time.value);
    });

  getCountDown()
  unBindCount.value = 0
}

let updateSloganId = null;
onMounted(() => {
  // 初始调用
  // updateSlogan();
  getConfig();
  // 设置一个定时器来定期调用 getConfig
  setInterval(getConfig, 300000);
  // 设置初始的定时器
  if (slogan_update_time.value > 0) {
    updateSloganId = setInterval(updateSlogan, slogan_update_time.value * 1000);
  } else {
    clearInterval(updateSloganId)
  }
  updateSlogan()

});

// 观察 slogan_update_time 的变化
watch(slogan_update_time, (newValue) => {
  if (updateSloganId) {
    clearInterval(updateSloganId);
  }
  // 根据新的时间值重新设置定时器
  if (newValue > 0) {
    updateSloganId = setInterval(updateSlogan, newValue * 1000);
  }
});
const count_down_list = ref('')

function getCountDown() {
  axios.get("/api/count-down/get-event").then((res) => {
    count_down_list.value = res.data.data.Res
    count_down_list.value.forEach((item) => {
      item.left = getDayDiff(item.Date); // 使用 item.Date 计算并设置 left
      item.weekday = daysToWeeksAndDays(item.left)
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

const props = defineProps({
  isOutside: {
    type: Boolean,
    default: false
  }
})


const unBindCount = ref(0)

watch(unBindCount, () => {
  if (unBindCount.value === 10) {
    axios.get('/api/device/un-bind')
    location.reload()
  }
})

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
</script>

<template>
  <div>
    <div class="w-full bg-transparent backdrop-blur-xl rounded-4 flex pr-2 mb-2">
      <div
        class="h-20 w-16 rounded-2 backdrop-blur-xl mt-2 mb-2 ml-2 flex items-center justify-center flex-grow"
        style="background-color: #3daee920"
        v-for="index in nowDate.length">
        <a class="text-white text-4xl  font-600">{{ nowDate[index - 1] }}</a>
      </div>
    </div>
    <div v-if="!props.isOutside"
         class=" w-full bg-transparent backdrop-blur-xl rounded-4 flex  flex-col pr-2 mb-2 pt-1 pb-1">
      <div
        class="h-16 rounded-2 backdrop-blur-xl ml-2 flex items-center justify-center flex-grow"
        v-for="index in count_down_list.length">
        <a class="text-white text-4xl font-500">距离</a>
        <a class="text-white text-4xl font-900">{{ count_down_list[index - 1].Event }}</a>
        <a class="text-white text-4xl font-500">还有</a>
        <a class="text-white text-6xl font-900">{{ count_down_list[index - 1].left }}</a>
        <a class="text-white text-4xl font-500">天({{ count_down_list[index - 1].weekday }})</a>
      </div>
    </div>
    <div class=" w-full bg-transparent backdrop-blur-xl rounded-4 flex pr-2 mb-2" @click="unBindCount++">
      <div class="w-16 rounded-2 backdrop-blur-xl mt-2 mb-2 ml-2 flex flex-col flex-grow p-2">
        <v-progress-linear
          class="text-purple-lighten-3"
          height="20"
          :model-value="dayPast"
          rounded
        ></v-progress-linear>
        <div class="flex justify-end mt-4">
          <a class="text-white text-2xl font-700">今日已过去{{ dayPast }}%</a>
        </div>
      </div>
    </div>
    <div class=" w-full bg-transparent backdrop-blur-xl rounded-4 flex pr-2 mb-2" @click="getConfig">
      <div class="w-16 rounded-2 backdrop-blur-xl mt-2 mb-2 ml-2 flex flex-col flex-grow p-2" @click="updateSlogan">
        <div class="flex justify-start mb-1">
          <a class="text-white text-3xl font-700">{{ slogan.text }}</a>
        </div>
        <div class="flex justify-end">
          <a class="text-white text-2xl font-700">{{ slogan.from }}</a>
        </div>
      </div>
    </div>


  </div>

</template>

<style scoped>

</style>
