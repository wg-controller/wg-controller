<script lang="ts" setup>
import { ref, onMounted } from "vue";

import { useStore } from "vuex";
import { key } from "../store";
import { VForm } from "vuetify/components";
import { DELETE_APIKey, GET_APIKeyInit, GET_APIKeys, PUT_APIKey } from "@/api/methods";
import type { APIKey, APIKeyInit, APIKeyWithToken } from "@/types/shared";
import { required } from "@/utils/validators";
const store = useStore(key);

onMounted(() => {
  Init();
});

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

const items = ref<APIKey[]>([]);
const headers = ref([
  { title: "Name", key: "name" },
  { title: "Expires", key: "expires" },
  { title: "Permissions", key: "attributes" },
  { title: "", key: "actions", align: "end", sortable: false }
] as const);

const search = ref("");
const loading = ref(true);

const keyDialog = ref(false);
const keyDialogEditMode = ref(false);
const keyBuffer = ref<APIKeyWithToken>();
const keyBufferCustomPermissions = ref(false);

async function NewKeyDialog() {
  let initKey: APIKeyInit;
  try {
    initKey = await GET_APIKeyInit();
  } catch (error: any) {
    console.error(error);
    store.state.SnackBarText = error;
    store.state.SnackBarError = true;
    store.state.SnackBarShow = true;
    return;
  }
  keyBuffer.value = {
    uuid: initKey.uuid,
    expiresUnixMillis: 31536000000,
    name: "",
    attributes: [],
    token: initKey.token
  };
  keyBufferCustomPermissions.value = false;
  keyDialogEditMode.value = false;
  keyDialog.value = true;
}

const keyForm = ref<VForm>();
async function ApplyKeyDialog() {
  if (keyForm.value == null) {
    console.error("loginForm is null");
    return;
  }

  const result = await keyForm.value.validate();
  if (result.valid) {
    // Set permissions based on type
    if (!keyBufferCustomPermissions.value) {
      keyBuffer.value!.attributes = ["wg-client"];
    }

    // Set actual expiry date
    if (keyBuffer.value!.expiresUnixMillis > 0) {
      keyBuffer.value!.expiresUnixMillis = Date.now() + keyBuffer.value!.expiresUnixMillis;
    }

    try {
      await PUT_APIKey(keyBuffer.value!);
    } catch (error: any) {
      console.error(error);
      store.state.SnackBarText = error;
      store.state.SnackBarError = true;
      store.state.SnackBarShow = true;
    } finally {
      keyDialog.value = false;
      Init();
    }
  } else {
    console.error("Form is not valid");
  }
}

function RemoveKey(key: APIKey) {
  store.state.ConfirmDialogTitle = "Remove " + key.name;
  store.state.ConfirmDialogText = "Are you sure you want to remove this API key?";
  store.state.ConfirmDialogCallback = async () => {
    try {
      await DELETE_APIKey(key.uuid);
    } catch (error: any) {
      console.error(error);
      store.state.SnackBarText = error;
      store.state.SnackBarError = true;
      store.state.SnackBarShow = true;
    } finally {
      Init();
    }
  };
  store.state.ConfirmDialogShow = true;
}

function CopyToClipboard(text: string) {
  navigator.clipboard.writeText(text).then(
    () => {
      store.state.SnackBarText = "Copied to clipboard";
      store.state.SnackBarError = false;
      store.state.SnackBarShow = true;
    },
    () => {
      store.state.SnackBarText = "Failed to copy to clipboard";
      store.state.SnackBarError = true;
      store.state.SnackBarShow = true;
    }
  );
}
</script>

<template>
  <v-container fluid max-width="1400">
    <v-row no-gutters class="d-flex align-center">
      <span class="text-h4">API Keys</span>
      <v-icon size="x-large" color="rgb(186,194,202)" class="ml-3"> mdi-key </v-icon>
    </v-row>
    <v-row no-gutters class="d-flex align-center justify-start mt-3 mb-3">
      <v-text-field
        v-model="search"
        label="Search"
        style="max-width: 300px"
        outlined
        flat
        density="compact"
        hide-details
        variant="solo-filled"
        append-inner-icon="mdi-magnify"
      />

      <v-btn
        color="secondary"
        variant="tonal"
        class="ml-2"
        height="40px"
        append-icon="mdi-plus"
        @click="NewKeyDialog()"
      >
        New Key
      </v-btn>

      <v-spacer />
    </v-row>

    <v-data-table
      :items="items"
      :headers="headers"
      no-data-text="No keys found"
      :items-per-page="-1"
      :search="search"
      :loading="loading"
      density="compact"
      style="border-radius: 5px; height: calc(100vh - 185px)"
    >
      <template #[`item.expires`]="{ item }">
        <span v-if="item.expiresUnixMillis > 0">{{
          new Date(item.expiresUnixMillis).toLocaleString()
        }}</span>
        <span v-else>never</span>
      </template>

      <template #[`item.attributes`]="{ item }">
        <v-chip v-for="attr in item.attributes" size="x-small" :key="attr" class="mr-1">
          {{ attr }}
        </v-chip>
      </template>

      <template #[`item.actions`]="{ item }">
        <v-menu open-on-click>
          <template #[`activator`]="{ props }">
            <v-icon v-bind="props" color="grey"> mdi-dots-horizontal </v-icon>
          </template>

          <v-list density="compact">
            <v-list-item class="d-flex flex-row" base-color="red" @click="RemoveKey(item)">
              <v-list-item-title>Remove</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>
      </template>
    </v-data-table>

    <v-dialog v-model="keyDialog" width="580">
      <v-card>
        <v-form ref="keyForm" @submit.prevent="ApplyKeyDialog">
          <v-card-title class="text-h6 ma-3"> New API Key </v-card-title>

          <v-text-field
            v-model="keyBuffer!.name"
            label="Name"
            variant="solo"
            flat
            bg-color="oddRow"
            density="compact"
            class="mx-7"
            :rules="[required]"
          />

          <v-select
            v-model="keyBuffer!.expiresUnixMillis"
            label="Expires"
            variant="solo"
            flat
            bg-color="oddRow"
            density="compact"
            class="mx-7 mt-2"
            :items="[
              { title: '1 day', value: 86400000 },
              { title: '1 week', value: 604800000 },
              { title: '1 month', value: 2592000000 },
              { title: '1 year', value: 31536000000 },
              { title: 'Never', value: 0 }
            ]"
          />

          <v-select
            v-model="keyBufferCustomPermissions"
            label="Type"
            variant="solo"
            flat
            bg-color="oddRow"
            density="compact"
            class="mx-7 mt-2"
            :items="[
              { title: 'WG Client', value: false },
              { title: 'Custom', value: true }
            ]"
          />

          <v-combobox
            v-if="keyBufferCustomPermissions"
            v-model="keyBuffer!.attributes"
            multiple
            label="Permissions"
            variant="solo"
            :return-object="false"
            flat
            bg-color="oddRow"
            density="compact"
            class="mx-7 mt-2"
            :items="[
              { title: 'Read Peers', value: 'read-peers' },
              { title: 'Write Peers', value: 'write-peers' },
              { title: 'Delete Peers', value: 'delete-peers' },
              { title: 'Read Users', value: 'read-users' },
              { title: 'Write Users', value: 'write-users' },
              { title: 'Delete Users', value: 'delete-users' }
            ]"
          />

          <v-text-field
            v-model="keyBuffer!.token"
            variant="solo"
            flat
            bg-color="oddRow"
            density="compact"
            class="mx-7 mt-2"
            :readonly="true"
            id="greyText"
          >
            <template #[`append`]>
              <v-btn
                size="small"
                class="mr-1"
                variant="flat"
                icon="mdi-content-copy"
                @click="CopyToClipboard(keyBuffer!.token)"
              />
            </template>
          </v-text-field>

          <div class="mx-13 mb-4">
            <v-icon color="grey" size="x-small" style="margin-bottom: 2px" class="ml-n5">
              mdi-information
            </v-icon>
            <span class="text-grey-darken-2">
              Store this key in a safe place. It will not be shown again.
            </span>
          </div>

          <v-card-actions class="mb-3 mr-5">
            <v-spacer />

            <v-btn color="secondary" variant="outlined" @click="keyDialog = false"> Cancel </v-btn>

            <v-btn color="secondary" type="submit" variant="flat"> Create </v-btn>
          </v-card-actions>
        </v-form>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<style>
#greyText {
  color: rgb(101, 101, 101);
}
</style>
