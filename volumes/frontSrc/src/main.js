import Vue from 'vue'
import App from './App.vue'

import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';
import 'xe-utils'
import VXETable from 'vxe-table'
import 'vxe-table/lib/style.css'
Vue.use(ElementUI);
Vue.use(VXETable);

Vue.config.productionTip = false

new Vue({
  render: h => h(App),
}).$mount('#app')

if ('serviceWorker' in navigator) {
  navigator.serviceWorker.register('/sw.js')
    .then(reg => console.log('註冊啦～SW registered!'))
    .catch(err => console.log('Error!', err));
}