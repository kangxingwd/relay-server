const apis = {
  // 登录
  login: {
    url: '/manage/login/login',
    method: 'post',
  },
  // 登出
  loginOut: {
    url: '/manage/login/loginOut',
    method: 'post',
  },
  // 更改密码
  update: {
    url: '/manage/login/update',
    method: 'post',
  },
  // 查询cpu，内存
  system: {
    url: '/manage/system',
    method: 'post',
  },
  // 根据种类查询设备列表
  dataType: {
    url: '/manage/device/dataType',
    method: 'post',
  },
  // 根据devId查询设备列表
  devId: {
    url: '/manage/device/devId',
    method: 'post',
  },
  // 查看某一个设备的角色信息
  role: {
    url: '/manage/device/role',
    method: 'post',
  },
  total: {
    url: '/manage/device/total',
    method: 'post',
  }
};

export default apis;
