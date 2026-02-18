import { request } from '@/utils/request'
import type { ApiResponse, PageData } from '@/types/api'

// BOM视图类型
export interface BOMView {
  id: number
  name: string
  root_material_id: number
  root_version: string
  is_exact: boolean
  status: string
  description: string
  created_at: string
  updated_at: string
  created_by: number
  root_material?: Material
  item_count?: number
}

// BOM项类型
export interface BOMItem {
  id: number
  bom_view_id: number
  parent_id: number | null
  material_id: number
  version: string
  quantity: number
  sort_order: number
  level: number
  material?: Material
  children?: BOMItem[]
}

// BOM树节点
export interface BOMTreeNode {
  id: number
  material_id: number
  item_id: string
  item_name: string
  version: string
  quantity: number
  sort_order: number
  level: number
  unit: string
  children?: BOMTreeNode[]
}

// 物料简要信息
export interface Material {
  id: number
  item_id: string
  item_name: string
  item_type: string
  version: string
  unit: string
}

// BOM列表查询参数
export interface BOMListParams {
  page?: number
  page_size?: number
  name?: string
  status?: string
  is_exact?: boolean
}

// 创建BOM视图请求
export interface CreateBOMViewRequest {
  name: string
  root_material_id: number
  root_version?: string
  is_exact?: boolean
  description?: string
}

// 更新BOM视图请求
export interface UpdateBOMViewRequest {
  name?: string
  status?: string
  description?: string
}

// 添加BOM项请求
export interface AddBOMItemRequest {
  parent_id?: number
  material_id: number
  version?: string
  quantity: number
  sort_order?: number
}

// 更新BOM项请求
export interface UpdateBOMItemRequest {
  quantity: number
  sort_order?: number
}

// 转换BOM类型请求
export interface ConvertBOMRequest {
  to_exact: boolean
}

// 获取BOM视图列表
export function getBOMList(params: BOMListParams): Promise<ApiResponse<PageData<BOMView>>> {
  return request.get('/bom/boms', { params })
}

// 创建BOM视图
export function createBOMView(data: CreateBOMViewRequest): Promise<ApiResponse<BOMView>> {
  return request.post('/bom/boms', data)
}

// 获取BOM视图详情
export function getBOMDetail(id: number): Promise<ApiResponse<BOMView>> {
  return request.get(`/bom/boms/${id}`)
}

// 更新BOM视图
export function updateBOMView(id: number, data: UpdateBOMViewRequest): Promise<ApiResponse<BOMView>> {
  return request.put(`/bom/boms/${id}`, data)
}

// 删除BOM视图
export function deleteBOMView(id: number): Promise<ApiResponse<null>> {
  return request.delete(`/bom/boms/${id}`)
}

// 获取BOM树形结构
export function getBOMTree(id: number): Promise<ApiResponse<BOMTreeNode>> {
  return request.get(`/bom/boms/${id}/tree`)
}

// 添加BOM项
export function addBOMItem(bomId: number, data: AddBOMItemRequest): Promise<ApiResponse<BOMItem>> {
  return request.post(`/bom/boms/${bomId}/items`, data)
}

// 更新BOM项
export function updateBOMItem(
  bomId: number,
  itemId: number,
  data: UpdateBOMItemRequest
): Promise<ApiResponse<BOMItem>> {
  return request.put(`/bom/boms/${bomId}/items/${itemId}`, data)
}

// 删除BOM项
export function deleteBOMItem(bomId: number, itemId: number): Promise<ApiResponse<null>> {
  return request.delete(`/bom/boms/${bomId}/items/${itemId}`)
}

// 转换BOM类型
export function convertBOMType(id: number, data: ConvertBOMRequest): Promise<ApiResponse<null>> {
  return request.post(`/bom/boms/${id}/convert`, data)
}

// 导出BOM
export function exportBOM(id: number): Promise<ApiResponse<{ download_url: string }>> {
  return request.get(`/bom/boms/${id}/export`)
}

// 导入BOM
export function importBOM(file: File): Promise<ApiResponse<{ success_count: number; fail_count: number }>> {
  const formData = new FormData()
  formData.append('file', file)
  return request.post('/bom/import', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
}
