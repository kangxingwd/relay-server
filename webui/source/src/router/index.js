import Vue from 'vue';
import Router from 'vue-router';
import {getActiveRouters} from '../utils/index';
import fullRoutes from '../router/fullpath';
// const login = () => import('@/views/login/Login');

Vue.use(Router);

let baseRoute = [
    // {
    //   path: '/',
    //   redirect: '/login'
    // },
    // {
    //   path: '/login',
    //   name: 'login',
    //   components: login,
    // },
    {
      path: '*',
  }];
const permission = {home: 1, parent: 1, child1: 1, child2: 1};
const activeRoutes = getActiveRouters(fullRoutes, permission);
if (activeRoutes.length !== 0) {
  baseRoute = activeRoutes.concat(baseRoute);
}
let router = new Router({
  mode: 'history',
  routes: baseRoute
});

export default router;
