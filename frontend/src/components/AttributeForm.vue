<template>
  <div class="attribute-form" v-if="schemas.length > 0">
    <template v-for="schema in schemas" :key="schema.attr_type">
      <el-divider v-if="schema.schema_config?.label" content-position="left">
        {{ schema.schema_config.label }}
      </el-divider>
      <el-divider v-else-if="schema.attr_type !== 'main'" content-position="left">
        {{ getAttrTypeLabel(schema.attr_type) }}
      </el-divider>

      <el-row :gutter="20">
        <el-col
          v-for="field in schema.schema_config?.fields"
          :key="field.key"
          :span="getFieldSpan(field.type)"
        >
          <el-form-item
            :label="field.label"
            :prop="getFieldProp(schema.attr_type, field.key)"
            :rules="getFieldRules(field)"
          >
            <!-- 文本输入 -->
            <el-input
              v-if="field.type === 'text'"
              v-model="modelValue[getFieldKey(schema.attr_type, field.key)]"
              :placeholder="`请输入${field.label}`"
              clearable
            />
            <!-- 数字输入 -->
            <div v-else-if="field.type === 'number'" class="number-input">
              <el-input-number
                v-model="modelValue[getFieldKey(schema.attr_type, field.key)]"
                :placeholder="`请输入${field.label}`"
                style="width: 100%"
              />
              <span v-if="field.unit" class="unit">{{ field.unit }}</span>
            </div>
            <!-- 下拉选择 -->
            <el-select
              v-else-if="field.type === 'select'"
              v-model="modelValue[getFieldKey(schema.attr_type, field.key)]"
              :placeholder="`请选择${field.label}`"
              clearable
              style="width: 100%"
            >
              <el-option
                v-for="option in field.options"
                :key="option"
                :label="option"
                :value="option"
              />
            </el-select>
            <!-- 日期选择 -->
            <el-date-picker
              v-else-if="field.type === 'date'"
              v-model="modelValue[getFieldKey(schema.attr_type, field.key)]"
              type="date"
              :placeholder="`请选择${field.label}`"
              style="width: 100%"
            />
            <!-- 多行文本 -->
            <el-input
              v-else-if="field.type === 'textarea'"
              v-model="modelValue[getFieldKey(schema.attr_type, field.key)]"
              type="textarea"
              :rows="3"
              :placeholder="`请输入${field.label}`"
            />
          </el-form-item>
        </el-col>
      </el-row>
    </template>
  </div>
  <el-empty v-else description="暂无属性配置" :image-size="60" />
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { AttributeSchema } from '@/api/attribute'

interface Props {
  schemas: AttributeSchema[]
  modelValue: Record<string, any>
  entityType: 'MATERIAL' | 'DOCUMENT'
}

const props = defineProps<Props>()
const emit = defineEmits<{
  'update:modelValue': [value: Record<string, any>]
}>()

// 获取属性类型标签
function getAttrTypeLabel(type: string) {
  const map: Record<string, string> = {
    main: '基本信息',
    description: '描述信息',
    specification: '规格参数',
    custom: '自定义属性'
  }
  return map[type] || type
}

// 获取字段占用列数
function getFieldSpan(type: string) {
  return type === 'textarea' ? 24 : 12
}

// 获取字段属性路径
function getFieldProp(attrType: string, key: string) {
  if (attrType === 'main') {
    return `attributes.${key}`
  }
  return `dynamicAttributes.${attrType}.${key}`
}

// 获取字段key（用于modelValue）
function getFieldKey(attrType: string, key: string) {
  if (attrType === 'main') {
    return `main_${key}`
  }
  return `${attrType}_${key}`
}

// 获取字段校验规则
function getFieldRules(field: any) {
  if (field.required) {
    return [{ required: true, message: `请输入${field.label}`, trigger: 'blur' }]
  }
  return []
}
</script>

<style scoped lang="scss">
.attribute-form {
  .number-input {
    display: flex;
    align-items: center;
    gap: 8px;

    .unit {
      color: var(--el-text-color-secondary);
      white-space: nowrap;
    }
  }

  :deep(.el-divider__text) {
    font-size: 14px;
    font-weight: 500;
    color: var(--el-text-color-regular);
  }
}
</style>
