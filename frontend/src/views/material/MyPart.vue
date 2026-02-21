<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>我的PART</span>
          <div class="header-actions">
            <el-switch
              v-model="showAttributes"
              active-text="显示属性"
              inactive-text=""
              class="attr-switch"
            />
            <el-button type="primary" @click="handleAdd">新增PART</el-button>
            <el-button type="warning" @click="handleBatchApprove" :disabled="selectedRows.length === 0">
              批量发起审批 ({{ selectedRows.length }})
            </el-button>
          </div>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="PART编码">
          <el-input v-model="searchForm.item_id" placeholder="请输入PART编码" clearable />
        </el-form-item>
        <el-form-item label="PART名称">
          <el-input v-model="searchForm.item_name" placeholder="请输入PART名称" clearable />
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

      <!-- 表格 -->
      <el-table
        :data="tableData"
        v-loading="loading"
        stripe
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="item_id" label="PART编码" width="150" />
        <el-table-column prop="item_name" label="PART名称" width="200" />
        <el-table-column prop="item_type" label="类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getItemTypeTag(row.item_type)">{{ getItemTypeText(row.item_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="version" label="版本" width="80" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" fixed="right" width="150">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleView(row)">查看</el-button>
            <el-button type="primary" link @click="handleEdit(row)" v-if="canEdit(row.status)">编辑</el-button>
            <el-button type="warning" link @click="handleInitiateWorkflow(row)" v-if="canInitiate(row.status)">发起审批</el-button>
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
    <el-dialog v-model="detailVisible" title="PART详情" width="700">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="PART编码">{{ currentMaterial?.item_id }}</el-descriptions-item>
        <el-descriptions-item label="PART名称">{{ currentMaterial?.item_name }}</el-descriptions-item>
        <el-descriptions-item label="类型">{{ getItemTypeText(currentMaterial?.item_type) }}</el-descriptions-item>
        <el-descriptions-item label="版本">{{ currentMaterial?.version }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(currentMaterial?.status)">{{ getStatusText(currentMaterial?.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ currentMaterial?.created_at }}</el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ currentMaterial?.description || '-' }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 编辑对话框 -->
    <el-dialog v-model="editVisible" title="编辑PART" width="700" destroy-on-close>
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="80px">
        <el-form-item label="PART编码">
          <el-input :value="formData.item_id" disabled />
        </el-form-item>
        <el-form-item label="PART名称" prop="item_name">
          <el-input v-model="formData.item_name" placeholder="请输入PART名称" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="请输入描述" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitEdit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 新增对话框 -->
    <el-dialog v-model="addVisible" title="新增PART" width="800" destroy-on-close>
      <el-form ref="addFormRef" :model="addFormData" :rules="addRules" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="PART编码" prop="item_id">
              <el-input v-model="addFormData.item_id" placeholder="请输入PART编码" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="PART名称" prop="item_name">
              <el-input v-model="addFormData.item_name" placeholder="请输入PART名称" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="类型" prop="item_type">
              <el-select v-model="addFormData.item_type" placeholder="请选择类型" style="width: 100%" @change="handleTypeChange">
                <el-option label="零件" value="PART" />
                <el-option label="组件" value="ASSEMBLY" />
                <el-option label="原材料" value="RAW_MATERIAL" />
                <el-option label="工具" value="TOOL" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="计量单位" prop="unit">
              <el-input v-model="addFormData.unit" placeholder="请输入计量单位" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="描述" prop="description">
          <el-input v-model="addFormData.description" type="textarea" :rows="3" placeholder="请输入描述" />
        </el-form-item>

        <!-- 动态属性表单 -->
        <template v-if="currentSchemas.length > 0">
          <AttributeForm
            :schemas="currentSchemas"
            v-model="attributeFormData"
            entity-type="MATERIAL"
          />
        </template>
      </el-form>
      <template #footer>
        <el-button @click="addVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitAdd" :loading="addSubmitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 批量发起审批对话框 -->
    <el-dialog v-model="batchApproveVisible" title="批量发起审批" width="600">
      <el-alert type="info" :closable="false" style="margin-bottom: 16px">
        <template #title>已选择 {{ selectedRows.length }} 个PART</template>
      </el-alert>
      <el-table :data="selectedRows" max-height="300" size="small">
        <el-table-column prop="item_id" label="PART编码" width="150" />
        <el-table-column prop="item_name" label="PART名称" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
      </el-table>
      <el-form ref="batchFormRef" :model="batchForm" :rules="batchRules" label-width="80px" style="margin-top: 16px">
        <el-form-item label="流程标题" prop="title">
          <el-input v-model="batchForm.title" placeholder="请输入流程标题" />
        </el-form-item>
        <el-form-item label="审批意见">
          <el-input v-model="batchForm.comment" type="textarea" :rows="3" placeholder="请输入审批意见（可选）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="batchApproveVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitBatchApprove" :loading="batchSubmitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 单个发起审批对话框 -->
    <el-dialog v-model="workflowVisible" title="发起审批" width="500">
      <el-form ref="workflowFormRef" :model="workflowForm" :rules="workflowRules" label-width="80px">
        <el-form-item label="PART编码">
          <el-input :value="currentMaterial?.item_id" disabled />
        </el-form-item>
        <el-form-item label="PART名称">
          <el-input :value="currentMaterial?.item_name" disabled />
        </el-form-item>
        <el-form-item label="流程标题" prop="title">
          <el-input v-model="workflowForm.title" placeholder="请输入流程标题" />
        </el-form-item>
        <el-form-item label="审批意见">
          <el-input v-model="workflowForm.comment" type="textarea" :rows="3" placeholder="请输入审批意见（可选）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="workflowVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitWorkflow" :loading="workflowSubmitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import {
  getMaterialList,
  getMaterial,
  createMaterial,
  updateMaterial,
  type Material,
  type MaterialDetail,
  type CreateMaterialRequest,
  type UpdateMaterialRequest,
} from '@/api/material'
import { getAllActiveSchemas, type AttributeSchema } from '@/api/attribute'
import { getActiveDefinition, initiateWorkflow, type WorkflowDefinition } from '@/api/workflow'
import { useUserStore } from '@/stores/user'
import AttributeForm from '@/components/AttributeForm.vue'

defineOptions({
  name: 'MyPart',
})

const userStore = useUserStore()
const loading = ref(false)
const submitting = ref(false)
const addSubmitting = ref(false)
const workflowSubmitting = ref(false)
const batchSubmitting = ref(false)
const tableData = ref<Material[]>([])
const selectedRows = ref<Material[]>([])
const detailVisible = ref(false)
const editVisible = ref(false)
const addVisible = ref(false)
const workflowVisible = ref(false)
const batchApproveVisible = ref(false)
const currentMaterial = ref<MaterialDetail | null>(null)
const formRef = ref<FormInstance>()
const addFormRef = ref<FormInstance>()
const workflowFormRef = ref<FormInstance>()
const batchFormRef = ref<FormInstance>()

// 属性相关
const showAttributes = ref(false)
const currentSchemas = ref<AttributeSchema[]>([])
const attributeFormData = ref<Record<string, any>>({})

// 流程定义
const workflowDefinition = ref<WorkflowDefinition | null>(null)

const searchForm = reactive({
  item_id: '',
  item_name: '',
  status: 'DRAFT',  // 默认显示拟制中
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0,
})

const formData = reactive<UpdateMaterialRequest & { item_id: string }>({
  item_id: '',
  item_name: '',
  description: '',
})

const rules: FormRules = {
  item_name: [
    { required: true, message: '请输入PART名称', trigger: 'blur' },
    { max: 200, message: 'PART名称不能超过200个字符', trigger: 'blur' },
  ],
}

// 新增表单
const addFormData = reactive<CreateMaterialRequest>({
  item_id: '',
  item_name: '',
  description: '',
  item_type: '',
  unit: '',
})

const addRules: FormRules = {
  item_id: [
    { required: true, message: '请输入PART编码', trigger: 'blur' },
    { max: 50, message: 'PART编码不能超过50个字符', trigger: 'blur' },
  ],
  item_name: [
    { required: true, message: '请输入PART名称', trigger: 'blur' },
    { max: 200, message: 'PART名称不能超过200个字符', trigger: 'blur' },
  ],
  item_type: [{ required: true, message: '请选择类型', trigger: 'change' }],
}

const workflowForm = reactive({
  title: '',
  comment: '',
})

const workflowRules: FormRules = {
  title: [{ required: true, message: '请输入流程标题', trigger: 'blur' }],
}

const batchForm = reactive({
  title: '',
  comment: '',
})

const batchRules: FormRules = {
  title: [{ required: true, message: '请输入流程标题', trigger: 'blur' }],
}

onMounted(() => {
  loadWorkflowDefinition()
  fetchData()
})

async function loadWorkflowDefinition() {
  try {
    const res = await getActiveDefinition('MATERIAL_APPROVAL')
    if (res.code === 0) {
      workflowDefinition.value = res.data
    }
  } catch (error) {
    console.error(error)
  }
}

async function fetchData() {
  loading.value = true
  try {
    const res = await getMaterialList({
      page: pagination.page,
      page_size: pagination.pageSize,
      created_by: userStore.userInfo?.id,  // 只查询我创建的
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
  searchForm.item_id = ''
  searchForm.item_name = ''
  searchForm.status = 'DRAFT'
  pagination.page = 1
  fetchData()
}

function handleSelectionChange(rows: Material[]) {
  selectedRows.value = rows
}

async function handleView(row: Material) {
  try {
    const res = await getMaterial(row.id)
    if (res.code === 0) {
      currentMaterial.value = res.data
      detailVisible.value = true
    }
  } catch (error) {
    console.error(error)
  }
}

async function handleEdit(row: Material) {
  try {
    const res = await getMaterial(row.id)
    if (res.code === 0) {
      currentMaterial.value = res.data
      formData.item_id = res.data.item_id
      formData.item_name = res.data.item_name
      formData.description = res.data.description
      editVisible.value = true
    }
  } catch (error) {
    console.error(error)
  }
}

async function handleSubmitEdit() {
  const valid = await formRef.value?.validate()
  if (!valid) return

  if (!currentMaterial.value) return

  submitting.value = true
  try {
    await updateMaterial(currentMaterial.value.id, {
      item_name: formData.item_name,
      description: formData.description,
    })
    ElMessage.success('更新成功')
    editVisible.value = false
    fetchData()
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

// 新增PART
function handleAdd() {
  addFormData.item_id = ''
  addFormData.item_name = ''
  addFormData.description = ''
  addFormData.item_type = ''
  addFormData.unit = ''
  attributeFormData.value = {}
  currentSchemas.value = []
  addVisible.value = true
}

async function handleTypeChange(typeCode: string) {
  if (typeCode) {
    try {
      const res = await getAllActiveSchemas('MATERIAL', typeCode)
      if (res.code === 0) {
        currentSchemas.value = res.data || []
      }
    } catch (error) {
      console.error(error)
      currentSchemas.value = []
    }
  } else {
    currentSchemas.value = []
  }
}

async function handleSubmitAdd() {
  const valid = await addFormRef.value?.validate()
  if (!valid) return

  addSubmitting.value = true
  try {
    // 构建属性数据
    const mainAttributes: Record<string, any> = {}
    const dynamicAttributes: Array<{ attr_type: string; attr_key: string; attr_value: string }> = []

    Object.entries(attributeFormData.value).forEach(([key, value]) => {
      if (value === undefined || value === null || value === '') return

      if (key.startsWith('main_')) {
        const attrKey = key.substring(5)
        mainAttributes[attrKey] = value
      } else {
        const underscoreIndex = key.indexOf('_')
        const attrType = key.substring(0, underscoreIndex)
        const attrKey = key.substring(underscoreIndex + 1)
        if (['description', 'specification', 'custom'].includes(attrType)) {
          dynamicAttributes.push({
            attr_type: attrType,
            attr_key: attrKey,
            attr_value: String(value)
          })
        }
      }
    })

    const createData: CreateMaterialRequest = {
      item_id: addFormData.item_id,
      item_name: addFormData.item_name,
      description: addFormData.description,
      item_type: addFormData.item_type,
      unit: addFormData.unit,
      attributes: JSON.stringify(mainAttributes),
      dynamic_attributes: dynamicAttributes.length > 0 ? dynamicAttributes : undefined,
    }

    await createMaterial(createData)
    ElMessage.success('创建成功')
    addVisible.value = false
    fetchData()
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    addSubmitting.value = false
  }
}

async function handleInitiateWorkflow(row: Material) {
  if (!workflowDefinition.value) {
    ElMessage.error('未找到PART审批流程，请先配置流程')
    return
  }

  try {
    const res = await getMaterial(row.id)
    if (res.code === 0) {
      currentMaterial.value = res.data
      workflowForm.title = `PART审批 - ${row.item_name}`
      workflowForm.comment = ''
      workflowVisible.value = true
    }
  } catch (error: any) {
    ElMessage.error(error.message || '获取信息失败')
  }
}

async function handleSubmitWorkflow() {
  const valid = await workflowFormRef.value?.validate()
  if (!valid) return

  if (!workflowDefinition.value || !currentMaterial.value) return

  workflowSubmitting.value = true
  try {
    await initiateWorkflow({
      definition_id: workflowDefinition.value.id,
      business_type: 'MATERIAL',
      business_id: currentMaterial.value.id,
      title: workflowForm.title,
      comment: workflowForm.comment,
    })
    ElMessage.success('发起审批成功')
    workflowVisible.value = false
    fetchData()
  } catch (error: any) {
    ElMessage.error(error.message || '发起审批失败')
  } finally {
    workflowSubmitting.value = false
  }
}

function handleBatchApprove() {
  if (selectedRows.value.length === 0) {
    ElMessage.warning('请选择要审批的PART')
    return
  }

  if (!workflowDefinition.value) {
    ElMessage.error('未找到PART审批流程，请先配置流程')
    return
  }

  // 检查是否有不可发起审批的PART
  const invalidItems = selectedRows.value.filter(item => !canInitiate(item.status))
  if (invalidItems.length > 0) {
    ElMessage.warning('只有拟制中或已驳回状态的PART可以发起审批')
    return
  }

  batchForm.title = `PART批量审批 - ${selectedRows.value.length}个`
  batchForm.comment = ''
  batchApproveVisible.value = true
}

async function handleSubmitBatchApprove() {
  const valid = await batchFormRef.value?.validate()
  if (!valid) return

  if (!workflowDefinition.value) return

  batchSubmitting.value = true
  let successCount = 0
  let failCount = 0

  for (const item of selectedRows.value) {
    try {
      await initiateWorkflow({
        definition_id: workflowDefinition.value.id,
        business_type: 'MATERIAL',
        business_id: item.id,
        title: batchForm.title,
        comment: batchForm.comment,
      })
      successCount++
    } catch {
      failCount++
    }
  }

  batchSubmitting.value = false
  batchApproveVisible.value = false

  if (successCount > 0) {
    ElMessage.success(`成功发起 ${successCount} 个审批`)
  }
  if (failCount > 0) {
    ElMessage.warning(`${failCount} 个PART发起审批失败`)
  }

  selectedRows.value = []
  fetchData()
}

function canEdit(status: string) {
  return status === 'DRAFT' || status === 'REJECTED'
}

function canInitiate(status: string) {
  return status === 'DRAFT' || status === 'REJECTED'
}

function getItemTypeText(type?: string) {
  const map: Record<string, string> = {
    PART: '零件',
    ASSEMBLY: '组件',
    RAW_MATERIAL: '原材料',
    TOOL: '工具',
  }
  return type ? map[type] || type : ''
}

function getItemTypeTag(type: string) {
  const map: Record<string, string> = {
    PART: '',
    ASSEMBLY: 'success',
    RAW_MATERIAL: 'warning',
    TOOL: 'info',
  }
  return map[type] || 'info'
}

function getStatusType(status?: string) {
  const map: Record<string, string> = {
    DRAFT: 'info',
    REVIEWING: 'warning',
    RELEASED: 'success',
    REJECTED: 'danger',
    OBSOLETE: '',
  }
  return status ? map[status] || 'info' : 'info'
}

function getStatusText(status?: string) {
  const map: Record<string, string> = {
    DRAFT: '拟制中',
    REVIEWING: '审核中',
    RELEASED: '已发布',
    REJECTED: '已驳回',
    OBSOLETE: '已失效',
  }
  return status ? map[status] || status : ''
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
