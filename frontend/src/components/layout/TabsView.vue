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
      <div class="dropdown-trigger">
        <el-icon><ArrowDown /></el-icon>
      </div>
      <template #dropdown>
        <el-dropdown-menu>
          <el-dropdown-item command="closeOther">
            <el-icon><Close /></el-icon>
            关闭其他
          </el-dropdown-item>
          <el-dropdown-item command="closeAll">
            <el-icon><FolderRemove /></el-icon>
            关闭所有
          </el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ArrowDown, Close, FolderRemove } from '@element-plus/icons-vue'
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
  padding: 0 16px;
  height: 40px;
  background: rgba(15, 23, 42, 0.6);
  border-bottom: 1px solid $border-color;
  position: relative;
  z-index: 5;

  /* 底部发光线 */
  &::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 1px;
    background: linear-gradient(90deg, transparent, rgba(0, 212, 255, 0.2), transparent);
  }

  :deep(.el-tabs) {
    flex: 1;
    height: 100%;

    .el-tabs__header {
      margin: 0;
      border-bottom: none;
      height: 100%;

      .el-tabs__nav-wrap {
        padding: 0;
        height: 100%;

        .el-tabs__nav-scroll {
          height: 100%;
        }

        .el-tabs__nav {
          border: none;
          display: flex;
          align-items: center;
          height: 100%;
          gap: 6px;

          .el-tabs__item {
            height: 28px;
            line-height: 26px;
            padding: 0 14px;
            border-radius: $border-radius-small;
            border: 1px solid transparent;
            background: rgba(30, 41, 59, 0.5);
            color: $text-secondary;
            font-size: 13px;
            transition: all $transition-fast;
            position: relative;

            &:hover {
              color: $text-primary;
              background: rgba(0, 212, 255, 0.1);
              border-color: rgba(0, 212, 255, 0.2);
            }

            &.is-active {
              background: rgba(0, 212, 255, 0.15);
              border-color: rgba(0, 212, 255, 0.3);
              color: $primary-color;

              /* 发光效果 */
              &::before {
                content: '';
                position: absolute;
                inset: -1px;
                border-radius: inherit;
                background: linear-gradient(135deg, rgba(0, 212, 255, 0.3), transparent);
                z-index: -1;
                filter: blur(4px);
              }
            }

            .el-icon {
              width: 14px;
              height: 14px;
              margin-left: 6px;
              border-radius: 50%;
              transition: all $transition-fast;

              &:hover {
                background: rgba(239, 68, 68, 0.8);
                color: #fff;
              }
            }
          }
        }
      }
    }
  }

  .tabs-dropdown {
    margin-left: 12px;

    .dropdown-trigger {
      width: 28px;
      height: 28px;
      display: flex;
      align-items: center;
      justify-content: center;
      border-radius: $border-radius-small;
      background: rgba(30, 41, 59, 0.5);
      border: 1px solid $border-color;
      color: $text-secondary;
      cursor: pointer;
      transition: all $transition-fast;

      &:hover {
        color: $primary-color;
        border-color: rgba(0, 212, 255, 0.3);
        background: rgba(0, 212, 255, 0.1);
      }
    }
  }
}
</style>
