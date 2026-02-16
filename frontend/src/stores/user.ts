import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, Role, Permission } from '@/types/api'
import { login as loginApi, getUser } from '@/api/user'

export const useUserStore = defineStore('user', () => {
  // 状态
  const token = ref<string>(localStorage.getItem('token') || '')
  const userInfo = ref<User | null>(null)
  const roles = ref<Role[]>([])
  const permissions = ref<string[]>([])

  // 计算属性
  const isLoggedIn = computed(() => !!token.value)
  const username = computed(() => userInfo.value?.username || '')
  const realName = computed(() => userInfo.value?.real_name || userInfo.value?.username || '')

  // 登录
  async function login(username: string, password: string) {
    const res = await loginApi({ username, password })
    if (res.code === 0) {
      token.value = res.data.token
      userInfo.value = res.data.user
      localStorage.setItem('token', res.data.token)

      // 提取角色和权限
      if (res.data.user.roles) {
        roles.value = res.data.user.roles
        permissions.value = res.data.user.roles
          .flatMap(role => role.permissions || [])
          .map((p: Permission) => p.code)
      }
    }
    return res
  }

  // 登出
  function logout() {
    token.value = ''
    userInfo.value = null
    roles.value = []
    permissions.value = []
    localStorage.removeItem('token')
  }

  // 获取用户信息
  async function fetchUserInfo() {
    if (!token.value) return null

    try {
      const res = await getUser(userInfo.value?.id || 0)
      if (res.code === 0) {
        userInfo.value = res.data
      }
      return res.data
    } catch {
      logout()
      return null
    }
  }

  // 检查权限
  function hasPermission(permission: string): boolean {
    return permissions.value.includes(permission) || permissions.value.includes('admin:all')
  }

  // 检查角色
  function hasRole(roleCode: string): boolean {
    return roles.value.some(role => role.code === roleCode)
  }

  return {
    token,
    userInfo,
    roles,
    permissions,
    isLoggedIn,
    username,
    realName,
    login,
    logout,
    fetchUserInfo,
    hasPermission,
    hasRole,
  }
})
