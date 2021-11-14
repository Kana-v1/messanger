import Vue from 'vue'
import App from './App.vue'
import { BootstrapVue } from 'bootstrap-vue'
import store from './store'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

console.log(store.state)
Vue.use(BootstrapVue)

Vue.config.productionTip = false
new Vue({
  store,
  render: h => h(App),  
}).$mount('#app')
