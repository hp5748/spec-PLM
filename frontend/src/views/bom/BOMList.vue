<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>BOM管理</span>
          <el-button type="primary" @click="handleAdd">新增BOM</el-button>
        </div>
      </template>

      <!-- 搜索表单 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="BOM名称">
          <el-input v-model="searchForm.name" placeholder="请输入BOM名称" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择状态" clearable>
            <el-option label="拟制中" value="DRAFT" />
            <el-option label="审核中" value="REVIEWING" />
            <el-option label="已发布" value="RELEASED" />
            <el-option label="已驳回" value="REJECTED" />
            <el-option label="已失效" value="OBSOLETE" />
          </el-select>
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="searchForm.is_exact" placeholder="请选择类型" clearable>
            <el-option label="精确BOM" :value="true" />
            <el-option label="非精确BOM" :value="false" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="BOM名称" min-width="180" />
        <el-table-column label="根物料" min-width="150">
          <template #default="{ row }">
            {{ row.root_material?.item_id }} - {{ row.root_material?.item_name }}
          </template>
        </el-table-column>
        <el-table-column prop="root_version" label="版本" width="80" />
        <el-table-column label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.is_exact ? 'success' : 'warning'">
              {{ row.is_exact ? '精确' : '非精确' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusTag(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="item_count" label="子项数" width="80" />
        <el-table-column label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="primary" link @click="handleViewTree(row)">查看结构</el-button>
            <el-button type="success" link @click="handleExport(row)">导出</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        class="pagination"
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="fetchData"
        @current-change="fetchData"
      />
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-form-item label="BOM名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入BOM名称" />
        </el-form-item>
        <el-form-item label="根物料" prop="root_material_id">
          <el-select
            v-model="formData.root_material_id"
            placeholder="请选择根物料"
            filterable
            remote
            :remote-method="searchMaterials"
            :loading="materialLoading"
            style="width: 100%"
          >
            <el-option
              v-for="item in materialOptions"
              :key="item.id"
              :label="`${item.item_id} - ${item.item_name}`"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="根物料版本" prop="root_version">
          <el-input v-model="formData.root_version" placeholder="精确BOM需要指定版本" />
        </el-form-item>
        <el-form-item label="BOM类型">
          <el-switch
            v-model="formData.is_exact"
            active-text="精确BOM"
            inactive-text="非精确BOM"
          />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="请输入描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import {
  getBOMList,
  createBOMView,
  deleteBOMView,
  exportBOM,
  type BOMView,
  type CreateBOMViewRequest
} from '@/api/bom'
import { searchMaterials as searchMaterialsApi, type Material } from '@/api/material'

const router = useRouter()

// 加载状态
const loading = ref(false)
const submitting = ref(false)
const materialLoading = ref(false)

// 数据
const tableData = ref<BOMView[]>([])
const materialOptions = ref<Material[]>([])

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 搜索表单
const searchForm = reactive({
  name: '',
  status: '',
  is_exact: undefined as boolean | undefined
})

// 对话框
const dialogVisible = ref(false)
const dialogTitle = computed(() => (isEdit.value ? '编辑BOM' : '新增BOM'))
const isEdit = ref(false)
const currentId = ref<number | null>(null)

// 表单
const formRef = ref<FormInstance>()
const formData = reactive<CreateBOMViewRequest>({
  name: '',
  root_material_id: 0,
  root_version: '',
  is_exact: false,
  description: ''
})

const rules: FormRules = {
  name: [
    { required: true, message: '请输入BOM名称', trigger: 'blur' },
    { max: 200, message: 'BOM名称不能超过200个字符', trigger: 'blur' }
  ],
  root_material_id: [{ required: true, message: '请选择根物料', trigger: 'change' }]
}

// 获取数据
async function fetchData() {
  loading.value = true
  try {
    const res = await getBOMList({
      page: pagination.page,
      page_size: pagination.pageSize,
      ...searchForm
    })
    if (res.code === 0) {
      tableData.value = res.data.list
      pagination.total = res.data.total
    }
  } catch (error: any) {
    ElMessage.error(error.message || '获取数据失败')
  } finally {
    loading.value = false
  }
}

// 搜索物料
async function searchMaterials(keyword: string) {
  if (!keyword) return
  materialLoading.value = true
  try {
    const res = await searchMaterialsApi(keyword, { page: 1, page_size: 50 })
    if (res.code === 0) {
      materialOptions.value = res.data.list
    }
  } catch (error) {
    console.error('搜索物料失败', error)
  } finally {
    materialLoading.value = false
  }
}

// 搜索
function handleSearch() {
  pagination.page = 1
  fetchData()
}

// 重置
function handleReset() {
  searchForm.name = ''
  searchForm.status = ''
  searchForm.is_exact = undefined
  handleSearch()
}

// 新增
function handleAdd() {
  isEdit.value = false
  currentId.value = null
  resetForm()
  dialogVisible.value = true
}

// 编辑
function handleEdit(row: BOMView) {
  router.push(`/bom/edit/${row.id}`)
}

// 查看结构
function handleViewTree(row: BOMView) {
  router.push(`/bom/edit/${row.id}`)
}

// 导出BOM
function handleExport(row: BOMView) {
  exportBOM(row.id)
  ElMessage.success('正在导出BOM...')
}

// 删除
async function handleDelete(row: BOMView) {
  try {
    await ElMessageBox.confirm('确定要删除该BOM吗？', '提示', {
      type: 'warning'
    })
    const res = await deleteBOMView(row.id)
    if (res.code === 0) {
      ElMessage.success('删除成功')
      fetchData()
    }
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

// 提交表单
async function handleSubmit() {
  const valid = await formRef.value?.validate()
  if (!valid) return

  submitting.value = true
  try {
    const res = await createBOMView(formData)
    if (res.code === 0) {
      ElMessage.success('创建成功')
      dialogVisible.value = false
      fetchData()
    }
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

// 重置表单
function resetForm() {
  formData.name = ''
  formData.root_material_id = 0
  formData.root_version = ''
  formData.is_exact = false
  formData.description = ''
  materialOptions.value = []
  formRef.value?.resetFields()
}

// 状态标签
function getStatusTag(status: string) {
  const map: Record<string, string> = {
    DRAFT: '',
    REVIEWING: 'warning',
    RELEASED: 'success',
    REJECTED: 'danger',
    OBSOLETE: 'info'
  }
  return map[status] || 'info'
}

function getStatusText(status: string) {
  const map: Record<string, string> = {
    DRAFT: '拟制中',
    REVIEWING: '审核中',
    RELEASED: '已发布',
    REJECTED: '已驳回',
    OBSOLETE: '已失效'
  }
  return map[status] || status
}

// 格式化日期
function formatDate(date: string) {
  if (!date) return ''
  return new Date(date).toLocaleString('zh-CN')
}

onMounted(() => {
  fetchData()
})
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
