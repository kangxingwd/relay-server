/**
 * 由于在focus和blur事件中，emite到父组件，父组件一直循环接受到事件，并非1次，无法解决
 * 弃用此组件 
 */
<template>
  <div class="el-input">
    <div class="title">{{infos.title}}</div>
    <div class="input-wrapper">
      <input v-model="infos.val"
             type="text"
             :placeholder="infos.placeholder"
             @input="handleInput($event)"
             @focus="handleFocus"
             @blur="handleBlur"
      >
    </div>
    <div v-if="show" class="tip">
      {{tip}}
    </div>
  </div>
</template>

<script>
export default {
  name: 'ElInput',
  props: {
    infos: {
      default() {
        return {
          msg: '*错误提示',
          reg: /\d/,
          name: 'oldPassword',
          title: '原始密码',
          val: '2344',
        }
      }
    },
  },
  data() {
    return {
      show: false,
      tip: '',
    }
  },
  methods: {
    handleFocus() {
      // this.$emit('focusEmit');
      // return;
      if(!this.infos.val) {
        this.show = true;
        this.tip = this.infos.noMsg;
      }
    },
    handleBlur() {
      // this.$emit('blurEmit');
      // return;
      if (this.infos.val) {
        if (!this.infos.reg.test(this.infos.val)) {
          this.show = true;
          this.tip = this.infos.errMsg;
        } else {
          this.show = false;
          this.tip = '';
          return false;
          // 对新密码特殊处理
          if (this.infos.name === 'newPassword') {
            if (this.infos.val !== this.$parent.infos[0].val) {
              this.show = false;
              this.tip = '';
            }
          }
          // 对确认密码特殊处理
          if (this.infos.name === 'vnewPassword') {
            if (this.infos.val === this.$parent.infos[1].val) {
              this.show = false;
              this.tip = '';
            }
          }
        }
      } else {
        this.show = true;
        this.tip = this.infos.noMsg;
      }
    },
    handleInput($event) {
      const val = $event.srcElement.value;
      this.$parent[this.infos.name] = val;
      this.$emit('inputEmit');
      // 为空处理
      if (!val) {
        this.show = true;
        return;
      }

      if (!this.infos.reg.test(val)) {
        this.show = true;
        this.tip = this.infos.errMsg;
      }

      // if (!this.infos.reg.test(val)) {
      //   this.show = true;
      //   this.tip = this.infos.errMsg;
      // } else {
      //   this.show = false;
      //   this.tip = '';
      //   // 对新密码特殊处理
      //   if (this.infos.name === 'newPassword') {
      //     if (val === this.$parent.infos[0].val) {
      //       this.show = true;
      //       this.tip = this.infos.specialMsg;
      //     }
      //   }
      // }
      // // 对确认密码特殊处理,不要正则验证
      // if (this.infos.name === 'vnewPassword') {
      //   if (val !== this.$parent.infos[1].val) {
      //     this.show = true;
      //     this.tip = this.infos.specialMsg;
      //   }
      // }
    }
  },
}
</script>

<style lang="less" scoped>
  .el-input{
    position: relative;
    .title{
      height: 40px;
      line-height: 40px;
      color: #646464;
      vertical-align: middle;
    }
    .input-wrapper{
      width:520px;
      height:50px;
      background:rgba(255,255,255,1);
      border:1px solid rgba(204,204,204,1);
      border-radius:2px;
      input{
        box-sizing: border-box;
        width: 100%;
        height: 48px;
        top: 1px;
        padding-left: 10px;
        border: none;
      }
    }
    .tip{
      position: absolute;
      top: 55px;
      right: 0;
      transform: translateX(110%);
      color: #e67878;
    }
  }
</style>

