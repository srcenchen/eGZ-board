<script setup>
import axios from "axios";

const slogan_update_time = ref()
const wall_url = ref()
// 初始化数据
onMounted(() => {
  axios.get("/api/setting/get-setting-value/?key=slogan_time&class=5").then(res => {
    slogan_update_time.value = res.data.data.Res.Value
  })
  axios.get("/api/setting/get-setting-value/?key=wall_url&class=5").then(res => {
    wall_url.value = res.data.data.Res.Value
  })
})

function setSetting(key, value, class_num) {
  axios.post("/api/setting/set-setting-value/", {
    key: key,
    class: class_num,
    value: value
  }).then(res => {

  })
}
</script>

<template>
  <div className="flex flex-col m-5">
    <h2>基础设置</h2>
    <h4 class="mt-3 mb-2">Slogan更新时间（秒）</h4>
    <v-text-field
      v-model="slogan_update_time"
      label="Slogan更新时间（秒）"
      outlined
      class="lg:w-1/3"
    ></v-text-field>
    <div>
      <v-btn variant="outlined" @click="setSetting('slogan_time', slogan_update_time, '5')">保存时间</v-btn>
    </div>

    <h4 class="mt-3 mb-2">壁纸链接</h4>
    <v-text-field
      v-model="wall_url"
      label="壁纸链接"
      outlined
      class="lg:w-1/3"
    ></v-text-field>
    <div>
      <v-btn variant="outlined" @click="setSetting('wall_url', wall_url, '5')">保存壁纸</v-btn>
    </div>
  </div>
</template>

<style scoped>

</style>
