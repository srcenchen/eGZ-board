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

function setSetting(key, value) {
  axios.post("/api/setting/set-setting-value/", {
    key: key,
    value: value
  }).then(res => {

  })
}

import {ref} from 'vue'

// 用于存储上传的文件
const file = ref(null)
// 用于存储返回的提示信息
const message = ref('')
const messageType = ref('success') // 'success' or 'error'

// 上传壁纸的方法
const uploadWallPaper = async () => {
  if (!file.value) {
    message.value = '请选择要上传的文件！'
    messageType.value = 'error'
    return
  }

  try {
    // 创建FormData对象
    const formData = new FormData()
    formData.append('file', file.value)

    // 发送POST请求到后端接口
    const response = await axios.post('/api/setting/upload-wall-paper', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })

    // 根据后端响应更新消息
    message.value = response.data.message || '壁纸设置成功！'
    messageType.value = 'success'
    wall_url.value = "/resource/upload/" + response.data.data.fileName
    setSetting('wall_url', wall_url.value)
  } catch (error) {
    message.value = '文件上传失败，请重试！'
    messageType.value = 'error'
    console.error('Upload Error:', error)
  }
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
      <v-btn variant="outlined" @click="setSetting('slogan_time', slogan_update_time)">保存时间</v-btn>
    </div>

    <h4 class="mt-3 mb-2">壁纸链接</h4>
    <v-text-field
      v-model="wall_url"
      label="壁纸链接"
      outlined
      class="lg:w-1/3"
    ></v-text-field>
    <!-- 上传课程表 -->
    <v-file-input
      v-model="file"
      label="上传壁纸"
      accept=".png, .jpg, .jpge, .webp"
      prepend-icon="mdi-paperclip"
      class="lg:w-1/3"
      @change="uploadWallPaper"
    ></v-file-input>

    <!-- 提示信息 -->
    <v-alert v-if="message" :type="messageType" class="lg:w-1/3" dismissible>
      {{ message }}
    </v-alert>
    <div>
      <v-btn variant="outlined" @click="setSetting('wall_url', wall_url)">保存壁纸</v-btn>
    </div>
  </div>
</template>

<style scoped>

</style>
