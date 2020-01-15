# vue-project

> A Vue.js project

## Build Setup

``` bash
# install dependencies
npm install

# serve with hot reload at localhost:8080
npm run dev

# build for production with minification
npm run build

# build for production and view the bundle analyzer report
npm run build --report
```
## 1、项目使用的插件

- less

- axios

- qs

- vuex

## 2、表格组件

props参数

columns: [
    title: '' // 表头
    data: '' // 对应字段
]

style: {
    width: '' // 宽度
    rowHeight: '60px' // 行高
}

## ElInput组件

props 

infos: {
    msg: '*错误提示',
    reg: /\d/,
    name: 'oldPassword',
    title: '原始密码',
    val: '2344',
}