<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>属性Schema管理</span>
          <el-button type="primary" @click="handleAdd">新增Schema</el-button>
        </div>
      </template>

      <!-- 搜索表单 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="实体类型">
          <el-select v-model="searchForm.entity_type" placeholder="全部" clearable @change="handleSearch">
            <el-option label="物料" value="MATERIAL" />
            <el-option label="文档" value="DOCUMENT" />
          </el-select>
        </el-form-item>
        <el-form-item label="类型编码">
          <el-input v-model="searchForm.type_code" placeholder="请输入类型编码" clearable />
        </el-form-item>
        <el-form-item label="属性类型">
          <el-select v-model="searchForm.attr_type" placeholder="全部" clearable @change="handleSearch">
            <el-option label="主属性" value="main" />
            <el-option label="描述属性" value="description" />
            <el-option label="规格属性" value="specification" />
            <el-option label="自定义属性" value="custom" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable @change="handleSearch">
            <el-option label="草稿" value="DRAFT" />
            <el-option label="已发布" value="RELEASED" />
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
        <el-table-column label="实体类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.entity_type === 'MATERIAL' ? 'success' : 'primary'">
              {{ row.entity_type === 'MATERIAL' ? '物料' : '文档' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="type_code" label="类型编码" width="120" />
        <el-table-column label="属性类型" width="120">
          <template #default="{ row }">
            {{ getAttrTypeLabel(row.attr_type) }}
          </template>
        </el-table-column>
        <el-table-column prop="schema_name" label="Schema名称" min-width="150" />
        <el-table-column prop="version" label="版本" width="80" />
        <el-table-column label="字段数" width="80">
          <template #default="{ row }">
            {{ row.schema_config?.fields?.length || 0 }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'RELEASED' ? 'success' : 'info'">
              {{ row.status === 'RELEASED' ? '已发布' : '草稿' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="激活" width="80">
          <template #default="{ row }">
            <el-tag v-if="row.is_active" type="success" size="small">是</el-tag>
            <el-tag v-else type="info" size="small">否</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="240" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button
              v-if="row.status === 'DRAFT'"
              type="success"
              link
              @click="handleRelease(row)"
            >发布</el-button>
            <el-button
              v-if="row.status === 'RELEASED' && !row.is_active"
              type="warning"
              link
              @click="handleActivate(row)"
            >激活</el-button>
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
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="800" destroy-on-close>
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="实体类型" prop="entity_type">
              <el-select
                v-model="formData.entity_type"
                placeholder="请选择实体类型"
                :disabled="isEdit"
                style="width: 100%"
              >
                <el-option label="物料" value="MATERIAL" />
                <el-option label="文档" value="DOCUMENT" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="类型编码" prop="type_code">
              <el-select
                v-model="formData.type_code"
                placeholder="请选择类型"
                :disabled="isEdit"
                style="width: 100%"
              >
                <template v-if="formData.entity_type === 'MATERIAL'">
                  <el-option label="零件 (PART)" value="PART" />
                  <el-option label="组件 (ASSEMBLY)" value="ASSEMBLY" />
                  <el-option label="原材料 (RAW_MATERIAL)" value="RAW_MATERIAL" />
                  <el-option label="工具 (TOOL)" value="TOOL" />
                </template>
                <template v-else-if="formData.entity_type === 'DOCUMENT'">
                  <el-option label="PDF文档" value="PDF" />
                  <el-option label="Word文档" value="DOC" />
                  <el-option label="CAD图纸" value="CAD" />
                  <el-option label="Excel表格" value="XLS" />
                  <el-option label="其他" value="OTHER" />
                </template>
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="属性类型" prop="attr_type">
              <el-select
                v-model="formData.attr_type"
                placeholder="请选择属性类型"
                :disabled="isEdit"
                style="width: 100%"
              >
                <el-option label="主属性" value="main" />
                <el-option label="描述属性" value="description" />
                <el-option label="规格属性" value="specification" />
                <el-option label="自定义属性" value="custom" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Schema名称" prop="schema_name">
              <el-input v-model="formData.schema_name" placeholder="请输入Schema名称" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="分组标签">
          <el-input v-model="formData.schema_config.label" placeholder="如：规格参数、描述信息" />
        </el-form-item>

        <!-- 字段配置 -->
        <el-divider content-position="left">字段配置</el-divider>
        <div class="field-list">
          <div v-for="(field, index) in formData.schema_config.fields" :key="index" class="field-item">
            <el-card shadow="never">
              <div class="field-header">
                <span>字段 {{ index + 1 }}</span>
                <el-button type="danger" link @click="removeField(index)">删除</el-button>
              </div>
              <el-row :gutter="12">
                <el-col :span="6">
                  <el-form-item label="键名" :prop="`schema_config.fields.${index}.key`" :rules="fieldRules.key">
                    <el-input v-model="field.key" placeholder="如：weight" />
                  </el-form-item>
                </el-col>
                <el-col :span="6">
                  <el-form-item label="标签" :prop="`schema_config.fields.${index}.label`" :rules="fieldRules.label">
                    <el-input v-model="field.label" placeholder="如：重量" />
                  </el-form-item>
                </el-col>
                <el-col :span="6">
                  <el-form-item label="类型">
                    <el-select v-model="field.type" placeholder="类型" style="width: 100%">
                      <el-option label="文本" value="text" />
                      <el-option label="数字" value="number" />
                      <el-option label="下拉选择" value="select" />
                      <el-option label="日期" value="date" />
                      <el-option label="多行文本" value="textarea" />
                    </el-select>
                  </el-form-item>
                </el-col>
                <el-col :span="6">
                  <el-form-item label="单位">
                    <el-input v-model="field.unit" placeholder="如：kg" />
                  </el-form-item>
                </el-col>
              </el-row>
              <el-row :gutter="12">
                <el-col :span="8">
                  <el-form-item label="必填">
                    <el-switch v-model="field.required" />
                  </el-form-item>
                </el-col>
                <el-col :span="16" v-if="field.type === 'select'">
                  <el-form-item label="选项">
                    <el-select
                      v-model="field.options"
                      multiple
                      filterable
                      allow-create
                      default-first-option
                      placeholder="输入后回车添加选项"
                      style="width: 100%"
                    />
                  </el-form-item>
                </el-col>
              </el-row>
            </el-card>
          </div>
          <el-button type="primary" plain @click="addField" style="width: 100%">+ 添加字段</el-button>
        </div>
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
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import {
  getAttributeSchemaList,
  createAttributeSchema,
  updateAttributeSchema,
  deleteAttributeSchema,
  releaseAttributeSchema,
  activateAttributeSchema,
  type AttributeSchema,
  type CreateAttributeSchemaRequest,
  type UpdateAttributeSchemaRequest,
  type FieldSchema
} from '@/api/attribute'

// 加载状态
const loading = ref(false)
const submitting = ref(false)

// 数据
const tableData = ref<AttributeSchema[]>([])

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 搜索表单
const searchForm = reactive({
  entity_type: '',
  type_code: '',
  attr_type: '',
  status: ''
})

// 对话框
const dialogVisible = ref(false)
const isEdit = ref(false)
const currentId = ref<number | null>(null)
const formRef = ref<FormInstance>()

const dialogTitle = computed(() => (isEdit.value ? '编辑Schema' : '新增Schema'))

// 默认字段
const createEmptyField = (): FieldSchema => ({
  key: '',
  label: '',
  type: 'text',
  unit: '',
  required: false,
  options: []
})

// 表单数据
const formData = reactive<CreateAttributeSchemaRequest>({
  entity_type: 'MATERIAL',
  type_code: '',
  attr_type: 'main',
  schema_name: '',
  schema_config: {
    label: '',
    fields: [createEmptyField()]
  }
})

// 表单规则
const rules: FormRules = {
  entity_type: [{ required: true, message: '请选择实体类型', trigger: 'change' }],
  type_code: [{ required: true, message: '请选择类型编码', trigger: 'change' }],
  attr_type: [{ required: true, message: '请选择属性类型', trigger: 'change' }]
}

// 字段规则
const fieldRules = {
  key: [{ required: true, message: '请输入键名', trigger: 'blur' }],
  label: [{ required: true, message: '请输入标签', trigger: 'blur' }]
}

// 获取数据
async function fetchData() {
  loading.value = true
  try {
    const res = await getAttributeSchemaList({
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

// 搜索
function handleSearch() {
  pagination.page = 1
  fetchData()
}

// 重置
function handleReset() {
  searchForm.entity_type = ''
  searchForm.type_code = ''
  searchForm.attr_type = ''
  searchForm.status = ''
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
function handleEdit(row: AttributeSchema) {
  isEdit.value = true
  currentId.value = row.id
  formData.entity_type = row.entity_type
  formData.type_code = row.type_code
  formData.attr_type = row.attr_type
  formData.schema_name = row.schema_name || ''
  formData.schema_config = {
    label: row.schema_config?.label || '',
    fields: row.schema_config?.fields?.length
      ? [...row.schema_config.fields]
      : [createEmptyField()]
  }
  dialogVisible.value = true
}

// 删除
async function handleDelete(row: AttributeSchema) {
  try {
    await ElMessageBox.confirm('确定要删除该Schema吗？', '提示', { type: 'warning' })
    const res = await deleteAttributeSchema(row.id)
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

// 发布
async function handleRelease(row: AttributeSchema) {
  try {
    await ElMessageBox.confirm('确定要发布该Schema吗？发布后将成为正式版本。', '提示', { type: 'warning' })
    const res = await releaseAttributeSchema(row.id)
    if (res.code === 0) {
      ElMessage.success('发布成功')
      fetchData()
    }
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '发布失败')
    }
  }
}

// 激活
async function handleActivate(row: AttributeSchema) {
  try {
    await ElMessageBox.confirm('确定要激活该版本吗？激活后将替换当前使用的版本。', '提示', { type: 'warning' })
    const res = await activateAttributeSchema(row.id)
    if (res.code === 0) {
      ElMessage.success('激活成功')
      fetchData()
    }
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '激活失败')
    }
  }
}

// 添加字段
function addField() {
  formData.schema_config.fields.push(createEmptyField())
}

// 删除字段
function removeField(index: number) {
  formData.schema_config.fields.splice(index, 1)
}

// 提交表单
async function handleSubmit() {
  const valid = await formRef.value?.validate()
  if (!valid) return

  // 过滤空字段
  const fields = formData.schema_config.fields.filter(f => f.key && f.label)

  submitting.value = true
  try {
    if (isEdit.value && currentId.value) {
      // 编辑
      const data: UpdateAttributeSchemaRequest = {
        schema_name: formData.schema_name,
        schema_config: {
          label: formData.schema_config.label,
          fields
        }
      }
      const res = await updateAttributeSchema(currentId.value, data)
      if (res.code === 0) {
        ElMessage.success('更新成功')
        dialogVisible.value = false
        fetchData()
      }
    } else {
      // 新增
      const data: CreateAttributeSchemaRequest = {
        entity_type: formData.entity_type,
        type_code: formData.type_code,
        attr_type: formData.attr_type,
        schema_name: formData.schema_name,
        schema_config: {
          label: formData.schema_config.label,
          fields
        }
      }
      const res = await createAttributeSchema(data)
      if (res.code === 0) {
        ElMessage.success('创建成功')
        dialogVisible.value = false
        fetchData()
      }
    }
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

// 重置表单
function resetForm() {
  formData.entity_type = 'MATERIAL'
  formData.type_code = ''
  formData.attr_type = 'main'
  formData.schema_name = ''
  formData.schema_config = {
    label: '',
    fields: [createEmptyField()]
  }
  formRef.value?.resetFields()
}

// 获取属性类型标签
function getAttrTypeLabel(type: string) {
  const map: Record<string, string> = {
    main: '主属性',
    description: '描述属性',
    specification: '规格属性',
    custom: '自定义属性'
  }
  return map[type] || type
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

  .field-list {
    .field-item {
      margin-bottom: 12px;

      .field-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 12px;
        font-weight: 500;
      }
    }
  }
}
</style>
