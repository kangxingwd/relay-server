import Vue from 'vue';
import Vuex from 'vuex';
import login from './modules/login';
import common from './modules/common';
import { mutations } from './mutations';
import actions from './actions';
import getters from './getters';

Vue.use(Vuex);

export default new Vuex.Store({
  state: {
    groups: []
  },
  actions,
  mutations,
  getters,
  modules: {
    login,
    common,
  },
});
