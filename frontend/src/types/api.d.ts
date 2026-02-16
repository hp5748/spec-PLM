// API响应类型
export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
  timestamp: number
}

// 分页数据类型
export interface PageData<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}

// 分页请求参数
export interface PageParams {
  page?: number
  page_size?: number
  sort?: string
  order?: 'asc' | 'desc'
}

// 用户类型
export interface User {
  id: number
  username: string
  email: string
  phone: string
  real_name: string
  avatar: string
  organization_id: number
  department_id: number
  organization?: Organization
  department?: Department
  status: 'active' | 'inactive' | 'locked'
  roles: Role[]
  created_at: string
  updated_at: string
}

// 组织类型
export interface Organization {
  id: number
  name: string
  code: string
  description: string
  logo: string
  status: 'active' | 'inactive'
  created_at: string
  updated_at: string
}

// 部门类型
export interface Department {
  id: number
  organization_id: number
  parent_id: number
  name: string
  code: string
  sort_order: number
  manager_id: number
  status: 'active' | 'inactive'
  parent?: Department
  children?: Department[]
  created_at: string
  updated_at: string
}

// 角色类型
export interface Role {
  id: number
  organization_id: number
  name: string
  code: string
  description: string
  is_system: boolean
  permissions: Permission[]
  created_at: string
  updated_at: string
}

// 权限类型
export interface Permission {
  id: number
  name: string
  code: string
  module: string
  description: string
}

// 登录请求
export interface LoginRequest {
  username: string
  password: string
}

// 登录响应
export interface LoginResponse {
  token: string
  user: User
}

// 路由元信息
declare module 'vue-router' {
  interface RouteMeta {
    title?: string
    icon?: string
    requiresAuth?: boolean
    permissions?: string[]
    roles?: string[]
  }
}
