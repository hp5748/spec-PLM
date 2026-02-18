import router from './index'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

// 白名单路由
const whiteList = ['/login', '/403', '/404']

router.beforeEach(async (to, _from, next) => {
  // 设置页面标题
  document.title = to.meta.title ? `${to.meta.title} - PLM` : 'PLM系统'

  const userStore = useUserStore()
  const token = userStore.token

  if (token) {
    if (to.path === '/login') {
      // 已登录，跳转到首页
      next({ path: '/' })
    } else {
      // 如果有token但没有用户信息，尝试恢复用户状态
      if (!userStore.userInfo) {
        const success = await userStore.initUser()
        if (!success) {
          // 恢复失败，跳转到登录页
          next({ path: '/login', query: { redirect: to.fullPath } })
          return
        }
      }

      // 检查权限
      if (to.meta.permissions && to.meta.permissions.length > 0) {
        const hasPermission = to.meta.permissions.every(permission =>
          userStore.hasPermission(permission)
        )

        if (!hasPermission) {
          ElMessage.error('没有权限访问该页面')
          next({ path: '/403' })
          return
        }
      }

      // 检查角色
      if (to.meta.roles && to.meta.roles.length > 0) {
        const hasRole = to.meta.roles.some(role => userStore.hasRole(role))

        if (!hasRole) {
          ElMessage.error('没有权限访问该页面')
          next({ path: '/403' })
          return
        }
      }

      next()
    }
  } else {
    // 未登录
    if (whiteList.includes(to.path)) {
      next()
    } else {
      next({ path: '/login', query: { redirect: to.fullPath } })
    }
  }
})

router.afterEach(to => {
  // 添加标签页
  if (to.meta.title && to.path !== '/login') {
    const appStore = useAppStore()
    appStore.addTab({
      path: to.path,
      title: to.meta.title as string,
      name: to.name as string,
    })
  }
})

// 需要在afterEach中使用appStore，但这里需要导入
import { useAppStore } from '@/stores/app'
