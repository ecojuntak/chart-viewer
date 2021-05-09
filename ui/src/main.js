import Vue from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify';
import router from './router';
import VueRouter from 'vue-router'
import VueGtag from "vue-gtag";

Vue.config.productionTip = false
Vue.use(VueRouter)
Vue.component('inspect', require('./pages/Inspect').default);

Vue.use(VueGtag, {
    config: { 
      id: process.env.VUE_APP_GOOGLE_ANALYTIC_TOKEN
    },
  }, 
  router
);

new Vue({
  vuetify,
  router: router,
  render: h => h(App)
}).$mount('#app')
