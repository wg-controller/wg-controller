<script lang="ts" setup>
import { useStore } from "vuex";
import { key } from "./store";
import { useRouter } from "vue-router";
import { POST_Logout } from "./api/methods";

const store = useStore(key);
const router = useRouter();

function Logout() {
  // Send logout request
  try {
    POST_Logout();
  } catch (error: any) {
    console.error(error);
    store.state.SnackBarText = "Error logging out";
    store.state.SnackBarError = true;
    store.state.SnackBarShow = true;
    return;
  }

  // Clear state
  store.state.LoggedIn = false;
  store.state.UserEmail = "";
  store.state.SnackBarText = "Logged out";
  store.state.SnackBarError = false;
  store.state.SnackBarShow = true;

  // Redirect to login
  router.push("/login");
}
</script>

<template>
  <v-app>
    <v-snackbar
      v-model="store.state.SnackBarShow"
      location="bottom left"
      :color="store.state.SnackBarError ? 'red-darken-2' : 'primary'"
      :timer="store.state.SnackBarError ? 'red-darken-4' : 'secondary'"
      :timeout="7000"
      max-width="400"
      close-on-content-click
      transition="slide-x-transition"
    >
      <v-icon class="mr-3"> mdi-information </v-icon>
      <span>{{ store.state.SnackBarText }}</span>
    </v-snackbar>
    <v-app-bar
      v-bind:class="{ invisible: !store.state.LoggedIn }"
      app
      density="compact"
      :elevation="0"
      color="secondary"
    >
      <template #[`prepend`]>
        <img class="ml-3 mr-4" style="width: 40px" src="./assets/Logo.png" />
        <v-tabs slider-color="primary">
          <v-tab to="/clients"> Clients </v-tab>
          <v-tab to="/apikeys"> API Keys </v-tab>
          <v-tab to="/users"> Users </v-tab>
        </v-tabs>
      </template>

      <template #[`append`]>
        <v-menu open-on-hover>
          <template v-slot:activator="{ props }">
            <span class="mr-1 ml-1">{{ store.state.UserEmail }}</span>
            <v-btn v-bind="props" icon size="small" class="mr-1">
              <v-icon class="mx-3" size="x-large"> mdi-account-circle </v-icon>
            </v-btn>
          </template>

          <v-list>
            <v-list-item class="d-flex flex-row" @click="Logout()" append-icon="mdi-logout">
              <v-list-item-title>Logout</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>
      </template>
    </v-app-bar>
    <v-main>
      <router-view />
    </v-main>

    <v-dialog v-model="store.state.ConfirmDialogShow" width="460" scrim="grey-darken-1">
      <v-card>
        <v-form
          ref="entryForm"
          @submit.prevent="
            (store.state.ConfirmDialogCallback(), (store.state.ConfirmDialogShow = false))
          "
        >
          <v-card-title class="text-h6 mb-1 ml-2 mt-4">
            {{ store.state.ConfirmDialogTitle }}
          </v-card-title>

          <v-card-text class="text-subtitle-1 ml-2">
            {{ store.state.ConfirmDialogText }}
          </v-card-text>

          <v-card-actions class="mb-3 mr-3">
            <v-spacer />

            <v-btn
              color="secondary"
              variant="outlined"
              @click="store.state.ConfirmDialogShow = false"
            >
              Cancel
            </v-btn>

            <v-btn color="secondary" type="submit" variant="flat"> Confirm </v-btn>
          </v-card-actions>
        </v-form>
      </v-card>
    </v-dialog>
  </v-app>
</template>

<style>
.html,
body {
  min-width: 750px !important;
}

.v-table__wrapper > table > tbody > tr:nth-of-type(odd) {
  background-color: rgb(var(--v-theme-oddRow)) !important;
}

.v-table__wrapper > table > tbody > tr:nth-of-type(even) {
  background-color: rgb(var(--v-theme-obtCard)) !important;
}

.v-table__wrapper > table > tbody td {
  border: none !important;
}

.v-table__wrapper > table > thead th {
  border: none !important;
}

.v-field__overlay {
  opacity: 0 !important;
}

.v-list {
  padding: 0px !important;
}

.invisible {
  display: none !important;
}
</style>
