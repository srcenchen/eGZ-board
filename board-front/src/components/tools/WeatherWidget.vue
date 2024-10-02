<template>
  <div class="w-full bg-transparent backdrop-blur-xl flex flex-col pl-4 pr-4 pb-2 rounded-4  mr-2">
    <p class="text-white justify-end items-end d-flex mt-2">观测时间：{{ updateTime }}</p>
    <!-- 当前温度与天气信息 -->
    <div class="current-weather d-flex align-center items-center justify-center">
      <img :src="getWeatherIcon(now.icon)" alt="weather icon" width="64" height="64" referrerpolicy="no-referrer"/>
      <h1 class="text-white ml-2 mr-2">{{ now.temp }}°</h1>
      <p class="text-white text-3xl">{{ now.text }}</p>
    </div>

    <!-- 天气详细信息：风力、湿度等 -->
    <div class="weather-details d-flex justify-space-between mt-2">
      <p class="text-white text-2xl">{{ now.windDir }} {{ now.windScale }}级</p>
      <p class="text-white text-2xl">湿度 {{ now.humidity }}%</p>
      <p class="text-white text-2xl">体感 {{ now.feelsLike }}°C</p>
      <p class="text-white text-2xl">LD {{ days }}</p>
    </div>
    <div class="weather-details d-flex justify-center mt-2">
      <p class="text-white text-2xl">{{ rain.summary }}</p>
    </div>


    <!-- 三天预报 -->
    <div class="forecast mt-4">
      <div v-for="day in forecast" :key="day.fxDate" class="d-flex justify-space-between align-center mb-2">
        <div>
          <p class="text-white text-2xl">{{ formatDate(day.fxDate) }}</p>
        </div>
        <div class="d-flex align-center">
          <img :src="getWeatherIcon(day.iconDay)" alt="day icon" width="32" height="32" referrerpolicy="no-referrer"/>
          <div class="temperature-range ml-4">
            <p class="text-white text-2xl">{{ day.tempMax }}° / {{ day.tempMin }}°</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import {onMounted, ref} from 'vue';

const updateTime = ref('');
const now = ref({});
const forecast = ref([]);
const rain = ref({});

// 获取当前天气信息
const fetchNowWeather = async () => {
  const response = await fetch('/api/proxy/get-weather?key=now');
  const data = await response.json();
  const weatherData = JSON.parse(data.data.Data.Value);
  updateTime.value = weatherData.updateTime;
  now.value = weatherData.now;
};

// 获取三天预报信息
const fetchForecast = async () => {
  const response = await fetch('/api/proxy/get-weather?key=3d');
  const data = await response.json();
  const forecastData = JSON.parse(data.data.Data.Value);
  forecast.value = forecastData.daily;
};

// 获取降水预报信息
const fetchRain = async () => {
  const response = await fetch('/api/proxy/get-weather?key=rain');
  const data = await response.json();
  rain.value = JSON.parse(data.data.Data.Value);
};

// 获取天气图标URL
const getWeatherIcon = (iconCode) => {
  return `https://a.hecdn.net/img/common/icon/202106d/${iconCode}.png`;
};

// 格式化日期
const formatDate = (date) => {
  const options = {weekday: 'short', month: 'long', day: 'numeric'};
  return new Date(date).toLocaleDateString('zh-CN', options);
};

// 组件挂载时请求数据
onMounted(() => {
  update()
});
const days = ref(0)

function update() {
  fetchNowWeather();
  fetchForecast();
  fetchRain()
  const startDate = new Date('2024-04-28');
  const currentDate = new Date();
  days.value = Math.floor((currentDate - startDate) / (1000 * 60 * 60 * 24));
}

setInterval(update, 60 * 1000)
</script>

<style scoped>

.current-weather {
  font-size: 2rem;
}

.weather-details {
  font-size: 0.9rem;
}

.forecast {
  font-size: 1rem;
}
</style>
