// 导入 Vue Router 实例和 Vuex Store 实例，用于路由管理和状态管理。
import router from './router'
import store from './store'

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
const whiteList = ['/login'] // no redirect whitelist


// 注册全局前置守卫，在路由导航之前执行相关操作。
router.beforeEach(async (to, from, next) => {
  // start progress bar
  NProgress.start()

  // set page title
  document.title = getPageTitle(to.meta.title)

  // determine whether the user has logged in
  // 判断用户是否已登录，根据用户登录状态执行不同的操作：
  const hasToken = getToken()


  // 如果已登录：
  //      如果要前往的页面是登录页面，直接重定向到首页。
  //      否则，判断用户是否已获取权限角色信息：
  //          如果已获取角色信息，继续路由导航。
  //          否则，调用 account/getInfo 获取用户信息，并根据角色生成可访问的路由表。
  // 如果未登录：
  //       如果要前往的页面在白名单中，直接继续路由导航。
  //       否则，重定向到登录页面，并携带当前页面路径作为重定向参数。

  if (hasToken) {
    if (to.path === '/login') {
      // if is logged in, redirect to the home page
      next({
        path: '/'
      })
      NProgress.done()
    } else {
      // determine whether the user has obtained his permission roles through getInfo
      const hasRoles = store.getters.roles && store.getters.roles.length > 0
      if (hasRoles) {
        next()
      } else {
        try {
          // get user info
          // note: roles must be a object array! such as: ['admin'] or ,['developer','editor']
          var roles = await store.dispatch('account/getInfo')

          // generate accessible routes map based on roles
          const accessRoutes = await store.dispatch('permission/generateRoutes', roles)

          // dynamically add accessible routes
          router.addRoutes(accessRoutes)

          // hack method to ensure that addRoutes is complete
          // set the replace: true, so the navigation will not leave a history record
          next({
            ...to,
            replace: true
          })
        } catch (error) {
          // remove token and go to login page to re-login
          await store.dispatch('account/resetToken')
          Message.error(error || 'Has Error')
          next(`/login?redirect=${to.path}`)
          NProgress.done()
        }
      }
    }
  } else {
    /* has no token*/

    if (whiteList.indexOf(to.path) !== -1) {
      // in the free login whitelist, go directly
      next()
    } else {
      // other pages that do not have permission to access are redirected to the login page.
      next(`/login?redirect=${to.path}`)
      NProgress.done()
    }
  }
})

router.afterEach(() => {
  // finish progress bar
  NProgress.done()
})
