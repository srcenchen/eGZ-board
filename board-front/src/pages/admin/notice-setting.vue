<script setup>
import axios from "axios";
import {Editor, Toolbar} from "@wangeditor/editor-for-vue";
import '@wangeditor/editor/dist/css/style.css'; // 引入 css
import {onBeforeUnmount, ref, shallowRef, onMounted} from 'vue';

// 编辑器实例，必须用 shallowRef
const editorRef = shallowRef();

// 内容 HTML
const valueHtml = ref('');

// 工具栏配置
const toolbarConfig = {};

// 编辑器配置，设置默认字号
const editorConfig = {
  placeholder: '请输入内容...',
  MENU_CONF: {
    fontSize: {
      // 配置可选的字号
      fontSizeList: ['12px', '14px', '16px', '18px', '24px', '32px', '48px'],
      defaultSize: '24px', // 设置默认字号为 24px
    },
  },
  // 默认字体大小为 24px
  style: {
    'font-size': '24px',
  }
};
onMounted(() => {
  axios.get("/api/setting/get-setting-value/?key=notice")
    .then((res) => {
      valueHtml.value = res.data.data.Res.Value;
      //console.log(slogan_update_time.value);
    });
})
// 组件销毁时，也及时销毁编辑器
onBeforeUnmount(() => {
  const editor = editorRef.value;
  if (editor == null) return;
  editor.destroy();
});

const handleCreated = (editor) => {
  editorRef.value = editor; // 记录 editor 实例，重要！
};

function submit() {
  axios.post("/api/setting/set-setting-value/", {
    key: 'notice',
    value: valueHtml.value
  }).then(res => {

  })
}
</script>

<template>
  <div class="flex flex-col m-5">
    <h2>基础设置</h2>
    <div style="border: 1px solid #ccc">
      <Toolbar
        style="border-bottom: 1px solid #ccc"
        :editor="editorRef"
        :defaultConfig="toolbarConfig"
        :mode="mode"
      />
      <Editor
        style="height: 500px; overflow-y: hidden"
        v-model="valueHtml"
        :defaultConfig="editorConfig"
        :mode="mode"
        @onCreated="handleCreated"
      />
    </div>
    <v-btn variant="outlined" @click="submit">提交</v-btn>
    <a>{{ valueHtml }}</a>
  </div>
</template>

<style scoped>

</style>
