/* eslint no-shadow: ["error", { "allow": ["state"] }] */
import types from '../types';

const state = {
    offset: 0,
    currentPage: 1,
    prePage: 10,
    reload: new Date().toString() || '',
};

const getters = {
  offset(state) {
    return state.offset;
  },
  reload(state) {
    return state.reload;
  },
  currentPage(state) {
    return state.currentPage;
  },
  prePage(state) {
    return state.prePage;
  },
};

const actions = {
    
};

const mutations = {
  [types.SETOFFSET](state, offset) {
    state.offset = offset;
  },
  [types.SETRELOAD](state) {
    state.reload = new Date().toString();
  },
  [types.SETCURRENTPAGE](state, curPage) {
    state.currentPage = curPage;
  },
  [types.SETPREPAGE](state, prePage) {
    state.prePage = prePage;
  },
};

export default {
    state,
    getters,
    actions,
    mutations
  };