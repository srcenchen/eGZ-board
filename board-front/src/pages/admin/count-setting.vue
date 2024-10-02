<script setup>
import axios from "axios";

const header = [
  {title: '事件', key: 'Event'},
  {title: '日期', key: 'Date'},
  {title: '操作', key: 'actions'},
]

const items = ref([])

function initData() {
  axios.get("/api/count-down/get-event").then(res => {
    items.value = res.data.data.Res
  })
}

// 初始化数据
onMounted(() => {
  initData()
})

const dialog = ref(false)
const date_choose = ref()
const event = ref('')

function addEvent() {
  if (date_choose && event.value !== '') {
    const year = date_choose.value.getFullYear();
    const month = date_choose.value.getMonth() + 1;
    const date = date_choose.value.getDate();
    const date_final = year + "-" + month + "-" + date;
    axios.post('/api/count-down/add-event', {
      event: event.value,
      date: date_final,
      type: 'count',
      during: 0
    })
      .then(res => {
        dialog.value = false
        initData()
        event.value = ''
      })
  }
}

function delEvent(id) {
  axios.post('/api/count-down/del-event', {
    id: id
  })
    .then(res => {
      initData()
    })
}
</script>

<template>
  <div className="flex flex-col m-5">
    <h2>倒计时管理</h2>
    <v-btn class="w-24 " variant="outlined" @click="dialog = true">新增倒计时</v-btn>
    <v-dialog
      v-model="dialog"
      persistent
      width="1024"
    >
      <v-card>
        <v-card-title>
          <span class="text-h5">添加倒计时</span>
        </v-card-title>
        <v-card-text>
          <v-container class=" flex-col items-center justify-center">
            <v-text-field label="日程名称" v-model="event"></v-text-field>
            <v-date-picker v-model="date_choose"></v-date-picker>
          </v-container>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            color="blue-darken-1"
            variant="text"
            @click="dialog = false"
          >
            关闭
          </v-btn>
          <v-btn
            color="blue-darken-1"
            variant="text"
            @click="addEvent"
          >
            保存
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-data-table :headers="header" :items="items">
      <template v-slot:item.actions="{item}">
        <v-icon
          size="small"
          class="me-2"
          @click="delEvent(item.Id)"
        >
          mdi-delete
        </v-icon>
      </template>
    </v-data-table>
  </div>
</template>

<style scoped>

</style>
