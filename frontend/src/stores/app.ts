import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore('app', () => {
  // 侧边栏折叠状态
  const sidebarCollapsed = ref(false)

  // 标签页列表
  const visitedTabs = ref<Array<{ path: string; title: string; name?: string }>>([])

  // 切换侧边栏
  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  // 添加标签页
  function addTab(tab: { path: string; title: string; name?: string }) {
    if (!visitedTabs.value.some(t => t.path === tab.path)) {
      visitedTabs.value.push(tab)
    }
  }

  // 关闭标签页
  function closeTab(path: string) {
    const index = visitedTabs.value.findIndex(t => t.path === path)
    if (index > -1) {
      visitedTabs.value.splice(index, 1)
    }
  }

  // 关闭其他标签页
  function closeOtherTabs(path: string) {
    visitedTabs.value = visitedTabs.value.filter(t => t.path === path)
  }

  // 关闭所有标签页
  function closeAllTabs() {
    visitedTabs.value = []
  }

  return {
    sidebarCollapsed,
    visitedTabs,
    toggleSidebar,
    addTab,
    closeTab,
    closeOtherTabs,
    closeAllTabs,
  }
})
