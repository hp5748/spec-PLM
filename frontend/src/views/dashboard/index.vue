<template>
  <div class="dashboard-container">
    <el-row :gutter="20">
      <!-- 用户信息卡片 -->
      <el-col :span="8">
        <el-card class="user-card tech-card">
          <template #header>
            <div class="card-header">
              <span class="header-title">个人信息</span>
              <div class="header-indicator"></div>
            </div>
          </template>
          <div class="user-info">
            <div class="avatar-wrapper">
              <el-avatar :size="80" :src="userStore.userInfo?.avatar" class="user-avatar">
                {{ userStore.realName?.charAt(0) }}
              </el-avatar>
              <div class="avatar-glow"></div>
            </div>
            <div class="user-detail">
              <h3 class="user-name">{{ userStore.realName }}</h3>
              <p class="user-account">@{{ userStore.userInfo?.username }}</p>
              <div class="role-tags">
                <el-tag
                  v-for="role in userStore.roles"
                  :key="role.id"
                  size="small"
                  type="primary"
                  class="role-tag"
                >
                  {{ role.name }}
                </el-tag>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 快捷入口 -->
      <el-col :span="8">
        <el-card class="tech-card">
          <template #header>
            <div class="card-header">
              <span class="header-title">快捷入口</span>
              <div class="header-indicator"></div>
            </div>
          </template>
          <div class="quick-links">
            <div class="quick-item" @click="$router.push('/user/list')">
              <div class="quick-icon-wrapper">
                <el-icon :size="28"><User /></el-icon>
              </div>
              <span>用户管理</span>
            </div>
            <div class="quick-item" @click="$router.push('/material/list')">
              <div class="quick-icon-wrapper">
                <el-icon :size="28"><Box /></el-icon>
              </div>
              <span>物料管理</span>
            </div>
            <div class="quick-item" @click="$router.push('/document/list')">
              <div class="quick-icon-wrapper">
                <el-icon :size="28"><Document /></el-icon>
              </div>
              <span>文档管理</span>
            </div>
            <div class="quick-item" @click="$router.push('/bom/list')">
              <div class="quick-icon-wrapper">
                <el-icon :size="28"><List /></el-icon>
              </div>
              <span>BOM管理</span>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 统计信息 -->
      <el-col :span="8">
        <el-card class="tech-card">
          <template #header>
            <div class="card-header">
              <span class="header-title">数据统计</span>
              <div class="header-indicator"></div>
            </div>
          </template>
          <el-row :gutter="16">
            <el-col :span="12">
              <div class="stat-item">
                <div class="stat-icon material-icon">
                  <el-icon :size="24"><Box /></el-icon>
                </div>
                <div class="stat-content">
                  <div class="stat-value">0</div>
                  <div class="stat-label">物料数量</div>
                </div>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="stat-item">
                <div class="stat-icon document-icon">
                  <el-icon :size="24"><Document /></el-icon>
                </div>
                <div class="stat-content">
                  <div class="stat-value">0</div>
                  <div class="stat-label">文档数量</div>
                </div>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="stat-item">
                <div class="stat-icon bom-icon">
                  <el-icon :size="24"><List /></el-icon>
                </div>
                <div class="stat-content">
                  <div class="stat-value">0</div>
                  <div class="stat-label">BOM数量</div>
                </div>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="stat-item">
                <div class="stat-icon task-icon">
                  <el-icon :size="24"><Clock /></el-icon>
                </div>
                <div class="stat-content">
                  <div class="stat-value">0</div>
                  <div class="stat-label">待办任务</div>
                </div>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
    </el-row>

    <!-- 待办事项 -->
    <el-row :gutter="20" class="mt-20">
      <el-col :span="12">
        <el-card class="tech-card">
          <template #header>
            <div class="card-header">
              <span class="header-title">待办事项</span>
              <el-button type="primary" link class="view-all-btn" @click="$router.push('/workflow/todo')">查看全部</el-button>
            </div>
          </template>
          <div v-if="todoList.length > 0" class="todo-list">
            <div
              v-for="item in todoList"
              :key="item.id"
              class="todo-item"
              @click="$router.push('/workflow/todo')"
            >
              <div class="todo-icon">
                <el-icon :size="18"><component :is="getBusinessIcon(item.business_type)" /></el-icon>
              </div>
              <div class="todo-info">
                <span class="todo-title">{{ item.title }}</span>
                <span class="todo-meta">
                  <el-tag size="small" :type="getBusinessTag(item.business_type)">{{ getBusinessText(item.business_type) }}</el-tag>
                  <span class="todo-time">{{ formatTodoTime(item.created_at) }}</span>
                </span>
              </div>
              <el-icon class="todo-arrow"><ArrowRight /></el-icon>
            </div>
          </div>
          <el-empty v-else description="暂无待办事项" class="tech-empty" />
        </el-card>
      </el-col>

      <!-- 最近访问 -->
      <el-col :span="12">
        <el-card class="tech-card">
          <template #header>
            <div class="card-header">
              <span class="header-title">最近访问</span>
              <el-button type="primary" link class="view-all-btn" @click="clearRecentVisits">
                清空记录
              </el-button>
            </div>
          </template>
          <div v-if="recentVisits.length > 0" class="recent-visits">
            <div
              v-for="item in recentVisits"
              :key="item.path + item.timestamp"
              class="visit-item"
              @click="$router.push(item.path)"
            >
              <div class="visit-icon">
                <el-icon :size="18"><component :is="getIcon(item.type)" /></el-icon>
              </div>
              <div class="visit-info">
                <span class="visit-title">{{ item.title }}</span>
                <span class="visit-time">{{ formatTime(item.timestamp) }}</span>
              </div>
              <el-icon class="visit-arrow"><ArrowRight /></el-icon>
            </div>
          </div>
          <el-empty v-else description="暂无访问记录" class="tech-empty" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import {
  User,
  Box,
  Document,
  List,
  FolderOpened,
  DataAnalysis,
  Clock,
  ArrowRight
} from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { getMyTodos, getTodoCount, type TodoItemResponse } from '@/api/workflow'

const userStore = useUserStore()

interface RecentVisit {
  path: string
  title: string
  type: string
  timestamp: number
}

const recentVisits = ref<RecentVisit[]>([])
const todoList = ref<TodoItemResponse[]>([])
const todoCount = ref(0)

onMounted(() => {
  loadRecentVisits()
  loadTodoList()
  // 添加一些模拟数据
  if (recentVisits.value.length === 0) {
    const mockVisits: RecentVisit[] = [
      { path: '/user/list', title: '用户列表', type: 'user', timestamp: Date.now() - 3600000 },
      { path: '/user/organization', title: '组织管理', type: 'user', timestamp: Date.now() - 7200000 },
      { path: '/user/department', title: '部门管理', type: 'user', timestamp: Date.now() - 10800000 },
    ]
    localStorage.setItem('recentVisits', JSON.stringify(mockVisits))
    loadRecentVisits()
  }
})

async function loadTodoList() {
  try {
    const res = await getMyTodos({ page: 1, page_size: 5 })
    if (res.code === 0) {
      todoList.value = res.data.list
      todoCount.value = res.data.total
    }
  } catch (error) {
    console.error(error)
  }
}

function loadRecentVisits() {
  const visits = localStorage.getItem('recentVisits')
  if (visits) {
    recentVisits.value = JSON.parse(visits)
  }
}

function clearRecentVisits() {
  recentVisits.value = []
  localStorage.removeItem('recentVisits')
}

function getIcon(type: string) {
  const iconMap: Record<string, unknown> = {
    user: User,
    material: Box,
    document: Document,
    bom: List,
    workflow: DataAnalysis,
    default: FolderOpened,
  }
  return iconMap[type] || iconMap.default
}

function formatTime(timestamp: number) {
  const now = Date.now()
  const diff = now - timestamp

  if (diff < 60000) {
    return '刚刚'
  } else if (diff < 3600000) {
    return `${Math.floor(diff / 60000)} 分钟前`
  } else if (diff < 86400000) {
    return `${Math.floor(diff / 3600000)} 小时前`
  } else {
    const date = new Date(timestamp)
    return `${date.getMonth() + 1}/${date.getDate()} ${date.getHours()}:${String(date.getMinutes()).padStart(2, '0')}`
  }
}

function getBusinessIcon(type: string) {
  const iconMap: Record<string, unknown> = {
    MATERIAL: Box,
    DOCUMENT: Document,
    BOM: List,
    default: FolderOpened,
  }
  return iconMap[type] || iconMap.default
}

function getBusinessText(type: string) {
  const textMap: Record<string, string> = {
    MATERIAL: '物料',
    DOCUMENT: '文档',
    BOM: 'BOM',
  }
  return textMap[type] || type
}

function getBusinessTag(type: string) {
  const tagMap: Record<string, string> = {
    MATERIAL: '',
    DOCUMENT: 'success',
    BOM: 'warning',
  }
  return tagMap[type] || 'info'
}

function formatTodoTime(timeStr: string) {
  const date = new Date(timeStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()

  if (diff < 60000) {
    return '刚刚'
  } else if (diff < 3600000) {
    return `${Math.floor(diff / 60000)} 分钟前`
  } else if (diff < 86400000) {
    return `${Math.floor(diff / 3600000)} 小时前`
  } else {
    return `${date.getMonth() + 1}/${date.getDate()}`
  }
}
</script>

<style scoped lang="scss">
@use '@/styles/variables.scss' as *;

.dashboard-container {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .header-title {
      font-weight: 600;
      color: $text-primary;
    }

    .header-indicator {
      width: 8px;
      height: 8px;
      background: $primary-color;
      border-radius: 50%;
      box-shadow: 0 0 10px $primary-color;
      animation: pulse 2s ease-in-out infinite;
    }

    .view-all-btn {
      color: $primary-color;
      font-size: 13px;

      &:hover {
        text-shadow: 0 0 8px rgba(0, 212, 255, 0.5);
      }
    }
  }
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
    transform: scale(1);
  }
  50% {
    opacity: 0.5;
    transform: scale(0.8);
  }
}

/* 用户卡片 */
.user-card {
  .user-info {
    display: flex;
    align-items: center;
  }

  .avatar-wrapper {
    position: relative;
    margin-right: 20px;

    .user-avatar {
      border: 3px solid rgba(0, 212, 255, 0.4);
      box-shadow: 0 0 20px rgba(0, 212, 255, 0.3);
    }

    .avatar-glow {
      position: absolute;
      inset: -4px;
      border-radius: 50%;
      background: $primary-gradient;
      opacity: 0.2;
      filter: blur(8px);
      animation: avatarGlow 3s ease-in-out infinite;
    }
  }

  .user-detail {
    flex: 1;

    .user-name {
      margin: 0 0 6px;
      font-size: 20px;
      font-weight: 600;
      color: $text-primary;
    }

    .user-account {
      margin: 0 0 10px;
      font-size: 14px;
      color: $text-secondary;
    }

    .role-tags {
      display: flex;
      gap: 8px;
      flex-wrap: wrap;
    }

    .role-tag {
      margin: 0;
    }
  }
}

@keyframes avatarGlow {
  0%, 100% {
    opacity: 0.2;
    transform: scale(1);
  }
  50% {
    opacity: 0.4;
    transform: scale(1.1);
  }
}

/* 快捷入口 */
.quick-links {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;

  .quick-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 20px 12px;
    cursor: pointer;
    border-radius: $border-radius-base;
    background: rgba(0, 212, 255, 0.03);
    border: 1px solid transparent;
    transition: all $transition-fast;

    &:hover {
      background: rgba(0, 212, 255, 0.08);
      border-color: rgba(0, 212, 255, 0.2);
      transform: translateY(-2px);

      .quick-icon-wrapper {
        background: $primary-gradient;
        box-shadow: 0 0 20px rgba(0, 212, 255, 0.4);
      }

      span {
        color: $primary-color;
      }
    }

    .quick-icon-wrapper {
      width: 52px;
      height: 52px;
      display: flex;
      align-items: center;
      justify-content: center;
      background: rgba(0, 212, 255, 0.1);
      border-radius: $border-radius-base;
      margin-bottom: 12px;
      transition: all $transition-fast;
      color: $primary-color;
    }

    span {
      font-size: 13px;
      color: $text-regular;
      transition: color $transition-fast;
    }
  }
}

/* 统计卡片 */
.stat-item {
  display: flex;
  align-items: center;
  padding: 16px;
  background: rgba(0, 212, 255, 0.03);
  border-radius: $border-radius-base;
  margin-bottom: 16px;
  border: 1px solid transparent;
  transition: all $transition-fast;

  &:hover {
    background: rgba(0, 212, 255, 0.06);
    border-color: rgba(0, 212, 255, 0.15);
  }

  .stat-icon {
    width: 48px;
    height: 48px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: $border-radius-base;
    margin-right: 14px;
    color: #fff;

    &.material-icon {
      background: linear-gradient(135deg, #00d4ff, #0ea5e9);
    }

    &.document-icon {
      background: linear-gradient(135deg, #10b981, #059669);
    }

    &.bom-icon {
      background: linear-gradient(135deg, #f59e0b, #d97706);
    }

    &.task-icon {
      background: linear-gradient(135deg, #8b5cf6, #7c3aed);
    }
  }

  .stat-content {
    flex: 1;

    .stat-value {
      font-size: 26px;
      font-weight: 700;
      color: $text-primary;
      line-height: 1.2;
    }

    .stat-label {
      font-size: 13px;
      color: $text-secondary;
      margin-top: 4px;
    }
  }
}

/* 最近访问 */
.recent-visits {
  .visit-item {
    display: flex;
    align-items: center;
    padding: 14px 16px;
    margin: 0 -16px;
    cursor: pointer;
    transition: all $transition-fast;
    border-left: 2px solid transparent;

    &:hover {
      background: rgba(0, 212, 255, 0.06);
      border-left-color: $primary-color;

      .visit-icon {
        background: rgba(0, 212, 255, 0.15);
        color: $primary-color;
      }

      .visit-arrow {
        opacity: 1;
        transform: translateX(0);
      }
    }

    .visit-icon {
      width: 36px;
      height: 36px;
      display: flex;
      align-items: center;
      justify-content: center;
      background: rgba(0, 212, 255, 0.08);
      border-radius: $border-radius-small;
      color: $primary-color;
      margin-right: 14px;
      transition: all $transition-fast;
    }

    .visit-info {
      flex: 1;
      display: flex;
      flex-direction: column;

      .visit-title {
        font-size: 14px;
        color: $text-primary;
        font-weight: 500;
      }

      .visit-time {
        font-size: 12px;
        color: $text-placeholder;
        margin-top: 4px;
      }
    }

    .visit-arrow {
      color: $text-placeholder;
      opacity: 0;
      transform: translateX(-8px);
      transition: all $transition-fast;
    }
  }
}

/* 待办列表 */
.todo-list {
  .todo-item {
    display: flex;
    align-items: center;
    padding: 14px 16px;
    margin: 0 -16px;
    cursor: pointer;
    transition: all $transition-fast;
    border-left: 2px solid transparent;

    &:hover {
      background: rgba(0, 212, 255, 0.06);
      border-left-color: $primary-color;

      .todo-icon {
        background: rgba(0, 212, 255, 0.15);
        color: $primary-color;
      }

      .todo-arrow {
        opacity: 1;
        transform: translateX(0);
      }
    }

    .todo-icon {
      width: 36px;
      height: 36px;
      display: flex;
      align-items: center;
      justify-content: center;
      background: rgba(139, 92, 246, 0.1);
      border-radius: $border-radius-small;
      color: #8b5cf6;
      margin-right: 14px;
      transition: all $transition-fast;
    }

    .todo-info {
      flex: 1;
      display: flex;
      flex-direction: column;

      .todo-title {
        font-size: 14px;
        color: $text-primary;
        font-weight: 500;
      }

      .todo-meta {
        display: flex;
        align-items: center;
        gap: 8px;
        margin-top: 6px;

        .todo-time {
          font-size: 12px;
          color: $text-placeholder;
        }
      }
    }

    .todo-arrow {
      color: $text-placeholder;
      opacity: 0;
      transform: translateX(-8px);
      transition: all $transition-fast;
    }
  }
}

/* 空状态 */
.tech-empty {
  padding: 30px 0;

  :deep(.el-empty__description) {
    color: $text-secondary;
  }
}
</style>
