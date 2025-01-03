<script lang="ts" setup>
import { ref } from 'vue'

import { useStore } from "vuex";
import { key } from "../store";
const store = useStore(key);

const items = ref([
    { 
        uuid: "1",
        email: "jack.wdgt@gmail.com",
        role: "admin",
        failedAttempts: 0,
    }
])
const headers = [
    { title: 'Email', value: 'email' },
    { title: 'Role', value: 'role' },
    { title: 'Failed Attempts', value: 'failedAttempts' },
    { title: 'Suspended', value: 'suspended' },
    { title: '', value: 'actions', align: 'end' },
]
const search = ref('')

function RemoveUser(uuid: string) {
    store.state.ConfirmDialogText = 'Are you sure you want to remove this user?'
    store.state.ConfirmDialogCallback = () => {
        console.log('Removing user with UUID:', uuid)
    }
    store.state.ConfirmDialogShow = true
}
</script>

<template>
    <v-container fluid max-width="1300">
        <v-row no-gutters class="d-flex align-center">
            <span class="text-h4">Users</span>
            <v-icon size="x-large" color="rgb(186,194,202)" class="ml-3">mdi-account-multiple</v-icon>
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
                >New User</v-btn>

            <v-spacer></v-spacer>
        </v-row>

        <v-data-table
        :items="items"
        :headers="headers"
        no-data-text="No users found"
        :items-per-page="-1"
        :search="search"
        density="compact"
        style="border-radius: 5px; height: calc(100vh - 185px)"
        >
        <template v-slot:item.actions="{ item }">
                <v-menu open-on-click origin="top" width="150">
                    <template v-slot:activator="{ props }">
                        <v-icon v-bind="props" color="grey">mdi-dots-horizontal</v-icon>
                    </template>

                    <v-list density="compact">
                        <v-list-item class="d-flex flex-row" @click="console.log(item)">
                            <v-list-item-title>Edit</v-list-item-title>
                        </v-list-item>
                        <v-list-item class="d-flex flex-row" @click="RemoveUser(item.uuid)" base-color="red">
                            <v-list-item-title>Remove</v-list-item-title>
                        </v-list-item>
                    </v-list>
                </v-menu>
            </template>
        </v-data-table>

    </v-container>
</template>