import { request } from '@/utils/request'
import type { ApiResponse, PageData, PageParams } from '@/types/api'

// 物料类型
export interface Material {
  id: number
  item_id: string
  item_name: string
  description: string
  item_type: string
  version: string
  unit: string
  status: string
  attributes: string
  created_at: string
  updated_at: string
  created_by: number
}

// 物料属性类型
export interface MaterialAttribute {
  id: number
  material_id: number
  attr_type: string
  attr_key: string
  attr_value: string
  unit: string
  sort_order: number
}

// 物料详情响应
export interface MaterialDetail extends Material {
  dynamicAttributes?: MaterialAttribute[]
}

// 物料列表查询参数
export interface MaterialListParams extends PageParams {
  item_id?: string
  item_name?: string
  item_type?: string
  status?: string
  keyword?: string
}

// 动态属性请求
export interface DynamicAttributeRequest {
  attr_type: string
  attr_key: string
  attr_value: string
  unit?: string
}

// 创建物料请求
export interface CreateMaterialRequest {
  item_id: string
  item_name: string
  description?: string
  item_type: string
  unit?: string
  attributes?: string
  dynamic_attributes?: DynamicAttributeRequest[]
}

// 更新物料请求
export interface UpdateMaterialRequest {
  item_name?: string
  description?: string
  item_type?: string
  unit?: string
  status?: string
  attributes?: string
  dynamic_attributes?: DynamicAttributeRequest[]
}

// 获取物料列表
export function getMaterialList(params: MaterialListParams): Promise<ApiResponse<PageData<Material>>> {
  return request.get('/material/materials', { params })
}

// 获取物料详情
export function getMaterial(id: number): Promise<ApiResponse<MaterialDetail>> {
  return request.get(`/material/materials/${id}`)
}

// 创建物料
export function createMaterial(data: CreateMaterialRequest): Promise<ApiResponse<Material>> {
  return request.post('/material/materials', data)
}

// 更新物料
export function updateMaterial(id: number, data: UpdateMaterialRequest): Promise<ApiResponse<Material>> {
  return request.put(`/material/materials/${id}`, data)
}

// 删除物料
export function deleteMaterial(id: number): Promise<ApiResponse<null>> {
  return request.delete(`/material/materials/${id}`)
}

// 搜索物料
export function searchMaterials(keyword: string, params?: PageParams): Promise<ApiResponse<PageData<Material>>> {
  return request.get('/material/materials/search', { params: { keyword, ...params } })
}

// 版本升级
export function incrementMaterialVersion(id: number): Promise<ApiResponse<Material>> {
  return request.post(`/material/materials/${id}/version`)
}
