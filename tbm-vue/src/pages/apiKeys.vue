<script lang="ts" setup>
import { ref } from 'vue'

import { useStore } from "vuex";
import { key } from "../store";
const store = useStore(key);

const items = ref([
    { 
        uuid: "1",
        name: "OBT Servers",
        expires: "2022-12-31",
    }
])
const headers = [
    { title: 'Name', value: 'name' },
    { title: 'Expires', value: 'expires' },
    { title: '', value: 'actions', align: 'end' },
]
const search = ref('')

function RemoveKey(uuid: string) {
    store.state.ConfirmDialogText = 'Are you sure you want to remove this API key?'
    store.state.ConfirmDialogCallback = () => {
        console.log('Removing API key with UUID:', uuid)
    }
    store.state.ConfirmDialogShow = true
}
</script>

<template>
    <v-container fluid max-width="1300">
        <v-row no-gutters class="d-flex align-center">
            <span class="text-h4">API Keys</span>
            <v-icon size="x-large" color="rgb(186,194,202)" class="ml-3">mdi-key</v-icon>
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
                >New Key</v-btn>

            <v-spacer></v-spacer>
        </v-row>

        <v-data-table
        :items="items"
        :headers="headers"
        no-data-text="No keys found"
        :items-per-page="-1"
        :search="search"
        density="compact"
        style="border-radius: 5px; height: calc(100vh - 185px)"
        >
            <template v-slot:item.actions="{ item }">
                <v-menu open-on-click origin="top">
                    <template v-slot:activator="{ props }">
                        <v-icon v-bind="props" color="grey">mdi-dots-horizontal</v-icon>
                    </template>

                    <v-list density="compact">
                        <v-list-item class="d-flex flex-row" @click="RemoveKey(item.uuid)" base-color="red">
                            <v-list-item-title>Remove</v-list-item-title>
                        </v-list-item>
                    </v-list>
                </v-menu>
            </template>
        </v-data-table>

    </v-container>
</template>