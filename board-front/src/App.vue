<template>
  <v-app>
    <v-main>
      <router-view v-if="isBind" />
      <div v-if="!isBind" class="flex h-screen background"
           :style="{ backgroundImage: `url(https://bing.img.run/1920x1080.php)` }">
        <div class=" w-full bg-transparent backdrop-blur-xl rounded-4 flex pr-2 "
             @click="getConfig">
          <div class="w-16 rounded-2 backdrop-blur-xlflex flex-col flex-grow p-2 items-center ">
            <div class="flex justify-center items-center h-full mb-1 flex-col">
              <h1 class="text-white">善良的孩子哟，可以告诉我，这是哪里吗？</h1>
              <h2 class="text-white">请谨慎选择，错选后事情将变得复杂起来哦</h2>
              <div class="flex mt-4">
                <h1 class="text-white">高</h1>

                <!-- 年级选择器 -->
                <v-select
                  v-model="selectedGrade"
                  :items="grades"
                  label="年级"
                  outlined
                  class="mx-2 text-white"
                ></v-select>

                <!-- 班级选择器 -->
                <v-select
                  v-model="selectedClass"
                  :items="classes"
                  label="班级"
                  outlined
                  class="mx-2 text-white"
                ></v-select>

                <h1 class="text-white">班</h1>
              </div>
              <h1 class="text-white mb-4">
                当前选择: 高{{ selectedGrade }} {{ selectedClass }} 班
              </h1>
              <v-btn @click="bindDevice">提交</v-btn>
            </div>
          </div>
        </div>

      </div>
    </v-main>
  </v-app>
</template>

<script setup>
import axios from "axios";

const isBind = ref(true)

function testBind() {
  axios.get('/api/device/is-bind').then(res => {
    isBind.value = res.data.data.isBind.DeviceID !== '';
  })
}

onMounted(() => {
  testBind()
})
const grades = ref(['一', '二', '三']); // 年级选项
const classes = ref(Array.from({length: 24}, (v, i) => i + 1)); // 1到22班
const selectedGrade = ref(null); // 当前选中的年级
const selectedClass = ref(null); // 当前选中的班级
// 计算属性，用于获取年级和班级的序号
const gradeIndex = computed(() => {
  return grades.value.indexOf(selectedGrade.value) + 1; // 返回年级的序号 (1-3)
});

const classIndex = computed(() => {
  return selectedClass.value; // 班级的序号 (1-22)
});

function bindDevice() {
  if (gradeIndex.value > 0 && classIndex.value > 0) {
    axios.post('/api/device/bind-device ', {
      grade: gradeIndex.value,
      class: classIndex.value
    }).then(res => {
      testBind()
    })
  }
}
</script>
