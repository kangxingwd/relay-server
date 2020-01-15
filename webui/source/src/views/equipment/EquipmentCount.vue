<template>
  <div class="equipment">
    <div class="top">
      <EquipmentNum v-for="(item, index) of equipmentNumInfos.infos"
                    :key="index"
                    :info="item"
                    :showActive="equipmentNumInfos.showActive === index"
                    @click.native="handleItemClick($event, index)"
      />
    </div>
    <div class="main">
      <div class="search">
        <span>设备ID</span>
        <input v-model="search" type="text">
        <button @click="handleSearch">搜索</button>
      </div>
      <ElTable @goPage="goPage" :datas='datas' :columns="columns" :pages="pages"></ElTable>
    </div>
  </div>
</template>

<script>
import ElTable from '@/components/Tabel';
import EquipmentNum from './EquipmentNum';
import RelayDialog from './RelayDialog';
import McloudDialog from './McloudDialog';
import DetailDialog from './DetailDialog';
import { mapGetters, mapActions } from "vuex";

export default {
  name: 'EquipmentCount',
  components: {
    ElTable,
    EquipmentNum,
    RelayDialog,
    McloudDialog,
    DetailDialog,
  },
  computed:mapGetters([
    'currentPage',
    'offset',
  ]),
  watch: {
    datas() {
      this.$nextTick(() => {
        const height = document.getElementsByClassName('right-fix')[0].scrollHeight;
        // const table = document.getElementsByClassName('tabel')[0].offsetHeight;
        // const val = table >= 646 ? height : 959; 
        // const val = h - 80;
        const val = height < window.innerHeight ? window.innerHeight : height;
        if(document.getElementsByClassName('right-fix') && height) {
          document.getElementsByClassName('left')[0].style.height = val + 'px';
        }
      });
    }
  },
  data() {
    return {
      equipmentNumInfos: {
        showActive: 0,
        infos:[
          {
            title: '设备总数',
            num: ''
          }, {
            title: '在线设备数',
            num: ''
          }, {
            title: '离线设备数',
            num: ''
          }, {
            title: 'Mcloud设备数',
            num: ''
          }, {
            title: 'Relay设备数',
            num: ''
          }
        ]
      },
      search: '',
      type: '',
      datas: [],
      columns: [
        {
          title: '序号',
          data: 'no',
        },
        {
          title: '设备ID',
          data: 'dev_id',
          // render(val, item) {
          //   return `<span class="m" data-index='1'>M</span><span class="r" data-index='2'>R</span>`;
          // },
        },
        {
          title: '硬件型号',
          data: 'firmware',
        },
        {
          title: 'MAC地址',
          data: 'mac',
        },
        {
          title: '申请角色',
          data: 'role',
          click: 'handleRoleClick',
          render(val) {
            if (val == 1) {
              return `<span class="r" data-index='2'>R</span>`;
            } else if (val == 2) {
              return `<span class="m" data-index='1'>M</span>`;
            } else if (val == 3) {
              return `<span class="m" data-index='1'>M</span><span class="r" data-index='2'>R</span>`;
            } else {
              return val;
            }
          }
        },
        {
          title: '操作',
          click: 'openDetailDialog',
          render() {
            return `<span style="color: #1E90FF; cursor: pointer;">详情</span>`
          }
        }
      ],
      pages: {
        num: 0,
        prePage: 10,
        pageSelect: [10,20,50,100],
      }
    }
  },
  methods: {
    goPage(cur, offset) {
      if (this.search) {
        this.handleSearch();
      } else {
        this.getList();
      }
    },
    handleItemClick($event, index) {
      this.equipmentNumInfos.showActive = index;
      const typeArr = ['all', 'online', 'offline', 'mcloud', 'relay'];
      this.type = typeArr[index];
      this.search = '';
      this.getList();
    },
    getList() {
      this.$http(this.$api.dataType, {
        start: this.offset,
        length: this.pages.prePage,
        dataType: this.type || 'all',
      }).then((res) => {
        if (res.status === 0) {
          const data = res.data;
          let start = (this.currentPage - 1) * this.pages.prePage;
          data.forEach((item) => {
            start += 1;
            item.no = start;
            item.firmware = item.ext.firmware;
          });
          this.datas = data;
          this.pages.num = res.recordsFiltered;
        } else {
          // 弹窗
          this.showConfirm(res.msg);
        }
      }).catch((err) => {
        throw new Error(err);
      });
    },
    async openRelayDialog(item, roleIndex) {
      const roleArr = ['mcloud', 'relay'];
      let res = await this.$http(this.$api.role, {
        devId: item.dev_id,
        role: roleArr[roleIndex - 1],
      });
      let info;

      if (res.status === 0) {
        info = res.data;
      } else {
        this.showConfirm(res.msg);
      }

      this.$store.commit('SETRELOAD');
      this.$Modal.confirm({
        scrollable: false,
        width: '760px',
        closable: true,
        render: (h) => {
          return h(RelayDialog, {
            props: {info},
            on: {
              goDevID: (id) => {
                this.search = id;
                this.handleSearch();
                this.equipmentNumInfos.showActive = -1;
              }
            },
          });
        },
        onOk: () => {},
      });
    },
    async openMcloudDialog(item, roleIndex) {
      const roleArr = ['mcloud', 'relay'];
      let res = await this.$http(this.$api.role, {
        devId: item.dev_id,
        role: roleArr[roleIndex - 1],
      })
      let info;
      if (res.status === 0) {
        info = res.data;
      } else {
        this.showConfirm(res.msg);
      }
      
      this.$store.commit('SETRELOAD');
      this.$Modal.confirm({
        scrollable: false,
        width: '760px',
        closable: true,
        render: (h) => {
          return h(McloudDialog, {
            props: {info},
            on: {
              goDevId: (id) => {
                this.search = id;
                this.handleSearch();
                this.equipmentNumInfos.showActive = -1;
              }
            },
          });
        },
        onOk: () => {},
      });
    },
    openDetailDialog(item) {
      this.$Modal.confirm({
        scrollable: false,
        width: '760px',
        closable: true,
        render: (h) => {
          return h(DetailDialog, {
            props: {info: item},
            on: {},
          });
        },
        onOk: () => {},
        onCancel: () => { this.$Modal.remove(); const body = document.getElementsByTagName('body'); body.removeClass('v-transfer-dom');}
      });
    },
    getTotal() {
      this.$http(this.$api.total)
      .then((res) => {
        if (res.status === 0) {
          const data = res.data;
          this.equipmentNumInfos.infos.forEach((item, index) => {
            switch (index) {
              case 0:
                item.num = data.allNum;
                break;
              case 1:
                item.num = data.onlineNum;
                break;
              case 2:
                item.num = data.OfflineNum;
                break;
              case 3:
                item.num = data.mcloudNum;
                break;
              case 4:
                item.num = data.relayNum;
                break;
              break;
              default:
                break;
            }
          });
        } else {
          this.showConfirm(res.msg);
        }
      }).catch((err) => {
        throw new Error(err);
      });
    },
    handleSearch() {
      this.equipmentNumInfos.showActive = -1;
      this.$http(this.$api.devId, {
        start: this.offset,
        length: this.pages.prePage,
        devId: this.search,
      }).then((res) => {
        if (res.status === 0) {
          const data = res.data;
          let start = (this.currentPage - 1) * this.pages.prePage;
          data.forEach((item) => {
            start += 1;
            item.no = start;
          });
          this.pages.num = res.recordsFiltered;
          this.datas = data;
        } else {
          this.showConfirm(res.msg);
        }
      }).catch((err) => {
        throw new Error(err);
      });
    },
    handleRoleClick(item, index) {
      // 1: m, 2: r
      if (index == 1) {
        this.openMcloudDialog(item, index);
      } else if (index == 2) {
        this.openRelayDialog(item, index);
      }
      
    },
  },
  mounted() {
    this.getList();
    this.getTotal();
  }
}
</script>


<style lang="less" scoped>
  .equipment /deep/ .equipment-num{
    margin-right: 10px;
  }
  .equipment{
    .top{
      display: flex;
      margin-bottom: 24px;
    }
    .main{
      margin: 10px;
      background: #ffffff;
      .search{
        height: 80px;
        line-height: 80px;
        margin-left: 30px;
        vertical-align: middle;
        input{
          box-sizing: border-box;
          width:252px;
          height:34px;
          background:rgba(255,255,255,1);
          border:1px solid rgba(228, 228, 228, 1);
          padding-left: 10px;
          border-radius:3px;
        }
        button{
          width:80px;
          height:34px;
          line-height: 34px;
          margin-left: 30px;
          color: #ffffff;
          border-radius:3px;
          background: #1E90FF;
        }
      }
    }
  }
</style>