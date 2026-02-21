<template>
  <div class="sidebar-container">
    <!-- Logo -->
    <div class="sidebar-logo">
      <div class="logo-icon">
        <img src="@/assets/logo.svg" alt="Logo" class="logo-img" />
      </div>
      <transition name="fade">
        <span v-show="!collapsed" class="logo-text neon-text">PLM</span>
      </transition>
    </div>

    <!-- 菜单 -->
    <el-menu
      :default-active="activeMenu"
      :collapse="collapsed"
      :collapse-transition="false"
      background-color="transparent"
      text-color="#94a3b8"
      active-text-color="#00d4ff"
      router
      class="tech-menu"
    >
      <template v-for="route in menuRoutes" :key="route.path">
        <!-- 有子菜单 -->
        <el-sub-menu v-if="route.children && route.children.length > 1" :index="route.path">
          <template #title>
            <div class="menu-item-content">
              <el-icon class="menu-icon">
                <component :is="route.meta?.icon || 'Menu'" />
              </el-icon>
              <span>{{ route.meta?.title }}</span>
            </div>
          </template>

          <el-menu-item
            v-for="child in route.children"
            :key="child.path"
            :index="`${route.path}/${child.path}`"
          >
            <el-icon class="menu-icon">
              <component :is="child.meta?.icon || 'Document'" />
            </el-icon>
            <span>{{ child.meta?.title }}</span>
          </el-menu-item>
        </el-sub-menu>

        <!-- 无子菜单 -->
        <el-menu-item v-else :index="getMenuIndex(route)">
          <div class="menu-item-content">
            <el-icon class="menu-icon">
              <component :is="route.meta?.icon || 'Menu'" />
            </el-icon>
            <span>{{ route.meta?.title }}</span>
          </div>
        </el-menu-item>
      </template>
    </el-menu>

    <!-- 底部装饰 -->
    <div class="sidebar-footer">
      <div class="tech-line"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, type RouteRecordRaw } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { privateRoutes } from '@/router'

const route = useRoute()
const appStore = useAppStore()

const collapsed = computed(() => appStore.sidebarCollapsed)
const activeMenu = computed(() => route.path)

// 过滤出需要显示的菜单
const menuRoutes = computed(() => {
  return privateRoutes.filter(r => r.meta?.title)
})

function getMenuIndex(route: RouteRecordRaw): string {
  if (route.children && route.children.length > 0) {
    const child = route.children[0]
    if (child) {
      return `${route.path}/${child.path}`.replace('//', '/')
    }
  }
  return route.path
}
</script>

<style scoped lang="scss">
@use '@/styles/variables.scss' as *;

.sidebar-container {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: linear-gradient(180deg, rgba(15, 23, 42, 0.95) 0%, rgba(15, 23, 42, 0.98) 100%);
  border-right: 1px solid $border-color;
  position: relative;

  /* 科技感网格背景 */
  &::before {
    content: '';
    position: absolute;
    inset: 0;
    background-image:
      linear-gradient(rgba(0, 212, 255, 0.03) 1px, transparent 1px),
      linear-gradient(90deg, rgba(0, 212, 255, 0.03) 1px, transparent 1px);
    background-size: 30px 30px;
    pointer-events: none;
  }
}

.sidebar-logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 0 16px;
  position: relative;
  z-index: 1;

  /* 底部发光线 */
  &::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 20px;
    right: 20px;
    height: 1px;
    background: linear-gradient(90deg, transparent, rgba(0, 212, 255, 0.5), transparent);
  }

  .logo-icon {
    width: 36px;
    height: 36px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: $primary-gradient;
    border-radius: $border-radius-base;
    box-shadow: 0 0 15px rgba(0, 212, 255, 0.4);

    .logo-img {
      width: 24px;
      height: 24px;
      filter: brightness(0) invert(1);
    }
  }

  .logo-text {
    font-size: 20px;
    font-weight: 700;
    letter-spacing: 2px;
    white-space: nowrap;
  }
}

.tech-menu {
  border-right: none;
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 8px 0;
  position: relative;
  z-index: 1;

  /* 自定义滚动条 */
  &::-webkit-scrollbar {
    width: 4px;
  }

  &::-webkit-scrollbar-thumb {
    background: rgba(0, 212, 255, 0.3);
    border-radius: 2px;
  }
}

.menu-item-content {
  display: flex;
  align-items: center;
  gap: 8px;
}

.menu-icon {
  font-size: 18px;
  transition: all $transition-fast;
}

:deep(.el-menu-item),
:deep(.el-sub-menu__title) {
  height: 48px;
  line-height: 48px;
  margin: 4px 12px;
  border-radius: $border-radius-base;
  transition: all $transition-fast;
  position: relative;

  &:hover {
    background: rgba(0, 212, 255, 0.08) !important;

    .menu-icon {
      color: $primary-color;
      transform: scale(1.1);
    }
  }
}

:deep(.el-menu-item.is-active) {
  background: rgba(0, 212, 255, 0.12) !important;
  color: $primary-color;
  font-weight: 500;

  /* 左侧指示条 */
  &::before {
    content: '';
    position: absolute;
    left: 0;
    top: 50%;
    transform: translateY(-50%);
    width: 3px;
    height: 60%;
    background: $primary-gradient;
    border-radius: 0 2px 2px 0;
    box-shadow: 0 0 10px rgba(0, 212, 255, 0.5);
  }

  .menu-icon {
    color: $primary-color;
    filter: drop-shadow(0 0 4px rgba(0, 212, 255, 0.5));
  }
}

:deep(.el-sub-menu) {
  .el-menu {
    background: transparent;
  }

  .el-menu-item {
    padding-left: 56px !important;
    min-width: auto;
  }

  &.is-active {
    .el-sub-menu__title {
      color: $primary-color;

      .menu-icon {
        color: $primary-color;
      }
    }
  }
}

/* 子菜单箭头 */
:deep(.el-sub-menu__icon-arrow) {
  color: $text-secondary;
  font-size: 12px;
}

.sidebar-footer {
  padding: 16px;
  position: relative;
  z-index: 1;

  .tech-line {
    height: 2px;
    background: linear-gradient(90deg, transparent, rgba(0, 212, 255, 0.3), transparent);
    border-radius: 1px;
    animation: pulse 2s ease-in-out infinite;
  }
}

@keyframes pulse {
  0%, 100% {
    opacity: 0.5;
  }
  50% {
    opacity: 1;
  }
}

/* 淡入淡出动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
