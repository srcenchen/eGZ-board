<template>
  <v-container>
    <!-- 上传课程表 -->
    <v-file-input
      v-model="file"
      label="上传课程表"
      accept=".json, .xls, .xlsx"
      prepend-icon="mdi-paperclip"
      @change="uploadSchedule"
    ></v-file-input>

    <!-- 提示信息 -->
    <v-alert v-if="message" :type="messageType" dismissible>
      {{ message }}
    </v-alert>

<!--    <v-data-table-->
<!--      :headers="headers"-->
<!--      :items="scheduleData"-->
<!--      item-value="index"-->
<!--      class="elevation-1"-->
<!--    >-->
<!--      <template v-slot:items="props">-->
<!--        <tr>-->
<!--          <td v-for="(header, index) in headers" :key="header.text">-->
<!--            {{ props.item[index] || '' }}-->
<!--          </td>-->
<!--        </tr>-->
<!--      </template>-->
<!--    </v-data-table>-->
  </v-container>
</template>

<script setup>
import { ref } from 'vue'
import axios from 'axios'

// 用于存储上传的文件
const file = ref(null)
// 用于存储返回的提示信息
const message = ref('')
const messageType = ref('success') // 'success' or 'error'

// 上传课程表的方法
const uploadSchedule = async () => {
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
    const response = await axios.post('/api/schedule/upload-schedule', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })

    // 根据后端响应更新消息
    message.value = response.data.message || '文件上传成功！'
    messageType.value = 'success'
  } catch (error) {
    message.value = '文件上传失败，请重试！'
    messageType.value = 'error'
    console.error('Upload Error:', error)
  }
}
</script>

<style scoped>
/* 你可以在这里根据需要添加样式 */
</style>
