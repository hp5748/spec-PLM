<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>角色管理</span>
          <el-button type="primary" @click="handleAdd">新增角色</el-button>
        </div>
      </template>

      <!-- 表格 -->
      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="code" label="角色编码" width="150" />
        <el-table-column prop="name" label="角色名称" width="150" />
        <el-table-column prop="description" label="描述" min-width="200" />
        <el-table-column label="权限" width="300">
          <template #default="{ row }">
            <el-tag
              v-for="perm in row.permissions"
              :key="perm.id"
              size="small"
              class="mr-10"
              type="info"
            >
              {{ perm.name }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="is_system" label="系统角色" width="100">
          <template #default="{ row }">
            <el-tag :type="row.is_system ? 'warning' : 'info'">
              {{ row.is_system ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" fixed="right" width="150">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)" :disabled="row.is_system">
              编辑
            </el-button>
            <el-button
              type="danger"
              link
              @click="handleDelete(row)"
              :disabled="row.is_system"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-form-item label="角色编码" prop="code">
          <el-input v-model="formData.code" placeholder="请输入角色编码" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="角色名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入角色名称" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="formData.description" type="textarea" rows="2" placeholder="请输入描述" />
        </el-form-item>
        <el-form-item label="权限配置" prop="permission_ids">
          <el-checkbox-group v-model="formData.permission_ids">
            <div v-for="group in permissionGroups" :key="group.name" class="permission-group">
              <div class="group-title">{{ group.name }}</div>
              <el-checkbox
                v-for="perm in group.permissions"
                :key="perm.id"
                :value="perm.id"
              >
                {{ perm.name }}
              </el-checkbox>
            </div>
          </el-checkbox-group>
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
import { getRoleList, createRole } from '@/api/user'
import type { Role, Permission } from '@/types/api'

defineOptions({
  name: 'RoleList',
})

interface PermissionGroup {
  name: string
  permissions: Permission[]
}

const loading = ref(false)
const submitting = ref(false)
const tableData = ref<Role[]>([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const currentId = ref<number | null>(null)
const formRef = ref<FormInstance>()

// 模拟权限列表（实际应从后端获取）
const allPermissions = ref<Permission[]>([
  { id: 1, name: '查看用户', code: 'user:view', module: '用户管理', created_at: '', updated_at: '' },
  { id: 2, name: '创建用户', code: 'user:create', module: '用户管理', created_at: '', updated_at: '' },
  { id: 3, name: '编辑用户', code: 'user:edit', module: '用户管理', created_at: '', updated_at: '' },
  { id: 4, name: '删除用户', code: 'user:delete', module: '用户管理', created_at: '', updated_at: '' },
  { id: 5, name: '查看物料', code: 'material:view', module: '物料管理', created_at: '', updated_at: '' },
  { id: 6, name: '创建物料', code: 'material:create', module: '物料管理', created_at: '', updated_at: '' },
  { id: 7, name: '编辑物料', code: 'material:edit', module: '物料管理', created_at: '', updated_at: '' },
  { id: 8, name: '删除物料', code: 'material:delete', module: '物料管理', created_at: '', updated_at: '' },
  { id: 9, name: '查看文档', code: 'document:view', module: '文档管理', created_at: '', updated_at: '' },
  { id: 10, name: '上传文档', code: 'document:upload', module: '文档管理', created_at: '', updated_at: '' },
  { id: 11, name: '查看BOM', code: 'bom:view', module: 'BOM管理', created_at: '', updated_at: '' },
  { id: 12, name: '创建BOM', code: 'bom:create', module: 'BOM管理', created_at: '', updated_at: '' },
  { id: 13, name: '查看流程', code: 'workflow:view', module: '流程管理', created_at: '', updated_at: '' },
  { id: 14, name: '审批流程', code: 'workflow:approve', module: '流程管理', created_at: '', updated_at: '' },
])

const dialogTitle = computed(() => (isEdit.value ? '编辑角色' : '新增角色'))

const formData = reactive({
  code: '',
  name: '',
  description: '',
  permission_ids: [] as number[],
})

const rules: FormRules = {
  code: [{ required: true, message: '请输入角色编码', trigger: 'blur' }],
  name: [{ required: true, message: '请输入角色名称', trigger: 'blur' }],
}

const permissionGroups = computed<PermissionGroup[]>(() => {
  const groups: Map<string, Permission[]> = new Map()
  allPermissions.value.forEach((perm) => {
    if (!groups.has(perm.module)) {
      groups.set(perm.module, [])
    }
    groups.get(perm.module)!.push(perm)
  })
  return Array.from(groups.entries()).map(([name, permissions]) => ({
    name,
    permissions,
  }))
})

onMounted(() => {
  fetchData()
})

async function fetchData() {
  loading.value = true
  try {
    const res = await getRoleList()
    if (res.code === 0) {
      tableData.value = res.data
    }
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

function resetForm() {
  formData.code = ''
  formData.name = ''
  formData.description = ''
  formData.permission_ids = []
  currentId.value = null
}

function handleAdd() {
  isEdit.value = false
  resetForm()
  dialogVisible.value = true
}

function handleEdit(row: Role) {
  isEdit.value = true
  currentId.value = row.id
  formData.code = row.code
  formData.name = row.name
  formData.description = row.description || ''
  formData.permission_ids = row.permissions?.map((p) => p.id) || []
  dialogVisible.value = true
}

async function handleSubmit() {
  const valid = await formRef.value?.validate()
  if (!valid) return

  submitting.value = true
  try {
    if (isEdit.value) {
      // TODO: 实现更新角色API
      ElMessage.success('更新成功')
    } else {
      await createRole(formData)
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

async function handleDelete(row: Role) {
  try {
    await ElMessageBox.confirm(`确定要删除角色 ${row.name} 吗？`, '提示', {
      type: 'warning',
    })
    // TODO: 实现删除角色API
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

  .mr-10 {
    margin-right: 10px;
    margin-bottom: 5px;
  }

  .permission-group {
    margin-bottom: 15px;

    .group-title {
      font-weight: bold;
      margin-bottom: 8px;
      color: #606266;
    }
  }
}
</style>
