<template>
  <div class="container">
    <img src="@/assets/prove.svg" alt="logo">
    <div class="title">基于联盟链的电子证书验证系统</div>
    <div class="form-group-row">
      <input type="text" class="form-control" v-model="certificateId" placeholder="请输入证书编号">
      <button class="btn-login" @click="openDialog">验证</button>
    </div>
    <el-dialog :visible.sync="dialogVisible" title="申请证书" @close="handleCloseDialog" :width="dialogWidth" class="custom-dialog">
      <el-form ref="applyForm" label-width="100px">
        <el-form-item label="身份ID">
          <el-input v-model="userId" placeholder="请输入身份ID" :style="{ width: 'calc(100% - 20px)' }"></el-input>
        </el-form-item>
        <el-form-item label="用户名">
          <el-input v-model="userName" placeholder="请输入用户名" :style="{ width: 'calc(100% - 20px)' }"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="verify">提交</el-button>
          <el-button @click="closeDialog">关闭</el-button>
        </el-form-item>
      </el-form>
    </el-dialog>
    <div class="cards-section">
      <div class="card-container">
        <div class="card" v-for="(link, index) in links" :key="index" @click="openLink(link.url)" :style="{ backgroundColor: link.backgroundColor, color: link.textColor }">
          <b>{{ link.title }}</b>
        </div>
      </div>
    </div>



    <el-dialog :visible.sync="resultDialogVisible" title="证书详情">
      <el-row v-for="(value, key) in resultData" :key="key">
        <el-col :span="6">{{ key }}</el-col>
        <el-col :span="18">
          <!-- 判断是否为 "authorityContactInfo"，如果是则展示里面的内容 -->
          <template v-if="key === 'authorityContactInfo'">
            <div>{{ value.address }}</div>
            <div>{{ value.email }}</div>
            <div>{{ value.phone }}</div>
          </template>
          <!-- 如果不是 "authorityContactInfo" 则直接展示值 -->
          <template v-else>
            {{ value }}
          </template>
        </el-col>
      </el-row>
    </el-dialog>



  </div>
</template>

<script>
import { prove_certificate } from "@/api/prove";
import { Message } from 'element-ui';

export default {
  data() {
    return {
      certificateId: '',
      userId: '',
      userName: '',
      dialogVisible: false,
      dialogWidth: '30%', // 根据需要调整宽度
      resultDialogVisible: false,
      resultData: {},
      links: [
        { title: 'CET-英语等级证书', url: 'http://localhost:9528/web/#/login', backgroundColor: '#ff7f50', textColor: '#fff' }, // 橙色
        { title: 'NCRE-计算机等级证书', url: 'http://example.com/2', backgroundColor: '#87CEEB', textColor: '#000' }, // 天蓝色
        { title: 'NECT-教师资格证书', url: 'http://example.com/2', backgroundColor: '#3CB371', textColor: '#fff' } // 草绿色
        // 根据需要添加更多网页卡片
      ]
    }
  },
  methods: {
    openDialog() {
      if (this.certificateId === '') {
        this.$message.error("请输入证书编号")
        return
      }
      this.dialogVisible = true;
    },
    handleCloseDialog() {
      this.dialogVisible = false;
    },

    verify() {
      let data = {
        "certID": this.certificateId,
        "holderID": this.userId,
        "holderName": this.userName
      }

      if (this.userId === '' || this.userName === '') {
        this.$message.error('请检查输入的身份ID和用户名');
        this.dialogVisible = false;
        return
      }

      // 模拟数据
      // data = {
      //   "certID": "cet.com-05674414-da22-11ee-8012-0242ac190004",
      //   "holderID": "001",
      //   "holderName": "hua"
      // }

      prove_certificate(data).then(res => {
        if (res.data === null) {
          Message({
            message: "验证失败, 请重新输入验证",
            type: "error",
            duration: 5 * 1000
          });
          this.clear_data()
          return
        }
        console.log(res.data)
        Message({
          message: "证书验证成功",
          type: "success",
          duration: 5 * 1000
        });

        this.resultData = res.data[0];
        this.clear_data()
        this.resultDialogVisible = true;
      })

    },
    openLink(url) {
      window.open(url, '_blank');
    },
    clear_data() {
      this.certificateId = "";
      this.userId = "";
      this.userName = "";
      this.dialogVisible = false;
    }
  }
}
</script>

<style scoped>
.container {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin: 50px;
}

.title {
  font-size: 30px;
  margin-bottom: 20px;
  font-weight: bold;
}

.card {
  margin: 20px; /* 调整卡片间距 */
  padding: 20px;
  border: none;
  cursor: pointer;
  width: 120px; /* 调整卡片大小 */
  height: 120px; /* 调整卡片大小 */
  text-align: center;
  border-radius: 10px;
  display: flex;
  justify-content: center;
  align-items: center;
  font-weight: bold; /* 加粗字体 */
}

.cards-section {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  margin-top: 20px;
}

.card-container {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-around;
  width: 100%;
}

.form-group-row {
  display: flex;
  justify-content: start;
  align-items: center;
  width: 30%;
  margin-bottom: 20px;
}

.form-group-row .form-control {
  flex-grow: 1;
  margin-right: 10px;
  padding: 10px;
  border: 1px solid #ccc;
  border-radius: 4px;
  box-sizing: border-box;
}

.form-group-row .btn-login {
  padding: 10px 20px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  color: #fff;
  background-color: #007bff;
}

.custom-dialog {
  margin-left: 0;
  margin-right: 0;
  border-radius: 10px;
}
</style>
