<template>
  <div class="login-container">
    <div class="login-form">
      <div class="title-container">
        <h3 class="title">基于联盟链的电子证书管理系统</h3>
      </div>
      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <select v-model="loginType" class="form-control">
            <option value="user">普通用户登录</option>
            <option value="admin">管理员登录</option>
          </select>
        </div>

        <div class="form-group">
          <input type="text" class="form-control" v-model="id" placeholder="请输入您的ID">
        </div>
        <div class="form-group">
          <input type="password" class="form-control" v-model="password" placeholder="请输入您的密码">
        </div>

        <div class="form-actions">
          <button type="button" class="btn-register" @click="GoToRegister">注册</button>
          <button type="submit" class="btn-login">登录</button>
        </div>
        <!-- 忘记密码链接 -->
        <div class="forgot-password">
          <router-link to="/forgot-password">忘记密码?</router-link>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import {adminLogin, userLogin} from '@/api/userInfo'
import { Message } from 'element-ui';
import {getReviewTypes} from "@/api/admin";

export default {
  name: 'Login',
  data() {
    return {
      id: "",
      password:"",
      loginType: "user",
      CERT_TYPE:[]
    }
  },
  watch: {

  },
  created() {
    // 在组件创建时立即执行逻辑
    const hasToken = window.localStorage.getItem('user_id')
    const role = window.localStorage.getItem('user_role')

    if (hasToken !== null && role !== null) {
      console.log("已登录: {user_id: "+hasToken+ " , user_role: "+ role+ " }")
      if (this.$route.path === '/login') {
        // 如果已登录且在访问登录页面，则根据角色重定向
        this.$router.push({ path: '/' + role })
      } else {
        // 对于其他路径，可以添加更多的角色检查和重定向逻辑
      }
    } else {
      // 未登录状态的处理逻辑
      console.log("未登录: ")
      if (!['/login', '/register', '/404'].includes(this.$route.path)) {
        this.$router.push(`/login`) // 重定向到登录页
      }
    }


    // 调用获取证书类型的方法
    this.loadCertTypes();

  },
  methods: {

    loadCertTypes() {
      // 调用后端获取证书类型的方法
      getReviewTypes()
          .then(res => {
            console.log("获取当前机构的所有证书类型: ")
            console.log(res);

            // 更新全局变量CERT_TYPE的值
            this.CERT_TYPE.splice(0, this.CERT_TYPE.length, ...res.data);

            //存储到本地
            window.localStorage.setItem('CET_TYPE', JSON.stringify(res.data));
            // const storedCETType = JSON.parse(window.localStorage.getItem('CET_TYPE'));

          }).catch(error => {
            console.log('加载证书类型错误:', error);
          });
    },

    handleLogin() {
      const loginData = {
        id: this.id,
        password: this.password,
        // type: this.loginType // 添加登录类型
      };
      console.log(loginData)

      const user_id=this.id;
      const user_role=this.loginType;

      if (this.loginType === 'user') {
        userLogin(loginData).then(res=> { // 调用用户登录API
          if (res.msg === 'Login_Success') { // 登录成功
            console.log("登录成功, 跳转主页")

            window.localStorage.setItem('user_id', user_id);
            window.localStorage.setItem('user_role', user_role);


            this.$router.push({ path: '/user' }); // 为普通用户跳转到用户主页

          } else if(res.msg === 'Login_Failed'){
            console.log("登录失败, 重新登录")
            Message({
              message: res.data,
              type: 'warning',
              duration: 4000 // 设置消息显示时间
            });

            setTimeout(() => {
              this.id=""
              this.password=""
            }, 4000);
          }
        }).catch(error => {
          console.log('登录错误:', error);
        });

      }else if (this.loginType === 'admin') {

        adminLogin(loginData).then(res=> { // 调用管理员登录API
              if (res.msg === 'Login_Success') { // 登录成功

                window.localStorage.setItem('user_id', user_id);
                window.localStorage.setItem('user_role', user_role);

                this.$router.push({ path: '/admin' }); // 为管理员跳转到管理员主页
              }else {
                // 这里可以处理登录失败的逻辑，比如显示错误消息
                console.error('登录失败:', res.data);
              }
        }).catch(error => {
          console.log('登录错误:', error);
        });

      }


    },

    GoToRegister(){
      this.$router.push('/register');
    },

  }
}
</script>

<style lang="scss" scoped>
$bg: #2d3a4b;
$dark_gray: #889aa4;
$light_gray: #eee;
$blue:#007bff;

.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: $bg;
}

.login-form {
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
  }

  .form-actions {
    display: flex;
    justify-content: space-between;

    .btn-login, .btn-register {
      padding: 10px 20px;
      border: none;
      border-radius: 4px;
      cursor: pointer;
      color: #fff;
    }

    .btn-login {
      background-color: $blue;
    }

    .btn-register {
      background-color: $light_gray;
      color: $bg;
    }
  }
}

.forgot-password {
  text-align: right;
  margin-top: 20px;

  a {
    color: $dark_gray; // 可以调整为您希望的颜色
    text-decoration: underline; // 添加下划线以表示可点击
    &:hover {
      color: darken(#007bff, 10%); // 鼠标悬停时的颜色
    }
  }
}
</style>
