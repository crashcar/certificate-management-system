<template>
  <div>
    <div class="toolbar">
      <div class="toolbar-right">
        <el-button type="primary" @click="fresh_check_certificates" plain>刷新</el-button>
      </div>
    </div>
    <!-- 表格部分 -->
    <el-table
        :data="tableData"
        style="width: 100%; margin-left: 10px; margin-right: 10px;"
        :highlight-current-row="true">
      <el-table-column
          prop="id"
          label="证书编号">
      </el-table-column>
      <el-table-column
          prop="uploaderId"
          label="申请者ID">
      </el-table-column>
      <el-table-column
          prop="uploaderName"
          label="申请者姓名">
      </el-table-column>
      <el-table-column
          prop="createdAt"
          label="申请日期">
      </el-table-column>
      <el-table-column
          label="操作">
        <template v-slot="scope">
          <el-button type="primary" icon="el-icon-edit" @click="showReviewDialog(scope.row.id)">审查</el-button>
          <el-button type="success" icon="el-icon-download" @click="downloadCertificate(scope.row.id)" round>下载</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 审查对话框 -->
    <el-dialog :visible.sync="reviewDialogVisible" title="审查证书">
      <el-radio-group v-model="reviewResult">
        <el-radio label="pass">通过</el-radio>
        <el-radio label="reject">拒绝</el-radio>
      </el-radio-group>
      <el-button type="primary" @click="submitReview">提交</el-button>
    </el-dialog>
  </div>
</template>

<script>
import { query_institute_showCertList } from "@/api/adminCert";

export default {
  name: 'Institute_Certificates',
  data() {
    return {
      tableData: [],
      reviewDialogVisible: false,
      reviewCertificateId: null,
      reviewResult: ''
    }
  },
  watch: {
    $route: {
      immediate: true
    }
  },
  created() {
    // 页面创建后立即执行查询方法
    this.fresh_check_certificates();
  },
  methods: {
    fresh_check_certificates() {
      console.log('fresh_check_certificates: ');
      // 这里可以添加查询证书的逻辑
      const user_id = parseInt(window.localStorage.getItem('user_id'))
      const queryCertData = {
        "adminID": user_id
      };
      // 发送查询请求
      query_institute_showCertList(queryCertData).then(res => {
        console.log("fresh_check_certificates(): POST 请求");
        console.log(res)
        // 将返回的证书数据赋值给表格数据
        if (res.data != null) {
          console.log("fresh_check_certificates(): 本机构存在待审核的证书")
          this.tableData = res.data;
        } else {
          console.log("fresh_check_certificates(): 本机构暂不存在待审核的证书")
        }
      }).catch(error => {
        console.log('查询机构待审核证书信息错误:', error);
      });
    },
    downloadCertificate(certificateId) {
      // 下载证书的逻辑
    },
    showReviewDialog(certificateId) {
      this.reviewCertificateId = certificateId;
      this.reviewDialogVisible = true;
    },
    submitReview() {
      // 提交审查结果的逻辑
      console.log('审查结果:', this.reviewResult);
      // 关闭审查对话框
      this.reviewDialogVisible = false;
    }
  }
}
</script>

<style scoped>
.toolbar {
  display: flex;
  position: relative;
  justify-content: flex-end;
}

.toolbar-right {
  text-align: right;
  margin: 5px 20px 5px 30px;
}

.el-button-group .el-button {
  margin-right: 10px; /* 添加右侧间距 */
}
</style>
