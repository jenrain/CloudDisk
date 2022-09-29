// 导入需要的库
import Vue from 'vue'
import App from './App.vue'
import router from './router/router'
import axios from 'axios'
import store from './store'

import ElementUI from 'element-ui';
import 'element-ui/lib/theme-chalk/index.css';

// 导入全局css
import './assets/css/base.css'

// 创建axios对象，设置全局服务端URL
// 然后把axios对象设置为vue实例的全局属性，这样就不用在每个组件中单独导入axios了
const instance = axios.create({
  baseURL: "/api",
})
Vue.prototype.$axios = instance;

Vue.use(ElementUI);
Vue.use(router);

new Vue({
  store,
  router,
  render: h => h(App),
}).$mount('#app')
