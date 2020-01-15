<template>
  <div class="confirm">
    <div class="title">
      提示
      <span class="close" @click="close"></span>
    </div>
    <div class="line"></div>
    <div class="content">
      <span>{{msg}}</span>
    </div>
    <button @click="close">确定</button>
  </div>
</template>

<script>
export default {
  name: 'ConfirmDialog',
  props: ['msg'],
  data() {
    return {
      // modal: true,
    }
  },
  methods: {
    close() {
      this.$Modal.remove();
    }
  }
}
</script>

<style lang="less" scoped>
  .confirm{
    .title{
      // 弹框body有padding: 16px; ui是60px 60-32=28px
      height: 28px;
      line-height: 28px;
      font-size:18px;
      color: #333333;
      text-align: left;
      vertical-align: middle;
      .close{
        position: absolute;
        display: inline-block;
        width: 30px;
        height: 30px;
        top: -5px;
        right: -10px;
        z-index: 9999999;
        cursor: pointer;
      }
    }
    .line{
      margin-top: 10px;
      height: 1px;
      border-bottom: 1px solid #EAEAEA
    }
    .content{
      display: flex;
      width: 100%;
      min-height: 150px;
      font-size: 16px;
      padding: 20px 0;
      justify-content: center;
      align-items: center;
    }
    button{
      position: relative;
      width:80px;
      height:34px;
      color: #ffffff;
      background:rgba(30,144,255,1);
      margin-bottom: 20px;
      border-radius:3px;
      left: 50%;
      transform: translateX(-50%);
    }
  }
</style>

