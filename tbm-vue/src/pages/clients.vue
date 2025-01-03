<script lang="ts" setup>
import { timeSinceSeconds, timeSinceString } from '@/utils/utils';
import { ref } from 'vue'

import { useStore } from "vuex";
import { key } from "../store";
const store = useStore(key);

const items = ref([
    { 
        uuid: "1",
        hostname: "Jack's MBP",
        lastSeenUnixMillis: 1735863146168,
        localTunAddress: "192.168.1.1",
        enabled: true
    },
    { 
        uuid: "2",
        hostname: "HD3-DOCKER-01",
        lastSeenUnixMillis: 1735863599253,
        localTunAddress: "192.168.1.2",
        enabled: true
    },
    { 
        uuid: "3",
        hostname: "HD4-DOCKER-01",
        lastSeenUnixMillis: 1735863246168,
        localTunAddress: "192.168.1.3",
        enabled: true
    },
])
const headers = [
    { title: 'Hostname', value: 'hostname' },
    { title: 'Last Seen', value: 'lastSeenUnixMillis' },
    { title: 'Address', value: 'localTunAddress' },
    { title: 'Subnets', value: 'subnets' },
    { title: '', value: 'actions', align: 'end' },
]
const search = ref('')

function RemoveClient(uuid: string) {
    store.state.ConfirmDialogText = 'Are you sure you want to remove this client?'
    store.state.ConfirmDialogCallback = () => {
        console.log('Removing client with UUID:', uuid)
    }
    store.state.ConfirmDialogShow = true
}
</script>

<template>
    <v-container fluid max-width="1300">
        <v-row no-gutters class="d-flex align-center">
            <span class="text-h4">Clients</span>
            <v-icon size="x-large" color="rgb(186,194,202)" class="ml-3">mdi-server-network</v-icon>
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
                >
            </v-text-field>

            <v-btn
                color="appBar"
                variant="tonal"
                class="ml-2"
                height="40px"
                @click=""
                append-icon="mdi-plus"
                >New Client</v-btn>

            <v-spacer></v-spacer>
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
            <template v-slot:item.lastSeenUnixMillis="{ item }">
                <div class="indicator" :class="{ green: (timeSinceSeconds(item.lastSeenUnixMillis) < 60) }"></div>
                {{ timeSinceString(item.lastSeenUnixMillis) }}
            </template>

            <template v-slot:item.subnets="{ item }">
                {{ item.subnets || 'N/A' }}
            </template>

            <template v-slot:item.actions="{ item }">
                <v-menu open-on-click origin="top">
                    <template v-slot:activator="{ props }">
                        <v-icon v-bind="props" color="grey">mdi-dots-horizontal</v-icon>
                    </template>

                    <v-list density="compact">
                        <v-list-item class="d-flex flex-row" @click="console.log(item)">
                            <v-list-item-title>Edit</v-list-item-title>
                        </v-list-item>
                        <v-list-item class="d-flex flex-row" @click="console.log(item)">
                            <v-list-item-title>Deploy TBM</v-list-item-title>
                        </v-list-item>
                        <v-list-item class="d-flex flex-row" @click="console.log(item)">
                            <v-list-item-title>Deploy Wireguard</v-list-item-title>
                        </v-list-item>
                        <v-list-item class="d-flex flex-row" @click="RemoveClient(item.uuid)" base-color="red">
                            <v-list-item-title>Remove</v-list-item-title>
                        </v-list-item>
                    </v-list>
                </v-menu>
            </template>
        </v-data-table>
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