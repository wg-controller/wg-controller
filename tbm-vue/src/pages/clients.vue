<script lang="ts" setup>
import { BytesString, timeSinceSeconds, timeSinceString } from '@/utils/utils';
import { ref } from 'vue'

import { useStore } from "vuex";
import { key } from "../store";
import { required, hostValidate, subnetsValidate, ipValidate } from '@/utils/validators';
const store = useStore(key);

function CopyToClipboard(text: string) {
    navigator.clipboard.writeText(text).then(() => {
        store.state.SnackBarText = 'Copied to clipboard'
        store.state.SnackBarError = false
        store.state.SnackBarShow = true
    }, () => {
        store.state.SnackBarText = 'Failed to copy to clipboard'
        store.state.SnackBarError = true
        store.state.SnackBarShow = true
    })
}

const items = ref([
    { 
        uuid: "1",
        hostname: "Jack's MBP",
        lastSeenUnixMillis: 1735863146168,
        localTunAddress: "192.168.1.1",
        enabled: true,
        transmitBytes: 123456789,
        receiveBytes: 98763354321,
    },
    { 
        uuid: "2",
        hostname: "HD3-DOCKER-01",
        lastSeenUnixMillis: 1735863599253,
        localTunAddress: "192.168.1.2",
        enabled: true,
        transmitBytes: 1256789,
        receiveBytes: 984321,
    },
    { 
        uuid: "3",
        hostname: "HD4-DOCKER-01",
        lastSeenUnixMillis: 1735863246168,
        localTunAddress: "192.168.1.3",
        enabled: true,
        transmitBytes: 123456789,
        receiveBytes: 987654321,
    },
])
const headers = ref([
    { title: 'Hostname', key: 'hostname' },
    { title: 'Last Seen', key: 'lastSeenUnixMillis' },
    { title: 'TX Bytes', key: 'transmitBytes' },
    { title: 'RX Bytes', key: 'receiveBytes' },
    { title: 'Address', key: 'localTunAddress', sortable: false },
    { title: 'Subnets', key: 'subnets', sortable: false },
    { title: '', key: 'actions', align: 'end', sortable: false },
] as const)
const search = ref('')

const clientWizard = ref(false)
const clientWizardType = ref('TBM Client')
const clientWizardPlatform = ref('Linux')
const clientWizardStep = ref(1)
const clientWizardInstallCMD = function() {
    switch (clientWizardPlatform.value) {
        case 'Linux':
            return 'sudo apt-get install net-tbm-client'
        case 'MacOS':
            return 'brew install net-tbm-client'
        case 'Windows':
            return 'choco install net-tbm-client'
        default:
            return 'sudo apt-get install net-tbm-client'
    }
}
const clientStartCMD = 'sudo systemctl start net-tbm-client'

function CopyWGConfig() {
    const el = document.getElementById('wgConfig')
    if (el) {
        CopyToClipboard(el.innerText)
    }
}

function NextClientWizardStep() {
    if (clientWizardStep.value < 3) {
        clientWizardStep.value++
    } else {
        if (clientWizardType.value === 'TBM Client') {
            clientWizard.value = false
        } else {
            console.log('Deploying Wireguard client')
            clientWizard.value = false
        }
    }
}

function NextButtonText() {
    if (clientWizardStep.value < 3) {
        return 'Next'
    } else {
        if (clientWizardType.value === 'TBM Client') {
            return 'Finish'
        } else {
            return 'Apply'
        }
    }
}

function NextButtonColor() {
    if (clientWizardStep.value < 3) {
        return ''
    } else {
        return 'secondary'
    }
}

function NextButtonEnabled() {
    // Validate input
    if (clientWizardStep.value === 2 && clientWizardType.value === 'Wireguard Client') {
        if (hostValidate(clientBuffer.value.hostname) !== true || subnetsValidate(clientBuffer.value.subnets) !== true) {
            return false
        }
    }
    return true
}

const clientDialog = ref(false)
const clientBuffer = ref({ hostname: '', subnets: [], tunnelAddress: '' })
const wgConfigDialog = ref(false)

const platforms = ref([
    'Linux',
    'MacOS',
    'Windows',
])

function EditClient(client: any) {
    clientBuffer.value = client
    clientDialog.value = true
}

function ExportWGConfig(client: any) {
    clientBuffer.value = client
    wgConfigDialog.value = true
}

function RemoveClient(client: any) {
    store.state.ConfirmDialogTitle = 'Remove ' + client.hostname
    store.state.ConfirmDialogText = 'Are you sure you want to remove this client?'
    store.state.ConfirmDialogCallback = () => {
        console.log('Removing client with UUID:', client.uuid)
    }
    store.state.ConfirmDialogShow = true
}

function NewClientWizardDialog() {
    clientBuffer.value = { hostname: '', subnets: [] }
    clientWizardStep.value = 1
    clientWizardType.value = 'TBM Client'
    clientWizardPlatform.value = 'Linux'
    clientWizard.value = true
}
</script>

<template>
  <v-container
    fluid
    max-width="1300"
  >
    <v-row
      no-gutters
      class="d-flex align-center"
    >
      <span class="text-h4">Clients</span>
      <v-icon
        size="x-large"
        color="rgb(186,194,202)"
        class="ml-3"
      >
        mdi-server-network
      </v-icon>
    </v-row>
    <v-row
      no-gutters
      class="d-flex align-center justify-start mt-3 mb-3"
    >
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
      density="compact"
      style="border-radius: 5px; height: calc(100vh - 185px)"
    >
      <template #[`header.localTunAddress`]>
        Address<v-icon
          size="x-small"
          color="grey"
          class="ml-1"
          style="margin-bottom: 2px"
        >
          mdi-information-outline
        </v-icon>
      </template>
      <template #[`item.lastSeenUnixMillis`]="{ item }">
        <div
          class="indicator"
          :class="{ green: (timeSinceSeconds(item.lastSeenUnixMillis) < 60) }"
        />
        {{ timeSinceString(item.lastSeenUnixMillis) }}
      </template>

      <template #[`item.transmitBytes`]="{ item }">
        {{ BytesString(item.transmitBytes) }}
      </template>

      <template #[`item.receiveBytes`]="{ item }">
        {{ BytesString(item.receiveBytes) }}
      </template>

      <template #[`item.actions`]="{ item }">
        <v-menu
          open-on-click
          origin="top"
        >
          <template #[`activator`]="{ props }">
            <v-icon
              v-bind="props"
              color="grey"
            >
              mdi-dots-horizontal
            </v-icon>
          </template>

          <v-list density="compact">
            <v-list-item
              class="d-flex flex-row"
              @click="EditClient(item)"
            >
              <v-list-item-title>Edit</v-list-item-title>
            </v-list-item>

            <v-list-item
              class="d-flex flex-row"
              @click="ExportWGConfig(item)"
            >
              <v-list-item-title>Export WG Config</v-list-item-title>
            </v-list-item>

            <v-list-item
              class="d-flex flex-row"
              base-color="red"
              @click="RemoveClient(item)"
            >
              <v-list-item-title>Remove</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>
      </template>
    </v-data-table>

    <v-dialog
      v-model="clientWizard"
      width="700"
    >
      <v-stepper
        v-model="clientWizardStep"
        :items="['Client Type', 'Options', 'Deploy']"
      >
        <template #item.1>
          <v-card
            title="Select Client Type"
            flat
          >
            <v-radio-group
              v-model="clientWizardType"
              row
              class="ml-5"
            >
              <v-radio
                key="TBM Client"
                label="TBM Client"
                value="TBM Client"
              />
              <v-radio
                key="Wireguard Client"
                label="Wireguard Client"
                value="Wireguard Client"
              />
            </v-radio-group>
            <div
              v-if="clientWizardType === 'TBM Client'"
              class="mx-13"
            >
              <v-icon
                color="grey"
                size="x-small"
                style="margin-bottom: 2px"
                class="ml-n5"
              >
                mdi-information
              </v-icon>
              <span class="text-grey-darken-2">
                A device running the TBM Client software that is managed by the TBM Server.

              </span>
            </div>
            <div
              v-if="clientWizardType === 'Wireguard Client'"
              class="mx-13"
            >
              <v-icon
                color="grey"
                size="x-small"
                style="margin-bottom: 2px"
                class="ml-n5"
              >
                mdi-information
              </v-icon>
              <span class="text-grey-darken-2">
                A standard Wireguard client or 3rd party device that is not managed by the TBM Server (no subnet route provisioning).
              </span>
            </div>
          </v-card>
        </template>

        <template #item.2>
          <v-card
            v-if="clientWizardType === 'TBM Client'"
            title="Select Platform"
            flat
          >
            <v-radio-group
              v-model="clientWizardPlatform"
              row
              class="ml-5"
            >
              <v-radio
                v-for="platform in platforms"
                :key="platform"
                :label="platform"
                :value="platform"
              />
            </v-radio-group>
          </v-card>
          <v-card
            v-else
            title="Client Options"
            flat
          >
            <v-text-field
              v-model="clientBuffer.hostname"
              :rules="[required, hostValidate]"
              label="Hostname"
              variant="solo"
              flat
              bg-color="oddRow"
              density="compact"
              class="mx-7"
            />

            <v-combobox
              v-model="clientBuffer.subnets"
              :rules="[subnetsValidate]"
              multiple
              chips
              label="Subnet CIDRs (optional)"
              variant="solo"
              flat
              bg-color="oddRow"
              class="mx-7"
            />
          </v-card>
        </template>

        <template #item.3>
          <v-card
            v-if="clientWizardType === 'TBM Client'"
            title="Steps To Deploy"
            flat
          >
            <h4 class="mx-7 mb-1 mt-2">
              Install Client
            </h4>
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

            <h4 class="mx-7 mb-1 mt-3">
              Start Client
            </h4>
            <v-code class="mx-7">
              {{ clientStartCMD }}
              <v-icon
                color="grey"
                size="x-small"
                @click="CopyToClipboard(clientStartCMD)"
              >
                mdi-content-copy
              </v-icon>
            </v-code>

            <div
              class="ml-13 mr-9 mt-6"
            >
              <v-icon
                color="grey"
                size="x-small"
                style="margin-bottom: 2px"
                class="ml-n5"
              >
                mdi-information
              </v-icon>
              <span class="text-grey-darken-2">
                TBM Clients will automatically appear in the Clients list as they connect to the server.
              </span>
            </div>
          </v-card>

          <v-card
            v-else
            title="Wireguard Config"
            flat
          >
            <v-card-text
              style="background-color: white; border-radius: 5px;"
              class="px-5 mx-7"
            >
              <v-icon
                style="position: absolute; right: 40px; top: 65px;"
                color="grey"
                @click="CopyWGConfig()"
              >
                mdi-content-copy
              </v-icon>
              <pre>
                <code id="wgConfig">

[Interface]
PrivateKey = 12345678
Address = 10.8.0.3/24
DNS = 1.1.1.1

[Peer]
PublicKey = 1234567
PresharedKey = 1234567
AllowedIPs = 0.0.0.0/0, ::/0
PersistentKeepalive = 0
Endpoint = 3.27.229.123:51820
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

    <v-dialog
      v-model="clientDialog"
      width="560"
    >
      <v-card>
        <v-form
          ref="entryForm"
        >
          <v-card-title
            class="text-h6 ma-3"
          >
            Edit Client
          </v-card-title>

          <v-text-field
            v-model="clientBuffer.hostname"
            label="Hostname"
            :rules="[required, hostValidate]"
            variant="solo"
            flat
            bg-color="oddRow"
            density="compact"
            class="mx-7"
          />

          <v-combobox
            v-model="clientBuffer.subnets"
            :rules="[subnetsValidate]"
            multiple
            chips
            label="Subnet CIDRs (optional)"
            variant="solo"
            flat
            bg-color="oddRow"
            class="mx-7"
          />

          <v-text-field
            v-model="clientBuffer.tunnelAddress"
            label="Tunnel Address"
            :rules="[required, ipValidate]"
            variant="solo"
            flat
            bg-color="oddRow"
            density="compact"
            class="mx-7"
          />

          <v-card-actions class="mb-3 mr-5">
            <v-spacer />

            <v-btn
              color="secondary"
              variant="outlined"
              @click="clientDialog = false"
            >
              Cancel
            </v-btn>

            <v-btn
              color="secondary"
              type="submit"
              variant="flat"
            >
              Apply
            </v-btn>
          </v-card-actions>
        </v-form>
      </v-card>
    </v-dialog>

    <v-dialog
      v-model="wgConfigDialog"
      width="700"
    >
      <v-card>
        <v-card-title
          class="text-h6 ma-3"
        >
          {{ clientBuffer.hostname }}
        </v-card-title>

        <v-card-text
          style="background-color: white; border-radius: 5px;"
          class="px-5 mx-7"
        >
          <v-icon
            style="position: absolute; right: 40px; top: 85px;"
            color="grey"
            @click="CopyWGConfig()"
          >
            mdi-content-copy
          </v-icon>
          <pre>
                <code id="wgConfig">

[Interface]
PrivateKey = 12345678
Address = 10.8.0.3/24
DNS = 1.1.1.1

[Peer]
PublicKey = 1234567
PresharedKey = 1234567
AllowedIPs = 0.0.0.0/0, ::/0
PersistentKeepalive = 0
Endpoint = 3.27.229.123:51820
                </code>
              </pre>
        </v-card-text>

        <v-card-actions class="mb-3 mr-5">
          <v-spacer />

          <v-btn
            color="secondary"
            variant="flat"
            @click="wgConfigDialog = false"
          >
            Close
          </v-btn>
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