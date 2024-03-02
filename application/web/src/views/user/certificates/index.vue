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
          prop="id"
          label="证书编号">
      </el-table-column>
      <el-table-column
          prop="name"
          label="证书名称">
      </el-table-column>
      <el-table-column
          prop="type"
          label="证书类型">
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
          <el-button type="success" icon="el-icon-download" @click="downloadCertificate(scope.row.id)">下载</el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import {query_user_certificates, user_download_certificate} from '@/api/userCert'
export default {
  name: 'My_Certificates',
  data() {
    return {
      // 表格数据示例，你可以根据实际情况修改或替换
      tableData: []
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

      console.log('query_my_certificates')
      // 这里可以添加查询证书的逻辑
      const user_id=window.localStorage.getItem('user_id')
      const queryCertData = {
        holderID: user_id,
        issuingAuthority: "CET",
      };

      // 发送查询请求
      query_user_certificates(queryCertData).then(res => {
        console.log("fresh_my_certificates(): POST 请求");
        console.log(res)
        // 将返回的证书数据赋值给表格数据
        if (res.data != null) {
          console.log("fresh_my_certificates(): 该用户在本机构存在证书")
          this.tableData = res.data;
        }else{
          console.log("fresh_my_certificates(): 该用户在本机构不存在证书")
        }
      }).catch(error => {
        console.log('查询用户证书信息错误:', error);
      });
    },
    downloadCertificate(certificateId) {
      console.log("certificateId: " + certificateId + "")
      const userId = window.localStorage.getItem('user_id');
      // 发送下载请求
      user_download_certificate(userId, certificateId).then(response => {
        const blob = new Blob([response.data], { type: 'application/pdf' });
        const url = window.URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.setAttribute('download', 'certificate.pdf');
        document.body.appendChild(link);
        link.click();
        window.URL.revokeObjectURL(url);
        document.body.removeChild(link);
      }).catch(error => {
        console.log('downloadCertificate: 下载证书错误:', error);
      });
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
