<template>
  <div class="tabel">
    <div class="thead">
      <div class="tr"
           :style="trStyle"
      >
        <div class="th"
             v-for="(item, index) of columns"
             :key="index"
        >{{item.title}}</div>
      </div>
    </div>
    <div class="tbody">
      <div class="tr no-datas" v-if="datas.length === 0">
        <div class="td">暂无数据</div>
      </div>
      <div class="tr"
           :style="trStyle"
           v-for="(item, index) of datas"
           :key="index"
      >
        <div class="td"
             v-for="(column, index) of columns"
             :key="index"
             v-html="column.render && column.render(item[column.data], item) || item[column.data]"
             @click="column.click && parentOperate($event, column.click, item)"
        >
          {{item[column.data]}}
        </div>
      </div>
    </div>
    <div class="tfoot">
      <select name="" id="" v-model="pages.prePage" @change="selectPageChange($event)">
        <option
          v-for="(item, index) of this.pages.pageSelect"
          :key="index"
          :value="item"
        >{{item}}</option>
      </select>
      <li class="prePage" v-if="isPrePage" @click="goPrePage">上一页</li>
      <li v-for="(item, index) of showPageBtn"
            :key="index"
            :class="[{'active': item === currentPage && item !== 0}, {'page': item !== 0}, {'notButton': item == 0}]"
            @click="goPage($event, item)"
      >
      <a v-if="item !== 0">{{item}}</a>
      <a v-else>...</a>
      </li>
      <li class="nexPage" v-if="isNextPage" @click="goNextPage">下一页</li>
      <span>共{{totalPage}}页</span>
      <span>
        到第 <input type="number" v-model="toPageIndex"> 页
      </span>
      <button @click="jumpPage">确定</button>
    </div>
  </div>
</template>

<script>
import { mapGetters, mapActions } from "vuex";

export default {
  name: 'Tabel',
  props: {
    columns: {
      default() {
        return [
          {
            title: '序号',
            data: 'no',
            click: '',
          },
          {
            title: '日期',
            data: 'date',
            click: 'test',
            render(val, item) {
              return `<span class="m" data-index='1'>M</span><span class="r" data-index='2'>R</span>`;
            },
          },
          {
            title: '硬件型号',
            data: 'id',
          }
        ]
      }
    },
    datas: {
      default() {
        return []
      }
    }, 
    styles: {
      default() {
        return {
          width: '120px',
          rowHeight: '50px',
        }
      }
    },
    pages: {
      default() {
        return {
          num: 145,
          prePage: 20,
          pageSelect: [10,20,50,100],
        }
      },
    }
  },
  computed: {
    ...mapGetters([
      'offset'
    ]),
    totalPage() {
      return Math.ceil(this.pages.num / this.pages.prePage);
    },
    currentPage() {
      const curPage = (this.offset && (this.offset / this.pages.prePage) + 1) || 1;
      this.$store.commit('SETCURRENTPAGE', curPage);
      this.toPageIndex = curPage;
      return curPage;
    },
    isPrePage() {
      return (this.offset !== 0) && this.pages.num;
    },
    isNextPage() {
      return (this.offset + this.pages.prePage < this.pages.num) && this.pages.num;
    },
    showPageBtn() {
      let pageNum = this.totalPage,
          index = this.currentPage,
          arr = [];
      if (pageNum <= 5) {
        for(let i = 1; i <= pageNum; i++) {
          arr.push(i);
        }
        return arr;
      }
      if (index <= 2) return [1,2,3,0,pageNum]
      if (index >= pageNum -1) return [1,0, pageNum -2, pageNum -1, pageNum]
      if (index === 3) return [1,2,3,4,0,pageNum]
      if (index === pageNum -2) return [1,0, pageNum-3, pageNum-2, pageNum-1, pageNum]
      return [1,0, index-1, index, index + 1, 0, pageNum]
    }
  },
  data() {
    return {
      trStyle: {
        lineHeight: this.styles.rowHeight,
      },
      toPageIndex: null,
    }
  },
  methods: {
    goPage($event, pageIndex) {
      if (pageIndex == 0) return;
      const offset = (pageIndex - 1) * this.pages.prePage;
      this.$store.commit('SETOFFSET', offset);
      this.emitGoPage(pageIndex, this.offset);
    },
    goPrePage() {
      if (!this.offset) return;
      const curPage = this.currentPage - 1;
      const offset = (curPage - 1) * this.pages.prePage;
      this.$store.commit('SETOFFSET', offset);
      this.emitGoPage(curPage, offset);
    },
    goNextPage() {
      const offset = this.currentPage * this.pages.prePage;
      this.$store.commit('SETOFFSET', offset);
      this.emitGoPage(this.currentPage + 1, offset);
    },
    emitGoPage(curPage, offset) {
      this.$emit('goPage', curPage, offset);
    },
    jumpPage() {
      const curPage = this.toPageIndex;
      if (curPage > this.totalPage || curPage <= 0) return;
      const offset = (curPage - 1) * this.pages.prePage;
      this.$store.commit('SETOFFSET', offset);
      this.emitGoPage(curPage, offset);
    },
    parentOperate($event, clickName, item) {
      if ($event.srcElement.tagName !== 'SPAN') return;
      const index = event.srcElement.dataset.index || -1;
      this.$parent[clickName](item, index);
    },
    selectPageChange($event) {
      if (this.$parent.pages && this.$parent.pages.prePage) {
        const prePage = $event.srcElement.value;
        this.$parent.pages.prePage = prePage;
        this.$store.commit('SETOFFSET', 0);
        this.emitGoPage(1, 0);
      }
    },
  },
  mounted () {
    this.toPageIndex = this.currentPage;
  }
}
</script>

<style lang="less" scoped>
  @import '~styles/mixins.less';
  .trs{
      display: flex;
      flex-direction: row;
      justify-content: space-between;
      border-bottom: 1px solid #EAEAEA;
      &:nth-child(odd){
        background: #ffffff;
      }
      &:nth-child(even){
        background: @tbLine-BC;
      }
      .th{
        flex: 1;
        text-align: center;
        vertical-align: middle;
      }
      .td{
        flex: 1;
        text-align: center;
        vertical-align: middle;
      }
    }
  .tabel{
    padding: 0 30px;
    border-radius: 4px;
    border-bottom: 1px solid #EAEAEA;
    .thead{
      .tr{
        .trs;
        background: @tbLine-BC !important;
        font-family:MicrosoftYaHei-Bold;
        font-weight:bold;
      }
    }
    .tbody{
      .tr{
        .trs;
      }
      .no-datas{
        text-align: center;
        height: 50px;
        line-height: 50px;
        vertical-align: middle;
      }
    }
    .tfoot{
      text-align: right;
      height: 84px;
      line-height: 84px;
      vertical-align: middle;
      select{
        width: 52px;
        height:36px;
        background:#ffffff;
        border:1px solid rgba(218,218,218,1);
        text-align: center;
        border-radius:4px;
      }
      input{
        width:36px;
        height:36px;
        margin: 0 10px;
        background:#ffffff;
        border:1px solid rgba(218,218,218,1);
        border-radius:4px;
        text-align: center;
      }
      button{
        width:65px;
        height:36px;
        line-height: 36px;
        background:#ffffff;
        border:1px solid rgba(218,218,218,1);
        border-radius:4px;
        &:hover{
          background: #1E90FF;
          color: #ffffff;
        }
      }
      span{
        margin: 0 10px;
      }
      .page{
        display: inline-block;
        width:36px;
        height:36px;
        line-height: 36px;
        background:#ffffff;
        margin: 0 5px;
        border:1px solid rgba(218,218,218,1);
        border-radius:4px;
        text-align: center;
        cursor: pointer;
        a{
          color: #000000;
        }
      }
      .prePage{
        .page;
        width: 60px;
        &:hover{
          background: #1E90FF;
          color: #ffffff;
        }
      }
      .nexPage{
        .prePage;
      }
      .notButton{
        display: inline-block;
        width:36px;
        height:36px;
        text-align: center;
      }
      .active{
        background: #1E90FF;
        a{
          color: #ffffff !important;
        }
      }
      .hover{
        background: #1E90FF;
        color: #ffffff;
      }
    }
  }
  .tabel /deep/ .m{
    display: inline-block;
    width:45px;
    height:30px;
    margin: 0 10px;
    line-height: 30px;
    color: #ffffff;
    background:#60D8A1;
    border-radius:15px;
    cursor: pointer;
  }
  .tabel /deep/ .r{
    display: inline-block;
    width:45px;
    height:30px;
    margin: 0 10px;
    line-height: 30px;
    color: #ffffff;
    border-radius:15px;
    background: #E67878;
    vertical-align: middle;
    cursor: pointer;
  }
</style>


