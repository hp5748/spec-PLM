<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <div class="left">
            <el-button @click="handleBack">返回</el-button>
            <span class="title">{{ bomInfo?.name || 'BOM编辑' }}</span>
          </div>
          <div class="right">
            <el-button type="primary" @click="handleAddItem" :disabled="!canEdit">添加子物料</el-button>
            <el-upload
              ref="uploadRef"
              :show-file-list="false"
              :before-upload="handleImportBefore"
              :http-request="handleImport"
              accept=".xlsx,.xls"
              :disabled="!canEdit"
            >
              <el-button :disabled="!canEdit">导入</el-button>
            </el-upload>
            <el-button type="success" @click="handleExport">导出</el-button>
            <el-button @click="handleConvert" :disabled="!canEdit">
              {{ bomInfo?.is_exact ? '转为非精确' : '转为精确' }}
            </el-button>
            <el-button type="primary" @click="handleSave" :loading="saving" :disabled="!canEdit">保存</el-button>
          </div>
        </div>
      </template>

      <!-- BOM基本信息 -->
      <el-descriptions :column="4" border class="bom-info">
        <el-descriptions-item label="根物料">{{ bomInfo?.root_material?.item_id }} - {{ bomInfo?.root_material?.item_name }}</el-descriptions-item>
        <el-descriptions-item label="版本">{{ bomInfo?.root_version || '-' }}</el-descriptions-item>
        <el-descriptions-item label="类型">
          <el-tag :type="bomInfo?.is_exact ? 'success' : 'warning'">
            {{ bomInfo?.is_exact ? '精确BOM' : '非精确BOM' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusTag(bomInfo?.status || '')">{{ getStatusText(bomInfo?.status || '') }}</el-tag>
        </el-descriptions-item>
      </el-descriptions>

      <!-- BOM树形结构 -->
      <div class="tree-container">
        <div class="tree-header">
          <span>BOM结构</span>
          <el-button type="primary" link @click="expandAll">展开全部</el-button>
          <el-button type="primary" link @click="collapseAll">折叠全部</el-button>
        </div>
        <el-tree
          ref="treeRef"
          :data="treeData"
          :props="treeProps"
          node-key="id"
          default-expand-all
          :expand-on-click-node="false"
          v-loading="treeLoading"
        >
          <template #default="{ node, data }">
            <div class="tree-node">
              <div class="node-info">
                <span class="item-id">{{ data.item_id }}</span>
                <span class="item-name">{{ data.item_name }}</span>
                <el-tag v-if="data.version" size="small" type="info">{{ data.version }}</el-tag>
                <span class="quantity">× {{ data.quantity }} {{ data.unit || '个' }}</span>
              </div>
              <div class="node-actions" v-if="data.level > 0 && canEdit">
                <el-button type="primary" link size="small" @click.stop="handleEditItem(data)">编辑</el-button>
                <el-button type="danger" link size="small" @click.stop="handleDeleteItem(data)">删除</el-button>
              </div>
            </div>
          </template>
        </el-tree>
        <el-empty v-if="!treeLoading && treeData.length === 0" description="暂无BOM结构，请添加子物料" />
      </div>
    </el-card>

    <!-- 添加/编辑子物料对话框 -->
    <el-dialog v-model="itemDialogVisible" :title="itemDialogTitle" width="500">
      <el-form ref="itemFormRef" :model="itemFormData" :rules="itemRules" label-width="100px">
        <el-form-item label="父节点">
          <el-select v-model="itemFormData.parent_id" placeholder="根节点" clearable style="width: 100%">
            <el-option
              v-for="item in parentOptions"
              :key="item.id"
              :label="`${item.item_id} - ${item.item_name}`"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="物料" prop="material_id">
          <el-select
            v-model="itemFormData.material_id"
            placeholder="请选择物料"
            filterable
            remote
            :remote-method="searchMaterials"
            :loading="materialLoading"
            style="width: 100%"
          >
            <el-option
              v-for="item in materialOptions"
              :key="item.id"
              :label="`${item.item_id} - ${item.item_name}`"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="版本" v-if="bomInfo?.is_exact">
          <el-input v-model="itemFormData.version" placeholder="请输入版本号" />
        </el-form-item>
        <el-form-item label="数量" prop="quantity">
          <el-input-number v-model="itemFormData.quantity" :min="0.0001" :precision="4" style="width: 100%" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="itemFormData.sort_order" :min="0" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="itemDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitItem" :loading="itemSubmitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import {
  getBOMDetail,
  getBOMTree,
  addBOMItem,
  updateBOMItem,
  deleteBOMItem,
  convertBOMType,
  exportBOM,
  importBOM,
  type BOMView,
  type BOMTreeNode,
  type AddBOMItemRequest,
  type UpdateBOMItemRequest
} from '@/api/bom'
import { searchMaterials as searchMaterialsApi, type Material } from '@/api/material'

const route = useRoute()
const router = useRouter()

// 状态
const loading = ref(false)
const treeLoading = ref(false)
const saving = ref(false)
const materialLoading = ref(false)
const itemSubmitting = ref(false)

// 数据
const bomId = computed(() => Number(route.params.id))
const bomInfo = ref<BOMView | null>(null)
const treeData = ref<BOMTreeNode[]>([])
const treeRef = ref()
const uploadRef = ref() // used in template
const materialOptions = ref<Material[]>([])
const parentOptions = ref<BOMTreeNode[]>([])

// 计算是否可编辑
const canEdit = computed(() => {
  const status = bomInfo.value?.status
  return status === 'DRAFT' || status === 'REJECTED'
})

// 树形配置
const treeProps = {
  children: 'children',
  label: 'item_name'
}

// 子物料对话框
const itemDialogVisible = ref(false)
const itemDialogTitle = computed(() => (isEditItem.value ? '编辑子物料' : '添加子物料'))
const isEditItem = ref(false)
const currentItem = ref<BOMTreeNode | null>(null)
const itemFormRef = ref<FormInstance>()
const itemFormData = reactive<AddBOMItemRequest & { id?: number }>({
  parent_id: undefined,
  material_id: 0,
  version: '',
  quantity: 1,
  sort_order: 0
})

const itemRules: FormRules = {
  material_id: [{ required: true, message: '请选择物料', trigger: 'change' }],
  quantity: [{ required: true, message: '请输入数量', trigger: 'blur' }]
}

// 获取BOM详情
async function fetchBOMDetail() {
  loading.value = true
  try {
    const res = await getBOMDetail(bomId.value)
    if (res.code === 0) {
      bomInfo.value = res.data
    }
  } catch (error: any) {
    ElMessage.error(error.message || '获取BOM详情失败')
  } finally {
    loading.value = false
  }
}

// 获取BOM树
async function fetchBOMTree() {
  treeLoading.value = true
  try {
    const res = await getBOMTree(bomId.value)
    if (res.code === 0) {
      // 根节点的children作为树数据
      treeData.value = res.data.children || []
      // 收集所有节点作为父节点选项
      collectParentOptions(res.data)
    }
  } catch (error: any) {
    ElMessage.error(error.message || '获取BOM结构失败')
  } finally {
    treeLoading.value = false
  }
}

// 收集父节点选项
function collectParentOptions(node: BOMTreeNode) {
  parentOptions.value = []
  if (node && node.level === 0) {
    // 根节点作为选项
    parentOptions.value.push({
      id: 0,
      material_id: node.material_id,
      item_id: node.item_id,
      item_name: node.item_name,
      version: node.version,
      quantity: node.quantity,
      sort_order: 0,
      level: 0,
      unit: node.unit,
      children: []
    })
    // 收集所有子节点
    collectChildrenOptions(node.children || [])
  }
}

function collectChildrenOptions(nodes: BOMTreeNode[]) {
  for (const node of nodes) {
    parentOptions.value.push(node)
    if (node.children && node.children.length > 0) {
      collectChildrenOptions(node.children)
    }
  }
}

// 搜索物料
async function searchMaterials(keyword: string) {
  if (!keyword) return
  materialLoading.value = true
  try {
    const res = await searchMaterialsApi(keyword, { page: 1, page_size: 50 })
    if (res.code === 0) {
      materialOptions.value = res.data.list
    }
  } catch (error) {
    console.error('搜索物料失败', error)
  } finally {
    materialLoading.value = false
  }
}

// 返回
function handleBack() {
  router.push('/bom/list')
}

// 展开全部
function expandAll() {
  const nodes = treeRef.value?.store?.nodesMap
  if (nodes) {
    for (const key in nodes) {
      nodes[key].expanded = true
    }
  }
}

// 折叠全部
function collapseAll() {
  const nodes = treeRef.value?.store?.nodesMap
  if (nodes) {
    for (const key in nodes) {
      nodes[key].expanded = false
    }
  }
}

// 添加子物料
function handleAddItem() {
  isEditItem.value = false
  currentItem.value = null
  resetItemForm()
  itemDialogVisible.value = true
}

// 编辑子物料
function handleEditItem(node: BOMTreeNode) {
  isEditItem.value = true
  currentItem.value = node
  itemFormData.id = node.id
  itemFormData.material_id = node.material_id
  itemFormData.version = node.version || ''
  itemFormData.quantity = node.quantity
  itemFormData.sort_order = node.sort_order
  // 设置父节点（需要根据树结构找到）
  itemDialogVisible.value = true
}

// 删除子物料
async function handleDeleteItem(node: BOMTreeNode) {
  try {
    await ElMessageBox.confirm('确定要删除该子物料吗？其所有子项也将被删除。', '提示', {
      type: 'warning'
    })
    const res = await deleteBOMItem(bomId.value, node.id)
    if (res.code === 0) {
      ElMessage.success('删除成功')
      fetchBOMTree()
    }
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除失败')
    }
  }
}

// 转换类型
async function handleConvert() {
  const toExact = !bomInfo.value?.is_exact
  const confirmText = toExact
    ? '转为精确BOM需要所有物料都有版本号，确定要转换吗？'
    : '转为非精确BOM将清空所有版本信息，确定要转换吗？'

  try {
    await ElMessageBox.confirm(confirmText, '提示', { type: 'warning' })
    const res = await convertBOMType(bomId.value, { to_exact: toExact })
    if (res.code === 0) {
      ElMessage.success('转换成功')
      fetchBOMDetail()
      fetchBOMTree()
    }
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '转换失败')
    }
  }
}

// 保存（目前是即时保存，这里可以留作后续批量保存用）
function handleSave() {
  ElMessage.success('BOM已自动保存')
}

// 导出BOM
function handleExport() {
  exportBOM(bomId.value)
  ElMessage.success('正在导出BOM...')
}

// 导入前检查
function handleImportBefore(file: File) {
  const isExcel = file.name.endsWith('.xlsx') || file.name.endsWith('.xls')
  if (!isExcel) {
    ElMessage.error('只能上传Excel文件')
    return false
  }
  return true
}

// 导入BOM
async function handleImport(options: { file: File }) {
  try {
    const res = await importBOM(bomId.value, options.file)
    if (res.code === 0) {
      const result = res.data
      if (result.fail_count > 0) {
        ElMessage.warning(`导入完成：成功${result.success_count}条，失败${result.fail_count}条`)
        if (result.errors && result.errors.length > 0) {
          console.warn('导入错误:', result.errors)
        }
      } else {
        ElMessage.success(`导入成功：共${result.success_count}条`)
      }
      fetchBOMTree()
    }
  } catch (error: any) {
    ElMessage.error(error.message || '导入失败')
  }
}

// 提交子物料
async function handleSubmitItem() {
  const valid = await itemFormRef.value?.validate()
  if (!valid) return

  itemSubmitting.value = true
  try {
    if (isEditItem.value && currentItem.value) {
      // 编辑
      const data: UpdateBOMItemRequest = {
        quantity: itemFormData.quantity,
        sort_order: itemFormData.sort_order
      }
      const res = await updateBOMItem(bomId.value, currentItem.value.id, data)
      if (res.code === 0) {
        ElMessage.success('更新成功')
        itemDialogVisible.value = false
        fetchBOMTree()
      }
    } else {
      // 新增
      const data: AddBOMItemRequest = {
        parent_id: itemFormData.parent_id,
        material_id: itemFormData.material_id,
        version: itemFormData.version,
        quantity: itemFormData.quantity,
        sort_order: itemFormData.sort_order
      }
      const res = await addBOMItem(bomId.value, data)
      if (res.code === 0) {
        ElMessage.success('添加成功')
        itemDialogVisible.value = false
        fetchBOMTree()
      }
    }
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    itemSubmitting.value = false
  }
}

// 重置子物料表单
function resetItemForm() {
  itemFormData.id = undefined
  itemFormData.parent_id = undefined
  itemFormData.material_id = 0
  itemFormData.version = ''
  itemFormData.quantity = 1
  itemFormData.sort_order = 0
  materialOptions.value = []
  itemFormRef.value?.resetFields()
}

// 状态标签
function getStatusTag(status: string) {
  const map: Record<string, string> = {
    DRAFT: '',
    REVIEWING: 'warning',
    RELEASED: 'success',
    REJECTED: 'danger',
    OBSOLETE: 'info'
  }
  return map[status] || 'info'
}

function getStatusText(status: string) {
  const map: Record<string, string> = {
    DRAFT: '拟制中',
    REVIEWING: '审核中',
    RELEASED: '已发布',
    REJECTED: '已驳回',
    OBSOLETE: '已失效'
  }
  return map[status] || status
}

onMounted(() => {
  if (bomId.value) {
    fetchBOMDetail()
    fetchBOMTree()
  }
})
</script>

<style scoped lang="scss">
.page-container {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .left {
      display: flex;
      align-items: center;
      gap: 12px;

      .title {
        font-size: 16px;
        font-weight: 600;
      }
    }

    .right {
      display: flex;
      gap: 8px;
    }
  }

  .bom-info {
    margin-bottom: 20px;
  }

  .tree-container {
    border: 1px solid var(--el-border-color);
    border-radius: 4px;

    .tree-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 12px 16px;
      border-bottom: 1px solid var(--el-border-color);
      background: var(--el-fill-color-light);
    }

    .el-tree {
      padding: 16px;
    }

    .tree-node {
      display: flex;
      justify-content: space-between;
      align-items: center;
      width: 100%;
      padding-right: 8px;

      .node-info {
        display: flex;
        align-items: center;
        gap: 8px;

        .item-id {
          font-family: monospace;
          color: var(--el-text-color-secondary);
        }

        .item-name {
          font-weight: 500;
        }

        .quantity {
          color: var(--el-color-primary);
        }
      }

      .node-actions {
        visibility: hidden;
      }

      &:hover .node-actions {
        visibility: visible;
      }
    }
  }
}
</style>
