<template>
  <div class="background h-screen" :style="{ backgroundImage: `url(${wall_url})` }" @click="update_url">
    <div class="flex-col zoom h-full">
      <ScheduleTable class="col-span-3"/>
      <ToolBox class="col-end-5 mr-4 mt-4"></ToolBox>
    </div>

  </div>
</template>

<script setup>
import ScheduleTable from "@/components/index/ScheduleTable.vue";
import ToolBox from "@/components/outside/ToolBox.vue";
import axios from "axios";

const zoom = ref(window.devicePixelRatio);
const wall_url = ref('')

function update_url() {
  axios.get("/api/setting/get-setting-value/?key=wall_url").then(res => {
    wall_url.value = res.data.data.Res.Value
  })
}

update_url()
</script>

<style scoped>
.background {
  background-size: cover;
}

.zoom {
  zoom: calc(1 / v-bind(zoom));
}
</style>
