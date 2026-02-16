import { request } from '@/utils/request'
import type {
  ApiResponse,
  PageData,
  PageParams,
  User,
  LoginRequest,
  LoginResponse,
  Organization,
  Department,
  Role,
} from '@/types/api'

// 用户登录
export function login(data: LoginRequest): Promise<ApiResponse<LoginResponse>> {
  return request.post('/user/auth/login', data)
}

// 用户注册
export function register(data: {
  username: string
  password: string
  email?: string
  phone?: string
  real_name?: string
}): Promise<ApiResponse<User>> {
  return request.post('/user/auth/register', data)
}

// 刷新Token
export function refreshToken(): Promise<ApiResponse<{ token: string }>> {
  return request.post('/user/auth/refresh')
}

// 获取用户列表
export function getUserList(params: PageParams): Promise<ApiResponse<PageData<User>>> {
  return request.get('/user/users', { params })
}

// 获取用户详情
export function getUser(id: number): Promise<ApiResponse<User>> {
  return request.get(`/user/users/${id}`)
}

// 创建用户
export function createUser(data: Partial<User> & { password: string; role_ids?: number[] }): Promise<ApiResponse<User>> {
  return request.post('/user/users', data)
}

// 更新用户
export function updateUser(id: number, data: Partial<User> & { role_ids?: number[] }): Promise<ApiResponse<User>> {
  return request.put(`/user/users/${id}`, data)
}

// 删除用户
export function deleteUser(id: number): Promise<ApiResponse<null>> {
  return request.delete(`/user/users/${id}`)
}

// 获取组织列表
export function getOrganizationList(): Promise<ApiResponse<Organization[]>> {
  return request.get('/user/organizations')
}

// 创建组织
export function createOrganization(data: Partial<Organization>): Promise<ApiResponse<Organization>> {
  return request.post('/user/organizations', data)
}

// 获取部门列表（树形）
export function getDepartmentList(organizationId?: string): Promise<ApiResponse<Department[]>> {
  const params = organizationId ? { organization_id: organizationId } : {}
  return request.get('/user/departments', { params })
}

// 创建部门
export function createDepartment(data: Partial<Department>): Promise<ApiResponse<Department>> {
  return request.post('/user/departments', data)
}

// 获取角色列表
export function getRoleList(): Promise<ApiResponse<Role[]>> {
  return request.get('/user/roles')
}

// 创建角色
export function createRole(data: Partial<Role> & { permission_ids?: number[] }): Promise<ApiResponse<Role>> {
  return request.post('/user/roles', data)
}
