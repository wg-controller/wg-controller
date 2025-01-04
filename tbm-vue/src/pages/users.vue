<script lang="ts" setup>
import { ref, onMounted } from 'vue'

import { useStore } from "vuex";
import { key } from "../store";
import { DELETE_Account, DELETE_AccountFailedAttempts, GET_Accounts, PATCH_Account, PUT_Account } from '@/api/methods';
import { VForm } from "vuetify/components";
import type { UserAccount, UserAccountWithPass } from '@/types/shared';
import { emailValidate, passwordValidate, required } from '@/utils/validators';
const store = useStore(key);

onMounted(() => {
    Init()
})

async function Init() {
    try {
        let val = await GET_Accounts()
        if (val != null) {
            items.value = val
        }
    } catch (error: any) {
        console.error(error)
        store.state.SnackBarText = "Error fetching users"
        store.state.SnackBarError = true
        store.state.SnackBarShow = true
    }
}

const items = ref<UserAccount[]>([])
const headers = ref([
    { title: 'Email', key: 'email' },
    { title: 'Role', key: 'role' },
    { title: 'Last Active', key: 'lastActiveUnixMillis' },
    { title: 'Failed Attempts', key: 'failedAttempts' },
    { title: 'Suspended', key: 'suspended' },
    { title: '', key: 'actions', align: 'end', sortable: false },
] as const)
const search = ref('')

const userDialog = ref(false)
const userDialogEditMode = ref(false)
const userBuffer = ref<UserAccount>()
const userPassword = ref('')
const userConfirmPassword = ref('')

async function RemoveUser(email: string) {
    store.state.ConfirmDialogTitle = 'Remove ' + email
    store.state.ConfirmDialogText = 'Are you sure you want to remove this user?'
    store.state.ConfirmDialogCallback = async () => {
        try {
            await DELETE_Account(email)
        } catch (error: any) {
            console.error(error)
            store.state.SnackBarText = error
            store.state.SnackBarError = true
            store.state.SnackBarShow = true
        } finally {
            Init()
        }
    }
    store.state.ConfirmDialogShow = true
}

async function ResetAttempts(email: string) {
    store.state.ConfirmDialogTitle = 'Reset Attempts'
    store.state.ConfirmDialogText = 'Are you sure you want to reset failed attempts for this user?'
    store.state.ConfirmDialogCallback = async () => {
        try {
            await DELETE_AccountFailedAttempts(email)
        } catch (error: any) {
            console.error(error)
            store.state.SnackBarText = "Error resetting attempts"
            store.state.SnackBarError = true
            store.state.SnackBarShow = true
        } finally {
            Init()
        }
    }
    store.state.ConfirmDialogShow = true
}

function NewUserDialog() {
    userDialogEditMode.value = false
    userBuffer.value = {
        email: '',
        role: 'user',
        failedAttempts: 0,
        lastActiveUnixMillis: 0,
    }
    userPassword.value = ''
    userConfirmPassword.value = ''
    userDialog.value = true
}

function EditUserDialog(user: UserAccount) {
    userDialogEditMode.value = true
    userBuffer.value = JSON.parse(JSON.stringify(user))
    userPassword.value = ''
    userConfirmPassword.value = ''
    userDialog.value = true
}

const userForm = ref<VForm>();
async function ApplyUserDialog() {
  if (userForm.value == null) {
    console.error("loginForm is null");
    return
  }

  if (userPassword.value != userConfirmPassword.value) {
      store.state.SnackBarText = "Passwords do not match"
      store.state.SnackBarError = true
      store.state.SnackBarShow = true
      return
  }

  let result = await userForm.value.validate();
  if (result.valid) {
    try {
        let account: UserAccountWithPass = {
          email: userBuffer.value!.email,
          role: userBuffer.value!.role,
          password: userPassword.value,
        }

        if (userDialogEditMode.value) {
            await PATCH_Account(account)
        } else {
            await PUT_Account(account)
        }
    } catch (error: any) {
        console.error(error)
        store.state.SnackBarText = "Error applying user"
        store.state.SnackBarError = true
        store.state.SnackBarShow = true
    } finally {
        userDialog.value = false
        Init()
    }
  } else {
    console.error("Form is not valid");
  }
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

      <template #[`item.lastActiveUnixMillis`]="{ item }">
        <span v-if="item.lastActiveUnixMillis > 0">{{ new Date(item.lastActiveUnixMillis).toLocaleString() }}</span>
        <span v-else>never</span>
      </template> 

      <template #[`item.suspended`]="{ item }">
        <v-chip
          v-if="item.failedAttempts >= 5"
          color="red"
          size="x-small"
        >SUSPENDED</v-chip>
      </template> 

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
              @click="EditUserDialog(item)"
            >
              <v-list-item-title>Edit</v-list-item-title>
            </v-list-item>
            <v-list-item
              class="d-flex flex-row"
              base-color="red"
              @click="ResetAttempts(item.email)"
            >
              <v-list-item-title>Reset Attempts</v-list-item-title>
            </v-list-item>
            <v-list-item
              class="d-flex flex-row"
              base-color="red"
              @click="RemoveUser(item.email)"
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
          ref="userForm"
          @submit.prevent="ApplyUserDialog"
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
            v-model="userBuffer!.email"
            label="Email"
            variant="solo"
            flat
            bg-color="oddRow"
            density="compact"
            class="mx-7"
            :rules="[emailValidate, required]"
            :disabled="userDialogEditMode"
          />

          <v-text-field
            v-model="userPassword"
            label="Password"
            variant="solo"
            flat
            bg-color="oddRow"
            density="compact"
            class="mx-7"
            type="password"
            :rules="[passwordValidate, required]"
          />

          <v-text-field
            v-model="userConfirmPassword"
            label="Confirm Password"
            variant="solo"
            flat
            bg-color="oddRow"
            density="compact"
            class="mx-7"
            type="password"
            :rules="[passwordValidate, required]"
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