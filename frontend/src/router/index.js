import { createRouter, createWebHistory } from "vue-router";
import Login from "../views/Login.vue";
import App from "../App.vue";
import { auth } from "../firebase";

const routes = [
  {
    path: "/login",
    name: "Login",
    component: Login,
  },
  {
    path: "/",
    name: "App",
    component: App,
    meta: {
      requiresAuth: true, // Route requires authentication
    },
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// Navigation guard to protect routes that require authentication
router.beforeEach((to, from, next) => {
  const requiresAuth = to.matched.some((record) => record.meta.requiresAuth);
  const currentUser = auth.currentUser;

  if (requiresAuth && !currentUser) {
    next("/login"); // Redirect to login if not authenticated
  } else if (to.path === "/login" && currentUser) {
    next("/"); // Redirect to home if logged in and trying to access login
  } else {
    next(); // Proceed
  }
});

export default router;