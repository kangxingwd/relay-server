<template>
  <div>
     <div class="top-nav">
      <span class="icon-user">
        <!-- <img :src="icon_user" alt=""> -->
      </span>
      <span style="vertical-align: bottom;">{{this.userName}}</span>
      <!-- <span class="icon-arrow"
            @click="handleLoginoutClick"
      >
        <img :src="icon_arrow" alt="">
      </span> -->
      <!-- <div class="login-out"
            v-if="showLoginout"
            @mouseenter="mouseenter"
            @mouseleave="mouseleave"
            @click="handleLoginout"
      >
        <div class="wrapper">
          <span class="icon-loginOut">
            <img :src="isMouseEnter ? icon_login_active : icon_login" alt="">
          </span>
          <span>退出</span>
        </div>
      </div> -->
      <div class="login-out"
            @mouseenter="mouseenter"
            @mouseleave="mouseleave"
            @click="handleLoginout"
      >
        <div class="wrapper">
          <span style="vertical-align:bottom;">退出</span>
          <span class="icon-loginOut">
            <img :src="icon_login_new" alt="">
          </span>
        </div>
      </div>
    </div>
    <Modal
        v-model="isShowModal"
        title="提示"
        @on-ok="ok"
        @on-cancel="cancel"
        width="400"
        class-name="vertical-center-modal my-modal">
        <p>确定退出？</p>
        </Modal>
    <div class="top-title">
      <span class="line"></span>
      <span class="title">{{titles[itemIndex]}}</span>
    </div>
  </div>
</template>

<script>
import { mapGetters, mapActions } from "vuex";

export default {
  name: 'TopNav',
  data() {
    return {
      icon_arrow: require('images/arrow-bottom.png'),
      icon_user: require('images/user.png'),
      // icon_login: require('images/backLogin.png'),
      // icon_login_active: require('images/backLogin-active.png'),
      icon_login_new: require('images/icon_login_new.png'),
      isMouseEnter: false,
      titles: ['系统状态', '设备统计', '修改密码'],
      showLoginout: false,
      isShowModal: false,
    }
  },
  computed:mapGetters([
    'itemIndex',
    'userName',
  ]),
  methods: {
    mouseenter() {
      this.isMouseEnter = true;
    },
    mouseleave() {
      this.isMouseEnter = false;
      setTimeout(() => {
        this.showLoginout = false;
      }, 3000);
    },
    handleLoginoutClick() {
      this.showLoginout = !this.showLoginout;
      setTimeout(() => {
        if (!this.isMouseEnter) {
          this.showLoginout = false;
        }
      }, 3000);
    },
    handleLoginout() {
      this.isShowModal = true;
      // this.$http(this.$api.loginOut)
      //   .then((res) => {
      //     if (res.status === 0) {
      //       this.$store.commit('DOLOGIN', false);
      //     } else {
      //       this.showConfirm(res.msg);
      //     }
      //   }).catch(() => {
      //     throw new Error('登出');
      //   });
    },
    ok() {
      this.$http(this.$api.loginOut)
        .then((res) => {
          if (res.status === 0) {
            this.$store.commit('DOLOGIN', false);
          } else {
            this.showConfirm(res.msg);
          }
        }).catch(() => {
          throw new Error('登出');
        });
    },
    cancel() {},
  },
}
</script>

<style lang="less" scoped>
  .top-nav /deep/ .ivu-modal-header-inner{
    font-size: 16px !important;
  }
  .my-modal{
    font-size: 16px !important;
    .ivu-modal-header-inner{
      font-size: 16px;
    }
    p{
      text-align: center;
      font-size: 16px;
    }
  }
  .top-nav{
    position: relative;
    height: 80px;
    line-height: 80px;
    margin-left: -20px;
    background: #ffffff;
    text-align: right;
    vertical-align: middle;
    color: #1e90ff;
    .icon-user{
      width: 16px;
      height: 16px;
      margin-right: 12px;
      vertical-align: middle;
    }
    .icon-arrow{
      margin: 0 33px 0 10px;
      cursor: pointer;
    }
    .login-out{
      // position: absolute;
      // right: 30px;
      position: relative;
      display: inline-block;
      width:100px;
      height:50px;
      margin-right: 30px;
      background:#ffffff;
      // box-shadow:0px 0px 6px 0px rgba(59,64,94,0.3);
      border-radius:2px;
      vertical-align: middle;
      cursor: pointer;
      transition: height 3s ;
      // &:hover{
      //   color: #1E90FF;
      // }
      .wrapper{
        position: absolute;
        width: 100px;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        text-align: center;
        .icon-loginOut{
          margin-right: 10px;
          vertical-align: middle;
          img{
            vertical-align: text-top;
          }
        }
      }
      // &::before{
      //   position: absolute;
      //   top: -44px;
      //   left: 50%;
      //   content: '\25B2';
      //   color: #ffffff;
      //   font-size: 12px;
      //   text-shadow: 0px -2px 8px rgba(59,64,94,0.3);
      // }
    }
  }
  .top-title{
    height: 60px;
    line-height: 60px;
    padding-left: 30px;
    vertical-align: middle;
    .line{
      display: inline-block;
      width:4px;
      height:18px;
      background:rgba(255,106,106,1);
    }
    .title{
      font-size: 20px;
      color: #8F9DF6;
      margin-left: 10px;
    }
  }
</style>


