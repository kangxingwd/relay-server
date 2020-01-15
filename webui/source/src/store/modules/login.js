/* eslint no-shadow: ["error", { "allow": ["state"] }] */
import types from '../types';
import { CookieStorage } from 'cookie-storage';
 
const cookieStorage = new CookieStorage();

const state = {
  isLogin: false,
  itemIndex: '0',
  userName: '',
};

const getters = {
  isLogin(state) {
    state.isLogin = cookieStorage.getItem('isLogin');
    return state.isLogin;
  },
  itemIndex(state) {
    state.itemIndex = cookieStorage.getItem('itemIndex');
    return state.itemIndex;
  },
  userName(state) {
    state.userName = cookieStorage.getItem('userName');
    return state.userName;
  },
};

const actions = {
  doLogin({ commit, state }){
    commit(types.DOLOGIN);
  },
};

const mutations = {
  [types.DOLOGIN](state, isLogin) {
    state.isLogin = isLogin;
    cookieStorage.setItem('isLogin', isLogin);
  },
  [types.SETITEMINDEX](state, itemIndex) {
    state.itemIndex = itemIndex;
    cookieStorage.setItem('itemIndex', itemIndex);
  },
  [types.SETUSERNAME](state, userName) {
    state.userName = userName;
    cookieStorage.setItem('userName', userName);
  },
};

export default {
  state,
  getters,
  actions,
  mutations
};