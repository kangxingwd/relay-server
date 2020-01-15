<template>
  <div class="mcloud-dialog" :key="new Date().toString()">
    <div class="title">
      MCloud信息
      <span class="close" @click="close"></span>
    </div>
    <div class="line"></div>
    <div class="main-top">
      <div class="main-top-left main-top-common">
        <div>IP地址：{{info.data.mcloud_ip}}</div>
        <div>心跳时间：{{info.data.heart_time}}</div>
      </div>
      <div class="main-top-right main-top-common">
        <div>连接状态：{{info.data.conn_state === 1 ? '在线' : '离线'}}</div>
        <div>已连接Relay个数：{{info.association.length}}</div>
      </div>
    </div>
    <div class="main-bottom">
      <div class="desc">
        <span class="icon"></span>
        <span class="icon-desc">已连接Relay</span>
      </div>
      <div class="table">
        <div class="thead">
          <div class="tr">
            <div class="th">序号</div>
            <div class="th">Relay设备ID</div>
          </div>
        </div>
        <div class="tbody">
          <div class="tr"
               v-for="(item, index) of info.association"
               :key="index"
          >
            <div class="th">{{index + 1}}</div>
            <div class="th" style="color: #1E90FF;cursor:pointer">
              <span @click="handleDevId($event, item.relay_id)">{{item.relay_id}}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>

export default {
  name: 'McloudDialog',
  props: ['info'],
  methods: {
    close() {
      this.$Modal.remove();
    },
    handleDevId($event, id) {
      this.$Modal.remove();
      this.$emit('goDevId', id);
    },
  }
}
</script>

<style lang="less" scoped>
  .mcloud-dialog{
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
    }
    .main-bottom{
      padding: 0 14px;
      .desc{
        line-height: 40px;
        vertical-align: middle;
        .icon{
          display: inline-block;
          width:4px;
          height:18px;
          background: #FF6A6A;
          vertical-align: middle;
        }
        .icon-desc{
          margin-left: 15px;
          color: #333333;
          font-size: 16px;
        }
      }
      .table{
        padding: 0 20px;
        margin-bottom: 10px;
        .trs{
          display: flex;
          flex: 1;
          height: 50px;
          border-bottom: 1px solid #EAEAEA;
          border-left: 1px solid #EAEAEA;
          .th{
            flex: 1;
            text-align: center;
            align-self: center;
            border-right: 1px solid #EAEAEA;
            height: 100%;
            line-height: 50px;
            vertical-align: middle;
          }
        }
        .thead{
          .tr{
            .trs;
            &:first-child{
              border-top: 1px solid #EAEAEA;
            }
            color: #666666;
            font-weight:bold;
            background: #F9F9F9;
          }
        }
        .tbody{
          overflow: auto;
          .tr{
            .trs;
          }
        }
      }
    }
  }
</style>

