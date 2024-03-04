<template>
  <div>
    <div class="toolbar">
      <div class="toolbar-right">
        <el-button type="primary" icon="el-icon-info" @click="fresh_check_certificates" plain>刷新</el-button>
        <el-button type="primary" icon="el-icon-circle-plus" @click="showApplyDialog" plain>证书上链申请</el-button>
      </div>
    </div>
    <!-- 表格部分 -->
    <el-table :data="tableData" style="width: 100%; margin-left: 10px; margin-right: 10px;" :highlight-current-row="true">
      <el-table-column prop="applicationID" label="申请编号"></el-table-column>
      <el-table-column prop="certType" label="证书类型"></el-table-column>
      <el-table-column prop="uploaderID" label="申请者ID"></el-table-column>
      <el-table-column prop="uploaderName" label="申请者姓名"></el-table-column>
      <el-table-column prop="createdAt" label="申请日期"></el-table-column>

      <el-table-column label="审核状态">
        <template v-slot="scope">
          <span v-if="scope.row.isProcessed && scope.row.isApproved">审核通过</span>
          <span v-else-if="scope.row.isProcessed && !scope.row.isApproved">审核失败</span>
          <span v-else>待审核</span>
        </template>
      </el-table-column>


      <el-table-column label="操作">
        <template v-slot="scope">
          <el-button v-if="!scope.row.isProcessed" type="primary" icon="el-icon-remove"  @click="cancelApplication(scope.row)">撤销申请</el-button>
          <el-button v-else type="success" icon="el-icon-error" @click="deleteRecord(scope.row)">删除记录</el-button>
          <el-button type="info" icon="el-icon-info" @click="seeInfoDetail(scope.row.applicationID)" round>详细</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 申请对话框 -->
    <el-dialog :visible.sync="dialogVisible" title="申请证书" @close="handleCloseDialog">
      <el-form ref="applyForm" label-width="100px">
        <el-form-item label="证书类型">
          <el-select v-model="cetType" placeholder="请选择证书类型">
            <el-option v-for="type in CET_TYPE" :key="type" :label="type" :value="type"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="上传证书">
          <el-upload class="upload-demo" :file-list="fileList" :auto-upload="false" :on-change="handleChange" action="/fake_upload_endpoint" list-type="text" :limit="1">
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

    <!-- 详细信息对话框 -->
    <el-dialog :visible.sync="detailDialogVisible" title="详细信息" width="50%">
      <el-row v-for="(value, key) in singleDetailInfo" :key="key" class="detail-row">
        <el-col :span="8">{{ key }}</el-col>
        <el-col :span="16">
          <span v-if="key === 'imageURL' && !value">暂无</span>
          <span v-else-if="key === 'imageURL' && value">
            <img :src="value" alt="证书图片" class="detail-image" @dblclick="openImageInNewTab(value)">
          </span>
          <span v-else>{{ value }}</span>
        </el-col>
      </el-row>
    </el-dialog>
  </div>
</template>

<script>
import {
  user_apply_certificate,
  user_query_checking_certificates,
  user_revoke_or_delete_apply_certificates, user_single_cert_certificates
} from "@/api/userCert";
import { getReviewTypes } from "@/api/admin";
import { Message , MessageBox } from 'element-ui';

export default {
  name: 'My_Certificates',
  data() {
    return {
      tableData: [],
      dialogVisible: false,
      CET_TYPE: [],
      cetType: "", // 选中的证书类型
      fileList: [],
      singleDetailInfo:{},
      detailDialogVisible: false, // 详细信息对话框可见性
    }
  },
  created() {
    const storedCETType = JSON.parse(window.localStorage.getItem('CET_TYPE'));
    if (storedCETType != null) {
      this.CET_TYPE = storedCETType;
    } else {
      this.load_CertTypes();
    }
    this.fresh_check_certificates()
  },
  watch: {
    $route: {
      immediate: true
    }
  },
  methods: {
    load_CertTypes() {
      getReviewTypes()
          .then(res => {
            console.log("apply/index(): 获取当前机构的所有证书类型: ")
            console.log(res);
            //存储到本地
            window.localStorage.setItem('CET_TYPE', JSON.stringify(res.data));
            this.CET_TYPE = res.data;
          }).catch(error => {
        console.log('加载证书类型错误:', error);
      });
    },
    openImageInNewTab(imageUrl) {
      window.open(imageUrl, '_blank');
    },

    // 展示证书审核状态
    fresh_check_certificates() {
      const user_id = window.localStorage.getItem('user_id');
      const data = {
        "queryType": "user",
        "userID": user_id,
        "adminID": 0
      }
      user_query_checking_certificates(data).then(res => {

        console.log("apply/index(): 获取当前用户的待审核证书: ")
        console.log(res);
        this.tableData = res.data;
      }).catch(error => {
        console.log('获取待审核证书错误:', error);
      })
    },
    showApplyDialog() {
      this.dialogVisible = true;
    },
    handleChange(file, fileList) {
      this.fileList = fileList;
    },
    confirmApply() {
      const user_id = window.localStorage.getItem('user_id');
      const user_name = window.localStorage.getItem('user_name'); // Assuming user_name is constant for now
      const certType = this.cetType; // 使用选中的证书类型
      const file = this.fileList.length > 0 ? this.fileList[0].raw : null;

      console.log(user_id, user_name, certType)

      if (!certType || !file) {
        // 检查是否有未填写的信息或者未上传的文件
        this.$message.error('请检查选择的证书类型和是否选择了上传的文件');
        return;
      }

      console.log('确认申请');
      this.dialogVisible = false;

      user_apply_certificate(user_id, user_name, certType, file)
          .then(response => {

            console.log(response);

            Message({
              message: "证书上链申请成功",
              type: "success",
              duration: 5 * 1000
            });
            this.fresh_check_certificates();
          }).catch(error => {
        console.error('申请证书错误:', error);
      });

      // 申请成功后清空表单数据
      this.clearFormData();
      //刷新待审核列表
      this.fresh_check_certificates();
    },
    handleRemove(file, fileList) {
      this.fileList = fileList;
    },
    clearFormData() {
      this.$refs.applyForm.resetFields();
      this.cetType=""
      this.fileList = [];
    },
    handleCloseDialog() {
      // 关闭对话框时重置表单数据
      this.clearFormData();
    },
    cancelApplication(row) {
      // 弹出确认对话框
      MessageBox.confirm('确定要撤销申请吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        // 用户点击了确定按钮
        console.log('撤销申请', row);
        const data={
          "applicationId": row.applicationID
        }
        console.log(data);
        user_revoke_or_delete_apply_certificates(data)
            .then(res=>{
              console.log("apply/index(): 撤销申请成功: ")
              console.log(res);
              Message({
                message: "撤销申请成功",
                type: "success",
                duration: 5 * 1000
              });
              this.fresh_check_certificates();
            })
      }).catch(() => {
        // 用户点击了取消按钮或者点击了对话框的关闭按钮
        console.log('取消撤销申请');
      });
    },
    deleteRecord(row) {
      // 弹出确认对话框
      MessageBox.confirm('确定要删除记录吗？', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        // 用户点击了确定按钮
        console.log('删除记录', row);
        const data={
          "applicationId": row.applicationID
        }
        console.log(data);
        user_revoke_or_delete_apply_certificates(data)
            .then(res=>{
              console.log("apply/index(): 删除记录功: ")
              console.log(res);
              Message({
                message: "记录删除成功",
                type: "success",
                duration: 5 * 1000
              });
              this.fresh_check_certificates();
            })
      }).catch(() => {
        // 用户点击了取消按钮或者点击了对话框的关闭按钮
        console.log('取消删除记录');
      });
    },
    seeInfoDetail(applicationID){
      console.log("seeInfoDetail(): ",applicationID)

      const data={
        "applicationId": applicationID
      }

      user_single_cert_certificates(data)
          .then(res=> {
            console.log("apply/index(): 获取申请详情: ")
            console.log(res.data)
            // 将 imageURL 调整为最后一项
            const { imageURL, ...rest } = res.data;
            const adjustedData = { ...rest, imageURL }
            //TODO测试
            // 将 localhost 替换为 10.201.102.119
            // adjustedData.imageURL = adjustedData.imageURL.replace('localhost', '10.201.102.119');

            this.singleDetailInfo = adjustedData;
            this.detailDialogVisible = true; // 打开详细信息对话框

          }).catch(() => {
        console.log('seeInfoDetail失败');
      });



    },
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

.detail-row {
  margin-bottom: 10px; /* 设置行之间的间距 */
}

.detail-image {
  max-width: 100%; /* 图片最大宽度为100% */
  cursor: pointer; /* 设置鼠标样式为指针，表示可点击 */
}
</style>
