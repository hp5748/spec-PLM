<template>
  <div class="header-container">
    <!-- 左侧：折叠按钮 -->
    <div class="header-left">
      <div class="collapse-btn" @click="toggleSidebar">
        <el-icon :size="20">
          <Fold v-if="!collapsed" />
          <Expand v-else />
        </el-icon>
      </div>
    </div>

    <!-- 右侧：用户信息 -->
    <div class="header-right">
      <el-dropdown @command="handleCommand" trigger="click">
        <div class="user-info">
          <el-avatar :size="32" :src="userStore.userInfo?.avatar" class="user-avatar">
            {{ userStore.realName?.charAt(0) }}
          </el-avatar>
          <span class="username">{{ userStore.realName }}</span>
          <el-icon class="arrow-icon"><ArrowDown /></el-icon>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="profile">
              <el-icon><User /></el-icon>
              个人中心
            </el-dropdown-item>
            <el-dropdown-item command="settings">
              <el-icon><Setting /></el-icon>
              设置
            </el-dropdown-item>
            <el-dropdown-item divided command="logout">
              <el-icon><SwitchButton /></el-icon>
              退出登录
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessageBox } from 'element-plus'
import { Fold, Expand, ArrowDown, User, Setting, SwitchButton } from '@element-plus/icons-vue'
import { useAppStore } from '@/stores/app'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const appStore = useAppStore()
const userStore = useUserStore()

const collapsed = computed(() => appStore.sidebarCollapsed)

function toggleSidebar() {
  appStore.toggleSidebar()
}

function handleCommand(command: string) {
  switch (command) {
    case 'profile':
      router.push('/profile')
      break
    case 'settings':
      router.push('/settings')
      break
    case 'logout':
      handleLogout()
      break
  }
}

async function handleLogout() {
  try {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    })

    userStore.logout()
    router.push('/login')
  } catch {
    // 取消退出
  }
}
</script>

<style scoped lang="scss">
@use '@/styles/variables.scss' as *;

.header-container {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  background: rgba(30, 41, 59, 0.8);
  backdrop-filter: blur(10px);
  position: relative;

  /* 底部发光边框 */
  &::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 1px;
    background: linear-gradient(90deg, transparent, rgba(0, 212, 255, 0.3), transparent);
  }
}

.header-left {
  display: flex;
  align-items: center;
}

.collapse-btn {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: $border-radius-base;
  cursor: pointer;
  color: $text-regular;
  transition: all $transition-fast;
  position: relative;

  /* 科技感边框 */
  &::before {
    content: '';
    position: absolute;
    inset: 0;
    border-radius: inherit;
    padding: 1px;
    background: transparent;
    -webkit-mask:
      linear-gradient(#fff 0 0) content-box,
      linear-gradient(#fff 0 0);
    mask:
      linear-gradient(#fff 0 0) content-box,
      linear-gradient(#fff 0 0);
    -webkit-mask-composite: xor;
    mask-composite: exclude;
    transition: all $transition-fast;
  }

  &:hover {
    color: $primary-color;
    background: rgba(0, 212, 255, 0.1);

    &::before {
      background: linear-gradient(135deg, rgba(0, 212, 255, 0.5), rgba(14, 165, 233, 0.5));
    }
  }

  &:active {
    transform: scale(0.95);
  }
}

.header-right {
  display: flex;
  align-items: center;
}

.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 6px 12px;
  border-radius: $border-radius-large;
  transition: all $transition-fast;
  position: relative;

  /* 悬浮时显示发光边框 */
  &::before {
    content: '';
    position: absolute;
    inset: 0;
    border-radius: inherit;
    padding: 1px;
    background: transparent;
    -webkit-mask:
      linear-gradient(#fff 0 0) content-box,
      linear-gradient(#fff 0 0);
    mask:
      linear-gradient(#fff 0 0) content-box,
      linear-gradient(#fff 0 0);
    -webkit-mask-composite: xor;
    mask-composite: exclude;
    transition: all $transition-fast;
  }

  &:hover {
    background: rgba(0, 212, 255, 0.08);

    &::before {
      background: linear-gradient(135deg, rgba(0, 212, 255, 0.4), rgba(59, 130, 246, 0.4));
    }

    .username {
      color: $primary-color;
    }

    .arrow-icon {
      color: $primary-color;
    }
  }

  .user-avatar {
    border: 2px solid rgba(0, 212, 255, 0.3);
    box-shadow: 0 0 10px rgba(0, 212, 255, 0.2);
  }

  .username {
    margin: 0 10px;
    font-size: 14px;
    font-weight: 500;
    color: $text-primary;
    transition: color $transition-fast;
  }

  .arrow-icon {
    font-size: 12px;
    color: $text-secondary;
    transition: all $transition-fast;
  }
}
</style>
