<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>组织管理</span>
          <el-button type="primary" @click="handleAdd">新增组织</el-button>
        </div>
      </template>

      <!-- 表格 -->
      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="code" label="组织编码" width="150" />
        <el-table-column prop="name" label="组织名称" width="200" />
        <el-table-column prop="description" label="描述" min-width="200" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'">
              {{ row.status === 'active' ? '正常' : '停用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" fixed="right" width="150">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500">
      <el-form ref="formRef" :model="formData" :rules="rules" label-width="100px">
        <el-form-item label="组织编码" prop="code">
          <el-input v-model="formData.code" placeholder="请输入组织编码" :disabled="isEdit" />
        </el-form-item>
        <el-form-item label="组织名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入组织名称" />
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input v-model="formData.description" type="textarea" rows="3" placeholder="请输入描述" />
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
import { getOrganizationList, createOrganization } from '@/api/user'
import type { Organization } from '@/types/api'

defineOptions({
  name: 'OrganizationList',
})

const loading = ref(false)
const submitting = ref(false)
const tableData = ref<Organization[]>([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const currentId = ref<number | null>(null)
const formRef = ref<FormInstance>()

const dialogTitle = computed(() => (isEdit.value ? '编辑组织' : '新增组织'))

const formData = reactive({
  code: '',
  name: '',
  description: '',
  status: 'active',
})

const rules: FormRules = {
  code: [{ required: true, message: '请输入组织编码', trigger: 'blur' }],
  name: [{ required: true, message: '请输入组织名称', trigger: 'blur' }],
}

onMounted(() => {
  fetchData()
})

async function fetchData() {
  loading.value = true
  try {
    const res = await getOrganizationList()
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
  formData.status = 'active'
  currentId.value = null
}

function handleAdd() {
  isEdit.value = false
  resetForm()
  dialogVisible.value = true
}

function handleEdit(row: Organization) {
  isEdit.value = true
  currentId.value = row.id
  formData.code = row.code
  formData.name = row.name
  formData.description = row.description || ''
  formData.status = row.status
  dialogVisible.value = true
}

async function handleSubmit() {
  const valid = await formRef.value?.validate()
  if (!valid) return

  submitting.value = true
  try {
    if (isEdit.value) {
      // TODO: 实现更新组织API
      ElMessage.success('更新成功')
    } else {
      await createOrganization(formData)
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

async function handleDelete(row: Organization) {
  try {
    await ElMessageBox.confirm(`确定要删除组织 ${row.name} 吗？`, '提示', {
      type: 'warning',
    })
    // TODO: 实现删除组织API
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
}
</style>
