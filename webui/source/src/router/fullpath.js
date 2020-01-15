import AccountCurity from 'views/account/AccountCurity';
import Equipment from 'views/equipment/EquipmentCount';
import Syetem from 'views/system/SystemStatus';
import Admin from 'views/Admin';
const login = () => import('@/views/login/Login');

export default [
  // {
  //   path: '/login',
  //   name: 'login',
  //   component: login,
  // },
  {
    path: '/admin',
    component: Admin,
    children: [
      {
        path: 'login',
        name: 'login',
        component: login,
      },
      {
        path: 'account',
        name: 'account',
        component: AccountCurity,
        hasChildren: false,
      },
      {
        path: 'equipment',
        name: 'equipment',
        component: Equipment,
        // redirect: '/admin/equipment',
      },
      {
        path: 'system',
        name: 'system',
        component: Syetem,
        // redirect: '/admin/syetem',
      },
    ]
  },
];
