import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

// 路由定义
const routes = [
  {
    path: '/login',
    component: () => import('@/views/login/index'),
    hidden: true
  },
  {
    path: '/register',
    component: () => import('@/views/register/index'),
    hidden: true
  },
  {
    path: '/prove',
    component: () => import('@/views/prove/index'),
    hidden: true
  },
  {
    path: '/user',
    redirect: '/user/certificates',
    component: () => import('@/views/user/layout/index'),
    meta: { role: 'user' },
    children: [
      {
        path: 'certificates',
        component: () => import('@/views/user/certificates/index'),
        name: 'certificates',
        meta: { title: '我的证书' }
      },
      {
        path: 'apply',
        component: () => import('@/views/user/apply/index'),
        name: 'apply',
        meta: { title: '证书申请' }
      },
      {
        path: 'profile',
        component: () => import('@/views/user/profile/index'),
        name: 'Profile',
        meta: { title: '个人信息页面' }
      }
    ]
  },

  {
    path: '/admin',
    redirect: 'admin/instituteCerts',
    component: () => import('@/views/admin/layout/index'),
    meta: { role: 'admin' },
    children: [
        {
          path: 'instituteCerts',
          component: () => import('@/views/admin/instituteCerts/index'),
          name: 'instituteCerts',
          meta: { title: '机构证书管理' }
        },
        {
          path: 'checkCerts',
          component: () => import('@/views/admin/checkCerts/index'),
          name: 'checkCerts',
          meta: { title: '机构证书审查' }
        },
      {
        path: 'profile',
        component: () => import('@/views/admin/profile/index'),
        name: 'Profile',
        meta: { title: '管理员个人信息页面' }
      }
    ]
  },
  {
    path: '/404',
    component: () => import('@/views/404'),
    hidden: true
  },
  { path: '*', redirect: '/404', hidden: true }
]

const router = new Router({
  base: '/web',
  routes
})

router.beforeEach((to, from, next) => {
  const userRole = localStorage.getItem('user_role')
  if(to.path==='/prove'){
    next()
  } else if (!userRole && to.path !== '/login' && to.path !== '/register') {
    next('/login')
  } else if (to.meta && to.meta.role && to.meta.role !== userRole) {
    next('/404')
  } else {
    next()
  }
})


export default router
