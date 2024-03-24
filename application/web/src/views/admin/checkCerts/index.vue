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
          prop="applicationID"
          label="证书编号">
      </el-table-column>
      <el-table-column
          prop="certType"
          label="证书类型">
      </el-table-column>
      <el-table-column
          prop="uploaderID"
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
          <el-button v-if="!scope.row.isProcessed" type="primary" icon="el-icon-s-check" @click="showReviewDialog(scope.row.applicationID, scope.row.certType)">审查</el-button>
          <el-button type="info" icon="el-icon-info" @click="seeInfoDetail(scope.row.applicationID)" round>详细</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 审查对话框 -->
    <el-dialog :visible.sync="reviewDialogVisible" title="审查证书">
      <el-radio-group v-model="reviewResult">
        <el-radio label="pass">通过</el-radio>
        <el-radio label="reject">拒绝</el-radio>
      </el-radio-group>
      <div v-if="reviewResult === 'reject'" style="margin-top: 10px;">
        <el-input v-model="rejectReason" placeholder="请输入拒绝理由"></el-input>
      </div>
      <div style="margin-top: 20px;">
        <el-button type="primary" @click="submitReview">提交</el-button>
      </div>
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
import {check_pass_single_Cert, check_reject_single_Cert, query_institute_showCertList} from "@/api/adminCert";
import {Message} from "element-ui";
import {user_single_cert_certificates} from "@/api/userCert";

export default {
  name: 'Institute_Certificates',
  data() {
    return {
      tableData: [],
      reviewDialogVisible: false,
      reviewCertificateId: -1,
      reviewResult: 'reject', // 默认选择拒绝
      rejectReason: '', // 拒绝理由
      issuingAuthority: '',
      singleDetailInfo: {},
      detailDialogVisible: false, // 详细信息对话框可见性
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
        "adminID": user_id,
        "queryType": "admin",
        "userID": "",
      };
      // 发送查询请求
      query_institute_showCertList(queryCertData).then(res => {
        console.log("fresh_check_certificates(): POST 请求");
        console.log(res)
        this.tableData = res.data.filter(item => !item.isProcessed);
      }).catch(error => {
        console.log('查询机构待审核证书信息错误:', error);
      });
    },
    seeInfoDetail(applicationID) {
      console.log("seeInfoDetail(): ", applicationID);

      const data = {
        "applicationId": applicationID
      };

      user_single_cert_certificates(data)
          .then(res => {
            console.log("apply/index(): 获取申请详情: ", res.data);

            // 处理获取的申请详情数据
            const { imageURL, ...rest } = res.data;
            const adjustedData = { ...rest, imageURL };
                        // 将 localhost 替换为 10.201.102.119
            // adjustedData.imageURL = adjustedData.imageURL.replace('localhost', '10.201.102.119');

            console.log("imageURL: "+imageURL)
            // 更新详情数据并打开详细信息对话框
            this.singleDetailInfo = adjustedData;
            this.detailDialogVisible = true;
          })
          .catch(() => {
            console.log('seeInfoDetail失败');
          });
    },
    showReviewDialog(certificateId, certType) {
      this.reviewCertificateId = certificateId;
      this.reviewDialogVisible = true;
      this.issuingAuthority=certType;
    },
    submitReview() {
      // 提交审查结果的逻辑
      const user_id=parseInt(window.localStorage.getItem('user_id'))

      // 如果是拒绝，则打印拒绝理由
      if (this.reviewResult === 'pass') {
        console.log('user_id:',user_id)
        console.log('审查结果:', this.reviewResult);
        console.log('审查的证书ID:', this.reviewCertificateId);
        console.log('审查机构',this.issuingAuthority)

        const passData={
          "certDBID": this.reviewCertificateId,
          "adminID": user_id,
          "issuingAuthority": this.issuingAuthority,
        }
        console.log( passData);
        check_pass_single_Cert(passData).then(res => {
          console.log("通过审核证书成功:", res);
          Message({
            message: "审核成功",
            type: 'warning',
            duration: 4000 // 设置消息显示时间
          });
          this.fresh_check_certificates();

        }).catch(error => {
          console.log('通过审核证书错误:', error);
        });

      }else if(this.reviewResult === 'reject'){

        console.log('拒绝理由:', this.rejectReason);

        if (this.rejectReason === '') {
          Message({
            message: "请填写拒绝审核理由",
            type: 'warning',
            duration: 4000 // 设置消息显示时间
          });
          return
        }

        const rejectData={
          "certDBID": this.reviewCertificateId,
          "adminID": user_id,
          "DenialReason": this.rejectReason
        }
        check_reject_single_Cert(rejectData).then(res => {

          console.log("拒绝审核证书成功:", res);
          Message({
            message: "审核成功",
            type: 'warning',
            duration: 4000 // 设置消息显示时间
          });

          this.fresh_check_certificates();

        }).catch(error => {
          console.log('通过审核证书错误:', error);
        });

      }

      this.reviewDialogVisible = false;
      this.clear_check_info();
    },

    clear_check_info(){
      this.rejectReason="";
      this.issuingAuthority="";
      this.reviewResult="reject";
      this.reviewCertificateId=-1;
    },

    openImageInNewTab(imageURL) {
      window.open(imageURL, '_blank');
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

.detail-row {
  margin-bottom: 10px; /* 设置行之间的间距 */
}

.detail-image {
  max-width: 100%; /* 图片最大宽度为100% */
}
</style>
