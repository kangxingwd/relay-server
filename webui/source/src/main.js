// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue';
import App from './App';
import router from './router';
import http from 'services/http';
import api from 'services/api';
import store from './store/index';
import { Circle, Modal, Message } from 'iview';
import 'iview/dist/styles/iview.css';
import { CookieStorage } from 'cookie-storage';
import { showConfirm } from '@/utils/index';

require('./assets/styles/index.less');

const cookieStorage = new CookieStorage();

Vue.component('ElCircle', Circle);
Vue.component('Modal', Modal);
Vue.component('Message', Message);

Vue.prototype.$http = http.httpServer;
Vue.prototype.$api = api;
Vue.prototype.$cookieStorage = cookieStorage;
Vue.prototype.$Modal = Modal;
Vue.prototype.showConfirm = showConfirm;

Vue.config.productionTip = false;

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  store,
  components: { App },
  template: '<App/>'
});
