<template>
  <div>
    <div class="toolbar">
      <div class="toolbar-right">
        <el-button type="primary" @click="editPasswordDialogVisible = true">修改密码</el-button>
      </div>
    </div>

    <div class="profile">
      <el-card>
        <h2 slot="header">管理员信息</h2>
        <el-row>
          <el-col :span="12">
            <el-form label-position="left" label-width="100px">
              <el-form-item label="管理员ID">{{ userInfo.id }}</el-form-item>
              <el-form-item label="所属机构">{{ userInfo.institute }}</el-form-item>
            </el-form>
          </el-col>
        </el-row>
        <div slot="footer">

          <el-button type="primary" @click="editPasswordDialogVisible = true">修改密码</el-button>
        </div>
      </el-card>


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
        id: 0, // 示例用户ID
        institute: "CET"
      },
      passwordData: {
        oldPassword: '', // 用于存放旧密码
        newPassword: '' // 用于存放新密码
      },
      editPasswordDialogVisible: false // 控制修改密码对话框的显示与隐藏
    };
  },
  created() {
    this.userInfo.id = parseInt(window.localStorage.getItem('user_id'))
  },

  methods: {
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
