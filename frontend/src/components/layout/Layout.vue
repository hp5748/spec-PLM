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
          <transition name="fade" mode="out-in">
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
}

.layout-aside {
  background-color: #304156;
  transition: width 0.3s;
  overflow: hidden;
}

.layout-main {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.layout-header {
  height: $header-height;
  padding: 0;
  background-color: $bg-color;
  border-bottom: 1px solid $border-color-lighter;
  box-shadow: $box-shadow-light;
}

.layout-content {
  flex: 1;
  background-color: $bg-color-page;
  overflow: auto;
  padding: 20px;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
