<template>
  <div>
    <div class="toolbar">
      <div class="toolbar-right">
        <el-button type="primary" @click="showApplyDialog" plain>申请</el-button>
      </div>
    </div>
    <!-- 表格部分 -->
    <el-table
        :data="tableData"
        style="width: 100%; margin-top: 20px;"
        :highlight-current-row="true">
      <el-table-column
          prop="name"
          label="证书名称">
      </el-table-column>
      <el-table-column
          prop="date"
          label="获得日期">
      </el-table-column>
      <el-table-column
          prop="issuer"
          label="颁发机构">
      </el-table-column>
      <el-table-column
          label="操作">
        <template v-slot="scope">
          <el-button type="text" @click="download_certificate(scope.row.id)">下载</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 申请对话框 -->
    <el-dialog :visible.sync="dialogVisible" title="申请证书" @close="handleCloseDialog">
      <el-form ref="applyForm" label-width="100px">
        <el-form-item label="证书颁发机构编号">
          <el-input v-model="institutionId"></el-input>
        </el-form-item>
        <el-form-item label="证书编号">
          <el-input v-model="certificateId"></el-input>
        </el-form-item>
        <el-form-item label="上传证书">
          <el-upload
              class="upload-demo"
              :file-list="fileList"
              :auto-upload="false"
              :on-change="handleChange"
              action="/fake_upload_endpoint"
              list-type="text"
              :limit="1">
            <el-button size="small" type="primary">点击上传</el-button>
            <template v-slot="{ file, fileList }">
              <div class="uploaded-file-name">{{ file.name }}</div>
              <el-button type="text" size="small" @click="handleRemove(file, fileList)">删除</el-button>
            </template>
          </el-upload>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="confirmApply">确认申请</el-button>
        </el-form-item>
      </el-form>
    </el-dialog>
  </div>
</template>

<script>
import { user_apply_certificate } from "@/api/userCert";

export default {
  name: 'My_Certificates',
  data() {
    return {
      tableData: [],
      dialogVisible: false,
      institutionId: '',
      certificateId: '',
      fileList: []
    }
  },
  watch: {
    $route: {
      immediate: true
    }
  },
  methods: {
    showApplyDialog() {
      this.dialogVisible = true;
    },
    download_certificate() {
    },
    handleChange(file, fileList) {
      this.fileList = fileList;
    },
    confirmApply() {
      const institutionId = this.institutionId;
      const certificateId = this.certificateId;
      const file = this.fileList.length > 0 ? this.fileList[0].raw : null;

      console.log(institutionId)
      console.log(certificateId)
      console.log(file)

      if (!institutionId || !certificateId || !file) {
        // 检查是否有未填写的信息或者未上传的文件
        this.$message.error('请填写完整信息并上传文件');
        return;
      }

      console.log('确认申请');
      this.dialogVisible = false;

      user_apply_certificate(institutionId, certificateId, file)
          .then(response => {
            console.log(response);
          })
          .catch(error => {
            console.error('申请证书错误:', error);
          });

      // 申请成功后清空表单数据
      this.clearFormData();
    },
    handleRemove(file, fileList) {
      this.fileList = fileList;
    },
    clearFormData() {
      this.$refs.applyForm.resetFields();
      this.institutionId = '';
      this.certificateId = '';
      this.fileList = [];
    },
    handleCloseDialog() {
      // 关闭对话框时重置表单数据
      this.clearFormData();
    }
  }
}
</script>

<style scoped>
.toolbar {
  display: flex;
  position: relative;
  justify-content: flex-end; /* 修改为右侧对齐 */
}

.toolbar-right {
  text-align: right; /* 调整文本对齐 */
  margin: 5px 20px 5px 30px; /* 调整间距 */
}

.upload-demo {
  display: flex;
}

.uploaded-file-name {
  margin-right: 10px; /* 设置文件名和删除按钮之间的距离 */
}
</style>
