import VueRouter from "vue-router";

import Inspect from "./pages/Inspect.vue";
import CompareChart from "./pages/CompareChart.vue";
import RenderManifest from "./pages/Render.vue";
import CompareManifest from "./pages/CompareManifest.vue"

const routes = [
  {
    path: "/",
    name: "Inspect",
    component: Inspect
  },
  {
    path: "/compare",
    name: "CompareChart",
    component: CompareChart
  },
  {
    path: "/render-manifest",
    name: "RenderManifest",
    component: RenderManifest
  },
  {
    path: "/compare-manifest",
    name: "CompareManifest",
    component: CompareManifest
  }
];

const router = new VueRouter({
  mode: "history",
  routes
});

export default router;