<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>部门管理</span>
          <el-button type="primary" @click="handleAdd(null)">新增部门</el-button>
        </div>
      </template>

      <!-- 组织筛选 -->
      <el-form :inline="true" class="filter-form">
        <el-form-item label="所属组织">
          <el-select v-model="selectedOrgId" placeholder="请选择组织" clearable @change="fetchData">
            <el-option v-for="org in organizations" :key="org.id" :label="org.name" :value="org.id" />
          </el-select>
        </el-form-item>
      </el-form>

      <!-- 树形表格 -->
      <el-table
        :data="tableData"
        v-loading="loading"
        row-key="id"
        :tree-props="{ children: 'children', hasChildren: 'hasChildren' }"
        default-expand-all
        stripe
      >
        <el-table-column prop="name" label="部门名称" min-width="200" />
        <el-table-column prop="code" label="部门编码" width="150" />
        <el-table-column prop="sort_order" label="排序" width="80" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'">
              {{ row.status === 'active' ? '正常' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" fixed="right" width="200">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleAdd(row)">添加子部门</el-button>
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-form-item label="所属组织" prop="organization_id">
          <el-select v-model="formData.organization_id" placeholder="请选择组织" style="width: 100%">
            <el-option v-for="org in organizations" :key="org.id" :label="org.name" :value="org.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="上级部门" prop="parent_id">
          <el-tree-select
            v-model="formData.parent_id"
            :data="departmentTree"
            :props="{ value: 'id', label: 'name', children: 'children' }"
            placeholder="请选择上级部门"
            clearable
            check-strictly
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="部门名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入部门名称" />
        </el-form-item>
        <el-form-item label="部门编码" prop="code">
          <el-input v-model="formData.code" placeholder="请输入部门编码" />
        </el-form-item>
        <el-form-item label="排序" prop="sort_order">
          <el-input-number v-model="formData.sort_order" :min="0" :max="999" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="formData.status">
            <el-radio value="active">正常</el-radio>
            <el-radio value="inactive">停用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { getDepartmentList, getOrganizationList, createDepartment } from '@/api/user'
import type { Department, Organization } from '@/types/api'

defineOptions({
  name: 'DepartmentList',
})

const loading = ref(false)
const submitting = ref(false)
const tableData = ref<Department[]>([])
const organizations = ref<Organization[]>([])
const selectedOrgId = ref<string>('')
const dialogVisible = ref(false)
const isEdit = ref(false)
const currentId = ref<number | null>(null)
const formRef = ref<FormInstance>()

const dialogTitle = computed(() => (isEdit.value ? '编辑部门' : '新增部门'))

const formData = reactive({
  organization_id: '' as string | number,
  parent_id: null as number | null,
  name: '',
  code: '',
  sort_order: 0,
  status: 'active',
})

const rules: FormRules = {
  organization_id: [{ required: true, message: '请选择所属组织', trigger: 'change' }],
  name: [{ required: true, message: '请输入部门名称', trigger: 'blur' }],
}

const departmentTree = computed(() => {
  return buildTree(tableData.value)
})

function buildTree(items: Department[]): any[] {
  const map = new Map<number, any>()
  const roots: any[] = []

  // 先创建所有节点
  items.forEach((item) => {
    map.set(item.id, { ...item })
  })

  // 构建树形结构
  items.forEach((item) => {
    const node = map.get(item.id)
    if (item.parent_id && map.has(item.parent_id)) {
      const parent = map.get(item.parent_id)
      if (!parent.children) {
        parent.children = []
      }
      parent.children.push(node)
    } else {
      roots.push(node)
    }
  })

  return roots
}

onMounted(async () => {
  await fetchOrganizations()
  fetchData()
})

async function fetchOrganizations() {
  try {
    const res = await getOrganizationList()
    if (res.code === 0) {
      organizations.value = res.data
      // 默认选择第一个组织
      if (organizations.value.length > 0 && !selectedOrgId.value) {
        selectedOrgId.value = String(organizations.value[0].id)
      }
    }
  } catch (error) {
    console.error(error)
  }
}

async function fetchData() {
  loading.value = true
  try {
    const res = await getDepartmentList(selectedOrgId.value ? String(selectedOrgId.value) : undefined)
    if (res.code === 0) {
      // 构建树形结构数据
      tableData.value = buildTree(res.data)
    }
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

function resetForm() {
  formData.organization_id = selectedOrgId.value || ''
  formData.parent_id = null
  formData.name = ''
  formData.code = ''
  formData.sort_order = 0
  formData.status = 'active'
  currentId.value = null
}

function handleAdd(parent: Department | null) {
  isEdit.value = false
  resetForm()
  if (parent) {
    formData.parent_id = parent.id
    if (!formData.organization_id && parent.organization_id) {
      formData.organization_id = parent.organization_id
    }
  }
  dialogVisible.value = true
}

function handleEdit(row: Department) {
  isEdit.value = true
  currentId.value = row.id
  formData.organization_id = row.organization_id
  formData.parent_id = row.parent_id || null
  formData.name = row.name
  formData.code = row.code || ''
  formData.sort_order = row.sort_order || 0
  formData.status = row.status
  dialogVisible.value = true
}

async function handleSubmit() {
  const valid = await formRef.value?.validate()
  if (!valid) return

  submitting.value = true
  try {
    if (isEdit.value) {
      // TODO: 实现更新部门API
      ElMessage.success('更新成功')
    } else {
      await createDepartment({
        organization_id: Number(formData.organization_id),
        parent_id: formData.parent_id,
        name: formData.name,
        code: formData.code,
        sort_order: formData.sort_order,
        status: formData.status,
      })
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchData()
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

async function handleDelete(row: Department) {
  try {
    await ElMessageBox.confirm(`确定要删除部门 ${row.name} 吗？`, '提示', {
      type: 'warning',
    })
    // TODO: 实现删除部门API
    ElMessage.success('删除成功')
    fetchData()
  } catch {
    // 取消删除
  }
}
</script>

<style scoped lang="scss">
.page-container {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .filter-form {
    margin-bottom: 20px;
  }
}
</style>
