import { request } from '@/utils/request'
import type { ApiResponse, PageData, PageParams } from '@/types/api'

// 文档类型
export interface Document {
  id: number
  doc_id: string
  doc_name: string
  doc_type: string
  file_path: string
  file_name: string
  file_size: number
  mime_type: string
  version: string
  status: string
  attributes: string
  created_at: string
  updated_at: string
  created_by: number
}

// 文档列表查询参数
export interface DocumentListParams extends PageParams {
  doc_id?: string
  doc_name?: string
  doc_type?: string
  status?: string
  keyword?: string
}

// 创建文档请求
export interface CreateDocumentRequest {
  doc_id: string
  doc_name: string
  doc_type: string
  attributes?: string
}

// 更新文档请求
export interface UpdateDocumentRequest {
  doc_name?: string
  doc_type?: string
  status?: string
  attributes?: string
}

// 获取文档列表
export function getDocumentList(params: DocumentListParams): Promise<ApiResponse<PageData<Document>>> {
  return request.get('/document/documents', { params })
}

// 获取文档详情
export function getDocument(id: number): Promise<ApiResponse<Document>> {
  return request.get(`/document/documents/${id}`)
}

// 上传文档
export function uploadDocument(formData: FormData): Promise<ApiResponse<Document>> {
  return request.post('/document/documents', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  })
}

// 更新文档
export function updateDocument(id: number, data: UpdateDocumentRequest): Promise<ApiResponse<Document>> {
  return request.put(`/document/documents/${id}`, data)
}

// 删除文档
export function deleteDocument(id: number): Promise<ApiResponse<null>> {
  return request.delete(`/document/documents/${id}`)
}

// 获取文档下载链接
export function getDocumentDownloadUrl(id: number): Promise<ApiResponse<{ download_url: string }>> {
  return request.get(`/document/documents/${id}/download`)
}

// 搜索文档
export function searchDocuments(keyword: string, params?: PageParams): Promise<ApiResponse<PageData<Document>>> {
  return request.get('/document/documents/search', { params: { keyword, ...params } })
}

// 格式化文件大小
export function formatFileSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}
