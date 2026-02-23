<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>流程查询</span>
          <el-tag type="info" size="large">共 {{ pagination.total }} 条</el-tag>
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
        <el-form-item label="审批节点">
          <el-select v-model="searchForm.node_name" placeholder="请选择节点" clearable filterable @change="handleSearch">
            <el-option v-for="node in allApprovalNodes" :key="node" :label="node" :value="node" />
          </el-select>
        </el-form-item>
        <el-form-item label="关键字">
          <el-input v-model="searchForm.keyword" placeholder="标题/业务编码" clearable @keyup.enter="handleSearch" style="width: 180px" />
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
        <el-table-column prop="business_type" label="业务类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getBusinessTypeTag(row.business_type)">
              {{ getBusinessTypeText(row.business_type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="business_code" label="业务编码" width="150" show-overflow-tooltip />
        <el-table-column prop="business_name" label="业务名称" width="180" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="current_node" label="当前节点" width="120">
          <template #default="{ row }">
            <span v-if="row.status === 'PENDING'">{{ row.current_node || '-' }}</span>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="current_assignee" label="当前处理人" width="150">
          <template #default="{ row }">
            <span>{{ getCurrentAssignee(row) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="initiator" label="发起人" width="100">
          <template #default="{ row }">
            {{ row.initiator?.real_name || row.initiator?.username || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="发起时间" width="170" />
        <el-table-column label="操作" fixed="right" width="100">
          <template #default="{ row }">
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
    <el-dialog v-model="detailVisible" title="流程详情" width="800">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="流程标题" :span="2">{{ currentInstance?.title }}</el-descriptions-item>
        <el-descriptions-item label="业务类型">{{ getBusinessTypeText(currentInstance?.business_type) }}</el-descriptions-item>
        <el-descriptions-item label="业务编码">{{ currentInstance?.business_code }}</el-descriptions-item>
        <el-descriptions-item label="业务名称" :span="2">{{ currentInstance?.business_name }}</el-descriptions-item>
        <el-descriptions-item label="业务状态">
          <el-tag :type="getBusinessStatusType(currentInstance?.business_status)">
            {{ getBusinessStatusText(currentInstance?.business_status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="流程状态">
          <el-tag :type="getStatusType(currentInstance?.status)">{{ getStatusText(currentInstance?.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="当前节点">
          <span v-if="currentInstance?.status === 'PENDING'">{{ currentInstance?.current_node || '-' }}</span>
          <span v-else>-</span>
        </el-descriptions-item>
        <el-descriptions-item label="当前处理人">{{ getCurrentAssignee(currentInstance) }}</el-descriptions-item>
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import {
  getInstanceList,
  getInstance,
  getHistories,
  type WorkflowInstanceResponse,
  type WorkflowHistory,
  type WorkflowDefinition,
} from '@/api/workflow'

defineOptions({
  name: 'WorkflowQueryList',
})

const loading = ref(false)
const tableData = ref<WorkflowInstanceResponse[]>([])
const detailVisible = ref(false)
const currentInstance = ref<WorkflowInstanceResponse | null>(null)
const histories = ref<WorkflowHistory[]>([])

const searchForm = reactive({
  type: '',
  status: '',
  node_name: '',
  keyword: '',
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0,
})

// 所有审批节点名称（用于筛选）
const allApprovalNodes = computed<string[]>(() => {
  const nodeSet = new Set<string>()
  for (const inst of tableData.value) {
    if (inst.definition?.config) {
      try {
        const config = JSON.parse(inst.definition.config)
        for (const node of config.nodes || []) {
          if (node.type === 'approval' && node.name) {
            nodeSet.add(node.name)
          }
        }
      } catch {
        // ignore
      }
    }
  }
  return Array.from(nodeSet).sort()
})

onMounted(() => {
  fetchData()
})

async function fetchData() {
  loading.value = true
  try {
    const res = await getInstanceList({
      page: pagination.page,
      page_size: pagination.pageSize,
      type: searchForm.type,
      status: searchForm.status,
      keyword: searchForm.keyword,
    })
    if (res.code === 0) {
      let list = res.data.list
      // 前端按节点名称筛选
      if (searchForm.node_name) {
        list = list.filter((inst: WorkflowInstanceResponse) => {
          if (!inst.definition?.config) return false
          try {
            const config = JSON.parse(inst.definition.config)
            return config.nodes?.some((n: any) => n.type === 'approval' && n.name === searchForm.node_name)
          } catch {
            return false
          }
        })
      }
      tableData.value = list
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
  searchForm.node_name = ''
  searchForm.keyword = ''
  pagination.page = 1
  fetchData()
}

async function handleView(row: WorkflowInstanceResponse) {
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

// 获取当前处理人
function getCurrentAssignee(row: WorkflowInstanceResponse | null): string {
  if (!row?.definition?.config || !row.current_node_id) {
    return '-'
  }

  try {
    const config = JSON.parse(row.definition.config)
    const currentNode = config.nodes?.find((n: any) => n.id === row.current_node_id)
    if (!currentNode?.assignee) {
      return '-'
    }

    const assignee = currentNode.assignee
    const typeMap: Record<string, string> = {
      role: '角色',
      department: '部门',
      user: '用户',
      initiator: '发起人',
    }
    return `${typeMap[assignee.type] || assignee.type}: ${assignee.value}`
  } catch {
    return '-'
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
