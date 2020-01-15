<template>
    <div class="navigation">
        <div class="title">中心服务器</div>
        <ul class="navItems">
          <li :class="[{'active': itemIndex == index}, 'navItem']"
              v-for="(item, index) of navItems"
              :key="index"
              @click="handleItemClick($event, index)"
          >
            <div class="icon">
              <img :src="item.icon_active" alt="">
            </div>
            <span>{{item.title}}</span>
          </li>
        </ul>
    </div>
</template>

<script>
import { mapGetters, mapActions } from "vuex";

export default {
  name: 'Navigation',
  data() {
    return {
      navItems: [
        {
          title: '系统状态',
          path: '/admin/system',
          icon: require('images/system-icon.png'),
          icon_active: require('images/system-icon-active.png'),
        },
        {
          title: '设备统计',
          path: '/admin/equipment',
          icon: require('images/equipment-icon.png'),
          icon_active: require('images/equipment-icon-active.png'),
        },
        {
          title: '账号安全',
          path: '/admin/account',
          icon: require('images/account-icon.png'),
          icon_active: require('images/account-icon-active.png'),
        }
      ],
    }
  },
  computed:mapGetters([
    'itemIndex'
  ]),
  methods: {
    // 状态点击
    handleItemClick($event, index) {
      this.$store.commit('SETITEMINDEX', index);
    },
  },
}
</script>


<style lang="less" scoped>
  .navigation{
    width: 240px;
    height: 100%;
    background: #011A46;
    .title{
      height: 80px;
      line-height: 80px;
      font-size: 32px;
      color:#ffffff;
      text-align: center;
      vertical-align: middle;
    }
    .navItems{
      .navItem{
        height: 60px;
        line-height: 60px;
        color: #ffffff;
        vertical-align: middle;
        cursor: pointer;
        .icon{
          display: inline-block;
          width: 22px;
          line-height: 60px;
          margin-left: 34px;
          vertical-align: -webkit-baseline-middle;
          img{
            width: 100%;
            height: 100%;
          }
        }
        span{
          margin-left: 14px;
        }
      }
      .active{
        color: #ffffff;
        background: #1E90FF;
      }
    }
  }
</style>