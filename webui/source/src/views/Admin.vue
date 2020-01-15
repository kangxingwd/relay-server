<template>
  <div class="admin">
    <template v-if="isLogin === 'true'">
      <div class="left">
        <div class="clearfix" style="clear:both;"></div>
        <Navigation/>
      </div>
      <div class="right-fix">
        <div class="right">
          <TopNav/>
          <!-- <router-view/> -->
          <SystemStatus v-if="itemIndex == '0'" />
          <EquipmentCount v-else-if="itemIndex == '1'"/>
          <AccountCurity v-else-if="itemIndex == '2'"/>
        </div>
      </div>
    </template>
    <template v-else>
      <Login />
    </template>
  </div>
</template>

<script>
import Login from 'views/login/Login';
import Navigation from 'layout/Navigation';
import TopNav from 'layout/TopNav';
import AccountCurity from 'views/account/AccountCurity';
import EquipmentCount from 'views/equipment/EquipmentCount';
import SystemStatus from 'views/system/SystemStatus';
import { mapGetters, mapActions } from "vuex";

export default {
  components: {
    Login,
    Navigation,
    TopNav,
    AccountCurity,
    EquipmentCount,
    SystemStatus,
  },
  computed:mapGetters([
    'isLogin',
    'itemIndex',
  ]),
}
</script>

<style lang='less' scoped>
@import '~styles/mixins.less';

.admin {
  overflow: auto;
  width: 100%;
  height: 100%;
  background: @main-BC;
  .left{
    position: relative;
    float: left;
    width: 240px;
    height: 100%;
  }
  .right-fix{
    float: right;
    width: 100%;
    margin-left: -240px;
    background: @main-BC;
    .right{
      margin-left: 260px;
      min-width: 1300px;
    }
  }
}
</style>
