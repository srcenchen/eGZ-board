<template>
  <div class="background h-screen" :style="{ backgroundImage: `url(${wall_url})` }" @click="getConfig">
    <div class="flex-col zoom h-full">
      <!-- 日程表组件 -->
      <ScheduleTable class="col-span-3"/>

      <!-- 布局分为三部分 -->
      <div class="flex w-full mt-2">
        <!-- 左侧空白部分 -->
        <div class="flex-1 "></div>

        <!-- 右侧的天气组件和工具箱 -->
        <div class="flex w-2/3 h-auto">
          <div class="flex-1 mr-2">
            <div class=" w-full bg-transparent backdrop-blur-xl rounded-4 flex pr-2 mb-2"
                 @click="getConfig" v-if="noticeSwitch">
              <div class="w-16 rounded-2 backdrop-blur-xl mt-2 mb-2 ml-2 flex flex-col flex-grow p-2">
                <div class="flex justify-start mb-1">
                  <a class="text-white font-700 w-full text-xl" v-html="noticeHTML"/>
                </div>
              </div>
            </div>
          </div>
          <ToolBox class="flex-3 mr-2"/>          <!-- ToolBox 占 3 部分 -->
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import ScheduleTable from "@/components/index/ScheduleTable.vue";
import ToolBox from "@/components/index/ToolBox.vue";
import axios from "axios";

const zoom = ref(window.devicePixelRatio);
const wall_url = ref('');

function update_url() {
  axios.get("/api/setting/get-setting-value/?key=wall_url&class=5").then(res => {
    wall_url.value = res.data.data.Res.Value;
  });
}

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
  update_url();
}

setInterval(getConfig, 60 * 1000)
onMounted(() => {
  getConfig()
})

</script>

<style scoped>
.background {
  background-size: cover;
}

/* 解决缩放问题 */
.zoom {
  zoom: calc(1 / v-bind(zoom));
}

/* 设置右侧的 flex 布局 */
.flex {
  display: flex;
}
</style>
