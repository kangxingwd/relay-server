<template>
  <div id="app">
    <!-- <template v-if="isLogin === 'true'">
      <div class="left">
        <Navigation/>
      </div>
      <div class="right-fix">
        <div class="right">
          <TopNav/>
          <router-view/>
        </div>
      </div>
    </template>
    <template v-else>
      <login />
    </template> -->
    <router-view/>
  </div>
</template>

<script>
import Login from 'views/login/Login';
import Navigation from 'layout/Navigation';
import TopNav from 'layout/TopNav';
import { mapGetters, mapActions } from "vuex";
import Admin from 'views/Admin';

export default {
  name: 'App',
  components: {
    Login,
    Navigation,
    TopNav,
    Admin,
  },
  computed:mapGetters([
    'isLogin',
  ]),
  methods: {
    getScrollTop(){  
      var scrollTop = 0;  
      if(document.documentElement && document.documentElement.scrollTop){  
          scrollTop=document.documentElement.scrollTop;  
      } else if(document.body){  
          scrollTop = document.body.scrollTop;  
      }  
      return scrollTop;  
    },
  },
};
</script>

<style lang='less' scoped>
@import '~styles/mixins.less';

#app {
  width: 100%;
  height: 100%;
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
    }
  }
}
</style>
