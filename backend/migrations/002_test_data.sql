-- PLM系统测试数据
-- 执行前请确保已运行 001_init_schema.sql

SET NAMES utf8mb4;
USE plm;

-- =============================================
-- 1. 组织数据
-- =============================================
INSERT INTO organizations (id, name, code, description, status, created_by) VALUES
(1, '示例科技有限公司', 'DEMO', 'PLM系统演示组织', 'active', 1);

-- =============================================
-- 2. 部门数据（树形结构）
-- =============================================
INSERT INTO departments (id, organization_id, parent_id, name, code, sort_order, status, created_by) VALUES
-- 一级部门
(1, 1, NULL, '总经办', 'GM', 1, 'active', 1),
(2, 1, NULL, '研发中心', 'RD', 2, 'active', 1),
(3, 1, NULL, '生产部', 'PROD', 3, 'active', 1),
(4, 1, NULL, '质量部', 'QC', 4, 'active', 1),
-- 二级部门 - 研发中心
(5, 1, 2, '机械设计部', 'RD-MECH', 1, 'active', 1),
(6, 1, 2, '电气设计部', 'RD-ELEC', 2, 'active', 1),
(7, 1, 2, '软件研发部', 'RD-SW', 3, 'active', 1),
-- 二级部门 - 生产部
(8, 1, 3, '机加工车间', 'PROD-MACH', 1, 'active', 1),
(9, 1, 3, '装配车间', 'PROD-ASSY', 2, 'active', 1),
-- 三级部门 - 机械设计部
(10, 1, 5, '结构设计组', 'RD-MECH-STR', 1, 'active', 1),
(11, 1, 5, '传动设计组', 'RD-MECH-TRN', 2, 'active', 1);

-- =============================================
-- 3. 属性Schema配置
-- =============================================

-- 物料 - 零件(PART) - 主属性
INSERT INTO attribute_schemas (entity_type, type_code, attr_type, schema_name, schema_config, version, status, is_active, created_by) VALUES
('MATERIAL', 'PART', 'main', '零件主属性', '{"fields": [{"key": "weight", "label": "重量", "type": "number", "unit": "kg", "required": false}, {"key": "material", "label": "材料", "type": "text", "required": false}, {"key": "supplier", "label": "供应商", "type": "text", "required": false}, {"key": "color", "label": "颜色", "type": "select", "options": ["红色", "蓝色", "黑色", "银色", "白色"], "required": false}]}', 'AA', 'RELEASED', TRUE, 1),
-- 物料 - 零件(PART) - 规格属性
('MATERIAL', 'PART', 'specification', '零件规格参数', '{"label": "规格参数", "fields": [{"key": "长度", "label": "长度", "type": "number", "unit": "mm", "required": false}, {"key": "宽度", "label": "宽度", "type": "number", "unit": "mm", "required": false}, {"key": "高度", "label": "高度", "type": "number", "unit": "mm", "required": false}, {"key": "公差", "label": "公差等级", "type": "text", "required": false}]}', 'AA', 'RELEASED', TRUE, 1),
-- 物料 - 零件(PART) - 描述属性
('MATERIAL', 'PART', 'description', '零件描述信息', '{"label": "描述信息", "fields": [{"key": "用途", "label": "用途说明", "type": "textarea", "required": false}, {"key": "备注", "label": "备注", "type": "textarea", "required": false}]}', 'AA', 'RELEASED', TRUE, 1),

-- 物料 - 组件(ASSEMBLY) - 主属性
('MATERIAL', 'ASSEMBLY', 'main', '组件主属性', '{"fields": [{"key": "weight", "label": "重量", "type": "number", "unit": "kg", "required": false}, {"key": "assembly_type", "label": "装配类型", "type": "select", "options": ["机械装配", "电气装配", "机电混合"], "required": false}]}', 'AA', 'RELEASED', TRUE, 1),
-- 物料 - 组件(ASSEMBLY) - 规格属性
('MATERIAL', 'ASSEMBLY', 'specification', '组件规格参数', '{"label": "规格参数", "fields": [{"key": "功率", "label": "功率", "type": "number", "unit": "W", "required": false}, {"key": "电压", "label": "电压", "type": "number", "unit": "V", "required": false}]}', 'AA', 'RELEASED', TRUE, 1);

-- =============================================
-- 4. 物料数据
-- =============================================
INSERT INTO materials (id, item_id, item_name, description, item_type, version, unit, status, attributes, created_by) VALUES
-- 零件
(1, 'P-001', '传动轴', '主传动轴组件', 'PART', 'AA', '件', 'RELEASED', '{"weight": 2.5, "material": "45#钢", "supplier": "华东钢铁", "color": "银色"}', 1),
(2, 'P-002', '轴承座', '深沟球轴承座', 'PART', 'AA', '件', 'RELEASED', '{"weight": 0.8, "material": "铸铁HT200", "supplier": "北方铸造"}', 1),
(3, 'P-003', '齿轮', '直齿圆柱齿轮', 'PART', 'AA', '件', 'RELEASED', '{"weight": 1.2, "material": "20CrMnTi", "supplier": "华东钢铁", "color": "银色"}', 1),
(4, 'P-004', '电机安装板', '电机安装固定板', 'PART', 'AA', '件', 'DRAFT', '{"weight": 3.0, "material": "Q235", "supplier": "本地钢厂"}', 1),
(5, 'P-005', '端盖', '密封端盖', 'PART', 'AA', '件', 'DRAFT', '{"weight": 0.3, "material": "铝合金", "color": "银色"}', 1),
-- 组件
(6, 'A-001', '传动组件', '包含传动轴、齿轮、轴承的完整传动组件', 'ASSEMBLY', 'AA', '套', 'RELEASED', '{"weight": 5.2, "assembly_type": "机械装配"}', 1),
(7, 'A-002', '电机驱动组件', '电机及安装配件', 'ASSEMBLY', 'AA', '套', 'REVIEWING', '{"weight": 15.0, "assembly_type": "机电混合"}', 1),
-- 原材料
(8, 'R-001', '钢板10mm', '10mm厚Q235钢板', 'RAW_MATERIAL', 'AA', '张', 'RELEASED', '{}', 1),
(9, 'R-002', '圆钢50mm', '50mm直径45#圆钢', 'RAW_MATERIAL', 'AA', '米', 'RELEASED', '{}', 1);

-- 物料动态属性
INSERT INTO material_attributes (material_id, attr_type, attr_key, attr_value, unit, sort_order) VALUES
-- P-001 传动轴 规格属性
(1, 'specification', '长度', '350', 'mm', 1),
(1, 'specification', '直径', '50', 'mm', 2),
(1, 'specification', '公差', 'h6', '', 3),
(1, 'description', '用途', '主传动系统核心部件', '', 1),
-- P-002 轴承座 规格属性
(2, 'specification', '长度', '80', 'mm', 1),
(2, 'specification', '宽度', '60', 'mm', 2),
(2, 'specification', '高度', '50', 'mm', 3),
-- P-003 齿轮 规格属性
(3, 'specification', '齿数', '36', '', 1),
(3, 'specification', '模数', '2', 'mm', 2),
(3, 'specification', '压力角', '20', '度', 3),
-- A-001 传动组件 规格属性
(6, 'specification', '功率', '5.5', 'kW', 1),
(6, 'specification', '转速', '1450', 'rpm', 2),
-- A-002 电机驱动组件 规格属性
(7, 'specification', '功率', '7.5', 'kW', 1),
(7, 'specification', '电压', '380', 'V', 2);

-- =============================================
-- 5. BOM视图数据
-- =============================================
INSERT INTO bom_views (id, name, root_material_id, root_version, is_exact, status, description, created_by) VALUES
(1, '传动组件BOM', 6, 'AA', TRUE, 'RELEASED', '传动组件的精确BOM清单', 1),
(2, '电机驱动组件BOM', 7, NULL, FALSE, 'DRAFT', '电机驱动组件的非精确BOM（设计阶段）', 1);

-- =============================================
-- 6. BOM结构数据
-- =============================================
-- 传动组件BOM (精确BOM)
INSERT INTO bom_items (bom_view_id, parent_id, material_id, version, quantity, sort_order, level) VALUES
-- 一级子件（直接挂在根节点下，parent_id为NULL表示根节点的子件）
(1, NULL, 1, 'AA', 1.0000, 1, 1),  -- 传动轴 x 1
(1, NULL, 2, 'AA', 2.0000, 2, 1),  -- 轴承座 x 2
(1, NULL, 3, 'AA', 2.0000, 3, 1),  -- 齿轮 x 2
(1, NULL, 5, 'AA', 2.0000, 4, 1);  -- 端盖 x 2

-- 电机驱动组件BOM (非精确BOM)
INSERT INTO bom_items (bom_view_id, parent_id, material_id, version, quantity, sort_order, level) VALUES
(2, NULL, 6, NULL, 1.0000, 1, 1),  -- 传动组件 x 1
(2, NULL, 4, NULL, 1.0000, 2, 1),  -- 电机安装板 x 1
(2, NULL, 5, NULL, 4.0000, 3, 1);  -- 端盖 x 4

-- =============================================
-- 7. 文档数据
-- =============================================
INSERT INTO documents (id, doc_id, doc_name, doc_type, file_path, file_name, file_size, mime_type, version, status, attributes, created_by) VALUES
(1, 'D-001', '传动轴零件图', 'CAD', '/documents/d-001/传动轴零件图.dwg', '传动轴零件图.dwg', 256000, 'application/acad', 'AA', 'RELEASED', '{"scale": "1:1", "sheet": "A3"}', 1),
(2, 'D-002', '传动轴加工工艺', 'PDF', '/documents/d-002/传动轴加工工艺.pdf', '传动轴加工工艺.pdf', 128000, 'application/pdf', 'AA', 'RELEASED', '{"pages": 5}', 1),
(3, 'D-003', '传动组件装配图', 'CAD', '/documents/d-003/传动组件装配图.dwg', '传动组件装配图.dwg', 384000, 'application/acad', 'AA', 'RELEASED', '{"scale": "1:2", "sheet": "A1"}', 1),
(4, 'D-004', '产品使用说明书', 'PDF', '/documents/d-004/产品使用说明书.pdf', '产品使用说明书.pdf', 512000, 'application/pdf', 'AA', 'RELEASED', '{"pages": 20}', 1);

-- =============================================
-- 8. 更新用户所属组织
-- =============================================
UPDATE users SET organization_id = 1, department_id = 1 WHERE id = 1;

SELECT '测试数据导入完成!' AS message;
