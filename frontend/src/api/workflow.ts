import { request } from '@/utils/request'
import type { ApiResponse, PageData, PageParams } from '@/types/api'

// 流程定义类型
export interface WorkflowDefinition {
  id: number
  name: string
  code: string
  type: string
  version: string
  config: string
  status: string
  description: string
  created_at: string
  updated_at: string
  created_by: number
}

// 流程实例类型
export interface WorkflowInstance {
  id: number
  definition_id: number
  business_type: string
  business_id: number
  title: string
  status: string
  initiator_id: number
  current_node: string
  comment: string
  created_at: string
  updated_at: string
  definition?: WorkflowDefinition
  initiator?: {
    id: number
    username: string
    real_name: string
  }
}

// 流程实例响应（含业务对象信息）
export interface WorkflowInstanceResponse extends WorkflowInstance {
  business_name: string
  business_code: string
  business_status: string
}

// 流程历史类型
export interface WorkflowHistory {
  id: number
  instance_id: number
  node_name: string
  action: string
  operator_id: number
  operator_name: string
  comment: string
  created_at: string
  operator?: {
    id: number
    username: string
    real_name: string
  }
}

// 待办事项响应
export interface TodoItemResponse extends WorkflowInstanceResponse {
  pending_action: string
}

// 流程定义列表查询参数
export interface DefinitionListParams extends PageParams {
  type?: string
  status?: string
  keyword?: string
}

// 流程实例列表查询参数
export interface InstanceListParams extends PageParams {
  type?: string
  status?: string
  keyword?: string
}

// 待办列表查询参数
export interface TodoListParams extends PageParams {
  type?: string    // 业务类型筛选
  status?: string
}

// 创建流程定义请求
export interface CreateDefinitionRequest {
  name: string
  code: string
  type: string
  config: string
  description?: string
}

// 更新流程定义请求
export interface UpdateDefinitionRequest {
  name?: string
  config?: string
  description?: string
}

// 发起流程请求
export interface InitiateWorkflowRequest {
  definition_id: number
  business_type: string
  business_id: number
  title: string
  comment?: string
}

// 同意审批请求
export interface ApproveRequest {
  comment?: string
}

// 驳回审批请求
export interface RejectRequest {
  comment: string
}

// 转交审批请求
export interface TransferRequest {
  target_user_id: number
  comment?: string
}

// ==================== 流程定义API ====================

// 获取流程定义列表
export function getDefinitionList(params: DefinitionListParams): Promise<ApiResponse<PageData<WorkflowDefinition>>> {
  return request.get('/workflow/workflows/definitions', { params })
}

// 获取流程定义详情
export function getDefinition(id: number): Promise<ApiResponse<WorkflowDefinition>> {
  return request.get(`/workflow/workflows/definitions/${id}`)
}

// 创建流程定义
export function createDefinition(data: CreateDefinitionRequest): Promise<ApiResponse<WorkflowDefinition>> {
  return request.post('/workflow/workflows/definitions', data)
}

// 更新流程定义
export function updateDefinition(id: number, data: UpdateDefinitionRequest): Promise<ApiResponse<WorkflowDefinition>> {
  return request.put(`/workflow/workflows/definitions/${id}`, data)
}

// 发布流程定义
export function releaseDefinition(id: number): Promise<ApiResponse<WorkflowDefinition>> {
  return request.post(`/workflow/workflows/definitions/${id}/release`)
}

// 删除流程定义
export function deleteDefinition(id: number): Promise<ApiResponse<null>> {
  return request.delete(`/workflow/workflows/definitions/${id}`)
}

// 获取指定类型的激活流程定义
export function getActiveDefinition(type: string): Promise<ApiResponse<WorkflowDefinition>> {
  return request.get('/workflow/workflows/definitions/active', { params: { type } })
}

// ==================== 流程实例API ====================

// 获取流程实例列表
export function getInstanceList(params: InstanceListParams): Promise<ApiResponse<PageData<WorkflowInstanceResponse>>> {
  return request.get('/workflow/workflows/instances', { params })
}

// 获取流程实例详情
export function getInstance(id: number): Promise<ApiResponse<WorkflowInstanceResponse>> {
  return request.get(`/workflow/workflows/instances/${id}`)
}

// 发起流程
export function initiateWorkflow(data: InitiateWorkflowRequest): Promise<ApiResponse<WorkflowInstanceResponse>> {
  return request.post('/workflow/workflows/instances', data)
}

// 获取流程历史
export function getHistories(instanceId: number): Promise<ApiResponse<WorkflowHistory[]>> {
  return request.get(`/workflow/workflows/instances/${instanceId}/history`)
}

// ==================== 待办API ====================

// 获取我的待办
export function getMyTodos(params: TodoListParams): Promise<ApiResponse<PageData<TodoItemResponse>>> {
  return request.get('/workflow/workflows/todos', { params })
}

// 获取待办数量
export function getTodoCount(): Promise<ApiResponse<{ count: number }>> {
  return request.get('/workflow/workflows/todos/count')
}

// ==================== 审批操作API ====================

// 同意审批
export function approveWorkflow(instanceId: number, data: ApproveRequest): Promise<ApiResponse<WorkflowInstanceResponse>> {
  return request.post(`/workflow/workflows/instances/${instanceId}/approve`, data)
}

// 驳回审批
export function rejectWorkflow(instanceId: number, data: RejectRequest): Promise<ApiResponse<WorkflowInstanceResponse>> {
  return request.post(`/workflow/workflows/instances/${instanceId}/reject`, data)
}

// 转交审批
export function transferWorkflow(instanceId: number, data: TransferRequest): Promise<ApiResponse<WorkflowInstanceResponse>> {
  return request.post(`/workflow/workflows/instances/${instanceId}/transfer`, data)
}

// 撤回流程
export function withdrawWorkflow(instanceId: number): Promise<ApiResponse<WorkflowInstanceResponse>> {
  return request.post(`/workflow/workflows/instances/${instanceId}/withdraw`)
}
