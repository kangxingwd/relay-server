<template>
  <div class="use-status">
    <div class="title">
      <span class="icon"></span>
      <span class="desc">{{cicleInfo.title}}</span>
    </div>
    <div class="main">
      <ElCircle
        :size="170"
        :trail-width="4"
        :stroke-width="5"
        :percent="cicleInfo.num"
        stroke-linecap="square"
        stroke-color="#43a3fb">
        <div class="demo-Circle-custom" v-if="cicleInfo.num === null">
          <div class="loading">读取中...</div>
        </div>
        <div class="demo-Circle-custom" v-else>
            <p v-if="cicleInfo.totalNum">
              {{cicleInfo.totalNum}}
            </p>
            <p v-else>{{cicleInfo.num}}%</p>
            <span>
                {{cicleInfo.desc}}
            </span>
        </div>
    </ElCircle>
    </div>
    <div class="bottom" v-if="cicleInfo.desc1 && cicleInfo.num">
      <div>
        <span>{{cicleInfo.desc1}}：</span>{{cicleInfo.val1}}
      </div>
      <div>
        <span>{{cicleInfo.desc2}}：</span>{{cicleInfo.val2}}
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'UseStatus',
  props: ['cicleInfo'],
}
</script>

<style lang="less" scoped>
.use-status /deep/ .ivu-chart-circle{
  left: 50%;
  margin-top: 52px;
  transform: translateX(-50%);
}
  .demo-Circle-custom{
    & h1{
        color: #3f414d;
        font-size: 28px;
        font-weight: normal;
    }
    & p{
        color: #333333;
        font-size: 40px;
        margin: 10px 0 15px;
    }
    & span{
        display: block;
        padding-top: 15px;
        color: #657180;
        font-size: 14px;
        &:before{
            content: '';
            display: block;
            width: 50px;
            height: 1px;
            margin: 0 auto;
            background: #e0e3e6;
            position: relative;
            top: -15px;
        };
    }
    .loading{
      font-size:28px;
      color: #FF6A6A;
    }
  }
  .use-status{
    display: inline-block;
    width:520px;
    height:520px;
    background:#ffffff;
    border-radius:4px;
    .title{
      height: 68px;
      line-height: 68px;
      vertical-align: middle;
      .icon{
        display: inline-block;
        width:4px;
        height:18px;
        margin-left: 40px;
        background: #FF6A6A;
        vertical-align: middle;
      }
      .desc{
        display: inline-block;
        margin-left: 15px;
        vertical-align: middle;
      }
    }
    .bottom{
      display: flex;
      height: 100px;
      margin-top: 30px;
      flex-direction: column;
      justify-content: space-around;
      align-items: center;
      >div{
        display: inline-block;
        &:first-child{
          span{
            display: inline-block;
            text-align: right;
          }
        }
        &:nth-child(2){
          span{
            display: inline-block;
            text-align: right;
          }
        }
      }
    }
  }
</style>
