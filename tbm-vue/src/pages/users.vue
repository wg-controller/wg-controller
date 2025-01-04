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
const headers = ref([
    { title: 'Email', key: 'email' },
    { title: 'Role', key: 'role' },
    { title: 'Failed Attempts', key: 'failedAttempts' },
    { title: 'Suspended', key: 'suspended' },
    { title: '', key: 'actions', align: 'end', sortable: false },
] as const)
const search = ref('')

const userDialog = ref(false)
const userDialogEditMode = ref(false)
const userBuffer = ref({ email: '', password: '', confirmPassword: '' })

function RemoveUser(user: any) {
    store.state.ConfirmDialogTitle = 'Remove ' + user.email
    store.state.ConfirmDialogText = 'Are you sure you want to remove this user?'
    store.state.ConfirmDialogCallback = () => {
        console.log('Removing user with UUID:', user.uuid)
    }
    store.state.ConfirmDialogShow = true
}

function ResetAttempts(uuid: string) {
    store.state.ConfirmDialogTitle = 'Reset Attempts'
    store.state.ConfirmDialogText = 'Are you sure you want to reset failed attempts for this user?'
    store.state.ConfirmDialogCallback = () => {
        console.log('Resetting attempts for user with UUID:', uuid)
    }
    store.state.ConfirmDialogShow = true
}

function NewUserDialog() {
    userDialogEditMode.value = false
    userBuffer.value = { email: '', password: '', confirmPassword: '' }
    userDialog.value = true
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
      <span class="text-h4">Users</span>
      <v-icon
        size="x-large"
        color="rgb(186,194,202)"
        class="ml-3"
      >
        mdi-account-multiple
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
        @click="NewUserDialog()"
      >
        New User
      </v-btn>

      <v-spacer />
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
      <template #[`item.actions`]="{ item }">
        <v-menu
          open-on-click
          origin="top"
          width="150"
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
              @click="console.log(item)"
            >
              <v-list-item-title>Edit</v-list-item-title>
            </v-list-item>
            <v-list-item
              class="d-flex flex-row"
              base-color="red"
              @click="ResetAttempts(item.uuid)"
            >
              <v-list-item-title>Reset Attempts</v-list-item-title>
            </v-list-item>
            <v-list-item
              class="d-flex flex-row"
              base-color="red"
              @click="RemoveUser(item)"
            >
              <v-list-item-title>Remove</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>
      </template>
    </v-data-table>


    <v-dialog
      v-model="userDialog"
      width="460"
    >
      <v-card>
        <v-form
          ref="entryForm"
          @submit.prevent="
            store.state.ConfirmDialogCallback(),
            (store.state.ConfirmDialogShow = false)
          "
        >
          <v-card-title
            v-if="userDialogEditMode"
            class="text-h6 ma-3"
          >
            Edit User
          </v-card-title>

          <v-card-title
            v-else
            class="text-h6 ma-3"
          >
            New User
          </v-card-title>

          <v-text-field
            v-model="userBuffer.email"
            label="Email"
            variant="solo"
            flat
            bg-color="oddRow"
            density="compact"
            class="mx-7"
          />

          <v-text-field
            v-model="userBuffer.password"
            label="Password"
            variant="solo"
            flat
            bg-color="oddRow"
            density="compact"
            class="mx-7"
          />

          <v-text-field
            v-model="userBuffer.confirmPassword"
            label="Confirm Password"
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
              @click="userDialog = false"
            >
              Cancel
            </v-btn>

            <v-btn
              v-if="!userDialogEditMode"
              color="secondary"
              type="submit"
              variant="flat"
            >
              Create
            </v-btn>
            <v-btn
              v-if="userDialogEditMode"
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
  </v-container>
</template>