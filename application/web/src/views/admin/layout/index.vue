<template>
  <el-container class="container">
    <el-aside class="aside" :width="isCollapse ? '64px' : '200px'">
      <!-- 侧边栏头部，包含logo和小组名称 -->
      <div class="aside-header">
        <el-avatar
            v-if="!isCollapse"
            src="https://xzj-pic-1306183757.cos.ap-nanjing.myqcloud.com/picgo/摘月白底.jpg"
            :size="38"
            alt=""
        ></el-avatar>
        <h1 class="aside-title" v-if="!isCollapse">证书系统</h1>
      </div>

      <el-menu
          unique-opened
          text-color="#fff"
          active-text-color="#ffd04b"
          router
          mode="vertical"
          :default-active="defaultActive"
      >
        <el-menu-item index="/admin/instituteCerts">
          <i class="el-icon-document"></i>
          <span slot="title" v-show="!isCollapse">机构证书库</span>
        </el-menu-item>
        <el-menu-item index="/admin/checkCerts">
          <i class="el-icon-edit"></i>
          <span slot="title" v-show="!isCollapse">证书审查</span>
        </el-menu-item>
        <el-menu-item index="/admin/profile">
          <i class="el-icon-user"></i>
          <span slot="title" v-show="!isCollapse">管理员信息页面</span>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="custom-header">
        <div class="left">
          <el-button type="text" class="collapse-btn" @click="toggleCollapse">
            <i class="el-icon-s-unfold icon-size" v-if="isCollapse"></i>
            <i class="el-icon-s-fold icon-size" v-else></i>
          </el-button>
        </div>
        <div class="right">
          <el-dropdown>
            <div class="usr">
              <el-avatar :size="31" :src="userInfo.avatar"></el-avatar>
              <span style="color: white">{{ userInfo.username }}</span>
              <i class="el-icon-arrow-down el-icon--right" style="color: white"></i>
            </div>
            <template>
              <el-dropdown-menu>
                <div style="padding: 0 10px;"> <!-- 添加自定义边距 -->
                  <el-button type="info" @click="logout" style="width: 100%;">退出登录</el-button>
                </div>
              </el-dropdown-menu>
            </template>

          </el-dropdown>
        </div>
      </el-header>
      <el-main>
        <router-view></router-view>
      </el-main>
    </el-container>

  </el-container>
</template>

<script>
export default {
  name: 'Layout',
  data() {
    return {
      isCollapse: false, // 控制侧边栏展开折叠
      defaultActive: '/admin/instituteCerts', // 默认选中的菜单项索引
      userInfo: {
        username: '管理员用户',
        avatar: 'https://xzj-pic-1306183757.cos.ap-nanjing.myqcloud.com/picgo/摘月白底.jpg',
      },
    };
  },
  methods: {
    toggleCollapse() {
      this.isCollapse = !this.isCollapse;
    },
    // 你可以添加一些方法
    logout() {
      // 退出登录逻辑
      console.log('退出登录');

      // 清空localStorage中的所有数据
      localStorage.clear();

      // 重定向到登录页
      this.$router.push('/login');
    }
  },
};
</script>

<style scoped>

.icon-size {
  font-size: 25px;
}
.container {
  display: flex;
  position: absolute;
  width: 100%;
  height: 100%;
}

.aside {
  overflow-y: hidden; /* 隐藏垂直滚动条 */
  background-color: #001529; /* 调整背景颜色 */
}

/* 侧边栏头部样式调整 */
.aside-header {
  display: flex;
  align-items: center;
  justify-content: center; /* 确保内容水平居中 */
  padding: 10px; /* 调整内边距 */
}

.custom-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #007bff; /* New background color */
  height: 50px; /* Adjusted height */
}

.left {
  display: flex;
  align-items: center;
  margin-left: 20px; /* Adjusted margin */
}

.right {
  display: flex;
  align-items: center;
  margin-right: 20px; /* Adjusted margin */
}

.collapse-btn {
  color: #fff; /* Button text color */
  cursor: pointer;
}

.aside-title {
  color: #fff; /* 字体颜色 */
  font-size: 16px; /* 字体大小 */
  white-space: nowrap; /* 确保文本不会换行 */
  overflow: hidden; /* 超出隐藏 */
  text-overflow: ellipsis; /* 显示省略符号 */
  margin: 2px; /* 移除外边距 */
}

.logout-btn {
  color: #fff; /* Button text color */
  border: none; /* Remove border */
  font-size: 14px; /* Adjust font size */
  cursor: pointer;
  padding: 0; /* Adjust padding */
}

.logout-btn:hover {
  text-decoration: underline; /* Underline on hover */
}

.logoCollapse {
  display: flex;
  justify-content: center;
  align-items: center;
  margin: 4px 0px 10px 0;
}

.logoNoCollapse {
  display: flex;
  justify-content: center;
  align-items: center;
  margin: 4px 0px 10px 16px;
}

.logoNoCollapse h1 {
  font-size: 18px; /* 调整字体大小为 18px */
  color: #ffffff; /* 字体颜色设置为白色 */
  white-space: nowrap; /* 保持文本在一行显示 */
  overflow: hidden; /* 超出部分隐藏 */
  text-overflow: ellipsis; /* 超出部分显示省略号 */
  margin: 0; /* 移除外边距 */
  padding: 0 10px; /* 添加内边距，确保文本不紧贴边缘 */
}

.el-main {
  background-color: #f9f9f9;
}

.avatar-container {
  margin-bottom: 20px;
  text-align: center;
}

.user {
  display: flex;
  justify-content: center;
  align-items: center;
  color: #fff;
  font-size: 16px;
}

i.el-icon {
  margin-right: 10px;
}

.user i {
  font-size: 14px;
  cursor: pointer;
}

.el-menu {
  border-right: 0;
  height: 100%;
  background-color: #001529;
}

.el-menu-item {
  background-color: #001529;
}

.el-menu-item:hover,
:deep(.el-menu-item:active) {
  background-color: #0960bd !important;
}

:deep(.el-sub-menu:hover) {
  background-color: #0c2135 !important;
}

:deep(.el-sub-menu__title:hover) {
  background-color: #0c2135 !important;
}

.el-menu-item.is-active {
  background-color: #0960bd !important;
}

.el-menu-vertical-demo:not(.el-menu--collapse) {
  width: 200px;
  min-height: 400px;
}

main.el-main {
  padding: 0px 0px;
}

:deep(header.el-header) {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #fff;
  height: 42px;
  margin-bottom: 15px;
}

header.el-header .left {
  display: flex;
  align-items: center;
}

header.el-header i.el-icon[data-v-8a54e678] {
  cursor: pointer;
}

header.el-header .right {
  display: flex;
  align-items: center;
}

header.el-header .right .usr {
  display: flex;
  align-items: center;
  margin-right: 4px;
}

span.el-avatar.el-avatar--circle {
  margin: 0 6px;
  font-size: 16px;
}
/* 调整头像和标题之间的距离，如果有必要 */
.el-avatar {
  margin-right: 8px;
}

header.el-header .right .usr span {
  margin: 0 6px;
  font-size: 16px;
}

header.el-header i.el-icon {
  color: #6a748b;
  font-size: 21px;
}

header.el-header h2 {
  font-size: 16px;
  font-weight: 500;
  color: #6a748b;
  margin-left: 10px;
}

.routerView {
  padding: 0px 15px;
}

.toggle-button {
  cursor: pointer;
  color: #fff; /* 根据需要调整颜色 */
}

.el-menu-item .el-icon-document {
  width: 16px; /* Adjusted width for the document icon */
}
</style>
