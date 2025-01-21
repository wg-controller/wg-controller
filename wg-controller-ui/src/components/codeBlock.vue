<script lang="ts" setup>
import { defineProps } from "vue";
const props = defineProps<{
  body: string;
  copyCallback: (body: string) => void;
}>();
</script>

<template>
  <v-row no-gutters class="d-flex align-center">
    <div class="codeDiv" style="width: calc(100% - 40px)">
      <slot>{{ props.body }}</slot>
    </div>
    <v-btn
      class="ml-3"
      @click="copyCallback(props.body)"
      icon="mdi-content-copy"
      density="compact"
      flat
    ></v-btn>
  </v-row>
</template>

<style scoped>
.codeDiv {
  font-family: monospace;
  font-size: 0.9em;
  font-weight: normal;
  line-height: 2.3em;
  background-color: rgb(245, 245, 245);
  padding-left: 15px;
  padding-right: 15px;
  margin: 0px;
  border-radius: 5px;
  overflow-x: scroll;
  overflow-y: hidden;
  white-space: nowrap;
}
</style>
