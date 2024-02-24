// 导入 Vue 库，以便创建 Vue 实例
import Vue from 'vue'

// 导入 normalize.css 库，用于提供现代化的 CSS 重置样式。
import 'normalize.css/normalize.css' // A modern alternative to CSS resets

// 导入 Element UI 库，用于提供基于 Vue 的 UI 组件。同时也可以导入对应的样式文件。
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
// import locale from 'element-ui/lib/locale/lang/en' // lang i18n

// 导入全局的 SCSS 样式文件，用于定义项目的全局样式。
import '@/styles/index.scss' // global css

// 导入根组件 App.vue。
import App from './App'
// 分别导入 Vuex 状态管理库和 Vue Router 路由配置。
// import store from './store'
import router from './router'

// 导入自定义的 icon 图标和权限控制模块。
import '@/icons' // icon
import '@/permission' // permission control

/**
 * If you don't want to use mock-server
 * you want to use MockJs for mock api
 * you can execute: mockXHR()
 *
 * Currently MockJs will be used in the production environment,
 * please remove it before going online ! ! !
 */
// if (process.env.NODE_ENV === 'production') {
//   const { mockXHR } = require('../mock')
//   mockXHR()
// }

// set ElementUI lang to EN
// Vue.use(ElementUI, { locale })
// 如果想要中文版 element-ui，按如下方式声明

// 全局注册 Element UI 插件，使其可以在 Vue 实例中使用 Element UI 的组件。
Vue.use(ElementUI)

// 关闭生产环境的提示信息，避免在控制台输出一些无用的提示。
Vue.config.productionTip = false

// 创建 Vue 实例，传入选项对象：
// el: '#app'：指定 Vue 实例挂载的元素，这里是一个具有 id="app" 的元素，通常是 index.html 中的根元素。
// router：注册 Vue Router 实例。
// store：注册 Vuex Store 实例。
// render: h => h(App)：渲染根组件 App.vue，这里使用了 ES6 的箭头函数语法。
new Vue({
  el: '#app',
  router,
  // store,
  render: h => h(App)
})
