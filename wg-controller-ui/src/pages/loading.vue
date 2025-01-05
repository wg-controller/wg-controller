<script setup lang="ts">
import { onMounted } from "vue";
import { useStore } from "vuex";
import { key } from "../store";
import { useRouter } from "vue-router";
import { POST_PreLogin } from "@/api/methods";

const store = useStore(key);
const router = useRouter();

onMounted(() => {
  PreLogin();
  store.state.Loaded = true;
});

async function PreLogin() {
  try {
    const resp = await POST_PreLogin();
    if (resp.status == 200) {
      // Set logged in state
      store.state.LoggedIn = true;

      // Get email from response
      let body = await resp.json();
      if ("email" in body) {
        store.state.UserEmail = body.email;
      } else {
        console.error("Email not found in response");
      }

      // Check for redirect URL
      if (router.currentRoute.value.query.redirect) {
        router.push(router.currentRoute.value.query.redirect as string);
      } else {
        router.push("/clients");
      }
    } else {
      store.state.SnackBarText = "Error connecting to server";
      store.state.SnackBarError = true;
      store.state.SnackBarShow = true;
    }
  } catch (error) {
    console.error(error);
    store.state.LoggedIn = false;
    router.push("/login");
  }
}
</script>

<template>
  <v-container fluid class="d-flex align-center justify-center" style="height: 100vh">
    <v-progress-circular indeterminate color="primary" size="64" />
  </v-container>
</template>
