<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>PART查询</span>
          <div class="header-actions">
            <el-switch
              v-model="showAttributes"
              active-text="显示属性"
              inactive-text=""
              class="attr-switch"
            />
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
        <el-form-item label="类型">
          <el-select v-model="searchForm.item_type" placeholder="请选择类型" clearable>
            <el-option label="零件" value="PART" />
            <el-option label="组件" value="ASSEMBLY" />
            <el-option label="原材料" value="RAW_MATERIAL" />
            <el-option label="工具" value="TOOL" />
          </el-select>
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
        <el-table-column prop="item_id" label="PART编码" width="150" />
        <el-table-column prop="item_name" label="PART名称" width="200" />
        <el-table-column prop="item_type" label="类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getItemTypeTag(row.item_type)">{{ getItemTypeText(row.item_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="version" label="版本" width="80" />
        <el-table-column prop="unit" label="单位" width="80" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <!-- 动态属性列 -->
        <el-table-column
          v-for="col in (showAttributes ? attributeColumns : [])"
          :key="col.prop"
          :prop="col.prop"
          :label="col.label"
          :width="col.width"
        >
          <template #default="{ row }">
            {{ getAttributeValue(row, col.prop) }}
          </template>
        </el-table-column>
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" fixed="right" width="250">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleView(row)">查看</el-button>
            <el-button type="primary" link @click="handleEdit(row)" v-if="canEdit(row.status)">编辑</el-button>
            <el-button type="warning" link @click="handleInitiateWorkflow(row)" v-if="canInitiate(row.status)">发起审批</el-button>
            <el-button type="danger" link @click="handleDelete(row)" v-if="canDelete(row.status)">删除</el-button>
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

    <!-- 新增/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="800" destroy-on-close>
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="PART编码" prop="item_id">
              <el-input v-model="formData.item_id" placeholder="请输入PART编码" :disabled="isEdit" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="PART名称" prop="item_name">
              <el-input v-model="formData.item_name" placeholder="请输入PART名称" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="类型" prop="item_type">
              <el-select v-model="formData.item_type" placeholder="请选择类型" style="width: 100%" @change="handleTypeChange">
                <el-option label="零件" value="PART" />
                <el-option label="组件" value="ASSEMBLY" />
                <el-option label="原材料" value="RAW_MATERIAL" />
                <el-option label="工具" value="TOOL" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="计量单位" prop="unit">
              <el-input v-model="formData.unit" placeholder="请输入计量单位" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="描述" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" placeholder="请输入描述" />
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
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" title="PART详情" width="800">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="PART编码">{{ currentMaterial?.item_id }}</el-descriptions-item>
        <el-descriptions-item label="PART名称">{{ currentMaterial?.item_name }}</el-descriptions-item>
        <el-descriptions-item label="类型">{{ getItemTypeText(currentMaterial?.item_type) }}</el-descriptions-item>
        <el-descriptions-item label="版本">{{ currentMaterial?.version }}</el-descriptions-item>
        <el-descriptions-item label="计量单位">{{ currentMaterial?.unit }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(currentMaterial?.status)">{{ getStatusText(currentMaterial?.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ currentMaterial?.description }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ currentMaterial?.created_at }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ currentMaterial?.updated_at }}</el-descriptions-item>
      </el-descriptions>

      <!-- 动态属性展示 -->
      <template v-if="detailSchemas.length > 0">
        <el-divider content-position="left">扩展属性</el-divider>
        <div v-for="schema in detailSchemas" :key="schema.attr_type">
          <el-divider v-if="schema.schema_config?.label" content-position="left" style="margin: 12px 0">
            {{ schema.schema_config.label }}
          </el-divider>
          <el-descriptions :column="2" border size="small">
            <el-descriptions-item
              v-for="field in schema.schema_config?.fields"
              :key="field.key"
              :label="field.label"
            >
              {{ getDetailAttributeValue(schema.attr_type, field.key, field.unit) }}
            </el-descriptions-item>
          </el-descriptions>
        </div>
      </template>

      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 发起审批对话框 -->
    <el-dialog v-model="workflowVisible" title="发起审批" width="500">
      <el-form ref="workflowFormRef" :model="workflowForm" :rules="workflowRules" label-width="80px">
        <el-form-item label="PART编码">
          <el-input :value="currentMaterial?.item_id" disabled />
        </el-form-item>
        <el-form-item label="PART名称">
          <el-input :value="currentMaterial?.item_name" disabled />
        </el-form-item>
        <el-form-item label="审批流程" prop="definition_id">
          <el-select v-model="workflowForm.definition_id" placeholder="请选择审批流程" style="width: 100%">
            <el-option
              v-for="def in workflowDefinitionList"
              :key="def.id"
              :label="def.name"
              :value="def.id"
            >
              <div style="display: flex; justify-content: space-between;">
                <span>{{ def.name }}</span>
                <span style="color: #909399; font-size: 12px;">{{ def.description }}</span>
              </div>
            </el-option>
          </el-select>
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
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import {
  getMaterialList,
  getMaterial,
  createMaterial,
  updateMaterial,
  deleteMaterial,
  type Material,
  type MaterialDetail,
  type CreateMaterialRequest,
  type UpdateMaterialRequest,
} from '@/api/material'
import { getAllActiveSchemas, type AttributeSchema } from '@/api/attribute'
import AttributeForm from '@/components/AttributeForm.vue'
import { getDefinitionList, initiateWorkflow, type WorkflowDefinition } from '@/api/workflow'

defineOptions({
  name: 'PartList',
})

const loading = ref(false)
const submitting = ref(false)
const tableData = ref<Material[]>([])
const dialogVisible = ref(false)
const detailVisible = ref(false)
const isEdit = ref(false)
const currentId = ref<number | null>(null)
const currentMaterial = ref<MaterialDetail | null>(null)
const formRef = ref<FormInstance>()
const hasSearched = ref(false)  // 标记是否已搜索过

// 属性相关
const showAttributes = ref(false)
const currentSchemas = ref<AttributeSchema[]>([])
const detailSchemas = ref<AttributeSchema[]>([])
const attributeFormData = ref<Record<string, any>>({})
const allSchemas = ref<AttributeSchema[]>([])

// 发起审批相关
const workflowVisible = ref(false)
const workflowSubmitting = ref(false)
const workflowFormRef = ref<FormInstance>()
const workflowDefinitionList = ref<WorkflowDefinition[]>([])
const workflowForm = reactive({
  definition_id: null as number | null,
  title: '',
  comment: '',
})

const workflowRules: FormRules = {
  definition_id: [{ required: true, message: '请选择审批流程', trigger: 'change' }],
  title: [{ required: true, message: '请输入流程标题', trigger: 'blur' }],
}

const searchForm = reactive({
  item_id: '',
  item_name: '',
  item_type: '',
  status: '',
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0,
})

const dialogTitle = computed(() => (isEdit.value ? '编辑PART' : '新增PART'))

const formData = reactive<CreateMaterialRequest & UpdateMaterialRequest>({
  item_id: '',
  item_name: '',
  description: '',
  item_type: '',
  unit: '',
})

const rules: FormRules = {
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

// 计算属性列
const attributeColumns = computed(() => {
  const columns: Array<{ prop: string; label: string; width?: number }> = []
  allSchemas.value.forEach(schema => {
    if (schema.attr_type === 'main' && schema.schema_config?.fields) {
      schema.schema_config.fields.slice(0, 3).forEach(field => {
        columns.push({
          prop: `main_${field.key}`,
          label: field.label,
          width: 120
        })
      })
    }
  })
  return columns
})

onMounted(() => {
  // 默认不查询，等用户点击搜索
  loadAllSchemas()
})

async function loadAllSchemas() {
  // 加载所有物料类型的Schema用于列表显示
  const types = ['PART', 'ASSEMBLY', 'RAW_MATERIAL', 'TOOL']
  for (const type of types) {
    try {
      const res = await getAllActiveSchemas('MATERIAL', type)
      if (res.code === 0 && res.data) {
        allSchemas.value.push(...res.data)
      }
    } catch (error) {
      console.error(error)
    }
  }
}

async function loadSchemasForType(typeCode: string) {
  try {
    const res = await getAllActiveSchemas('MATERIAL', typeCode)
    if (res.code === 0) {
      currentSchemas.value = res.data || []
    }
  } catch (error) {
    console.error(error)
    currentSchemas.value = []
  }
}

async function fetchData() {
  loading.value = true
  try {
    const res = await getMaterialList({
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
  searchForm.item_id = ''
  searchForm.item_name = ''
  searchForm.item_type = ''
  searchForm.status = ''
  pagination.page = 1
  hasSearched.value = false
  tableData.value = []
  pagination.total = 0
}

function resetForm() {
  formData.item_id = ''
  formData.item_name = ''
  formData.description = ''
  formData.item_type = ''
  formData.unit = ''
  attributeFormData.value = {}
  currentId.value = null
}

function handleAdd() {
  isEdit.value = false
  resetForm()
  currentSchemas.value = []
  dialogVisible.value = true
}

async function handleTypeChange(typeCode: string) {
  if (typeCode) {
    await loadSchemasForType(typeCode)
  } else {
    currentSchemas.value = []
  }
}

async function handleView(row: Material) {
  try {
    const res = await getMaterial(row.id)
    if (res.code === 0) {
      currentMaterial.value = res.data
      // 加载详情页的Schema
      if (res.data.item_type) {
        const schemaRes = await getAllActiveSchemas('MATERIAL', res.data.item_type)
        if (schemaRes.code === 0) {
          detailSchemas.value = schemaRes.data || []
        }
      }
      detailVisible.value = true
    }
  } catch (error) {
    console.error(error)
  }
}

async function handleEdit(row: Material) {
  isEdit.value = true
  currentId.value = row.id
  formData.item_id = row.item_id
  formData.item_name = row.item_name
  formData.description = row.description
  formData.item_type = row.item_type
  formData.unit = row.unit

  // 加载Schema
  await loadSchemasForType(row.item_type)

  // 加载物料详情获取属性
  try {
    const res = await getMaterial(row.id)
    if (res.code === 0) {
      const detail = res.data
      attributeFormData.value = {}

      // 解析主属性
      if (detail.attributes) {
        const mainAttrs = typeof detail.attributes === 'string' ? JSON.parse(detail.attributes) : detail.attributes
        Object.entries(mainAttrs).forEach(([key, value]) => {
          attributeFormData.value[`main_${key}`] = value
        })
      }

      // 解析动态属性
      if (detail.dynamicAttributes && detail.dynamicAttributes.length > 0) {
        detail.dynamicAttributes.forEach(attr => {
          attributeFormData.value[`${attr.attr_type}_${attr.attr_key}`] = attr.attr_value
        })
      }
    }
  } catch (error) {
    console.error(error)
  }

  dialogVisible.value = true
}

async function handleSubmit() {
  const valid = await formRef.value?.validate()
  if (!valid) return

  submitting.value = true
  try {
    // 构建属性数据
    const mainAttributes: Record<string, any> = {}
    const dynamicAttributes: Array<{ attr_type: string; attr_key: string; attr_value: string; unit?: string }> = []

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

    if (isEdit.value) {
      const updateData: UpdateMaterialRequest = {
        item_name: formData.item_name,
        description: formData.description,
        item_type: formData.item_type,
        unit: formData.unit,
        attributes: JSON.stringify(mainAttributes),
        dynamic_attributes: dynamicAttributes.length > 0 ? dynamicAttributes : undefined,
      }
      await updateMaterial(currentId.value!, updateData)
      ElMessage.success('更新成功')
    } else {
      const createData: CreateMaterialRequest = {
        item_id: formData.item_id,
        item_name: formData.item_name,
        description: formData.description,
        item_type: formData.item_type,
        unit: formData.unit,
        attributes: JSON.stringify(mainAttributes),
        dynamic_attributes: dynamicAttributes.length > 0 ? dynamicAttributes : undefined,
      }
      await createMaterial(createData)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchData()
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

async function handleDelete(row: Material) {
  try {
    await ElMessageBox.confirm(`确定要删除PART ${row.item_name} 吗？`, '提示', {
      type: 'warning',
    })

    await deleteMaterial(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch {
    // 取消删除
  }
}

function canEdit(status: string) {
  return status === 'DRAFT' || status === 'REJECTED'
}

function canDelete(status: string) {
  return status === 'DRAFT' || status === 'REJECTED'
}

function canInitiate(status: string) {
  return status === 'DRAFT' || status === 'REJECTED'
}

// 发起审批
async function handleInitiateWorkflow(row: Material) {
  try {
    // 获取物料详情
    const res = await getMaterial(row.id)
    if (res.code === 0) {
      currentMaterial.value = res.data
    }

    // 获取物料审批流程定义列表
    const defRes = await getDefinitionList({
      type: 'MATERIAL_APPROVAL',
      status: 'RELEASED',
      page: 1,
      page_size: 100,
    })
    if (defRes.code === 0 && defRes.data.list.length > 0) {
      workflowDefinitionList.value = defRes.data.list
    } else {
      ElMessage.error('未找到PART审批流程，请先配置流程')
      return
    }

    // 设置默认值
    workflowForm.definition_id = null
    workflowForm.title = `PART审批 - ${row.item_name}`
    workflowForm.comment = ''
    workflowVisible.value = true
  } catch (error: any) {
    ElMessage.error(error.message || '获取流程信息失败')
  }
}

async function handleSubmitWorkflow() {
  const valid = await workflowFormRef.value?.validate()
  if (!valid) return

  if (!workflowForm.definition_id || !currentMaterial.value) return

  workflowSubmitting.value = true
  try {
    await initiateWorkflow({
      definition_id: workflowForm.definition_id,
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

// 获取列表中的属性值
function getAttributeValue(row: Material, prop: string) {
  if (prop.startsWith('main_')) {
    const key = prop.substring(5)
    if (row.attributes) {
      const attrs = typeof row.attributes === 'string' ? JSON.parse(row.attributes) : row.attributes
      return attrs[key] || '-'
    }
  }
  return '-'
}

// 获取详情中的属性值
function getDetailAttributeValue(attrType: string, key: string, unit?: string) {
  const material = currentMaterial.value
  if (!material) return '-'

  if (attrType === 'main') {
    if (material.attributes) {
      const attrs = typeof material.attributes === 'string' ? JSON.parse(material.attributes) : material.attributes
      const value = attrs[key]
      return value ? (unit ? `${value} ${unit}` : value) : '-'
    }
  } else {
    const attr = material.dynamicAttributes?.find(a => a.attr_type === attrType && a.attr_key === key)
    if (attr) {
      return attr.unit ? `${attr.attr_value} ${attr.unit}` : attr.attr_value
    }
  }
  return '-'
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
