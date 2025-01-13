<script lang="ts" setup>
import { BytesString, timeSinceSeconds, timeSinceString } from "@/utils/utils";
import { ref, onMounted, onBeforeUnmount } from "vue";
import { VForm } from "vuetify/components";
import { useStore } from "vuex";
import { key } from "../store";
import { required, hostValidate, subnetsValidate, ipValidate } from "@/utils/validators";
import VueQrcode from "vue-qrcode";

import type { Peer, PeerInit, ServerInfo } from "@/types/shared";
import {
  DELETE_Peer,
  GET_PeerInit,
  GET_Peers,
  PUT_Peer,
  PATCH_Peer,
  GET_ServerInfo
} from "@/api/methods";
const store = useStore(key);

const fetchInterval = ref<ReturnType<typeof setInterval>>();

onMounted(() => {
  Init(true);
  fetchInterval.value = setInterval(() => {
    Init(false);
  }, 10000);
});

onBeforeUnmount(() => {
  if (fetchInterval.value) {
    clearInterval(fetchInterval.value);
  }
});

async function Init(showLoading: boolean) {
  if (showLoading) {
    loading.value = true;
  }
  try {
    const val = await GET_Peers();
    if (val != null) {
      items.value = val;
    } else {
      items.value = [];
    }
  } catch (error: any) {
    console.error(error);
    store.state.SnackBarText = error;
    store.state.SnackBarError = true;
    store.state.SnackBarShow = true;
  } finally {
    loading.value = false;
  }
}

const items = ref<Peer[]>([]);
const headers = ref([
  { title: "Hostname", key: "hostname" },
  { title: "Last Seen", key: "lastSeenUnixMillis" },
  { title: "TX Bytes", key: "transmitBytes" },
  { title: "RX Bytes", key: "receiveBytes" },
  { title: "Remote Address", key: "remoteTunAddress", sortable: false },
  { title: "Remote Subnets", key: "remoteSubnets", sortable: false },
  { title: "Enabled", key: "enabled" },
  { title: "", key: "actions", align: "end", sortable: false }
] as const);

const search = ref("");
const loading = ref(true);

const clientWizard = ref(false);
const clientWizardType = ref("Managed Client");
const clientWizardPlatform = ref("Linux");
const clientWizardStep = ref(1);
const clientWizardInstallCMD = function () {
  switch (clientWizardPlatform.value) {
    case "Linux":
      return "sudo apt-get install wg-controller-client";
    case "MacOS":
      return "brew install wg-controller-client";
    case "Windows":
      return "choco install wg-controller-client";
    default:
      return "sudo apt-get install wg-controller-client";
  }
};
const clientStartCMD = "sudo systemctl start wg-controller-client";

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

function GenerateWGConfig(): string {
  return `
[Interface]
PrivateKey = ${clientBuffer.value!.privateKey}
Address = ${clientBuffer.value!.remoteTunAddress + serverInfo.value!.netmask}
DNS = ${serverInfo.value!.nameServers.join(", ")}

[Peer]
PublicKey = ${serverInfo.value!.publicKey}
PresharedKey = ${clientBuffer.value!.preSharedKey}
AllowedIPs = ${clientBuffer.value!.allowedSubnets.join(", ")}
PersistentKeepalive = ${clientBuffer.value!.keepAliveSeconds}
Endpoint = ${serverInfo.value!.publicEndpoint}
`;
}

function CopyWGConfig() {
  const el = document.getElementById("wgConfig");
  if (el) {
    CopyToClipboard(el.innerText);
  }
}

function DownloadWGConfig() {
  const el = document.getElementById("wgConfig");
  if (el) {
    const blob = new Blob([el.innerText], { type: "text/plain" });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = clientBuffer.value!.hostname + ".conf";
    a.click();
    window.URL.revokeObjectURL(url);
  }
}

function DownloadQR() {
  const imgEl = document.getElementById("qrCode");
  if (!imgEl) return;

  // Create a temporary link element
  const a = document.createElement("a");
  // Use the `src` of the image as the `href`
  a.href = (imgEl as HTMLImageElement).src;
  // Provide a default download filename
  a.download = clientBuffer.value!.hostname + ".png";
  // Simulate the click
  a.click();
}

async function NextClientWizardStep() {
  if (clientWizardStep.value < 3) {
    clientWizardStep.value++;
  } else {
    if (clientWizardType.value === "Managed Client") {
      clientWizard.value = false;
    } else {
      // Apply changes to Wireguard client
      try {
        await PUT_Peer(clientBuffer.value!);
      } catch (error: any) {
        console.error(error);
        store.state.SnackBarText = error;
        store.state.SnackBarError = true;
        store.state.SnackBarShow = true;
        return;
      } finally {
        clientWizard.value = false;
        Init(true);
      }
    }
  }
}

function NextButtonText() {
  if (clientWizardStep.value < 3) {
    return "Next";
  } else {
    if (clientWizardType.value === "Managed Client") {
      return "Finish";
    } else {
      return "Apply";
    }
  }
}

function NextButtonColor() {
  if (clientWizardStep.value < 3) {
    return "";
  } else {
    return "secondary";
  }
}

function NextButtonEnabled() {
  if (clientBuffer.value == undefined) {
    return false;
  }

  // Validate input
  if (clientWizardStep.value === 2 && clientWizardType.value === "Wireguard Client") {
    if (
      hostValidate(clientBuffer.value.hostname) !== true ||
      subnetsValidate(clientBuffer.value.remoteSubnets) !== true ||
      subnetsValidate(clientBuffer.value.allowedSubnets) !== true
    ) {
      return false;
    }
  }
  return true;
}

const clientDialog = ref(false);
const clientBuffer = ref<Peer>();
const serverInfo = ref<ServerInfo>();
const wgConfigDialog = ref(false);
const wgConfigQRDialog = ref(false);

const platforms = ref(["Linux", "MacOS", "Windows"]);

function StartEditClient(client: Peer) {
  clientBuffer.value = JSON.parse(JSON.stringify(client));
  clientDialog.value = true;
}

async function ApplyEditClient() {
  try {
    await PATCH_Peer(clientBuffer.value!);
  } catch (error: any) {
    console.error(error);
    store.state.SnackBarText = error;
    store.state.SnackBarError = true;
    store.state.SnackBarShow = true;
    return;
  } finally {
    Init(true);
    clientDialog.value = false;
  }
}

async function ToggleClientEnabled(client: Peer) {
  try {
    client.enabled = !client.enabled;
    await PATCH_Peer(client);
  } catch (error: any) {
    console.error(error);
    store.state.SnackBarText = error;
    store.state.SnackBarError = true;
    store.state.SnackBarShow = true;
    return;
  } finally {
    Init(true);
  }
}

async function ExportWGConfig(client: Peer, qrCode: boolean) {
  try {
    const ServerInfoVal = await GET_ServerInfo();
    if (ServerInfoVal != null) {
      serverInfo.value = ServerInfoVal;
    }
  } catch (error: any) {
    console.error(error);
    store.state.SnackBarText = error;
    store.state.SnackBarError = true;
    store.state.SnackBarShow = true;
    return;
  } finally {
    clientBuffer.value = client;
    if (qrCode) {
      wgConfigQRDialog.value = true;
    } else {
      wgConfigDialog.value = true;
    }
  }
}

function RemoveClient(client: Peer) {
  store.state.ConfirmDialogTitle = "Remove " + client.hostname;
  store.state.ConfirmDialogText = "Are you sure you want to remove this client?";
  store.state.ConfirmDialogCallback = () => {
    try {
      DELETE_Peer(client.uuid);
    } catch (error: any) {
      console.error(error);
      store.state.SnackBarText = error;
      store.state.SnackBarError = true;
      store.state.SnackBarShow = true;
      return;
    } finally {
      Init(true);
    }
  };
  store.state.ConfirmDialogShow = true;
}

async function NewClientWizardDialog() {
  // Get initial value for a potential new client
  let InitPeer: PeerInit;
  try {
    InitPeer = await GET_PeerInit();
    const ServerInfoVal = await GET_ServerInfo();
    if (ServerInfoVal != null) {
      serverInfo.value = ServerInfoVal;
    }
  } catch (error: any) {
    console.error(error);
    store.state.SnackBarText = error;
    store.state.SnackBarError = true;
    store.state.SnackBarShow = true;
    return;
  }

  clientBuffer.value = <Peer>{
    uuid: InitPeer.uuid,
    hostname: "",
    enabled: true,
    privateKey: InitPeer.privateKey,
    publicKey: InitPeer.publicKey,
    preSharedKey: InitPeer.preSharedKey,
    keepAliveSeconds: 15,
    localTunAddress: "",
    remoteTunAddress: InitPeer.remoteTunAddress,
    remoteSubnets: [],
    allowedSubnets: ["0.0.0.0/0"],
    lastSeenUnixMillis: 0,
    lastIPAddress: "",
    transmitBytes: 0,
    receiveBytes: 0,
    attributes: []
  };
  clientWizardStep.value = 1;
  clientWizardType.value = "Managed Client";
  clientWizardPlatform.value = "Linux";
  clientWizard.value = true;
}
</script>

<template>
  <v-container fluid max-width="1400">
    <v-row no-gutters class="d-flex align-center">
      <span class="text-h4">Clients</span>
      <v-icon size="x-large" color="rgb(186,194,202)" class="ml-3"> mdi-server-network </v-icon>
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
        @click="NewClientWizardDialog()"
      >
        New Client
      </v-btn>

      <v-spacer />
    </v-row>
    <v-data-table
      :items="items"
      :headers="headers"
      no-data-text="No clients found"
      :items-per-page="-1"
      :search="search"
      :loading="loading"
      density="compact"
      style="border-radius: 5px; height: calc(100vh - 185px)"
    >
      <template #[`item.lastSeenUnixMillis`]="{ item }">
        <div
          class="indicator"
          :class="{ green: timeSinceSeconds(item.lastSeenUnixMillis) < 120 }"
        />
        {{ timeSinceString(item.lastSeenUnixMillis) }}
      </template>

      <template #[`item.transmitBytes`]="{ item }">
        {{ BytesString(item.transmitBytes) }}
      </template>

      <template #[`item.receiveBytes`]="{ item }">
        {{ BytesString(item.receiveBytes) }}
      </template>

      <template #[`item.remoteSubnets`]="{ item }">
        <v-chip v-for="subnet in item.remoteSubnets" :key="subnet" class="mr-1" size="x-small">
          {{ subnet }}
        </v-chip>
      </template>

      <template #[`item.enabled`]="{ item }">
        <v-switch
          v-model="item.enabled"
          flat
          color="primary"
          density="compact"
          hide-details
          class="ml-2"
          @click="ToggleClientEnabled(item)"
        ></v-switch>
      </template>

      <template #[`item.actions`]="{ item }">
        <v-menu open-on-click origin="top">
          <template #[`activator`]="{ props }">
            <v-icon v-bind="props" color="grey"> mdi-dots-horizontal </v-icon>
          </template>

          <v-list density="compact">
            <v-list-item class="d-flex flex-row" @click="StartEditClient(item)">
              <v-list-item-title>Edit</v-list-item-title>
            </v-list-item>

            <v-list-item class="d-flex flex-row" @click="ExportWGConfig(item, false)">
              <v-list-item-title>Export WG Config</v-list-item-title>
            </v-list-item>

            <v-list-item class="d-flex flex-row" @click="ExportWGConfig(item, true)">
              <v-list-item-title>Export WG QR-Code</v-list-item-title>
            </v-list-item>

            <v-list-item class="d-flex flex-row" base-color="red" @click="RemoveClient(item)">
              <v-list-item-title>Remove</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>
      </template>
    </v-data-table>

    <v-dialog v-model="clientWizard" width="700">
      <v-stepper v-model="clientWizardStep" :items="['Client Type', 'Options', 'Deploy']">
        <template #item.1>
          <v-card title="Select Client Type" flat>
            <v-radio-group v-model="clientWizardType" row class="ml-5">
              <v-radio key="Managed Client" label="Managed Client" value="Managed Client" />
              <v-radio key="Wireguard Client" label="Wireguard Client" value="Wireguard Client" />
            </v-radio-group>
            <div v-if="clientWizardType === 'Managed Client'" class="mx-13">
              <v-icon color="grey" size="x-small" style="margin-bottom: 2px" class="ml-n5">
                mdi-information
              </v-icon>
              <span class="text-grey-darken-2"> A device running the client software. </span>
            </div>
            <div v-if="clientWizardType === 'Wireguard Client'" class="mx-13">
              <v-icon color="grey" size="x-small" style="margin-bottom: 2px" class="ml-n5">
                mdi-information
              </v-icon>
              <span class="text-grey-darken-2">
                A standard wireguard client or 3rd party device. Requires manual configuration.
              </span>
            </div>
          </v-card>
        </template>

        <template #item.2>
          <v-card v-if="clientWizardType === 'Managed Client'" title="Select Platform" flat>
            <v-radio-group v-model="clientWizardPlatform" row class="ml-5">
              <v-radio
                v-for="platform in platforms"
                :key="platform"
                :label="platform"
                :value="platform"
              />
            </v-radio-group>
          </v-card>
          <v-card v-else title="Client Options" flat>
            <v-text-field
              v-model="clientBuffer!.hostname"
              :rules="[required, hostValidate]"
              label="Hostname"
              variant="solo"
              flat
              bg-color="oddRow"
              density="compact"
              class="mx-7"
            />

            <v-text-field
              v-model="clientBuffer!.remoteTunAddress"
              label="Remote Tun Address"
              :rules="[required, ipValidate]"
              variant="solo"
              flat
              bg-color="oddRow"
              density="compact"
              class="mx-7"
            />

            <v-row no-gutters>
              <v-combobox
                v-model="clientBuffer!.remoteSubnets"
                :rules="[subnetsValidate, required]"
                multiple
                chips
                label="Remote Subnets"
                variant="solo"
                flat
                bg-color="oddRow"
                class="mx-7"
              />

              <v-tooltip
                text="Subnets on the client side that are available to the server"
                location="top"
                transition="none"
                close-delay="0"
              >
                <template #activator="{ props }">
                  <v-icon class="mt-4 mr-6" color="grey" v-bind="props"> mdi-help-circle </v-icon>
                </template>
              </v-tooltip>
            </v-row>

            <v-row no-gutters>
              <v-combobox
                v-model="clientBuffer!.allowedSubnets"
                :rules="[subnetsValidate, required]"
                multiple
                chips
                label="Allowed Subnets"
                variant="solo"
                flat
                bg-color="oddRow"
                class="mx-7"
              />
              <v-tooltip
                text="Subnets the client is allowed to access via the server"
                location="top"
                transition="none"
                close-delay="0"
              >
                <template #activator="{ props }">
                  <v-icon class="mt-4 mr-6" color="grey" v-bind="props"> mdi-help-circle </v-icon>
                </template>
              </v-tooltip>
            </v-row>
          </v-card>
        </template>

        <template #item.3>
          <v-card v-if="clientWizardType === 'Managed Client'" title="Steps To Deploy" flat>
            <h4 class="mx-7 mb-1 mt-2">Install Client</h4>
            <v-code class="mx-7">
              {{ clientWizardInstallCMD() }}
              <v-icon
                color="grey"
                size="x-small"
                @click="CopyToClipboard(clientWizardInstallCMD())"
              >
                mdi-content-copy
              </v-icon>
            </v-code>

            <h4 class="mx-7 mb-1 mt-3">Start Client</h4>
            <v-code class="mx-7">
              {{ clientStartCMD }}
              <v-icon color="grey" size="x-small" @click="CopyToClipboard(clientStartCMD)">
                mdi-content-copy
              </v-icon>
            </v-code>

            <div class="ml-13 mr-9 mt-6">
              <v-icon color="grey" size="x-small" style="margin-bottom: 2px" class="ml-n5">
                mdi-information
              </v-icon>
              <span class="text-grey-darken-2">
                Managed Clients will automatically appear in the Clients list as they connect to the
                server.
              </span>
            </div>
          </v-card>

          <v-card v-else title="Wireguard Config" flat>
            <v-card-text style="background-color: white; border-radius: 5px" class="px-5 mx-7">
              <v-icon
                style="position: absolute; right: 40px; top: 65px"
                color="grey"
                @click="CopyWGConfig()"
              >
                mdi-content-copy
              </v-icon>
              <pre>
                <code id="wgConfig">
{{ GenerateWGConfig() }}
                </code>
              </pre>
            </v-card-text>
          </v-card>
        </template>
        <template #next>
          <v-btn
            :disabled="!NextButtonEnabled()"
            :color="NextButtonColor()"
            @click="NextClientWizardStep()"
          >
            {{ NextButtonText() }}
          </v-btn>
        </template>
      </v-stepper>
    </v-dialog>

    <v-dialog v-model="clientDialog" width="650">
      <v-card>
        <v-form ref="entryForm" @submit.prevent="ApplyEditClient()">
          <v-card-title class="text-h6 ma-3"> Edit Client </v-card-title>

          <v-text-field
            v-model="clientBuffer!.hostname"
            label="Hostname"
            :rules="[required, hostValidate]"
            variant="solo"
            flat
            bg-color="oddRow"
            density="compact"
            class="mx-7"
          />

          <v-text-field
            v-model="clientBuffer!.remoteTunAddress"
            label="Remote Tun Address"
            :rules="[required, ipValidate]"
            variant="solo"
            flat
            bg-color="oddRow"
            density="compact"
            class="mx-7"
          />

          <v-row no-gutters>
            <v-combobox
              v-model="clientBuffer!.remoteSubnets"
              :rules="[subnetsValidate, required]"
              multiple
              chips
              label="Remote Subnets"
              variant="solo"
              flat
              bg-color="oddRow"
              class="mx-7"
            />

            <v-tooltip
              text="Subnets on the client side that are available to the server"
              location="top"
              transition="none"
              close-delay="0"
            >
              <template #activator="{ props }">
                <v-icon class="mt-4 mr-10" color="grey" v-bind="props"> mdi-help-circle </v-icon>
              </template>
            </v-tooltip>
          </v-row>

          <v-row no-gutters>
            <v-combobox
              v-model="clientBuffer!.allowedSubnets"
              :rules="[subnetsValidate, required]"
              multiple
              chips
              label="Allowed Subnets"
              variant="solo"
              flat
              bg-color="oddRow"
              class="mx-7"
            />
            <v-tooltip
              text="Subnets the client is allowed to access via the server"
              location="top"
              transition="none"
              close-delay="0"
            >
              <template #activator="{ props }">
                <v-icon class="mt-4 mr-10" color="grey" v-bind="props"> mdi-help-circle </v-icon>
              </template>
            </v-tooltip>
          </v-row>

          <div class="ml-13 mr-9 mt-2 mb-2">
            <v-icon color="grey" size="x-small" style="margin-bottom: 2px" class="ml-n5">
              mdi-information
            </v-icon>
            <span class="text-grey-darken-2">
              Standard wireguard clients require a config export before some changes will take
              effect.
            </span>
          </div>

          <v-card-actions class="mb-3 mr-5">
            <v-spacer />

            <v-btn color="secondary" variant="outlined" @click="clientDialog = false">
              Cancel
            </v-btn>

            <v-btn color="secondary" type="submit" variant="flat"> Apply </v-btn>
          </v-card-actions>
        </v-form>
      </v-card>
    </v-dialog>

    <v-dialog v-model="wgConfigDialog" width="700">
      <v-card>
        <v-card-title class="text-h6 ma-3">
          {{ clientBuffer!.hostname }}
        </v-card-title>

        <v-card-text style="background-color: white; border-radius: 5px" class="px-5 mx-7">
          <v-icon
            style="position: absolute; right: 40px; top: 85px"
            color="grey"
            @click="CopyWGConfig()"
          >
            mdi-content-copy
          </v-icon>
          <pre>
                <code id="wgConfig">
{{ GenerateWGConfig() }}
                </code>
              </pre>
        </v-card-text>

        <v-card-actions class="mb-3 mr-5 mt-2">
          <v-spacer />

          <v-btn color="secondary" variant="outlined" @click="wgConfigDialog = false">
            Close
          </v-btn>

          <v-btn color="secondary" variant="flat" @click="DownloadWGConfig()"> Download </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="wgConfigQRDialog" width="600">
      <v-card>
        <v-card-title class="text-h6 ma-3">
          {{ clientBuffer!.hostname }}
        </v-card-title>

        <v-card-text class="d-flex align-center justify-center" style="height: 290px">
          <vueQrcode
            id="qrCode"
            class="ma-3"
            :color="{
              dark: '#ebebeb',
              light: '#333333'
            }"
            type="image/png"
            :margin="1"
            :value="GenerateWGConfig()"
          />
        </v-card-text>

        <v-card-actions class="mb-3 mr-5 mt-2">
          <v-spacer />

          <v-btn color="secondary" variant="outlined" @click="wgConfigQRDialog = false">
            Close
          </v-btn>

          <v-btn color="secondary" variant="flat" @click="DownloadQR()"> Download </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<style>
.indicator {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  display: inline-block;
  margin-right: 5px;
  background-color: grey;
}

.green {
  background-color: rgb(69, 197, 69);
}
</style>
