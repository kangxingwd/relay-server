<template>
  <div class="detail-dialog">
    <div class="title">
      详情
      <span class="close" @click="close"></span>
    </div>
    <div class="line"></div>
    <div class="main-top">
      <div class="main-top-left main-top-common">
        <div>设备ID：{{info.dev_id}}</div>
        <div>硬件型号：{{info.ext.firmware}}</div>
        <div>MAC地址：{{info.mac}}</div>
        <div>域名：{{info.hostname}}</div>
        <div>申请角色：{{roleDesc}}</div>
      </div>
      <div class="main-top-right main-top-common">
        <div>软件型号：{{info.ext.software}}</div>
        <div>软件版本：{{info.ext.release}}</div>
        <div>入网时间：{{info.join_time}}</div>
        <div>厂家信息：{{info.vendor}}</div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'DetailDialog',
  props: ['info'],
  computed: {
    roleDesc() {
        const role = this.info.role ;
        if(role == 1) {
          return 'Relay'
        } else if (role == 2) {
          return 'Mcloud'
        } else if(role == 3) {
          return 'Mcloud+Relay'
        }
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
  .detail-dialog{
    .title{
      // 弹框body有padding: 16px; ui是60px 60-32=28px
      height: 28px;
      line-height: 28px;
      font-size:18px;
      color: #333333;
      text-align: center;
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
      margin-top: 16px;
      border-bottom: 1px solid #EAEAEA;
    }
    .main-top{
      padding: 20px 50px;
      display: flex;
      .main-top-common{
        display: flex;
        width: 50%;
        flex-direction: column;
        justify-content: space-around;
        div{
          line-height: 36px;
        }
      }
      .main-top-right{
        margin-bottom: 36px;
      }
    }
  }
</style>

