// store.ts
import type { InjectionKey } from "vue";
import { createStore, Store } from "vuex";
import type { StoreStruct } from "@/types/types";

// define injection key
export const key: InjectionKey<Store<StoreStruct>> = Symbol();

export const store = createStore<StoreStruct>({
  state() {
    return {
      ConfirmDialogShow: false,
      ConfirmDialogText: "",
      ConfirmDialogCallback: () => {},
    };
  },
});
