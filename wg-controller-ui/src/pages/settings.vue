<script lang="ts" setup>
import { ref, onMounted } from "vue";

import { useStore } from "vuex";
import { key } from "../store";
const store = useStore(key);

onMounted(() => {
  //Init();
});

/*
async function Init() {
  loading.value = true;
  try {
    const val = await GET_APIKeys();
    if (val != null) {
      items.value = val;
    } else {
      items.value = [];
    }
  } catch (error: any) {
    console.error(error);
    store.state.SnackBarText = "Error fetching API keys";
    store.state.SnackBarError = true;
    store.state.SnackBarShow = true;
  } finally {
    loading.value = false;
  }
}
  */

const loading = ref(true);
</script>

<template>
  <v-container fluid max-width="900">
    <v-row no-gutters class="d-flex align-center">
      <span class="text-h4">Settings</span>
      <v-icon size="x-large" color="rgb(186,194,202)" class="ml-3"> mdi-wrench </v-icon>
    </v-row>

    <v-container fluid class="mx-2">
      <v-row no-gutters>
        <span class="text-h6 mt-5">Server Info</span>
      </v-row>

      <v-row no-gutters class="pl-3">
        <span class="text-subtitle-1 mt-3">WG Endpoint</span>
        <v-spacer></v-spacer>
        <span class="text-subtitle-1 mt-3 text-grey-darken-1">localhost:554</span>
      </v-row>

      <v-row no-gutters class="pl-3">
        <span class="text-subtitle-1 mt-3">API Endpoint</span>
        <v-spacer></v-spacer>
        <span class="text-subtitle-1 mt-3 text-grey-darken-1">https://localhost:51820</span>
      </v-row>

      <v-row no-gutters class="pl-3">
        <span class="text-subtitle-1 mt-3">Tun Network</span>
        <v-spacer></v-spacer>
        <span class="text-subtitle-1 mt-3 text-grey-darken-1">http://localhost:51820</span>
      </v-row>

      <v-row no-gutters class="pl-3">
        <span class="text-subtitle-1 mt-3">Server Version</span>
        <v-spacer></v-spacer>
        <v-chip color="orange-darken-2" size="small" class="text-subtitle-1 mt-3 mx-2"
          >update available</v-chip
        >
        <span class="text-subtitle-1 mt-3 text-grey-darken-1">v0.0.3</span>
      </v-row>

      <v-row no-gutters class="pl-3">
        <span class="text-subtitle-1 mt-3">UI Version</span>
        <v-spacer></v-spacer>
        <v-chip color="green" size="small" class="text-subtitle-1 mt-3 mx-2">latest</v-chip>
        <span class="text-subtitle-1 mt-3 text-grey-darken-1">v0.0.4</span>
      </v-row>

      <v-row no-gutters>
        <span class="text-h6 mt-5">Defaults</span>
      </v-row>

      <v-row no-gutters>
        <span class="text-h6 mt-5">Monitoring</span>
      </v-row>

      <v-row no-gutters>
        <span class="text-h6 mt-5">About</span>
      </v-row>
    </v-container>
  </v-container>
</template>
