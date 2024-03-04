<template>
  <div>
    <div class="toolbar">
      <div class="toolbar-right">
        <el-button type="primary" @click="fresh_my_certificates" plain>刷新</el-button>
      </div>
    </div>
    <!-- 表格部分 -->
    <el-table
        :data="tableData"
        style="width: 100%; margin-left: 10px; margin-right: 10px;"
        :highlight-current-row="true">
      <el-table-column
          prop="certID"
          label="证书编号">
      </el-table-column>
      <el-table-column
          prop="certType"
          label="证书类型">
      </el-table-column>
      <el-table-column
          prop="holderID"
          label="持有者ID">
      </el-table-column>
      <el-table-column
          prop="holderName"
          label="持有者姓名">
      </el-table-column>
      <el-table-column
          prop="expiryDate"
          label="到期日期">
      </el-table-column>
      <!-- 添加操作列 -->
      <el-table-column
          label="操作">
        <template v-slot="scope">
          <el-button type="info" icon="el-icon-info" @click="showDetails(scope.row)" round>详情</el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 详情对话框 -->
    <el-dialog :visible.sync="detailsDialogVisible" title="证书详情">
      <el-row v-for="(value, key) in selectedCertificate" :key="key">
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
import { query_user_certificates } from '@/api/userCert';

export default {
  name: 'My_Certificates',
  data() {
    return {
      tableData: [],
      detailsDialogVisible: false,
      selectedCertificate: {}
    }
  },
  watch: {
    $route: {
      immediate: true
    }
  },
  created() {
    // 页面创建后立即执行查询方法
    this.fresh_my_certificates();
  },
  methods: {
    fresh_my_certificates() {
      console.log('query_my_certificates');
      // 这里可以添加查询证书的逻辑
      const user_id = window.localStorage.getItem('user_id');
      const queryCertData = {
        holderID: user_id,
        issuingAuthority: "CET",
      };

      // 发送查询请求
      query_user_certificates(queryCertData)
          .then(res => {
            console.log("fresh_my_certificates(): POST 请求");
            console.log(res);
            // 将返回的证书数据赋值给表格数据
            if (res.data != null) {
              console.log("fresh_my_certificates(): 该用户在本机构存在证书");
              this.tableData = res.data;
            } else {
              console.log("fresh_my_certificates(): 该用户在本机构不存在证书");
            }
          })
          .catch(error => {
            console.log('查询用户证书信息错误:', error);
          });
    },
    showDetails(certificate) {
      this.selectedCertificate = certificate;
      this.detailsDialogVisible = true;
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
</style>
