<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>文档列表</span>
          <el-button type="primary" @click="handleUpload">上传文档</el-button>
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
          <el-select v-model="searchForm.doc_type" placeholder="请选择" clearable>
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
          <el-select v-model="searchForm.status" placeholder="请选择" clearable>
            <el-option label="草稿" value="DRAFT" />
            <el-option label="已检入" value="CHECKED_IN" />
            <el-option label="已发布" value="RELEASED" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

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
        <el-table-column label="操作" fixed="right" width="200">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleDownload(row)">下载</el-button>
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
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

    <!-- 上传对话框 -->
    <el-dialog v-model="uploadDialogVisible" title="上传文档" width="600">
      <el-form ref="uploadFormRef" :model="uploadForm" :rules="uploadRules" label-width="100px">
        <el-form-item label="文档编码" prop="doc_id">
          <el-input v-model="uploadForm.doc_id" placeholder="请输入文档编码" />
        </el-form-item>
        <el-form-item label="文档名称" prop="doc_name">
          <el-input v-model="uploadForm.doc_name" placeholder="请输入文档名称" />
        </el-form-item>
        <el-form-item label="文档类型" prop="doc_type">
          <el-select v-model="uploadForm.doc_type" placeholder="请选择文档类型" style="width: 100%">
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
        <el-form-item label="选择文件" prop="file">
          <el-upload
            ref="uploadRef"
            :auto-upload="false"
            :limit="1"
            :on-change="handleFileChange"
            :on-exceed="handleExceed"
            drag
          >
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">
              拖拽文件到此处或 <em>点击上传</em>
            </div>
            <template #tip>
              <div class="el-upload__tip">支持上传100MB以内的文件</div>
            </template>
          </el-upload>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="uploadDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitUpload" :loading="uploading">确定上传</el-button>
      </template>
    </el-dialog>

    <!-- 编辑对话框 -->
    <el-dialog v-model="editDialogVisible" title="编辑文档" width="500">
      <el-form ref="editFormRef" :model="editForm" :rules="editRules" label-width="100px">
        <el-form-item label="文档名称" prop="doc_name">
          <el-input v-model="editForm.doc_name" placeholder="请输入文档名称" />
        </el-form-item>
        <el-form-item label="文档类型" prop="doc_type">
          <el-select v-model="editForm.doc_type" placeholder="请选择文档类型" style="width: 100%">
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
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitEdit" :loading="editing">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules, type UploadProps, type UploadUserFile } from 'element-plus'
import { UploadFilled } from '@element-plus/icons-vue'
import {
  getDocumentList,
  uploadDocument,
  updateDocument,
  deleteDocument,
  getDocumentDownloadUrl,
  formatFileSize,
  type Document,
} from '@/api/document'

defineOptions({
  name: 'DocumentList',
})

const loading = ref(false)
const uploading = ref(false)
const editing = ref(false)
const tableData = ref<Document[]>([])
const uploadDialogVisible = ref(false)
const editDialogVisible = ref(false)
const currentId = ref<number | null>(null)
const uploadFormRef = ref<FormInstance>()
const editFormRef = ref<FormInstance>()
const uploadRef = ref()
const selectedFile = ref<File | null>(null)

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

const uploadForm = reactive({
  doc_id: '',
  doc_name: '',
  doc_type: '',
})

const editForm = reactive({
  doc_name: '',
  doc_type: '',
})

const uploadRules: FormRules = {
  doc_id: [
    { required: true, message: '请输入文档编码', trigger: 'blur' },
    { max: 50, message: '文档编码不能超过50个字符', trigger: 'blur' },
  ],
  doc_name: [
    { required: true, message: '请输入文档名称', trigger: 'blur' },
    { max: 200, message: '文档名称不能超过200个字符', trigger: 'blur' },
  ],
  doc_type: [{ required: true, message: '请选择文档类型', trigger: 'change' }],
}

const editRules: FormRules = {
  doc_name: [
    { required: true, message: '请输入文档名称', trigger: 'blur' },
  ],
  doc_type: [{ required: true, message: '请选择文档类型', trigger: 'change' }],
}

onMounted(() => {
  fetchData()
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
  fetchData()
}

function handleReset() {
  searchForm.doc_id = ''
  searchForm.doc_name = ''
  searchForm.doc_type = ''
  searchForm.status = ''
  pagination.page = 1
  fetchData()
}

function handleUpload() {
  uploadForm.doc_id = ''
  uploadForm.doc_name = ''
  uploadForm.doc_type = ''
  selectedFile.value = null
  uploadDialogVisible.value = true
}

const handleFileChange: UploadProps['onChange'] = (uploadFile) => {
  selectedFile.value = uploadFile.raw || null
  // 自动填充文档名称
  if (!uploadForm.doc_name && uploadFile.name) {
    uploadForm.doc_name = uploadFile.name.replace(/\.[^/.]+$/, '')
  }
}

const handleExceed: UploadProps['onExceed'] = () => {
  ElMessage.warning('只能上传一个文件')
}

async function handleSubmitUpload() {
  const valid = await uploadFormRef.value?.validate()
  if (!valid) return

  if (!selectedFile.value) {
    ElMessage.warning('请选择要上传的文件')
    return
  }

  uploading.value = true
  try {
    const formData = new FormData()
    formData.append('file', selectedFile.value)
    formData.append('doc_id', uploadForm.doc_id)
    formData.append('doc_name', uploadForm.doc_name)
    formData.append('doc_type', uploadForm.doc_type)

    const res = await uploadDocument(formData)
    if (res.code === 0) {
      ElMessage.success('上传成功')
      uploadDialogVisible.value = false
      fetchData()
    }
  } catch (error: any) {
    ElMessage.error(error.message || '上传失败')
  } finally {
    uploading.value = false
  }
}

function handleEdit(row: Document) {
  currentId.value = row.id
  editForm.doc_name = row.doc_name
  editForm.doc_type = row.doc_type
  editDialogVisible.value = true
}

async function handleSubmitEdit() {
  const valid = await editFormRef.value?.validate()
  if (!valid) return

  editing.value = true
  try {
    await updateDocument(currentId.value!, {
      doc_name: editForm.doc_name,
      doc_type: editForm.doc_type,
    })
    ElMessage.success('更新成功')
    editDialogVisible.value = false
    fetchData()
  } catch (error: any) {
    ElMessage.error(error.message || '更新失败')
  } finally {
    editing.value = false
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

async function handleDelete(row: Document) {
  try {
    await ElMessageBox.confirm(`确定要删除文档 ${row.doc_name} 吗？`, '提示', {
      type: 'warning',
    })

    await deleteDocument(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch {
    // 取消删除
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
    CHECKED_IN: 'warning',
    RELEASED: 'success',
  }
  return map[status] || 'info'
}

function getStatusText(status: string) {
  const map: Record<string, string> = {
    DRAFT: '草稿',
    CHECKED_IN: '已检入',
    RELEASED: '已发布',
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
