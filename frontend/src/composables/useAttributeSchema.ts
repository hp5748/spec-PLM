import { ref, watch } from 'vue'
import { getAllActiveSchemas, type AttributeSchema } from '@/api/attribute'

/**
 * 属性Schema加载Hook
 * @param entityType 实体类型 MATERIAL/DOCUMENT
 * @param typeCode 类型编码
 */
export function useAttributeSchema(entityType: 'MATERIAL' | 'DOCUMENT', typeCode: string | undefined) {
  const schemas = ref<AttributeSchema[]>([])
  const loading = ref(false)

  async function loadSchemas(code: string) {
    if (!code) {
      schemas.value = []
      return
    }

    loading.value = true
    try {
      const res = await getAllActiveSchemas(entityType, code)
      if (res.code === 0) {
        schemas.value = res.data || []
      }
    } catch (error) {
      console.error('加载属性Schema失败', error)
      schemas.value = []
    } finally {
      loading.value = false
    }
  }

  // 监听typeCode变化
  watch(() => typeCode, (newCode) => {
    if (newCode) {
      loadSchemas(newCode)
    }
  }, { immediate: true })

  return {
    schemas,
    loading,
    loadSchemas
  }
}

/**
 * 解析属性数据用于表单
 * @param mainAttributes 主属性JSON字符串或对象
 * @param dynamicAttributes 动态属性数组
 * @param schemas Schema列表
 */
export function parseAttributesToForm(
  mainAttributes: string | Record<string, any> | undefined,
  dynamicAttributes: Array<{ attr_type: string; attr_key: string; attr_value: string; unit?: string }> | undefined,
  schemas: AttributeSchema[]
): Record<string, any> {
  const result: Record<string, any> = {}

  // 解析主属性
  if (mainAttributes) {
    const mainObj = typeof mainAttributes === 'string' ? JSON.parse(mainAttributes || '{}') : mainAttributes
    Object.entries(mainObj).forEach(([key, value]) => {
      result[`main_${key}`] = value
    })
  }

  // 解析动态属性
  if (dynamicAttributes && dynamicAttributes.length > 0) {
    dynamicAttributes.forEach(attr => {
      result[`${attr.attr_type}_${attr.attr_key}`] = attr.attr_value
    })
  }

  return result
}

/**
 * 将表单数据转换为主属性和动态属性
 * @param formData 表单数据
 * @param schemas Schema列表
 */
export function formToAttributes(
  formData: Record<string, any>,
  schemas: AttributeSchema[]
): {
  mainAttributes: Record<string, any>
  dynamicAttributes: Array<{ attr_type: string; attr_key: string; attr_value: string; unit?: string }>
} {
  const mainAttributes: Record<string, any> = {}
  const dynamicAttributes: Array<{ attr_type: string; attr_key: string; attr_value: string; unit?: string }> = []

  // 遍历表单数据
  Object.entries(formData).forEach(([key, value]) => {
    if (value === undefined || value === null || value === '') return

    // 主属性
    if (key.startsWith('main_')) {
      const attrKey = key.substring(5)
      mainAttributes[attrKey] = value
    }
    // 动态属性
    else if (key.includes('_')) {
      const underscoreIndex = key.indexOf('_')
      const attrType = key.substring(0, underscoreIndex)
      const attrKey = key.substring(underscoreIndex + 1)

      // 检查是否是有效的动态属性类型
      if (['description', 'specification', 'custom'].includes(attrType)) {
        // 查找对应的schema获取单位
        const schema = schemas.find(s => s.attr_type === attrType)
        const field = schema?.schema_config?.fields?.find(f => f.key === attrKey)

        dynamicAttributes.push({
          attr_type: attrType,
          attr_key: attrKey,
          attr_value: String(value),
          unit: field?.unit
        })
      }
    }
  })

  return { mainAttributes, dynamicAttributes }
}

/**
 * 获取列表显示的属性列配置
 * @param schemas Schema列表
 * @param showAttributes 是否显示属性
 */
export function getAttributeColumns(schemas: AttributeSchema[], showAttributes: boolean) {
  if (!showAttributes || schemas.length === 0) {
    return []
  }

  const columns: Array<{ prop: string; label: string; width?: number }> = []

  schemas.forEach(schema => {
    // 只显示主属性的前几个字段
    if (schema.attr_type === 'main' && schema.schema_config?.fields) {
      schema.schema_config.fields.slice(0, 3).forEach(field => {
        columns.push({
          prop: `main_${field.key}`,
          label: field.label,
          width: 120
        })
      })
    }
  })

  return columns
}
