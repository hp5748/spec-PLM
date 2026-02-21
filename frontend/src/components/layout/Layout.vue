<template>
  <el-container class="layout-container">
    <!-- 侧边栏 -->
    <el-aside :width="appStore.sidebarCollapsed ? '64px' : '220px'" class="layout-aside">
      <AppSidebar />
    </el-aside>

    <el-container class="layout-main">
      <!-- 顶部Header -->
      <el-header class="layout-header">
        <AppHeader />
      </el-header>

      <!-- 标签页 -->
      <TabsView v-if="appStore.visitedTabs.length > 0" />

      <!-- 主内容区 -->
      <el-main class="layout-content">
        <router-view v-slot="{ Component }">
          <transition name="fade-slide" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { useAppStore } from '@/stores/app'
import AppSidebar from './AppSidebar.vue'
import AppHeader from './AppHeader.vue'
import TabsView from './TabsView.vue'

const appStore = useAppStore()
</script>

<style scoped lang="scss">
@use '@/styles/variables.scss' as *;

.layout-container {
  height: 100vh;
  width: 100%;
  background: $bg-color-page;
}

.layout-aside {
  background: transparent;
  transition: width $transition-duration $ease-out-cubic;
  overflow: hidden;
  flex-shrink: 0;
}

.layout-main {
  display: flex;
  flex-direction: column;
  height: 100%;
  min-width: 0;
  position: relative;

  /* 科技感背景网格 */
  &::before {
    content: '';
    position: absolute;
    inset: 0;
    background-image:
      linear-gradient(rgba(0, 212, 255, 0.015) 1px, transparent 1px),
      linear-gradient(90deg, rgba(0, 212, 255, 0.015) 1px, transparent 1px);
    background-size: 40px 40px;
    pointer-events: none;
    z-index: 0;
  }
}

.layout-header {
  height: $header-height;
  padding: 0;
  background: rgba(30, 41, 59, 0.9);
  border-bottom: 1px solid $border-color;
  position: relative;
  z-index: 10;

  /* 底部发光线 */
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

.layout-content {
  flex: 1;
  background: transparent;
  overflow: auto;
  padding: 20px;
  position: relative;
  z-index: 1;
}

/* 页面切换动画 */
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.25s $ease-out-cubic;
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}
</style>
