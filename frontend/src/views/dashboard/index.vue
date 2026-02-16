<template>
  <div class="dashboard-container">
    <el-row :gutter="20">
      <!-- 用户信息卡片 -->
      <el-col :span="8">
        <el-card class="user-card">
          <template #header>
            <div class="card-header">
              <span>个人信息</span>
            </div>
          </template>
          <div class="user-info">
            <el-avatar :size="80" :src="userStore.userInfo?.avatar">
              {{ userStore.realName?.charAt(0) }}
            </el-avatar>
            <div class="user-detail">
              <h3>{{ userStore.realName }}</h3>
              <p>{{ userStore.userInfo?.username }}</p>
              <p>
                <el-tag
                  v-for="role in userStore.roles"
                  :key="role.id"
                  size="small"
                  type="primary"
                  class="role-tag"
                >
                  {{ role.name }}
                </el-tag>
              </p>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 快捷入口 -->
      <el-col :span="8">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>快捷入口</span>
            </div>
          </template>
          <div class="quick-links">
            <div class="quick-item" @click="$router.push('/user/list')">
              <el-icon :size="30"><User /></el-icon>
              <span>用户管理</span>
            </div>
            <div class="quick-item" @click="$router.push('#')">
              <el-icon :size="30"><Box /></el-icon>
              <span>物料管理</span>
            </div>
            <div class="quick-item" @click="$router.push('#')">
              <el-icon :size="30"><Document /></el-icon>
              <span>文档管理</span>
            </div>
            <div class="quick-item" @click="$router.push('#')">
              <el-icon :size="30"><List /></el-icon>
              <span>BOM管理</span>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 统计信息 -->
      <el-col :span="8">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>数据统计</span>
            </div>
          </template>
          <el-row :gutter="20">
            <el-col :span="12">
              <div class="stat-item">
                <div class="stat-value">0</div>
                <div class="stat-label">物料数量</div>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="stat-item">
                <div class="stat-value">0</div>
                <div class="stat-label">文档数量</div>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="stat-item">
                <div class="stat-value">0</div>
                <div class="stat-label">BOM数量</div>
              </div>
            </el-col>
            <el-col :span="12">
              <div class="stat-item">
                <div class="stat-value">0</div>
                <div class="stat-label">待办任务</div>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
    </el-row>

    <!-- 待办事项 -->
    <el-row :gutter="20" class="mt-20">
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>待办事项</span>
              <el-button type="primary" link>查看全部</el-button>
            </div>
          </template>
          <el-empty description="暂无待办事项" />
        </el-card>
      </el-col>

      <!-- 最近访问 -->
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>最近访问</span>
              <el-button type="primary" link @click="clearRecentVisits">清空记录</el-button>
            </div>
          </template>
          <div v-if="recentVisits.length > 0" class="recent-visits">
            <div
              v-for="item in recentVisits"
              :key="item.path + item.timestamp"
              class="visit-item"
              @click="$router.push(item.path)"
            >
              <el-icon :size="18"><component :is="getIcon(item.type)" /></el-icon>
              <div class="visit-info">
                <span class="visit-title">{{ item.title }}</span>
                <span class="visit-time">{{ formatTime(item.timestamp) }}</span>
              </div>
            </div>
          </div>
          <el-empty v-else description="暂无访问记录" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { User, Box, Document, List, FolderOpened, DataAnalysis } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()

interface RecentVisit {
  path: string
  title: string
  type: string
  timestamp: number
}

const recentVisits = ref<RecentVisit[]>([])

onMounted(() => {
  loadRecentVisits()
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
  const iconMap: Record<string, any> = {
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
</script>

<style scoped lang="scss">
@use '@/styles/variables.scss' as *;

.dashboard-container {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
}

.user-card {
  .user-info {
    display: flex;
    align-items: center;

    .user-detail {
      margin-left: 20px;

      h3 {
        margin: 0 0 5px;
        font-size: 18px;
      }

      p {
        margin: 0 0 5px;
        color: $text-secondary;
        font-size: 14px;
      }

      .role-tag {
        margin-right: 5px;
      }
    }
  }
}

.quick-links {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 15px;

  .quick-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 15px;
    cursor: pointer;
    border-radius: 8px;
    transition: all 0.3s;

    &:hover {
      background-color: $bg-color-page;
      color: $primary-color;
    }

    span {
      margin-top: 8px;
      font-size: 14px;
    }
  }
}

.stat-item {
  text-align: center;
  padding: 15px 0;

  .stat-value {
    font-size: 28px;
    font-weight: bold;
    color: $primary-color;
  }

  .stat-label {
    font-size: 14px;
    color: $text-secondary;
    margin-top: 5px;
  }
}

.recent-visits {
  .visit-item {
    display: flex;
    align-items: center;
    padding: 12px 0;
    border-bottom: 1px solid $border-color-light;
    cursor: pointer;
    transition: all 0.3s;

    &:last-child {
      border-bottom: none;
    }

    &:hover {
      background-color: $bg-color-page;
      padding-left: 10px;
      margin-left: -10px;
      margin-right: -10px;
      padding-right: 10px;
    }

    .el-icon {
      color: $primary-color;
    }

    .visit-info {
      margin-left: 12px;
      flex: 1;
      display: flex;
      justify-content: space-between;
      align-items: center;

      .visit-title {
        font-size: 14px;
        color: $text-primary;
      }

      .visit-time {
        font-size: 12px;
        color: $text-secondary;
      }
    }
  }
}
</style>
