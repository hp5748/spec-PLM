import request from '@/utils/request'

// 字段Schema定义
export interface FieldSchema {
  key: string
  label: string
  type: 'text' | 'number' | 'select' | 'date' | 'textarea'
  unit?: string
  required?: boolean
  options?: string[]
}

// Schema配置
export interface SchemaConfig {
  label?: string
  fields?: FieldSchema[]
}

// 属性Schema
export interface AttributeSchema {
  id: number
  entity_type: 'MATERIAL' | 'DOCUMENT'
  type_code: string
  attr_type: 'main' | 'description' | 'specification' | 'custom'
  schema_name: string
  schema_config: SchemaConfig
  version: string
  status: 'DRAFT' | 'RELEASED'
  is_active: boolean
  created_at: string
  updated_at: string
  created_by: number
}

// 创建请求
export interface CreateAttributeSchemaRequest {
  entity_type: 'MATERIAL' | 'DOCUMENT'
  type_code: string
  attr_type: 'main' | 'description' | 'specification' | 'custom'
  schema_name?: string
  schema_config: SchemaConfig
}

// 更新请求
export interface UpdateAttributeSchemaRequest {
  schema_name?: string
  schema_config: SchemaConfig
}

// 查询参数
export interface AttributeSchemaQuery {
  page?: number
  page_size?: number
  entity_type?: string
  type_code?: string
  attr_type?: string
  status?: string
}

// 获取属性Schema列表
export function getAttributeSchemaList(params?: AttributeSchemaQuery) {
  return request.get<{ list: AttributeSchema[]; total: number }>('/attribute/attribute-schemas', { params })
}

// 获取属性Schema详情
export function getAttributeSchema(id: number) {
  return request.get<AttributeSchema>(`/attribute/attribute-schemas/${id}`)
}

// 创建属性Schema
export function createAttributeSchema(data: CreateAttributeSchemaRequest) {
  return request.post<AttributeSchema>('/attribute/attribute-schemas', data)
}

// 更新属性Schema
export function updateAttributeSchema(id: number, data: UpdateAttributeSchemaRequest) {
  return request.put<AttributeSchema>(`/attribute/attribute-schemas/${id}`, data)
}

// 删除属性Schema
export function deleteAttributeSchema(id: number) {
  return request.delete(`/attribute/attribute-schemas/${id}`)
}

// 发布属性Schema
export function releaseAttributeSchema(id: number) {
  return request.post(`/attribute/attribute-schemas/${id}/release`)
}

// 激活属性Schema版本
export function activateAttributeSchema(id: number) {
  return request.post(`/attribute/attribute-schemas/${id}/activate`)
}

// 获取当前激活的Schema
export function getActiveSchema(entityType: string, typeCode: string, attrType: string) {
  return request.get<AttributeSchema | null>('/attribute/attribute-schemas/active', {
    params: { entity_type: entityType, type_code: typeCode, attr_type: attrType }
  })
}

// 获取指定类型的所有激活Schema
export function getAllActiveSchemas(entityType: string, typeCode: string) {
  return request.get<AttributeSchema[]>('/attribute/attribute-schemas/all-active', {
    params: { entity_type: entityType, type_code: typeCode }
  })
}
