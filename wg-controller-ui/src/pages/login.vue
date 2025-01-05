<script lang="ts" setup>
import { ref } from "vue";
import { emailValidate, passwordValidate } from "@/utils/validators";
import { VForm } from "vuetify/components";
import router from "@/router";
import { POST_Login } from "@/api/methods";
import type { LoginBody } from "@/types/shared";

// Types
import { useStore } from "vuex";
import { key } from "../store";

const store = useStore(key);

const email = ref("");
const password = ref("");
const showPassword = ref(false);
const loading = ref(false);

// Refs
const loginForm = ref<VForm>();

async function onSubmit() {
  if (loginForm.value == null) {
    console.error("loginForm is null");
    return;
  }

  let result = await loginForm.value.validate();
  if (result.valid) {
    try {
      loading.value = true;
      let body: LoginBody = {
        email: email.value,
        password: password.value
      };
      let resp = await POST_Login(body);
      if (resp.status == 200) {
        // Get email from response
        let body = await resp.json();
        if ("email" in body) {
          store.state.UserEmail = body.email;
        } else {
          console.error("Email not found in response");
        }

        // Set logged in state
        store.state.LoggedIn = true;

        // Redirect to clients
        router.push("/clients");
      } else {
        store.state.SnackBarText = "Login error";
        store.state.SnackBarError = true;
        store.state.SnackBarShow = true;
      }
    } catch (error: any) {
      console.error(error);
      store.state.SnackBarText = error;
      store.state.SnackBarError = true;
      store.state.SnackBarShow = true;
    } finally {
      loading.value = false;
    }
  } else {
    console.error("Form is not valid");
  }
}
</script>

<template>
  <v-container fluid class="d-flex align-center justify-center" style="height: 100vh">
    <v-card style="min-width: 500px" class="pa-1">
      <v-icon alt="logo" color="primary" class="mx-auto mt-3 d-block" :size="70">mdi-vpn</v-icon>
      <v-card-title class="text-center">Welcome, please login.</v-card-title>
      <v-card-text>
        <v-form ref="loginForm" @submit.prevent="onSubmit" validate-on="blur">
          <v-text-field
            v-model="email"
            label="Email"
            required
            density="comfortable"
            :rules="[emailValidate]"
            class="mb-1"
          ></v-text-field>

          <v-text-field
            v-model="password"
            label="Password"
            required
            density="comfortable"
            :append-inner-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
            :type="showPassword ? 'text' : 'password'"
            @click:append-inner="() => (showPassword = !showPassword)"
            :rules="[passwordValidate]"
            class="mb-5"
          ></v-text-field>

          <v-btn block color="primary" type="submit" :loading="loading"> Login </v-btn>
        </v-form>
      </v-card-text>
    </v-card>
  </v-container>
</template>
