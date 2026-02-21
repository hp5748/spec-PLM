<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>我的待办</span>
          <el-tag type="warning" size="large">待处理 {{ pagination.total }} 项</el-tag>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="业务类型">
          <el-select v-model="searchForm.type" placeholder="请选择业务类型" clearable @change="handleSearch">
            <el-option label="PART" value="MATERIAL" />
            <el-option label="文档" value="DOCUMENT" />
            <el-option label="BOM" value="BOM" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择状态" clearable @change="handleSearch">
            <el-option label="待审批" value="PENDING" />
            <el-option label="已通过" value="APPROVED" />
            <el-option label="已驳回" value="REJECTED" />
            <el-option label="已取消" value="CANCELLED" />
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
        <el-table-column prop="title" label="流程标题" min-width="200" show-overflow-tooltip />
        <el-table-column prop="business_type" label="业务类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getBusinessTypeTag(row.business_type)">
              {{ getBusinessTypeText(row.business_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="business_code" label="业务编码" width="150">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewBusiness(row)">
              {{ row.business_code }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column prop="business_name" label="业务名称" width="200" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="initiator" label="发起人" width="120">
          <template #default="{ row }">
            {{ row.initiator?.real_name || row.initiator?.username || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="发起时间" width="180" />
        <el-table-column label="操作" fixed="right" width="200">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleView(row)">查看</el-button>
            <el-button type="success" link @click="handleApprove(row)">审批</el-button>
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
    <el-dialog v-model="detailVisible" title="流程详情" width="800">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="流程标题" :span="2">{{ currentInstance?.title }}</el-descriptions-item>
        <el-descriptions-item label="业务类型">{{ getBusinessTypeText(currentInstance?.business_type) }}</el-descriptions-item>
        <el-descriptions-item label="业务编码">
          <el-button type="primary" link @click="handleViewBusiness(currentInstance)">
            {{ currentInstance?.business_code }}
          </el-button>
        </el-descriptions-item>
        <el-descriptions-item label="业务名称" :span="2">{{ currentInstance?.business_name }}</el-descriptions-item>
        <el-descriptions-item label="业务状态">
          <el-tag :type="getBusinessStatusType(currentInstance?.business_status)">
            {{ getBusinessStatusText(currentInstance?.business_status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="流程状态">
          <el-tag :type="getStatusType(currentInstance?.status)">{{ getStatusText(currentInstance?.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="发起人">{{ currentInstance?.initiator?.real_name || currentInstance?.initiator?.username }}</el-descriptions-item>
        <el-descriptions-item label="发起时间">{{ currentInstance?.created_at }}</el-descriptions-item>
        <el-descriptions-item label="发起意见" :span="2">{{ currentInstance?.comment || '-' }}</el-descriptions-item>
      </el-descriptions>

      <!-- 流程历史 -->
      <el-divider content-position="left">审批历史</el-divider>
      <el-timeline v-if="histories.length > 0">
        <el-timeline-item
          v-for="history in histories"
          :key="history.id"
          :type="getHistoryType(history.action)"
          :timestamp="history.created_at"
          placement="top"
        >
          <el-card shadow="hover">
            <div class="history-item">
              <span class="action">{{ getActionText(history.action) }}</span>
              <span class="operator">操作人：{{ history.operator_name || history.operator?.real_name || history.operator?.username }}</span>
              <div v-if="history.comment" class="comment">意见：{{ history.comment }}</div>
            </div>
          </el-card>
        </el-timeline-item>
      </el-timeline>
      <el-empty v-else description="暂无审批历史" />

      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 审批对话框 -->
    <el-dialog v-model="approveVisible" title="审批操作" width="500">
      <el-form ref="approveFormRef" :model="approveForm" :rules="approveRules" label-width="80px">
        <el-form-item label="审批结果">
          <el-radio-group v-model="approveForm.result">
            <el-radio value="approve">同意</el-radio>
            <el-radio value="reject">驳回</el-radio>
            <el-radio value="transfer">转交</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item
          v-if="approveForm.result === 'reject'"
          label="驳回原因"
          prop="comment"
        >
          <el-input
            v-model="approveForm.comment"
            type="textarea"
            :rows="3"
            placeholder="请输入驳回原因"
          />
        </el-form-item>
        <el-form-item v-else-if="approveForm.result === 'transfer'" label="转交用户" prop="target_user_id">
          <el-select v-model="approveForm.target_user_id" placeholder="请选择转交用户" style="width: 100%">
            <el-option
              v-for="user in userList"
              :key="user.id"
              :label="user.real_name || user.username"
              :value="user.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item v-else label="审批意见">
          <el-input
            v-model="approveForm.comment"
            type="textarea"
            :rows="3"
            placeholder="请输入审批意见（可选）"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="approveVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitApprove" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import {
  getMyTodos,
  getInstance,
  getHistories,
  approveWorkflow,
  rejectWorkflow,
  transferWorkflow,
  type TodoItemResponse,
  type WorkflowHistory,
} from '@/api/workflow'
import { getUserList, type User } from '@/api/user'

const router = useRouter()

defineOptions({
  name: 'TodoList',
})

const loading = ref(false)
const submitting = ref(false)
const tableData = ref<TodoItemResponse[]>([])
const detailVisible = ref(false)
const approveVisible = ref(false)
const currentInstance = ref<TodoItemResponse | null>(null)
const histories = ref<WorkflowHistory[]>([])
const userList = ref<User[]>([])
const approveFormRef = ref<FormInstance>()

const searchForm = reactive({
  type: '',
  status: '',
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0,
})

const approveForm = reactive({
  result: 'approve',
  comment: '',
  target_user_id: null as number | null,
})

const approveRules: FormRules = {
  comment: [
    { required: true, message: '请输入驳回原因', trigger: 'blur' },
  ],
  target_user_id: [
    { required: true, message: '请选择转交用户', trigger: 'change' },
  ],
}

onMounted(() => {
  fetchData()
  loadUserList()
})

async function loadUserList() {
  try {
    const res = await getUserList({ page: 1, page_size: 1000 })
    if (res.code === 0) {
      userList.value = res.data.list
    }
  } catch (error) {
    console.error(error)
  }
}

async function fetchData() {
  loading.value = true
  try {
    const res = await getMyTodos({
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
  searchForm.type = ''
  searchForm.status = ''
  pagination.page = 1
  fetchData()
}

async function handleView(row: TodoItemResponse) {
  try {
    const res = await getInstance(row.id)
    if (res.code === 0) {
      currentInstance.value = res.data
      // 加载流程历史
      const historyRes = await getHistories(row.id)
      if (historyRes.code === 0) {
        histories.value = historyRes.data
      }
      detailVisible.value = true
    }
  } catch (error) {
    console.error(error)
  }
}

// 跳转到业务对象详情
function handleViewBusiness(row: TodoItemResponse | null) {
  if (!row) return

  detailVisible.value = false

  switch (row.business_type) {
    case 'MATERIAL':
      router.push('/part/query')
      break
    case 'DOCUMENT':
      router.push('/document/query')
      break
    case 'BOM':
      router.push('/bom/list')
      break
  }
}

async function handleApprove(row: TodoItemResponse) {
  try {
    const res = await getInstance(row.id)
    if (res.code === 0) {
      currentInstance.value = res.data
      approveForm.result = 'approve'
      approveForm.comment = ''
      approveForm.target_user_id = null
      approveVisible.value = true
    }
  } catch (error) {
    console.error(error)
  }
}

async function handleSubmitApprove() {
  if (approveForm.result === 'reject' || approveForm.result === 'transfer') {
    const valid = await approveFormRef.value?.validate()
    if (!valid) return
  }

  if (!currentInstance.value) return

  submitting.value = true
  try {
    if (approveForm.result === 'approve') {
      await approveWorkflow(currentInstance.value.id, { comment: approveForm.comment })
      ElMessage.success('审批通过')
    } else if (approveForm.result === 'reject') {
      await rejectWorkflow(currentInstance.value.id, { comment: approveForm.comment })
      ElMessage.success('已驳回')
    } else if (approveForm.result === 'transfer') {
      await transferWorkflow(currentInstance.value.id, {
        target_user_id: approveForm.target_user_id!,
        comment: approveForm.comment,
      })
      ElMessage.success('已转交')
    }
    approveVisible.value = false
    fetchData()
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

function getBusinessTypeText(type?: string) {
  const map: Record<string, string> = {
    MATERIAL: 'PART',
    DOCUMENT: '文档',
    BOM: 'BOM',
  }
  return type ? map[type] || type : ''
}

function getBusinessTypeTag(type: string) {
  const map: Record<string, string> = {
    MATERIAL: '',
    DOCUMENT: 'success',
    BOM: 'warning',
  }
  return map[type] || 'info'
}

function getStatusType(status?: string) {
  const map: Record<string, string> = {
    DRAFT: 'info',
    PENDING: 'warning',
    APPROVED: 'success',
    REJECTED: 'danger',
    CANCELLED: 'info',
  }
  return status ? map[status] || 'info' : 'info'
}

function getStatusText(status?: string) {
  const map: Record<string, string> = {
    DRAFT: '拟制中',
    PENDING: '待审批',
    APPROVED: '已通过',
    REJECTED: '已驳回',
    CANCELLED: '已取消',
  }
  return status ? map[status] || status : ''
}

function getBusinessStatusType(status?: string) {
  const map: Record<string, string> = {
    DRAFT: 'info',
    REVIEWING: 'warning',
    RELEASED: 'success',
    REJECTED: 'danger',
    OBSOLETE: '',
  }
  return status ? map[status] || 'info' : 'info'
}

function getBusinessStatusText(status?: string) {
  const map: Record<string, string> = {
    DRAFT: '拟制中',
    REVIEWING: '审核中',
    RELEASED: '已发布',
    REJECTED: '已驳回',
    OBSOLETE: '已失效',
  }
  return status ? map[status] || status : ''
}

function getHistoryType(action: string) {
  const map: Record<string, string> = {
    SUBMIT: 'primary',
    APPROVE: 'success',
    REJECT: 'danger',
    TRANSFER: 'warning',
    CANCEL: 'info',
    WITHDRAW: 'info',
  }
  return map[action] || 'info'
}

function getActionText(action: string) {
  const map: Record<string, string> = {
    SUBMIT: '提交审批',
    APPROVE: '审批通过',
    REJECT: '审批驳回',
    TRANSFER: '转交',
    CANCEL: '取消',
    WITHDRAW: '撤回',
  }
  return map[action] || action
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

  .history-item {
    .action {
      font-weight: bold;
      margin-right: 16px;
    }

    .operator {
      color: #909399;
    }

    .comment {
      margin-top: 8px;
      color: #606266;
    }
  }
}
</style>
