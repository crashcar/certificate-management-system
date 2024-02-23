<template>
  <div class="register-container">
    <div class="register-form">
      <div class="title-container">
        <h3 class="title">注册新用户</h3>
      </div>
      <form @submit.prevent="handleRegister">
        <div class="form-group">
          <input type="text" class="form-control" v-model="id" placeholder="请输入您的用户id">
        </div>
        <div class="form-group">
          <input type="password" class="form-control" v-model="password" placeholder="请输入您的密码">
        </div>
        <div class="form-group">
          <input type="text" class="form-control" v-model="phone" placeholder="请输入您的电话号码">
        </div>
        <div class="form-group">
          <input type="email" class="form-control" v-model="email" placeholder="请输入您的邮箱">
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
import { userRegister } from '@/api/userInfo'

export default {
  name: 'Register',
  data() {
    return {
      id: '',
      password: '',
      phone: '',
      email: '',
    }
  },
  methods: {
    handleRegister() {
      const registerData = {
        id: this.id,
        password: this.password,
        phone: this.phone,
        email: this.email,
      };
      console.log(registerData)
      userRegister(registerData).then(res => {

        console.log(res)
        if (res.data === 'Register_Success') {
          this.$router.push('/login');
        } else {
          // 处理注册失败逻辑
          console.error('注册失败:', res.data);
        }
      }).catch(error => {
        console.error('注册错误:', error);
      });
    },

    GoToLogin(){
      this.$router.push('/login');
    },
  }
}
</script>

<style lang="scss" scoped>
$bg: #2d3a4b; // 背景颜色
$dark_gray: #889aa4; // 深灰色，用于按钮等
$light_gray: #eee; // 浅灰色，用于另一个按钮
$text_color: #5a5a5a; // 文本颜色
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

    input {
      width: 100%;
      padding: 10px;
      border: 1px solid $dark_gray;
      border-radius: 4px;
      box-sizing: border-box; // 确保内边距不影响最终尺寸
      color: $text_color; // 输入框内的文本颜色

      // 输入框获得焦点时的样式
      &:focus {
        border-color: lighten($dark_gray, 10%);
        outline: none; // 移除焦点边框
      }
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

// 响应式设计调整，适应更小的屏幕
@media (max-width: 500px) {
  .register-form {
    padding: 20px;
    max-width: 90%;
  }
}

</style>