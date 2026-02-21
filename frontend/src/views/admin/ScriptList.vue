<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>脚本管理</span>
          <el-button type="primary" @click="handleAdd">新增脚本</el-button>
        </div>
      </template>

      <!-- 搜索表单 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="脚本类型">
          <el-select v-model="searchForm.type" placeholder="请选择" clearable>
            <el-option label="JavaScript" value="JAVASCRIPT" />
            <el-option label="SQL" value="SQL" />
          </el-select>
        </el-form-item>
        <el-form-item label="触发时机">
          <el-select v-model="searchForm.trigger_type" placeholder="请选择" clearable>
            <el-option label="发起前" value="BEFORE_SUBMIT" />
            <el-option label="发起后" value="AFTER_SUBMIT" />
            <el-option label="审批前" value="BEFORE_APPROVE" />
            <el-option label="审批后" value="AFTER_APPROVE" />
            <el-option label="驳回前" value="BEFORE_REJECT" />
            <el-option label="驳回后" value="AFTER_REJECT" />
          </el-select>
        </el-form-item>
        <el-form-item label="业务类型">
          <el-select v-model="searchForm.business_type" placeholder="请选择" clearable>
            <el-option label="物料" value="MATERIAL" />
            <el-option label="文档" value="DOCUMENT" />
            <el-option label="BOM" value="BOM" />
            <el-option label="全部" value="ALL" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择" clearable>
            <el-option label="草稿" value="DRAFT" />
            <el-option label="已发布" value="RELEASED" />
            <el-option label="已废弃" value="OBSOLETE" />
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
        <el-table-column prop="name" label="脚本名称" min-width="150" />
        <el-table-column prop="code" label="脚本编码" width="150" />
        <el-table-column label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.type === 'JAVASCRIPT' ? 'warning' : 'success'">
              {{ row.type }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="触发时机" width="100">
          <template #default="{ row }">
            {{ getTriggerText(row.trigger_type) }}
          </template>
        </el-table-column>
        <el-table-column label="业务类型" width="100">
          <template #default="{ row }">
            {{ getBusinessText(row.business_type) }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusTag(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="timeout" label="超时(ms)" width="100" />
        <el-table-column label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="primary" link @click="handleTest(row)">测试</el-button>
            <el-button type="info" link @click="handleViewLogs(row)">日志</el-button>
            <el-button
              v-if="row.status === 'DRAFT'"
              type="success"
              link
              @click="handleRelease(row)"
            >发布</el-button>
            <el-button
              v-if="row.status === 'RELEASED'"
              type="warning"
              link
              @click="handleObsolete(row)"
            >废弃</el-button>
            <el-button
              v-if="row.status !== 'RELEASED'"
              type="danger"
              link
              @click="handleDelete(row)"
            >删除</el-button>
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
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="900" :close-on-click-modal="false">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="脚本名称" prop="name">
              <el-input v-model="formData.name" placeholder="请输入脚本名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="脚本编码" prop="code">
              <el-input v-model="formData.code" placeholder="请输入脚本编码" :disabled="isEdit" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="脚本类型" prop="type">
              <el-select v-model="formData.type" placeholder="请选择" :disabled="isEdit" style="width: 100%">
                <el-option label="JavaScript" value="JAVASCRIPT" />
                <el-option label="SQL" value="SQL" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="触发时机" prop="trigger_type">
              <el-select v-model="formData.trigger_type" placeholder="请选择" :disabled="isEdit" style="width: 100%">
                <el-option label="发起前" value="BEFORE_SUBMIT" />
                <el-option label="发起后" value="AFTER_SUBMIT" />
                <el-option label="审批前" value="BEFORE_APPROVE" />
                <el-option label="审批后" value="AFTER_APPROVE" />
                <el-option label="驳回前" value="BEFORE_REJECT" />
                <el-option label="驳回后" value="AFTER_REJECT" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="业务类型" prop="business_type">
              <el-select v-model="formData.business_type" placeholder="请选择" :disabled="isEdit" style="width: 100%">
                <el-option label="物料" value="MATERIAL" />
                <el-option label="文档" value="DOCUMENT" />
                <el-option label="BOM" value="BOM" />
                <el-option label="全部" value="ALL" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="超时时间">
              <el-input-number v-model="formData.timeout" :min="1000" :max="60000" :step="1000" style="width: 100%" />
              <span class="form-tip">毫秒（默认5000）</span>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="描述">
          <el-input v-model="formData.description" type="textarea" :rows="2" placeholder="请输入描述" />
        </el-form-item>
        <el-form-item label="脚本内容" prop="content">
          <div class="code-editor-wrapper">
            <div class="code-editor-tips" v-if="formData.type === 'JAVASCRIPT'">
              <p>可用变量：businessType, businessID, context</p>
              <p>可用API：dbQuery(sql), dbQueryOne(sql), getMaterial(id), getDocument(id), getBOM(id)</p>
              <p>入口函数：定义main()函数返回结果，返回 { success: true/false, message: '...' }</p>
            </div>
            <div class="code-editor-tips" v-else-if="formData.type === 'SQL'">
              <p>只允许SELECT查询</p>
              <p>可用变量：${businessType}, ${businessID}, ${context.xxx}</p>
            </div>
            <el-input
              v-model="formData.content"
              type="textarea"
              :rows="12"
              placeholder="请输入脚本内容"
              class="code-editor"
            />
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>

    <!-- 测试对话框 -->
    <el-dialog v-model="testDialogVisible" title="测试脚本" width="800">
      <el-form :model="testForm" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="业务类型">
              <el-select v-model="testForm.business_type" placeholder="请选择" style="width: 100%">
                <el-option label="物料" value="MATERIAL" />
                <el-option label="文档" value="DOCUMENT" />
                <el-option label="BOM" value="BOM" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="业务ID">
              <el-input-number v-model="testForm.business_id" :min="1" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="上下文">
          <el-input
            v-model="testForm.contextStr"
            type="textarea"
            :rows="3"
            placeholder="JSON格式的上下文数据，如 { &quot;key&quot;: &quot;value&quot; }"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="testDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="executeTest" :loading="testLoading">执行测试</el-button>
      </template>

      <!-- 测试结果 -->
      <div v-if="testResult" class="test-result">
        <el-divider>测试结果</el-divider>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="执行状态">
            <el-tag :type="testResult.success ? 'success' : 'danger'">
              {{ testResult.success ? '成功' : '失败' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="耗时">{{ testResult.duration }}ms</el-descriptions-item>
        </el-descriptions>
        <div v-if="testResult.error" class="test-error">
          <strong>错误信息：</strong>
          <pre>{{ testResult.error }}</pre>
        </div>
        <div v-if="testResult.output" class="test-output">
          <strong>输出结果：</strong>
          <pre>{{ JSON.stringify(testResult.output, null, 2) }}</pre>
        </div>
      </div>
    </el-dialog>

    <!-- 日志对话框 -->
    <el-dialog v-model="logDialogVisible" title="执行日志" width="1000">
      <el-table :data="logData" v-loading="logLoading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="script_name" label="脚本名称" width="150" />
        <el-table-column label="触发时机" width="100">
          <template #default="{ row }">
            {{ getTriggerText(row.trigger_type) }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.success ? 'success' : 'danger'">
              {{ row.success ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="duration" label="耗时(ms)" width="100" />
        <el-table-column label="执行时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.executed_at) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewLogDetail(row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        class="pagination"
        v-model:current-page="logPagination.page"
        v-model:page-size="logPagination.pageSize"
        :total="logPagination.total"
        layout="total, prev, pager, next"
        @current-change="fetchLogs"
      />
    </el-dialog>

    <!-- 日志详情对话框 -->
    <el-dialog v-model="logDetailVisible" title="日志详情" width="800">
      <el-descriptions :column="2" border v-if="logDetail">
        <el-descriptions-item label="脚本名称">{{ logDetail.script_name }}</el-descriptions-item>
        <el-descriptions-item label="触发时机">{{ getTriggerText(logDetail.trigger_type) }}</el-descriptions-item>
        <el-descriptions-item label="业务类型">{{ getBusinessText(logDetail.business_type) }}</el-descriptions-item>
        <el-descriptions-item label="业务ID">{{ logDetail.business_id }}</el-descriptions-item>
        <el-descriptions-item label="执行状态">
          <el-tag :type="logDetail.success ? 'success' : 'danger'">
            {{ logDetail.success ? '成功' : '失败' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="耗时">{{ logDetail.duration }}ms</el-descriptions-item>
        <el-descriptions-item label="执行时间" :span="2">{{ formatDate(logDetail.executed_at) }}</el-descriptions-item>
      </el-descriptions>
      <div v-if="logDetail?.input" class="log-section">
        <strong>输入参数：</strong>
        <pre>{{ formatJSON(logDetail.input) }}</pre>
      </div>
      <div v-if="logDetail?.output" class="log-section">
        <strong>输出结果：</strong>
        <pre>{{ formatJSON(logDetail.output) }}</pre>
      </div>
      <div v-if="logDetail?.error_msg" class="log-section error">
        <strong>错误信息：</strong>
        <pre>{{ logDetail.error_msg }}</pre>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import {
  getScriptList,
  createScript,
  updateScript,
  deleteScript,
  releaseScript,
  obsoleteScript,
  testScript,
  getExecutionLogs,
  getExecutionLogDetail,
  type Script,
  type CreateScriptRequest,
  type UpdateScriptRequest,
  type ScriptTestResult,
  type ScriptExecutionLog
} from '@/api/script'

// 加载状态
const loading = ref(false)
const submitting = ref(false)
const testLoading = ref(false)
const logLoading = ref(false)

// 数据
const tableData = ref<Script[]>([])

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 搜索表单
const searchForm = reactive({
  type: '',
  trigger_type: '',
  business_type: '',
  status: '',
  keyword: ''
})

// 对话框
const dialogVisible = ref(false)
const dialogTitle = computed(() => (isEdit.value ? '编辑脚本' : '新增脚本'))
const isEdit = ref(false)
const currentId = ref<number | null>(null)

// 表单
const formRef = ref<FormInstance>()
const formData = reactive<CreateScriptRequest & { timeout: number }>({
  name: '',
  code: '',
  type: 'JAVASCRIPT',
  trigger_type: 'BEFORE_APPROVE',
  business_type: 'ALL',
  content: '',
  description: '',
  timeout: 5000
})

const rules: FormRules = {
  name: [
    { required: true, message: '请输入脚本名称', trigger: 'blur' },
    { max: 200, message: '名称不能超过200个字符', trigger: 'blur' }
  ],
  code: [
    { required: true, message: '请输入脚本编码', trigger: 'blur' },
    { max: 50, message: '编码不能超过50个字符', trigger: 'blur' }
  ],
  type: [{ required: true, message: '请选择脚本类型', trigger: 'change' }],
  trigger_type: [{ required: true, message: '请选择触发时机', trigger: 'change' }],
  business_type: [{ required: true, message: '请选择业务类型', trigger: 'change' }],
  content: [{ required: true, message: '请输入脚本内容', trigger: 'blur' }]
}

// 测试对话框
const testDialogVisible = ref(false)
const testForm = reactive({
  type: 'JAVASCRIPT' as 'JAVASCRIPT' | 'SQL',
  content: '',
  business_type: 'MATERIAL',
  business_id: 1,
  contextStr: '{}'
})
const testResult = ref<ScriptTestResult | null>(null)

// 日志对话框
const logDialogVisible = ref(false)
const logData = ref<ScriptExecutionLog[]>([])
const logPagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
  scriptId: 0
})

// 日志详情
const logDetailVisible = ref(false)
const logDetail = ref<ScriptExecutionLog | null>(null)

// 获取数据
async function fetchData() {
  loading.value = true
  try {
    const res = await getScriptList({
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
  searchForm.type = ''
  searchForm.trigger_type = ''
  searchForm.business_type = ''
  searchForm.status = ''
  searchForm.keyword = ''
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
function handleEdit(row: Script) {
  isEdit.value = true
  currentId.value = row.id
  formData.name = row.name
  formData.code = row.code
  formData.type = row.type as 'JAVASCRIPT' | 'SQL'
  formData.trigger_type = row.trigger_type
  formData.business_type = row.business_type
  formData.content = row.content
  formData.description = row.description || ''
  formData.timeout = row.timeout || 5000
  dialogVisible.value = true
}

// 测试
function handleTest(row: Script) {
  testForm.type = row.type as 'JAVASCRIPT' | 'SQL'
  testForm.content = row.content
  testForm.business_type = row.business_type === 'ALL' ? 'MATERIAL' : row.business_type
  testForm.business_id = 1
  testForm.contextStr = '{}'
  testResult.value = null
  testDialogVisible.value = true
}

// 执行测试
async function executeTest() {
  testLoading.value = true
  testResult.value = null
  try {
    let context = {}
    if (testForm.contextStr) {
      try {
        context = JSON.parse(testForm.contextStr)
      } catch {
        ElMessage.warning('上下文JSON格式不正确')
        return
      }
    }

    const res = await testScript({
      type: testForm.type,
      content: testForm.content,
      business_type: testForm.business_type,
      business_id: testForm.business_id,
      context
    })
    if (res.code === 0) {
      testResult.value = res.data
    }
  } catch (error: any) {
    ElMessage.error(error.message || '测试失败')
  } finally {
    testLoading.value = false
  }
}

// 查看日志
function handleViewLogs(row: Script) {
  logPagination.scriptId = row.id
  logPagination.page = 1
  logDialogVisible.value = true
  fetchLogs()
}

// 获取日志
async function fetchLogs() {
  logLoading.value = true
  try {
    const res = await getExecutionLogs({
      script_id: logPagination.scriptId,
      page: logPagination.page,
      page_size: logPagination.pageSize
    })
    if (res.code === 0) {
      logData.value = res.data.list
      logPagination.total = res.data.total
    }
  } catch (error: any) {
    ElMessage.error(error.message || '获取日志失败')
  } finally {
    logLoading.value = false
  }
}

// 查看日志详情
async function handleViewLogDetail(row: ScriptExecutionLog) {
  try {
    const res = await getExecutionLogDetail(row.id)
    if (res.code === 0) {
      logDetail.value = res.data
      logDetailVisible.value = true
    }
  } catch (error: any) {
    ElMessage.error(error.message || '获取日志详情失败')
  }
}

// 发布
async function handleRelease(row: Script) {
  try {
    await ElMessageBox.confirm('确定要发布该脚本吗？发布后将在流程中自动执行。', '提示', {
      type: 'warning'
    })
    const res = await releaseScript(row.id)
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

// 废弃
async function handleObsolete(row: Script) {
  try {
    await ElMessageBox.confirm('确定要废弃该脚本吗？废弃后将不再执行。', '提示', {
      type: 'warning'
    })
    const res = await obsoleteScript(row.id)
    if (res.code === 0) {
      ElMessage.success('废弃成功')
      fetchData()
    }
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '废弃失败')
    }
  }
}

// 删除
async function handleDelete(row: Script) {
  try {
    await ElMessageBox.confirm('确定要删除该脚本吗？', '提示', {
      type: 'warning'
    })
    const res = await deleteScript(row.id)
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
    if (isEdit.value && currentId.value) {
      const data: UpdateScriptRequest = {
        name: formData.name,
        content: formData.content,
        description: formData.description,
        timeout: formData.timeout
      }
      const res = await updateScript(currentId.value, data)
      if (res.code === 0) {
        ElMessage.success('更新成功')
        dialogVisible.value = false
        fetchData()
      }
    } else {
      const res = await createScript(formData)
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
  formData.name = ''
  formData.code = ''
  formData.type = 'JAVASCRIPT'
  formData.trigger_type = 'BEFORE_APPROVE'
  formData.business_type = 'ALL'
  formData.content = ''
  formData.description = ''
  formData.timeout = 5000
  formRef.value?.resetFields()
}

// 状态标签
function getStatusTag(status: string) {
  const map: Record<string, string> = {
    DRAFT: '',
    RELEASED: 'success',
    OBSOLETE: 'info'
  }
  return map[status] || 'info'
}

function getStatusText(status: string) {
  const map: Record<string, string> = {
    DRAFT: '草稿',
    RELEASED: '已发布',
    OBSOLETE: '已废弃'
  }
  return map[status] || status
}

// 触发时机文本
function getTriggerText(type: string) {
  const map: Record<string, string> = {
    BEFORE_SUBMIT: '发起前',
    AFTER_SUBMIT: '发起后',
    BEFORE_APPROVE: '审批前',
    AFTER_APPROVE: '审批后',
    BEFORE_REJECT: '驳回前',
    AFTER_REJECT: '驳回后'
  }
  return map[type] || type
}

// 业务类型文本
function getBusinessText(type: string) {
  const map: Record<string, string> = {
    MATERIAL: '物料',
    DOCUMENT: '文档',
    BOM: 'BOM',
    ALL: '全部'
  }
  return map[type] || type
}

// 格式化日期
function formatDate(date: string) {
  if (!date) return ''
  return new Date(date).toLocaleString('zh-CN')
}

// 格式化JSON
function formatJSON(str: string) {
  try {
    return JSON.stringify(JSON.parse(str), null, 2)
  } catch {
    return str
  }
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

  .form-tip {
    margin-left: 10px;
    color: var(--el-text-color-secondary);
    font-size: 12px;
  }

  .code-editor-wrapper {
    width: 100%;

    .code-editor-tips {
      margin-bottom: 10px;
      padding: 10px;
      background: var(--el-fill-color-light);
      border-radius: 4px;
      font-size: 12px;
      color: var(--el-text-color-secondary);

      p {
        margin: 4px 0;
      }
    }

    .code-editor {
      :deep(textarea) {
        font-family: 'Fira Code', 'Consolas', monospace;
        font-size: 13px;
        line-height: 1.5;
      }
    }
  }

  .test-result {
    margin-top: 20px;

    .test-error {
      margin-top: 10px;
      padding: 10px;
      background: var(--el-color-danger-light-9);
      border-radius: 4px;

      pre {
        margin: 5px 0 0 0;
        white-space: pre-wrap;
        word-break: break-all;
      }
    }

    .test-output {
      margin-top: 10px;
      padding: 10px;
      background: var(--el-fill-color-light);
      border-radius: 4px;

      pre {
        margin: 5px 0 0 0;
        white-space: pre-wrap;
        word-break: break-all;
      }
    }
  }

  .log-section {
    margin-top: 15px;
    padding: 10px;
    background: var(--el-fill-color-light);
    border-radius: 4px;

    &.error {
      background: var(--el-color-danger-light-9);
    }

    pre {
      margin: 5px 0 0 0;
      white-space: pre-wrap;
      word-break: break-all;
      max-height: 200px;
      overflow: auto;
    }
  }
}
</style>
