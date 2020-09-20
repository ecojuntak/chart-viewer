import Vue from 'vue'
import App from './App.vue'
import vuetify from './plugins/vuetify';
import router from './router';
import VueRouter from 'vue-router'
import VueAnalytics from 'vue-analytics';

Vue.config.productionTip = false
Vue.use(VueRouter)
Vue.component('inspect', require('./pages/Inspect').default);

Vue.use(VueAnalytics, {
  id: process.env.VUE_APP_GOOGLE_ANALITIC_TOKEN,
  router
})

new Vue({
  vuetify,
  router: router,
  render: h => h(App)
}).$mount('#app')
