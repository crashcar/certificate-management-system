// 导入 Vue Router 实例和 Vuex Store 实例，用于路由管理和状态管理。
import router from './router'
// import store from './store'

// 导入 Element UI 的 Message 组件和 NProgress 进度条库，用于显示提示信息和页面加载进度条。
import { Message } from 'element-ui'
import NProgress from 'nprogress' // progress bar
import 'nprogress/nprogress.css' // progress bar style

// 导入自定义的工具函数，用于从 Cookie 中获取用户 token 和设置页面标题。
import { getToken } from '@/utils/auth' // get token from cookie
import getPageTitle from '@/utils/get-page-title'

// 配置 NProgress，设置是否显示加载进度条中的旋转图标。
NProgress.configure({ showSpinner: false }) // NProgress Configuration

// 定义一个白名单数组，表示不需要登录即可访问的页面路径。
const whiteList = ['/login',"/register","/prove","/404"] // no redirect whitelist


// 注册全局前置守卫，在路由导航之前执行相关操作。
router.beforeEach(async (to, from, next) => {
  // 开始进度条...
  NProgress.start()

  // 设置页面标题...
  document.title = getPageTitle(to.meta.title)

  const hasToken = window.localStorage.getItem('user_id')
  const role = window.localStorage.getItem('user_role')

  // 已登录状态的处理逻辑
  if (hasToken!==null && role!==null) {
    console.log("已登录: {user_id: "+hasToken+ " , user_role: "+ role+ " }")
    if (to.path === '/login') {
      // 如果已登录且在访问登录页面，则根据角色重定向
      next({ path: '/' + role })
      NProgress.done()
    } else {
      // 对于其他路径，可以添加更多的角色检查和重定向逻辑
      next() // 直接放行
      NProgress.done()
    }
  } else {
    // 未登录状态的处理逻辑
    console.log("permision.js: user/admin未登录")
    if (whiteList.indexOf(to.path) !== -1) {
      console.log("permision.js: 白名单路径放行: ", to.path)
      next() // 白名单路径放行
      NProgress.done()
    } else {
      next(`/login`) // 重定向到登录页
      NProgress.done()
    }
  }
})

router.afterEach(() => {
  // finish progress bar
  NProgress.done()
})
