// Composables
import { createRouter, createWebHashHistory } from "vue-router";
import { useStore } from "vuex";
import { key } from "../store";

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      name: "Loading",
      path: "/",
      component: () => import("@/pages/loading.vue")
    },
    {
      name: "Login",
      path: "/login",
      component: () => import("@/pages/login.vue")
    },
    {
      name: "Clients",
      path: "/clients",
      component: () => import("@/pages/clients.vue")
    },
    {
      name: "Users",
      path: "/users",
      component: () => import("@/pages/users.vue")
    },
    {
      name: "API Keys",
      path: "/apikeys",
      component: () => import("@/pages/apiKeys.vue")
    },
    {
      name: "Settings",
      path: "/settings",
      component: () => import("@/pages/settings.vue")
    }
  ]
});

router.beforeEach((to, from, next) => {
  const store = useStore(key);
  // Allow access to loading page
  if (to.path === "/") {
    next();
    return;
  }

  // Check if user is logged in
  if (store.state.LoggedIn) {
    // Dont allow access to login page
    if (to.path === "/login") {
      next("/clients");
      return;
    }

    // Allow access to other pages
    next();
    return;
  } else {
    // Allow access to login page
    if (to.path === "/login" && store.state.Loaded) {
      next();
      return;
    }

    // Redirect to loading page
    next("/?redirect=" + to.path);
  }
});

// Workaround for https://github.com/vitejs/vite/issues/11804
router.onError((err, to) => {
  if (err?.message?.includes?.("Failed to fetch dynamically imported module")) {
    if (!localStorage.getItem("vuetify:dynamic-reload")) {
      console.log("Reloading page to fix dynamic import error");
      localStorage.setItem("vuetify:dynamic-reload", "true");
      location.assign(to.fullPath);
    } else {
      console.error("Dynamic import error, reloading page did not fix it", err);
    }
  } else {
    console.error(err);
  }
});

router.isReady().then(() => {
  localStorage.removeItem("vuetify:dynamic-reload");
});

export default router;
