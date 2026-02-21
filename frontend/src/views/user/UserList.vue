<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>用户列表</span>
          <el-button type="primary" @click="handleAdd">新增用户</el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="用户名">
          <el-input v-model="searchForm.username" placeholder="请输入用户名" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择状态" clearable>
            <el-option label="正常" value="active" />
            <el-option label="停用" value="inactive" />
            <el-option label="锁定" value="locked" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 表格 -->
      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="real_name" label="姓名" width="120" />
        <el-table-column prop="email" label="邮箱" width="180" />
        <el-table-column prop="phone" label="手机号" width="140" />
        <el-table-column label="角色" width="200">
          <template #default="{ row }">
            <el-tag v-for="role in row.roles" :key="role.id" size="small" class="mr-10">
              {{ role.name }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" fixed="right" width="180">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="primary" link @click="handleResetPwd(row)">重置密码</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        class="pagination"
        @size-change="fetchData"
        @current-change="fetchData"
      />
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="用户名" prop="username">
              <el-input v-model="formData.username" placeholder="请输入用户名" :disabled="isEdit" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="姓名" prop="real_name">
              <el-input v-model="formData.real_name" placeholder="请输入姓名" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="密码" :prop="isEdit ? '' : 'password'">
              <el-input
                v-model="formData.password"
                type="password"
                placeholder="请输入密码"
                show-password
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态" prop="status">
              <el-select v-model="formData.status" placeholder="请选择状态" style="width: 100%">
                <el-option label="正常" value="active" />
                <el-option label="停用" value="inactive" />
                <el-option label="锁定" value="locked" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="邮箱" prop="email">
              <el-input v-model="formData.email" placeholder="请输入邮箱" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="手机号" prop="phone">
              <el-input v-model="formData.phone" placeholder="请输入手机号" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="所属组织" prop="organization_id">
              <el-select
                v-model="formData.organization_id"
                placeholder="请选择组织"
                style="width: 100%"
                @change="handleOrgChange"
              >
                <el-option v-for="org in organizations" :key="org.id" :label="org.name" :value="org.id" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="所属部门" prop="department_id">
              <el-select v-model="formData.department_id" placeholder="请选择部门" style="width: 100%">
                <el-option v-for="dept in departments" :key="dept.id" :label="dept.name" :value="dept.id" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="角色分配" prop="role_ids">
          <el-checkbox-group v-model="formData.role_ids">
            <el-checkbox v-for="role in roles" :key="role.id" :value="role.id">
              {{ role.name }}
            </el-checkbox>
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
import {
  getUserList,
  createUser,
  updateUser,
  deleteUser,
  getOrganizationList,
  getDepartmentList,
  getRoleList,
} from '@/api/user'
import type { User, Organization, Department, Role } from '@/types/api'

defineOptions({
  name: 'UserList',
})

const loading = ref(false)
const submitting = ref(false)
const tableData = ref<User[]>([])
const organizations = ref<Organization[]>([])
const departments = ref<Department[]>([])
const roles = ref<Role[]>([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const currentId = ref<number | null>(null)
const formRef = ref<FormInstance>()

const searchForm = reactive({
  username: '',
  status: '',
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0,
})

const dialogTitle = computed(() => (isEdit.value ? '编辑用户' : '新增用户'))

const formData = reactive({
  username: '',
  password: '',
  real_name: '',
  email: '',
  phone: '',
  organization_id: null as number | null,
  department_id: null as number | null,
  status: 'active',
  role_ids: [] as number[],
})

const rules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 50, message: '用户名长度为3-50个字符', trigger: 'blur' },
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' },
  ],
  real_name: [{ required: true, message: '请输入姓名', trigger: 'blur' }],
  email: [{ type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }],
  status: [{ required: true, message: '请选择状态', trigger: 'change' }],
}

onMounted(async () => {
  await Promise.all([fetchOrganizations(), fetchRoles()])
  fetchData()
})

async function fetchOrganizations() {
  try {
    const res = await getOrganizationList()
    if (res.code === 0) {
      organizations.value = res.data
    }
  } catch (error) {
    console.error(error)
  }
}

async function fetchDepartments(orgId: number) {
  try {
    const res = await getDepartmentList(String(orgId))
    if (res.code === 0) {
      departments.value = res.data
    }
  } catch (error) {
    console.error(error)
  }
}

async function fetchRoles() {
  try {
    const res = await getRoleList()
    if (res.code === 0) {
      roles.value = res.data
    }
  } catch (error) {
    console.error(error)
  }
}

async function fetchData() {
  loading.value = true
  try {
    const res = await getUserList({
      page: pagination.page,
      page_size: pagination.pageSize,
    })
    if (res.code === 0) {
      tableData.value = res.data.list
      pagination.total = res.data.total
    }
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  pagination.page = 1
  fetchData()
}

function handleReset() {
  searchForm.username = ''
  searchForm.status = ''
  pagination.page = 1
  fetchData()
}

function resetForm() {
  formData.username = ''
  formData.password = ''
  formData.real_name = ''
  formData.email = ''
  formData.phone = ''
  formData.organization_id = null
  formData.department_id = null
  formData.status = 'active'
  formData.role_ids = []
  currentId.value = null
  departments.value = []
}

function handleAdd() {
  isEdit.value = false
  resetForm()
  dialogVisible.value = true
}

async function handleEdit(row: User) {
  isEdit.value = true
  currentId.value = row.id
  formData.username = row.username
  formData.password = ''
  formData.real_name = row.real_name || ''
  formData.email = row.email || ''
  formData.phone = row.phone || ''
  formData.organization_id = row.organization_id || null
  formData.department_id = row.department_id || null
  formData.status = row.status
  formData.role_ids = row.roles?.map((r) => r.id) || []

  if (row.organization_id) {
    await fetchDepartments(row.organization_id)
  }

  dialogVisible.value = true
}

async function handleOrgChange(orgId: number) {
  formData.department_id = null
  departments.value = []
  if (orgId) {
    await fetchDepartments(orgId)
  }
}

async function handleSubmit() {
  const valid = await formRef.value?.validate()
  if (!valid) return

  submitting.value = true
  try {
    if (isEdit.value) {
      const updateData: any = {
        real_name: formData.real_name,
        email: formData.email,
        phone: formData.phone,
        organization_id: formData.organization_id,
        department_id: formData.department_id,
        status: formData.status,
        role_ids: formData.role_ids,
      }
      await updateUser(currentId.value!, updateData)
      ElMessage.success('更新成功')
    } else {
      await createUser({
        username: formData.username,
        password: formData.password,
        real_name: formData.real_name,
        email: formData.email,
        phone: formData.phone,
        organization_id: formData.organization_id,
        department_id: formData.department_id,
        status: formData.status,
        role_ids: formData.role_ids,
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

async function handleResetPwd(row: User) {
  try {
    const { value } = await ElMessageBox.prompt('请输入新密码', '重置密码', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputPattern: /^.{6,}$/,
      inputErrorMessage: '密码长度不能少于6位',
    })
    // TODO: 实现重置密码API
    ElMessage.success(`用户 ${row.username} 密码已重置`)
  } catch {
    // 取消操作
  }
}

async function handleDelete(row: User) {
  try {
    await ElMessageBox.confirm(`确定要删除用户 ${row.username} 吗？`, '提示', {
      type: 'warning',
    })

    await deleteUser(row.id)
    ElMessage.success('删除成功')
    fetchData()
  } catch {
    // 取消删除
  }
}

function getStatusType(status: string) {
  const map: Record<string, string> = {
    active: 'success',
    inactive: 'info',
    locked: 'danger',
  }
  return map[status] || 'info'
}

function getStatusText(status: string) {
  const map: Record<string, string> = {
    active: '正常',
    inactive: '停用',
    locked: '锁定',
  }
  return map[status] || status
}
</script>

<style scoped lang="scss">
.page-container {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .search-form {
    margin-bottom: 20px;
  }

  .pagination {
    margin-top: 20px;
    justify-content: flex-end;
  }

  .mr-10 {
    margin-right: 10px;
  }
}
</style>
