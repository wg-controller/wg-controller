<script lang="ts" setup>
import { useStore } from "vuex";
import { key } from "./store";

const store = useStore(key);
</script>

<template>
  <v-app>
    <v-app-bar app density="compact" :elevation="0" color="appBar">
      <template v-slot:prepend>
        <v-icon class="mx-3" size="x-large" color="primary">mdi-vpn</v-icon>
        <v-tabs slider-color="primary">
          <v-tab to="/clients">Clients</v-tab>
          <v-tab to="/apikeys">API Keys</v-tab>
          <v-tab to="/users">Users</v-tab>
          
        </v-tabs>
      </template>

      <template v-slot:append>
        <v-icon class="mx-3" size="x-large">mdi-account-circle</v-icon>
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
            store.state.ConfirmDialogCallback(),
              (store.state.ConfirmDialogShow = false)
          "
        >
          <v-card-title class="text-h5 mb-5"> Please Confirm </v-card-title>

          <v-card-text class="text-subtitle-1 mx-2">
            {{ store.state.ConfirmDialogText }}
          </v-card-text>

          <v-card-actions>
            <v-spacer></v-spacer>

            <v-btn color="primary" @click="store.state.ConfirmDialogShow = false">
              Cancel
            </v-btn>

            <v-btn color="primary" type="submit"> Confirm </v-btn>
          </v-card-actions>
        </v-form>
      </v-card>
    </v-dialog>
  </v-app>
</template>

<style>
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
</style>