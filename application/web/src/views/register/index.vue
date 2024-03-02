<template>
  <div class="register-container">
    <div class="register-form">
      <div class="title-container">
        <h3 class="title">注册新用户</h3>
      </div>
      <form @submit.prevent="handleRegister">
        <div class="form-group">
          <select v-model="userType" class="form-control" @change="handleUserTypeChange">
            <option value="user">普通用户注册</option>
            <option value="admin">管理员注册</option>
          </select>
        </div>
        <div v-if="userType === 'user'">
          <div class="form-group">
            <input type="text" class="form-control" v-model="id" placeholder="请输入您的用户id">
          </div>
          <div class="form-group">
            <input type="text" class="form-control" v-model="realname" placeholder="请输入您的用户名">
          </div>
          <div class="form-group">
            <input type="password" class="form-control" v-model="password" placeholder="请输入您的密码">
          </div>
          <div class="form-group">
            <input type="email" class="form-control" v-model="email" placeholder="请输入您的邮箱">
          </div>
        </div>

        <div v-else-if="userType === 'admin'">
          <!-- 管理员注册的输入框 -->
          <!-- 假设管理员注册只需要密码和选择机构类型 -->
          <div class="form-group">
            <el-select v-model="certType" class="form-control admin-select" placeholder="请选择证书类型">
              <el-option
                  v-for="type in CERT_TYPE"
                  :key="type"
                  :label="type"
                  :value="type">
              </el-option>
            </el-select>
          </div>

          <div class="form-group">
            <input type="password" class="form-control" v-model="password" placeholder="请输入管理员密码">
          </div>
        </div>
        <div class="form-actions">
          <button type="button" class="btn-back" @click="GoToLogin">返回登录</button>
          <button type="button" class="btn-register" @click="handleRegister">注册</button>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import {adminRegister, userRegister} from '@/api/userInfo'
import { Message } from 'element-ui';
import {getReviewTypes} from "@/api/admin";

export default {
  name: 'Register',
  data() {
    return {
      id: '',
      password: '',
      realname: '',
      email: '',
      userType: 'user', // 默认选择普通用户注册
      certType:'',
      CERT_TYPE: [] // 管理员注册时的机构类型，默认展示CET
    }
  },
  created() {
    const storedCETType = JSON.parse(window.localStorage.getItem('CET_TYPE'));
    if (storedCETType != null) {
      console.log("register: 管理员可注册证书类型已存在")
      this.CERT_TYPE = storedCETType;
    }else{
      console.log("register: 查询管理员可注册证书类型")
      this.loadCertTypes_register();
    }
  },
  methods: {
    loadCertTypes_register() {
      getReviewTypes()
          .then(res => {
            console.log("register/index(): 获取当前机构的所有证书类型: ")
            console.log(res);
            //存储到本地
            window.localStorage.setItem('CET_TYPE', JSON.stringify(res.data));
            this.CERT_TYPE=res.data;
          }).catch(error => {
        console.log('加载证书类型错误:', error);
      });
    },
    handleRegister() {

      if(this.userType==="user"){
        this.handUserRegister();
      }else if(this.userType==="admin"){
        this.handleAdminRegister();
      }

    },
    handUserRegister() {
      const registerData = {
        id: this.id,
        realname: this.realname,
        password: this.password,
        email: this.email,
      };
      console.log(registerData);
      userRegister(registerData).then(res => {
        console.log("handleRegister: ", res);
        if (res.msg === 'Register_Success') {
          this.$router.push('/login');
        } else if(res.msg === 'Register_already') {
          console.log('register/index.vue: 注册失败，该用户已存在, 请返回登录');
          Message({
            message: '用户已注册, 即将返回登录页面',
            type: 'warning',
            duration: 4000
          });
          setTimeout(() => {
            this.$router.push('/login');
          }, 5000);
        }
      }).catch(error => {
        console.error('register/index.vue: 注册错误:', error);
      });
    },

    handleAdminRegister() {
      const registerData = {
        reviewType: this.certType,
        password: this.password,
      };

      if(!this.certType || !this.password){
        console.log("register/index.vue: 注册失败，请填写完整信息");
        Message({
          message: '请填写完整信息',
          type: 'warning',
          duration: 3000
        });
      }

      adminRegister(registerData).then(
          res => {
            if(res.msg==="Admin_Register_Success"){
              console.log("register/index.vue: admin注册成功");
              console.log(res);
              Message({
                message: '管理员id为: ' + res.data,
                type: 'warning',
                duration: 5000
              });
              setTimeout(() => {
                this.$router.push('/login');
              }, 5000);
            }
          }
      ).catch(error=>{
        console.log("register/index.vue: admin注册失败");
        console.log(error);
      });
    },

    handleUserTypeChange() {
      // 切换用户类型时清空已有的数据
      this.id = '';
      this.password = '';
      this.realname = '';
      this.email = '';
      if (this.userType === 'admin') {
        this.certType = '';
      }
    },

    GoToLogin(){
      this.$router.push('/login');
    },
  }
}
</script>

<style lang="scss" scoped>
$bg: #2d3a4b;
$dark_gray: #889aa4;
$light_gray: #eee;
$text_color: #5a5a5a;
$blue:#007bff;

.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: $bg;
}

.register-form {
  background: #fff;
  padding: 40px;
  border-radius: 8px;
  box-shadow: 0 4px 8px rgba(0,0,0,0.1);
  width: 100%;
  max-width: 400px;

  .title-container {
    margin-bottom: 20px;

    .title {
      color: $bg;
      text-align: center;
    }
  }

  .form-group {
    margin-bottom: 20px;

    input, select {
      width: 100%;
      padding: 10px;
      border: 1px solid $dark_gray;
      border-radius: 4px;
      box-sizing: border-box; // Ensure padding does not affect the final size
    }

    .el-select.admin-select {
      width: 100%;
      .el-input__inner {
        padding: 10px;
        border-radius: 4px;
        border: 1px solid $dark_gray;
      }

      .el-input__suffix .el-input__suffix-inner {
        color: $dark_gray;
      }
    }

    .el-option.selected {
      color: $text_color;
    }
  }

  .form-actions {
    display: flex;
    justify-content: space-between;

    .btn-back, .btn-register {
      padding: 10px 20px;
      border: none;
      border-radius: 4px;
      cursor: pointer;
      color: #fff;
    }

    .btn-register {
      background-color: $blue;
    }

    .btn-back {
      background-color: $light_gray;
      color: $bg;
    }
  }
}

@media (max-width: 500px) {
  .register-form {
    padding: 20px;
    max-width: 90%;
  }
}
</style>
