<template>
  <div>
    <div class="toolbar">
      <div class="toolbar-right">
        <el-button type="primary" @click="editInfoDialogVisible = true">修改信息</el-button>
        <el-button type="primary" @click="editPasswordDialogVisible = true">修改密码</el-button>
      </div>
    </div>

    <div class="profile">
      <el-card>
        <h2 slot="header">个人信息</h2>
        <el-row>
          <el-col :span="12">
            <el-form label-position="left" label-width="100px">
              <el-form-item label="用户ID">{{ userInfo.id }}</el-form-item>
              <el-form-item label="用户姓名">{{ userInfo.realname}}</el-form-item>
              <el-form-item label="手机号码">{{ userInfo.phone }}</el-form-item>
              <el-form-item label="邮箱">{{ userInfo.email }}</el-form-item>
            </el-form>
          </el-col>
        </el-row>
        <div slot="footer">
          <el-button type="primary" @click="editInfoDialogVisible = true">修改信息</el-button>
          <el-button type="primary" @click="editPasswordDialogVisible = true">修改密码</el-button>
        </div>
      </el-card>

      <!-- 修改信息对话框 -->
      <el-dialog :visible.sync="editInfoDialogVisible" title="修改信息" width="30%">
        <el-form :model="formData" label-width="80px">
          <el-form-item label="手机号码" prop="phone">
            <el-input v-model="formData.phone" placeholder="请输入新手机号码"></el-input>
          </el-form-item>
          <el-form-item label="邮箱" prop="email">
            <el-input v-model="formData.email" placeholder="请输入新邮箱"></el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="updateInfo">修改信息</el-button>
          </el-form-item>
        </el-form>
      </el-dialog>

      <!-- 修改密码对话框 -->
      <el-dialog :visible.sync="editPasswordDialogVisible" title="修改密码" width="30%">
        <el-form :model="passwordData" label-width="80px">
          <el-form-item label="旧密码" prop="oldPassword">
            <el-input v-model="passwordData.oldPassword" type="password" placeholder="请输入旧密码"></el-input>
          </el-form-item>
          <el-form-item label="新密码" prop="newPassword">
            <el-input v-model="passwordData.newPassword" type="password" placeholder="请输入新密码"></el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="updatePassword">修改密码</el-button>
          </el-form-item>
        </el-form>
      </el-dialog>


    </div>

  </div>
</template>

<script>
export default {
  data() {
    return {
      userInfo: {
        id: "", // 示例用户ID
        realname: '', // 示例姓名
        phone: '1234567890', // 示例手机号码
        email: 'example@example.com' // 示例邮箱
      },
      formData: {
        phone: '', // 用于存放修改后的手机号码
        email: '' // 用于存放修改后的邮箱
      },
      passwordData: {
        oldPassword: '', // 用于存放旧密码
        newPassword: '' // 用于存放新密码
      },
      editInfoDialogVisible: false, // 控制修改信息对话框的显示与隐藏
      editPasswordDialogVisible: false // 控制修改密码对话框的显示与隐藏
    };
  },
  created() {
    this.userInfo.id =window.localStorage.getItem("user_id")
    this.userInfo.realname =window.localStorage.getItem("user_name")
  },


  methods: {
    updateInfo() {
      // 这里实现更新用户信息的逻辑，可调用后端接口进行更新
      // 示例代码：this.formData.phone 和 this.formData.email 分别是修改后的手机号码和邮箱
      console.log('更新信息', this.formData);
      // 关闭对话框
      this.editInfoDialogVisible = false;
    },
    updatePassword() {
      // 这里实现更新密码的逻辑，可调用后端接口进行更新
      // 示例代码：this.passwordData.oldPassword 和 this.passwordData.newPassword 分别是旧密码和新密码
      console.log('更新密码', this.passwordData);
      // 关闭对话框
      this.editPasswordDialogVisible = false;
    }
  }
};
</script>

<style scoped>
.profile {
  margin: 5px;
}

.toolbar {
  display: flex;
  position: relative;
  justify-content: flex-end;
}

.toolbar-right {
  text-align: right;
  margin: 5px 5px 5px 5px;
}
</style>
