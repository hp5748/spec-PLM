import { request } from '@/utils/request'
import type { ApiResponse, PageData } from '@/types/api'

// 脚本类型
export interface Script {
  id: number
  name: string
  code: string
  type: 'JAVASCRIPT' | 'SQL'
  trigger_type: string
  business_type: string
  content: string
  description: string
  status: string
  timeout: number
  created_at: string
  updated_at: string
  created_by: number
}

// 脚本执行日志
export interface ScriptExecutionLog {
  id: number
  script_id: number
  script_name: string
  instance_id: number
  business_type: string
  business_id: number
  trigger_type: string
  input: string
  output: string
  success: boolean
  error_msg: string
  duration: number
  executed_at: string
  executed_by: number
}

// 脚本列表查询参数
export interface ScriptListParams {
  page?: number
  page_size?: number
  type?: string
  trigger_type?: string
  business_type?: string
  status?: string
  keyword?: string
}

// 创建脚本请求
export interface CreateScriptRequest {
  name: string
  code: string
  type: 'JAVASCRIPT' | 'SQL'
  trigger_type: string
  business_type: string
  content: string
  description?: string
  timeout?: number
}

// 更新脚本请求
export interface UpdateScriptRequest {
  name?: string
  content?: string
  description?: string
  timeout?: number
}

// 测试脚本请求
export interface ScriptTestRequest {
  type: 'JAVASCRIPT' | 'SQL'
  content: string
  business_type?: string
  business_id?: number
  context?: Record<string, any>
}

// 测试脚本结果
export interface ScriptTestResult {
  success: boolean
  output: Record<string, any>
  error: string
  duration: number
  logs: string[]
}

// 执行脚本结果
export interface ExecuteScriptResult {
  success: boolean
  output: Record<string, any>
  error: string
  log_id: number
}

// 获取脚本列表
export function getScriptList(params: ScriptListParams): Promise<ApiResponse<PageData<Script>>> {
  return request.get('/script/scripts', { params })
}

// 获取脚本详情
export function getScriptDetail(id: number): Promise<ApiResponse<Script>> {
  return request.get(`/script/scripts/${id}`)
}

// 创建脚本
export function createScript(data: CreateScriptRequest): Promise<ApiResponse<Script>> {
  return request.post('/script/scripts', data)
}

// 更新脚本
export function updateScript(id: number, data: UpdateScriptRequest): Promise<ApiResponse<Script>> {
  return request.put(`/script/scripts/${id}`, data)
}

// 删除脚本
export function deleteScript(id: number): Promise<ApiResponse<null>> {
  return request.delete(`/script/scripts/${id}`)
}

// 发布脚本
export function releaseScript(id: number): Promise<ApiResponse<Script>> {
  return request.post(`/script/scripts/${id}/release`)
}

// 废弃脚本
export function obsoleteScript(id: number): Promise<ApiResponse<Script>> {
  return request.post(`/script/scripts/${id}/obsolete`)
}

// 测试脚本
export function testScript(data: ScriptTestRequest): Promise<ApiResponse<ScriptTestResult>> {
  return request.post('/script/scripts/test', data)
}

// 获取执行日志列表
export function getExecutionLogs(params: {
  script_id?: number
  page?: number
  page_size?: number
}): Promise<ApiResponse<PageData<ScriptExecutionLog>>> {
  return request.get('/script-logs', { params })
}

// 获取执行日志详情
export function getExecutionLogDetail(id: number): Promise<ApiResponse<ScriptExecutionLog>> {
  return request.get(`/script-logs/${id}`)
}
