import types from './types';
// import router from '../router/index';
import $api from 'services/api';
import http from 'services/http';
import { showConfirm } from '@/utils/index';

const $http = http.httpServer;

const actions = {
  /**
   * 登录
   * @param {Object} - comnmit, state
   */
  doLoginAsync({commit, state}, info){
    let isDoLogin = false;
    let p = new Promise((resolve, reject) => {
      $http($api.login, info).then((res) => {
        if (res.status === 0) {
          isDoLogin = true;
          resolve(isDoLogin);
        } else {
          // 弹框
          showConfirm(res.msg);
        }
      }).catch(() => {
        throw new Error('登录失败');
      });
      
    });
    p.then(() => {
        commit(types.DOLOGIN, isDoLogin);
        commit(types.SETITEMINDEX, 0);
        commit(types.SETUSERNAME, info.name);
    }).catch((err) => {
        throw new Error(err);
    });
  },
};

export default actions;
