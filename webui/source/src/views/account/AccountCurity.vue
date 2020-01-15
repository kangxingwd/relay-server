<template>
  <div class="account">
    <div class="main">
      <!-- <ElInput v-for="(item, index) of infos"
               :key="index"
               :infos="item"
               @inputEmit="inputEmit"
               @focusEmit="focusEmit"
               @blurEmit="blurEmit"
               :ref="item.name"
      /> -->
      <div class="el-input" v-for="(info, index) of infos"
          :key="index"
      >
        <div class="title">{{info.title}}</div>
        <div class="input-wrapper">
          <input v-model="info.val"
                type="password"
                :placeholder="info.placeholder"
                @input="handleInput($event, index, info.name)"
                @focus="handleFocus($event, index, info.name)"
                @blur="handleBlur($event, index, info.name)"
          >
        </div>
        <div v-if="info.show" class="tip">
          {{info.tip}}
        </div>
      </div>
      <button @click="handleBtnClick">确认修改</button>
    </div>
  </div>
</template>

<script>
import ElInput from './ElInput';
export default {
  name: 'AccountCurity',
  components: {
    ElInput,
  },
  computed: {
    canSubmit() {
      let tag = true;
      this.infos.every((item) => {
        if (item.show) {
          tag = false;
          return false;
        }
        return true;
      });
      return tag;
    }
  },
  data() {
    return {
      oldPassword: '',
      newPassword: '',
      vnewPassword: '',
      show: true,
      tip: '',
      infos: [
        {
          noMsg: '请输入密码',
          errMsg: '密码输入错误',
          specialMsg: '',
          tip: '',
          reg: /^[a-zA-Z0-9_-]{5,20}$/,
          name: 'oldPassword',
          title: '原密码',
          val: '',
          placeholder: '请输入原密码',
          show: false,
        },
        {
          noMsg: '请输入密码',
          errMsg: '限制5-20个字符，允许输入字母、数字以及英文特殊字符-和_',
          specialMsg: '新密码不能与原密码相同',
          tip: '',
          reg: /^[a-zA-Z0-9_-]{5,20}$/,
          name: 'newPassword',
          title: '新密码',
          val: '',
          placeholder: '设置新密码',
          show: false,
        },
        {
          noMsg: '请输入密码',
          errMsg: '',
          specialMsg: '两次输入不一致',
          tip: '',
          reg: /^[a-zA-Z0-9_-]{5,20}$/,
          name: 'vnewPassword',
          title: '确认密码',
          val: '',
          placeholder: '确认密码',
          show: false,
        }
      ]
    }
  },
  methods: {
    handleInput($event, index, type) {
      const val = $event.srcElement.value;
      const infos = this.infos;
      const info = this.infos[index];
      info.val = val;
      this[type] = val;
      if (!val) {
        // 没有数据
        info.tip = info.noMsg;
        info.show = true;
      } else {
        // 有数据
        // 新密码特殊处理
        if(infos[0].val === infos[1].val) {
          infos[1].tip = infos[0].val && infos[1].specialMsg;
          infos[1].show = true;
        } else if (!infos[1].reg.test(infos[1].val)) {
          infos[1].tip = infos[1].val && infos[1].errMsg;
          infos[1].show = true;
        } else {
          infos[1].tip = ''
          infos[1].show = false;
        }
        // 确认密码特殊处理
        if(infos[1].val !== infos[2].val) {
          infos[2].tip = infos[2].val && infos[2].specialMsg;
          infos[2].show = true;
        } else {
          infos[2].tip = ''
          infos[2].show = false;
        }
        if (type === 'vnewPassword') return;
        if(!info.reg.test(val)) {
          // 未通过验证
          info.tip = info.errMsg;
          info.show = true;
        } else {
          // 通过验证
          if (type === 'newPassword') {
            if(infos[0].val === val) {
              info.tip = info.specialMsg;
              info.show = true;
            } else {
              info.tip = '';
              info.show = false;
            }
          }
          // 一般情况
          else {
            info.tip = info.errMsg;
            info.show = false;
          }
        }
      }
    },
    handleFocus($event, index, type) {
      const infos = this.infos;
      const info = this.infos[index];
      if(!info.val) {
        info.tip = info.noMsg;
        info.show = true;
      }
    },
    handleBlur($event, index, type) {
      const infos = this.infos;
      const info = this.infos[index];
      if (!info.val) {
        info.tip = info.noMsg;
        info.show = true;
      }
    },
    /** 
     * 修改密码
     */
    handleBtnClick() {
      let tag = this.oldPassword && this.newPassword && this.vnewPassword;
      if (this.canSubmit && tag) {
        this.$http(this.$api.update, {
          oldPassword: this.oldPassword,
          newPassword: this.newPassword,
        }).then((res) => {
          if(res.status === 0) {
            this.showConfirm('修改成功');
          } else {
            this.showConfirm(res.msg);
          }
        }).catch((err) => {
          throw new Error(err);
        });
      }
      
    },
  },
}
</script>


<style lang="less" scoped>
  .account{
    position: relative;
    min-height: 800px;
    margin: 0 30px 0 0;
    background: #ffffff;
    border-radius:4px;
    min-width: 600px;
    .main{
      position: absolute;
      top: 117px;
      left: 50%;
      transform: translateX(-50%);
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
          right: -10px;
          transform: translateX(100%);
          color: #e67878;
        }
      }
      button{
        width: 100%;
        height: 54px;
        margin-top: 50px;
        font-size:18px;
        color: #ffffff;
        background: #1E90FF;
        border-radius:4px;
        cursor: pointer;
      }
    }
  }
</style>

