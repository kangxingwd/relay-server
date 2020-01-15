// noinspection JSAnnotator
import axios from 'axios';
import qs from 'qs';

const instance = axios.create({
  path: '',
  timeout: '10000'
});

instance.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded';

// 添加请求拦截器
instance.interceptors.request.use(function (config) {
  // 在发送请求之前做些什么
  return config;
}, function (error) {
  // 对请求错误做些什么
  return Promise.reject(error);
});

// 添加响应拦截器
instance.interceptors.response.use(function (response) {
  // 对响应数据做点什么
  return response;
}, function (error) {
  // 对响应错误做点什么
  return Promise.reject(error);
});

const httpServer = (opts, data) => {
  const common = {};
  const httpDefaultConfig = {
    url: opts.url,
    method: opts.method,
    timeout: 5000,
    params: Object.assign(common, data),
    data: qs.stringify(Object.assign(common, data))
  };
  if (opts.method === 'get') {
    delete httpDefaultConfig.data;
  } else if (opts.method === 'post') {
    delete httpDefaultConfig.params;
  }
  const promise = new Promise((resolve, reject) => {
    instance(httpDefaultConfig).then((res) => {
      if (res.data) {
        resolve(res.data);
      } else {
        resolve(res);
      }
    }).catch((error) => {
      reject(error);
    });
  });
  return promise;
};

export default {
  httpServer,
  instance,
};
