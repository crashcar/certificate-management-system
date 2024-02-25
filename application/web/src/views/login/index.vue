<template>
  <div class="login-container">
    <div class="login-form">
      <div class="title-container">
        <h3 class="title">基于联盟链的电子证书管理系统</h3>
      </div>
      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <input type="text" class="form-control" v-model="id" placeholder="请输入您的用户ID">
        </div>
        <div class="form-group">
          <input type="password" class="form-control" v-model="password" placeholder="请输入您的密码">
        </div>
        <div class="form-actions">

          <button type="button" class="btn-register" @click="GoToRegister">注册用户</button>
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
import { userLogin } from '@/api/userInfo'

export default {
  name: 'Login',
  data() {
    return {
      loading: false,

      id: "",
      password:"",
    }
  },
  watch: {
    $route: {
      immediate: true
    }
  },
  created() {
   
  },
  methods: {
    handleLogin() {
      this.loading = true; // 开始加载
      const loginData = {
        id: this.id,
        password: this.password,
      };
      console.log(loginData)
      userLogin(loginData).then(res=> { // 调用登录API
        this.loading = false; // 停止加载
        if (res.msg === 'Login_Success') { // 登录成功
          console.log("登录成功, 跳转主页")

          const user_id=res.data.user_id;
          const user_role=res.data.user_role;

          window.localStorage.setItem('user_id', user_id);
          window.localStorage.setItem('user_role', user_role);

          // 根据角色跳转到相应的页面
          if (user_role === 'admin') {
            this.$router.push({ path: '/admin' }); // 为管理员跳转到管理员主页
          } else if (user_role === 'user') {
            this.$router.push({ path: '/user' }); // 为普通用户跳转到用户主页
          } else {
            // 可以处理其他角色或默认跳转逻辑
            this.$router.push({ path: '/' });
          }
        } else {
          // 这里可以处理登录失败的逻辑，比如显示错误消息
          console.error('登录失败:', res.data);
        }
      }).catch(error => {
        this.loading = false; // 出现错误，停止加载
        console.log('登录错误:', error);
      });
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

    input {
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
