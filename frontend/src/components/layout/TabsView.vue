<template>
  <div class="tabs-view">
    <el-tabs
      v-model="activeTab"
      type="card"
      closable
      @tab-remove="closeTab"
      @tab-click="clickTab"
    >
      <el-tab-pane
        v-for="tab in visitedTabs"
        :key="tab.path"
        :label="tab.title"
        :name="tab.path"
      />
    </el-tabs>

    <el-dropdown class="tabs-dropdown" @command="handleCommand">
      <el-button type="primary" size="small">
        <el-icon><ArrowDown /></el-icon>
      </el-button>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item command="closeOther">关闭其他</el-dropdown-item>
          <el-dropdown-item command="closeAll">关闭所有</el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ArrowDown } from '@element-plus/icons-vue'
import { useAppStore } from '@/stores/app'

const router = useRouter()
const route = useRoute()
const appStore = useAppStore()

const visitedTabs = computed(() => appStore.visitedTabs)
const activeTab = computed({
  get: () => route.path,
  set: () => {},
})

function clickTab(tab: { props: { name: string } }) {
  router.push(tab.props.name)
}

function closeTab(path: string) {
  appStore.closeTab(path)

  // 如果关闭的是当前页面，跳转到最后一个标签
  if (route.path === path) {
    const tabs = appStore.visitedTabs
    const lastTab = tabs[tabs.length - 1]
    if (lastTab) {
      router.push(lastTab.path)
    } else {
      router.push('/dashboard')
    }
  }
}

function handleCommand(command: string) {
  switch (command) {
    case 'closeOther':
      appStore.closeOtherTabs(route.path)
      break
    case 'closeAll':
      appStore.closeAllTabs()
      router.push('/dashboard')
      break
  }
}
</script>

<style scoped lang="scss">
@use '@/styles/variables.scss' as *;

.tabs-view {
  display: flex;
  align-items: center;
  padding: 0 10px;
  background-color: $bg-color;
  border-bottom: 1px solid $border-color-lighter;

  :deep(.el-tabs) {
    flex: 1;
    height: 36px;

    .el-tabs__header {
      margin: 0;
      border-bottom: none;

      .el-tabs__nav-wrap {
        padding: 0;

        .el-tabs__nav-scroll {
          height: 36px;
        }

        .el-tabs__nav {
          border: none;

          .el-tabs__item {
            height: 30px;
            line-height: 30px;
            margin-right: 5px;
            border-radius: 2px;
            border: 1px solid $border-color-lighter;

            &.is-active {
              background-color: $primary-color;
              border-color: $primary-color;
              color: #fff;
            }
          }
        }
      }
    }
  }

  .tabs-dropdown {
    margin-left: 10px;
  }
}
</style>
