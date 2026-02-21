<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>文档查询</span>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="文档编码">
          <el-input v-model="searchForm.doc_id" placeholder="请输入文档编码" clearable />
        </el-form-item>
        <el-form-item label="文档名称">
          <el-input v-model="searchForm.doc_name" placeholder="请输入文档名称" clearable />
        </el-form-item>
        <el-form-item label="文档类型">
          <el-select v-model="searchForm.doc_type" placeholder="请选择文档类型" clearable>
            <el-option label="PDF" value="PDF" />
            <el-option label="Word文档" value="DOC" />
            <el-option label="Excel表格" value="XLS" />
            <el-option label="PPT" value="PPT" />
            <el-option label="CAD图纸" value="CAD" />
            <el-option label="图片" value="IMAGE" />
            <el-option label="压缩包" value="ARCHIVE" />
            <el-option label="其他" value="OTHER" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择状态" clearable>
            <el-option label="拟制中" value="DRAFT" />
            <el-option label="审核中" value="REVIEWING" />
            <el-option label="已发布" value="RELEASED" />
            <el-option label="已驳回" value="REJECTED" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 提示信息 -->
      <el-alert v-if="!hasSearched" type="info" :closable="false" show-icon style="margin-bottom: 16px">
        <template #title>请设置筛选条件后点击"搜索"按钮查询数据</template>
      </el-alert>

      <!-- 表格 -->
      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="doc_id" label="文档编码" width="150" />
        <el-table-column prop="doc_name" label="文档名称" width="200" show-overflow-tooltip />
        <el-table-column prop="doc_type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getDocTypeTag(row.doc_type)">{{ row.doc_type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="version" label="版本" width="80" />
        <el-table-column prop="file_size" label="大小" width="100">
          <template #default="{ row }">
            {{ formatFileSize(row.file_size) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" fixed="right" width="150">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleDownload(row)">下载</el-button>
            <el-button type="primary" link @click="handleView(row)">查看</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        class="pagination"
        @size-change="fetchData"
        @current-change="fetchData"
      />
    </el-card>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" title="文档详情" width="700">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="文档编码">{{ currentDocument?.doc_id }}</el-descriptions-item>
        <el-descriptions-item label="文档名称">{{ currentDocument?.doc_name }}</el-descriptions-item>
        <el-descriptions-item label="文档类型">{{ currentDocument?.doc_type }}</el-descriptions-item>
        <el-descriptions-item label="版本">{{ currentDocument?.version }}</el-descriptions-item>
        <el-descriptions-item label="文件大小">{{ formatFileSize(currentDocument?.file_size) }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(currentDocument?.status || '')">{{ getStatusText(currentDocument?.status || '') }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ currentDocument?.created_at }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ currentDocument?.updated_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
        <el-button type="primary" @click="currentDocument && handleDownload(currentDocument)">下载</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  getDocumentList,
  getDocument,
  getDocumentDownloadUrl,
  formatFileSize,
  type Document,
  type DocumentDetail,
} from '@/api/document'

defineOptions({
  name: 'DocumentQuery',
})

const loading = ref(false)
const tableData = ref<Document[]>([])
const hasSearched = ref(false)
const detailVisible = ref(false)
const currentDocument = ref<DocumentDetail | null>(null)

const searchForm = reactive({
  doc_id: '',
  doc_name: '',
  doc_type: '',
  status: '',
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0,
})

onMounted(() => {
  // 默认不查询
})

async function fetchData() {
  loading.value = true
  try {
    const res = await getDocumentList({
      page: pagination.page,
      page_size: pagination.pageSize,
      ...searchForm,
    })
    if (res.code === 0) {
      tableData.value = res.data.list
      pagination.total = res.data.total
    }
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  pagination.page = 1
  hasSearched.value = true
  fetchData()
}

function handleReset() {
  searchForm.doc_id = ''
  searchForm.doc_name = ''
  searchForm.doc_type = ''
  searchForm.status = ''
  pagination.page = 1
  hasSearched.value = false
  tableData.value = []
  pagination.total = 0
}

async function handleView(row: Document) {
  try {
    const res = await getDocument(row.id)
    if (res.code === 0) {
      currentDocument.value = res.data
      detailVisible.value = true
    }
  } catch (error) {
    console.error(error)
  }
}

async function handleDownload(row: Document) {
  try {
    const res = await getDocumentDownloadUrl(row.id)
    if (res.code === 0 && res.data.download_url) {
      window.open(res.data.download_url, '_blank')
    }
  } catch (error: any) {
    ElMessage.error(error.message || '获取下载链接失败')
  }
}

function getDocTypeTag(type: string) {
  const map: Record<string, string> = {
    PDF: 'danger',
    DOC: 'primary',
    XLS: 'success',
    PPT: 'warning',
    CAD: 'info',
    IMAGE: '',
    ARCHIVE: 'info',
    OTHER: 'info',
  }
  return map[type] || 'info'
}

function getStatusType(status: string) {
  const map: Record<string, string> = {
    DRAFT: 'info',
    REVIEWING: 'warning',
    RELEASED: 'success',
    REJECTED: 'danger',
  }
  return map[status] || 'info'
}

function getStatusText(status: string) {
  const map: Record<string, string> = {
    DRAFT: '拟制中',
    REVIEWING: '审核中',
    RELEASED: '已发布',
    REJECTED: '已驳回',
  }
  return map[status] || status
}
</script>

<style scoped lang="scss">
.page-container {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .search-form {
    margin-bottom: 20px;
  }

  .pagination {
    margin-top: 20px;
    justify-content: flex-end;
  }
}
</style>
