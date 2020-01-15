<template>
  <div class="relay-dialog" :key="new Date().toString()">
    <div class="title">
      Relay信息
      <span class="close" @click="close"></span>
    </div>
    <div class="line"></div>
    <div class="main-top">
      <div class="main-top-left main-top-common">
        <div>公网IP地址：{{info.data.public_ip}}</div>
        <div>wan IP地址：{{info.data.wan_ip}}</div>
        <div>心跳时间：{{info.data.heart_time}}</div>
        <div>连接状态：{{info.data.conn_state === 1 ? '在线' : '离线'}}</div>
      </div>
      <div class="main-top-right main-top-common">
        <div>负载状态：{{info.data.overload === 0 ? '无负载' : '负载'}}</div>
        <div>cpu：192.2555</div>
        <div>mem：192.2555</div>
        <div>已服务Mcloud个数：{{info.data.mcloud_num}}</div>
      </div>
    </div>
    <div class="main-bottom">
      <div class="desc">
        <span class="icon"></span>
        <span class="icon-desc">已服务MCloud</span>
      </div>
      <div class="table">
        <div class="thead">
          <div class="tr">
            <div class="th">序号</div>
            <div class="th">MCloud设备ID</div>
          </div>
        </div>
        <div class="tbody" :style="[info.association.length > 5 ? tbodyStyle : '']">
          <!-- <div class="tr">
            <div class="th">序号</div>
            <div class="th">MCloud设备ID</div>
          </div> -->
          <div class="tr"
               v-for="(item, index) of info.association"
               :key="index"
          >
            <div class="th">{{index + 1}}</div>
            <div class="th" style="color: #1E90FF;cursor:pointer">
              <span @click="handleIdClick($event, item.mcloud_id)">{{item.mcloud_id}}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { mapGetters, mapActions } from "vuex";

export default {
  name:'RelayDialog',
  props: ['info'],
  data() {
    return {
      tbodyStyle: {
        width: '670px',
      }
    }
  },
  methods: {
    test() {
      alert('heheh');
    },
    // 点击设备ID跳转
    handleIdClick($event, id) {
      this.$Modal.remove();
      this.$emit('goDevID', id);
    },
    close() {
      this.$Modal.remove();
    }
  },
}
</script>

<style lang="less" scoped>
  .relay-dialog{
    background:#ffffff;
    border-radius:4px;
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
            color: #666666;
            font-weight:bold;
            &:first-child{
              border-top: 1px solid #EAEAEA;
            }
          }
        }
        .tbody{
          max-height: 250px;
          overflow: auto;
          .tr{
            .trs;
          }
        }
      }
    }
  }
</style>

