import VueRouter from "vue-router";

import Home from "./pages/Home.vue";
import Compare from "./pages/Compare.vue";
import RenderManifest from "./pages/Render.vue";

const routes = [
  {
    path: "/",
    name: "Home",
    component: Home
  },
  {
    path: "/compare",
    name: "Compare",
    component: Compare
  },
  {
    path: "/render-manifest",
    name: "RenderManifest",
    component: RenderManifest
  }
];

const router = new VueRouter({
  mode: "history",
  routes
});

export default router;