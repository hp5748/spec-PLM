<template>
  <div class="sidebar-container">
    <!-- Logo -->
    <div class="sidebar-logo">
      <img src="@/assets/logo.svg" alt="Logo" class="logo-img" />
      <span v-show="!collapsed" class="logo-text">PLM</span>
    </div>

    <!-- 菜单 -->
    <el-menu
      :default-active="activeMenu"
      :collapse="collapsed"
      :collapse-transition="false"
      background-color="#304156"
      text-color="#bfcbd9"
      active-text-color="#409eff"
      router
    >
      <template v-for="route in menuRoutes" :key="route.path">
        <!-- 有子菜单 -->
        <el-sub-menu v-if="route.children && route.children.length > 1" :index="route.path">
          <template #title>
            <el-icon>
              <component :is="route.meta?.icon || 'Menu'" />
            </el-icon>
            <span>{{ route.meta?.title }}</span>
          </template>

          <el-menu-item
            v-for="child in route.children"
            :key="child.path"
            :index="`${route.path}/${child.path}`"
          >
            <el-icon>
              <component :is="child.meta?.icon || 'Document'" />
            </el-icon>
            <span>{{ child.meta?.title }}</span>
          </el-menu-item>
        </el-sub-menu>

        <!-- 无子菜单 -->
        <el-menu-item v-else :index="getMenuIndex(route)">
          <el-icon>
            <component :is="route.meta?.icon || 'Menu'" />
          </el-icon>
          <span>{{ route.meta?.title }}</span>
        </el-menu-item>
      </template>
    </el-menu>
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
.sidebar-container {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.sidebar-logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #2b2f3a;
  overflow: hidden;

  .logo-img {
    width: 32px;
    height: 32px;
  }

  .logo-text {
    margin-left: 12px;
    font-size: 18px;
    font-weight: bold;
    color: #fff;
    white-space: nowrap;
  }
}

.el-menu {
  border-right: none;
  flex: 1;
  overflow-y: auto;
}

:deep(.el-menu-item),
:deep(.el-sub-menu__title) {
  height: 50px;
  line-height: 50px;
}

:deep(.el-menu-item.is-active) {
  background-color: #263445 !important;
}

:deep(.el-menu-item:hover),
:deep(.el-sub-menu__title:hover) {
  background-color: #263445 !important;
}
</style>
