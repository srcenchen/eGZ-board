<template>
  <div class="flex flex-col w-178">
    <BaseData :is-outside="false"/>
    <WeatherWidget/>
  </div>
</template>

<style scoped>

</style>
<script setup>
import BaseData from "@/components/tools/BaseData.vue";
import WeatherWidget from "@/components/tools/WeatherWidget.vue";
</script>
