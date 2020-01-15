<template>
  <div class="login">
    <div class="login-container">
      <div class="top">
        <span>中心服务器管理系统</span>
        <span>Central server management system</span>
      </div>
      <div class="bottom">
        <div class="el-input-wrapper">
          <div class="el-input">
            <div class="icon"></div>
            <input v-model="info.name" @input="handleInput($event, 'user')" type="text" placeholder="请输入账号">
          </div>
          <div v-if="tips.user.show" class="tip">{{tips.user.msg}}</div>
        </div>
        <div class="el-input-wrapper">
          <div class="el-input">
            <div class="icon"></div>
            <input v-model="info.password" v-on:keyup.13="handleLoginClick" @input="handleInput($event, 'password')"  type="password" placeholder="请输入密码">
          </div>
          <div v-if="tips.password.show" class="tip">{{tips.password.msg}}</div>
        </div>
        <button @click="handleLoginClick">登录</button>
      </div>
    </div>
  </div>
</template>

<script>
import { mapGetters, mapActions } from "vuex";

export default {
  name: 'login',
  computed:mapGetters([
    'isLogin',
  ]),
  data() {
    return {
      tips: {
        user: {
          show: false,
          reg: /.*/,
          msg: '请输入账号',
        },
        password: {
          show: false,
          reg: /^[a-zA-Z0-9_-]{5,20}$/,
          msg: '请输入密码',
        },
      },
      info: {
        name: '',
        password: '',
      }
    }
  },
  beforeRouteEnter (to, from, next) {
    next(vm => {
      if (vm.isLogin === 'true') {
        vm.$store.commit('DOLOGIN', false);
      }
    })
  },
  methods:mapActions(['doLoginAsync']),
  methods: {
    /** 
     * 点击登录按钮 
     */
    handleLoginClick() {
      const check = () => {
        if (!this.info.name) {
          this.showConfirm('请输入账户名');
          return false;
        } else if (!this.tips.user.reg.test(this.info.name)) {
          this.showConfirm('账户名格式不正确');
          return false;
        }
        if(!this.info.password) {
          this.showConfirm('请输入密码');
          return false;
        } else if (!this.tips.password.reg.test(this.info.password)) {
          this.showConfirm('密码格式不正确');
          return false;
        }
        return true;
      }
      if (!this.tips.user.show && !this.tips.password.show && this.info.name && this.info.password) {
        this.$store.dispatch('doLoginAsync', this.info);
      }
    },
    /** 
     * 账号和密码校验 
     */
    handleInput($event, type){
      const val = $event.srcElement.value;
      // 为空处理
      if (!val) {
        this.tips[type].show = true;
        return;
      }
      if (type === 'user' || type === 'password') {
        if (!val) { // !this.tips[type].reg.test(val)
          this.tips[type].show = true;
          return;
        }
        this.tips[type].show = false;
      }
    },
  }
};
</script>

<style lang="less" scoped>
  .login{
    width: 100%;
    height: 100%;
    background: url('~images/background.png') no-repeat;
    background-size:100% 100%;
    .login-container{
      position: absolute;
      width: 536px;
      height: 437px;
      border-radius:4px;
      top: 42%;
      left: 50%;
      transform: translate(-50%, -50%);
      background: #ffffff;
      .top{
        display: flex;
        flex-direction: column;
        height: 122px;
        justify-content: center;
        align-items: center;
        span{
          &:first-child{
            font-size: 26px;
            color: #1E90FF;
          }
          &:last-child{
            font-size: 12px;
            color: #959595;
          }
        }
      }
      .bottom{
        height: 100%;
        padding: 22px 43px 0 43px;
        .el-input-wrapper{
          height: 100px;
          .el-input{
            position: relative;
            height:50px;
            border:1px solid rgba(223, 223, 223, 1);
            border-radius:2px;
            .icon{
              position: absolute;
              width: 24px;
              height: 24px;
              margin: 13px 0 13px 20px;
            }
            input{
              width: 100%;
              height: 100%;
              top: 1px;
              padding-left: 58px;
              box-sizing: border-box;
              border: none;
            }
          }
          .tip{
            height: 30px;
            line-height: 30px;
            font-size: 14px;
            color: red;
            vertical-align: middle;
          }
          &:first-child{
            .icon{
              background: url('~images/user-icon.png') no-repeat; 
            }
          }
          &:nth-child(2){
            .icon{
              background: url('~images/password-icon.png') no-repeat; 
            }
          }
        }
        button{
          width:452px;
          height:56px;
          font-size:22px;
          color: #ffffff;
          background:#1E90FF;
          border-radius:2px;
          margin-bottom: 48px;
          cursor: pointer;
        }
      }
    }
  }
</style>
