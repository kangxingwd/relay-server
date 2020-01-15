<template>
  <div class="equipment-num" :style="[showActive ? styleObj : '']">
    <span class="desc">{{info.title}}</span>
    <span class="num" :style="[showActive ? whiteColor : '']">{{info.num}}</span>
  </div>  
</template>

<script>
export default {
  name: 'EquipmentNum',
  props: ['info', 'showActive'],
  data() {
    return {
      styleObj: {
        background: 'linear-gradient(155deg,rgba(138,196,255,1),rgba(84,169,255,1))',
        boxShadow: '0px 0px 10px 0px rgba(114,151,255,0.6)',
        color: '#ffffff',
      },
      whiteColor: {
        color: '#ffffff'
      },
    }
  },
}
</script>

<style lang="less" scoped>
  .equipment-num{
    display: flex;
    flex-direction: column;
    justify-content: space-around;
    align-items: center;
    flex: 1;
    height: 120px;
    padding: 20px;
    color: #888888;
    background: #ffffff;
    border-radius:4px;
    cursor: pointer;
    .desc{
      font-size: 14px;
    }
    .num{
      color: #1E90FF;
      font-size:30px;
    }
  }
</style>

