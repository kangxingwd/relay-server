<template>
  <div class="system">
    <div class="main">
      <UseStatus v-for="(item, index) of circleDatas"
                :key="index"
                :cicleInfo="item"
      ></UseStatus>
    </div>
  </div>
</template>

<script>
import UseStatus from './UseStatus';

export default {
  name: 'SystemStatus',
  components: {
    UseStatus,
  },
  data() {
    return {
      circleDatas: [{
        title: '内存使用率',
        desc: '内存使用率',
        num: null,
        desc1: '可用内存',
        val1: '',
        desc2: '总内存',
        val2: ''
      }, {
        title: 'CPU使用率',
        desc: 'CPU使用率',
        num: null,
        desc1: '',
        val1: '',
        desc2: '',
        val2: ''
      }, {
        title: '连接数统计',
        desc: '连接数',
        num: null,
        totalNum: null,
        desc1: '已确认连接数',
        val1: '',
        desc2: '待关闭连接数',
        val2: ''
      }],
      timer: null,
    }
  },
  methods: {
    init() {
      this.$http(this.$api.system)
        .then((res) => {
          if (res.status === 0) {
            // const disk = res.data.diskSpace;
            const system = res.data.system;
            const tcp = res.data.tcpStatus;
            // 内存使用率
            // const discPercent = Math.ceil(((disk.size - disk.free) / disk.size) * 100);
            // 可用内存
            // const unusedDisk = Math.ceil(disk.free / 1024 / 1024);
            // 总内存
            // const totalDisk = Math.ceil(disk.size / 1024 / 1024);
            // 连接数百分比
            const freememPercentage = system.freememPercentage * 100;
            // 总连接数
            const totalmem = system.totalmen;
            // freemem 待关闭连接数
            const freemem = system.freemem;
            // 已确认连接数
            const established = tcp.ESTABLISHED;
            // 待关闭连接数
            const TIME_WAIT = tcp.TIME_WAIT;
            // 总连接数
            const tcpNum = established + TIME_WAIT;
            this.circleDatas[0].val1 = `${system.freememPercentage}MB`;
            this.circleDatas[0].val2 = `${system.totalmem}MB`;
            this.circleDatas[2].val1 = established;
            this.circleDatas[2].val2 = TIME_WAIT;
            this.circleDatas[2].totalNum = tcpNum;
            this.circleDatas[2].num = Math.ceil((established / tcpNum).toFixed(4) * 100);
            this.circleDatas[0].num = system.freememPercentage;
            this.circleDatas[1].num = system.cpu;
          } else {
            this.showConfirm(res.msg);
          }
        }).catch((err) => {
          throw new Error(err);
        });
    }
  },
  mounted () {
    this.init();
    this.timer = setInterval(this.init, 3000);
  },
  beforeDestroy () {
    clearInterval(this.timer);
  }
}
</script>


<style lang="less" scoped>
  .system /deep/ .use-status{
    &:first-child{
      margin-right: 30px;
    }
    &:nth-child(2){
      margin-right: 30px;
    }
  }
  .system{
    .main{
      display: flex;
    }
  }
</style>