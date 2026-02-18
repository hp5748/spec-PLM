-- PLM系统初始化脚本
-- 设置客户端字符集
SET NAMES utf8mb4;
SET CHARACTER SET utf8mb4;

-- 创建数据库（如果不存在）
CREATE DATABASE IF NOT EXISTS plm CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE plm;

-- 组织表
CREATE TABLE IF NOT EXISTS organizations (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL COMMENT '组织名称',
    code VARCHAR(50) NOT NULL UNIQUE COMMENT '组织编码',
    description TEXT COMMENT '描述',
    logo VARCHAR(255) COMMENT 'Logo',
    status ENUM('active', 'inactive') DEFAULT 'active' COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by BIGINT NOT NULL DEFAULT 0,
    deleted_at TIMESTAMP NULL,
    INDEX idx_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='组织表';

-- 部门表
CREATE TABLE IF NOT EXISTS departments (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    organization_id BIGINT NOT NULL COMMENT '组织ID',
    parent_id BIGINT COMMENT '父部门ID',
    name VARCHAR(100) NOT NULL COMMENT '部门名称',
    code VARCHAR(50) COMMENT '部门编码',
    sort_order INT DEFAULT 0 COMMENT '排序',
    manager_id BIGINT COMMENT '部门负责人ID',
    status ENUM('active', 'inactive') DEFAULT 'active' COMMENT '状态',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by BIGINT NOT NULL DEFAULT 0,
    deleted_at TIMESTAMP NULL,
    INDEX idx_organization (organization_id),
    INDEX idx_parent (parent_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='部门表';

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    password VARCHAR(255) NOT NULL COMMENT '密码(bcrypt)',
    email VARCHAR(100) COMMENT '邮箱',
    phone VARCHAR(20) COMMENT '手机号',
    real_name VARCHAR(50) COMMENT '真实姓名',
    avatar VARCHAR(255) COMMENT '头像URL',
    organization_id BIGINT COMMENT '组织ID',
    department_id BIGINT COMMENT '部门ID',
    status ENUM('active', 'inactive', 'locked') DEFAULT 'active' COMMENT '状态',
    last_login_at TIMESTAMP NULL COMMENT '最后登录时间',
    login_fail_count INT DEFAULT 0 COMMENT '登录失败次数',
    locked_until TIMESTAMP NULL COMMENT '锁定截止时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by BIGINT NOT NULL DEFAULT 0,
    deleted_at TIMESTAMP NULL,
    INDEX idx_username (username),
    INDEX idx_organization (organization_id),
    INDEX idx_department (department_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 角色表
CREATE TABLE IF NOT EXISTS roles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    organization_id BIGINT COMMENT '组织ID（为空则全局角色）',
    name VARCHAR(50) NOT NULL COMMENT '角色名称',
    code VARCHAR(50) NOT NULL COMMENT '角色编码',
    description TEXT COMMENT '描述',
    is_system BOOLEAN DEFAULT FALSE COMMENT '是否系统内置角色',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by BIGINT NOT NULL DEFAULT 0,
    deleted_at TIMESTAMP NULL,
    INDEX idx_organization (organization_id),
    UNIQUE KEY uk_org_code (organization_id, code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

-- 权限表
CREATE TABLE IF NOT EXISTS permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL COMMENT '权限名称',
    code VARCHAR(100) NOT NULL UNIQUE COMMENT '权限编码',
    module VARCHAR(50) NOT NULL COMMENT '所属模块',
    description TEXT COMMENT '描述',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_module (module)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='权限表';

-- 用户角色关联表
CREATE TABLE IF NOT EXISTS user_roles (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL COMMENT '用户ID',
    role_id BIGINT NOT NULL COMMENT '角色ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_role (user_id, role_id),
    INDEX idx_user (user_id),
    INDEX idx_role (role_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户角色关联表';

-- 角色权限关联表
CREATE TABLE IF NOT EXISTS role_permissions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    role_id BIGINT NOT NULL COMMENT '角色ID',
    permission_id BIGINT NOT NULL COMMENT '权限ID',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_role_permission (role_id, permission_id),
    INDEX idx_role (role_id),
    INDEX idx_permission (permission_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色权限关联表';

-- 物料表
CREATE TABLE IF NOT EXISTS materials (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    item_id VARCHAR(50) NOT NULL COMMENT '物料编码',
    item_name VARCHAR(200) NOT NULL COMMENT '物料名称',
    description TEXT COMMENT '描述',
    item_type ENUM('PART', 'ASSEMBLY', 'RAW_MATERIAL', 'TOOL') NOT NULL COMMENT '物料类型',
    version VARCHAR(2) NOT NULL DEFAULT 'AA' COMMENT '版本号',
    unit VARCHAR(20) COMMENT '计量单位',
    status ENUM('DRAFT', 'REVIEWING', 'RELEASED', 'REJECTED', 'OBSOLETE') DEFAULT 'DRAFT' COMMENT '状态',
    attributes JSON COMMENT '主属性JSON',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by BIGINT NOT NULL DEFAULT 0,
    deleted_at TIMESTAMP NULL,
    UNIQUE KEY uk_item_version (item_id, version),
    INDEX idx_item_id (item_id),
    INDEX idx_item_type (item_type),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='物料表';

-- 物料动态属性表
CREATE TABLE IF NOT EXISTS material_attributes (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    material_id BIGINT NOT NULL COMMENT '物料ID',
    attr_type VARCHAR(50) NOT NULL COMMENT '属性类型',
    attr_key VARCHAR(100) NOT NULL COMMENT '属性键',
    attr_value TEXT COMMENT '属性值',
    unit VARCHAR(20) COMMENT '单位',
    sort_order INT DEFAULT 0 COMMENT '排序',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_material (material_id),
    INDEX idx_material_type (material_id, attr_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='物料动态属性表';

-- 文档表
CREATE TABLE IF NOT EXISTS documents (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    doc_id VARCHAR(50) NOT NULL COMMENT '文档编码',
    doc_name VARCHAR(200) NOT NULL COMMENT '文档名称',
    doc_type VARCHAR(50) NOT NULL COMMENT '文档类型',
    file_path VARCHAR(500) NOT NULL COMMENT '文件存储路径',
    file_name VARCHAR(255) COMMENT '原始文件名',
    file_size BIGINT COMMENT '文件大小（字节）',
    mime_type VARCHAR(100) COMMENT 'MIME类型',
    version VARCHAR(2) NOT NULL DEFAULT 'AA' COMMENT '版本号',
    status ENUM('DRAFT', 'CHECKED_IN', 'RELEASED') DEFAULT 'DRAFT' COMMENT '状态',
    attributes JSON COMMENT '扩展属性JSON',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by BIGINT NOT NULL DEFAULT 0,
    deleted_at TIMESTAMP NULL,
    UNIQUE KEY uk_doc_version (doc_id, version),
    INDEX idx_doc_id (doc_id),
    INDEX idx_doc_type (doc_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文档表';

-- 初始化系统角色
INSERT INTO roles (name, code, description, is_system, created_by) VALUES
('系统管理员', 'ADMIN', '系统配置、用户管理、全部数据访问', TRUE, 0),
('设计工程师', 'DESIGNER', '创建/编辑物料、BOM、文档', TRUE, 0),
('审批人员', 'REVIEWER', '审批流程、查看数据', TRUE, 0),
('普通用户', 'VIEWER', '只读访问数据', TRUE, 0);

-- 初始化权限
INSERT INTO permissions (name, code, module, description) VALUES
-- 用户管理
('查看用户', 'user:view', 'user', '用户列表和详情'),
('创建用户', 'user:create', 'user', '新建用户'),
('编辑用户', 'user:edit', 'user', '修改用户信息'),
('删除用户', 'user:delete', 'user', '删除用户'),
-- 角色管理
('查看角色', 'role:view', 'role', '角色列表和详情'),
('管理角色', 'role:manage', 'role', '创建/编辑/删除角色'),
-- 物料管理
('查看物料', 'material:view', 'material', '物料列表和详情'),
('创建物料', 'material:create', 'material', '新建物料'),
('编辑物料', 'material:edit', 'material', '修改物料信息'),
('删除物料', 'material:delete', 'material', '删除物料'),
('发布物料', 'material:release', 'material', '发布物料'),
-- 文档管理
('查看文档', 'document:view', 'document', '文档列表和详情'),
('上传文档', 'document:upload', 'document', '上传新文档'),
('编辑文档', 'document:edit', 'document', '修改文档信息'),
('删除文档', 'document:delete', 'document', '删除文档'),
('下载文档', 'document:download', 'document', '下载文件'),
-- BOM管理
('查看BOM', 'bom:view', 'bom', 'BOM列表和详情'),
('创建BOM', 'bom:create', 'bom', '新建BOM视图'),
('编辑BOM', 'bom:edit', 'bom', '修改BOM结构'),
('删除BOM', 'bom:delete', 'bom', '删除BOM视图'),
('导出BOM', 'bom:export', 'bom', '导出Excel'),
-- 流程管理
('查看流程', 'workflow:view', 'workflow', '流程列表和详情'),
('发起流程', 'workflow:initiate', 'workflow', '发起审批'),
('审批流程', 'workflow:approve', 'workflow', '同意/驳回/转交'),
-- 系统管理
('系统管理', 'admin:all', 'admin', '全部权限');

-- 初始化管理员用户（密码：admin123）
-- bcrypt hash for 'admin123'
INSERT INTO users (username, password, real_name, status, created_by) VALUES
('admin', '$2a$10$DcvP67Sgy2GVFMWe8ZC5WOOYdNIar128GNPII0flHcCrp9.NUWapy', '系统管理员', 'active', 0);

-- 给管理员分配ADMIN角色
INSERT INTO user_roles (user_id, role_id) VALUES (1, 1);

-- 给ADMIN角色分配所有权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 1, id FROM permissions;

-- 给DESIGNER角色分配权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 2, id FROM permissions WHERE code IN (
    'material:view', 'material:create', 'material:edit', 'material:delete',
    'document:view', 'document:upload', 'document:edit', 'document:download',
    'bom:view', 'bom:create', 'bom:edit', 'bom:export',
    'workflow:view', 'workflow:initiate'
);

-- 给REVIEWER角色分配权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 3, id FROM permissions WHERE code IN (
    'material:view', 'material:release',
    'document:view', 'document:download',
    'bom:view',
    'workflow:view', 'workflow:approve'
);

-- 给VIEWER角色分配权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 4, id FROM permissions WHERE code IN (
    'material:view', 'document:view', 'bom:view', 'workflow:view'
);

-- 版本历史表
CREATE TABLE IF NOT EXISTS version_histories (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    entity_type VARCHAR(50) NOT NULL COMMENT '实体类型(MATERIAL/DOCUMENT)',
    entity_id BIGINT NOT NULL COMMENT '实体ID',
    old_version VARCHAR(2) COMMENT '原版本',
    new_version VARCHAR(2) NOT NULL COMMENT '新版本',
    change_reason TEXT COMMENT '变更原因',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by BIGINT NOT NULL,
    INDEX idx_entity (entity_type, entity_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='版本历史表';

-- BOM视图表
CREATE TABLE IF NOT EXISTS bom_views (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(200) NOT NULL COMMENT 'BOM视图名称',
    root_material_id BIGINT NOT NULL COMMENT '根物料ID',
    root_version VARCHAR(2) COMMENT '根物料版本(精确BOM必填)',
    is_exact BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否精确BOM',
    status ENUM('DRAFT', 'REVIEWING', 'RELEASED', 'REJECTED', 'OBSOLETE') DEFAULT 'DRAFT' COMMENT '状态',
    description TEXT COMMENT '描述',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by BIGINT NOT NULL DEFAULT 0,
    deleted_at TIMESTAMP NULL,
    UNIQUE KEY uk_root_material (root_material_id),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='BOM视图表';

-- BOM关系表
CREATE TABLE IF NOT EXISTS bom_items (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    bom_view_id BIGINT NOT NULL COMMENT 'BOM视图ID',
    parent_id BIGINT COMMENT '父节点ID(为空则为根节点)',
    material_id BIGINT NOT NULL COMMENT '物料ID',
    version VARCHAR(2) COMMENT '物料版本(精确BOM必填)',
    quantity DECIMAL(10,4) NOT NULL DEFAULT 1 COMMENT '数量',
    sort_order INT DEFAULT 0 COMMENT '排序',
    level INT DEFAULT 1 COMMENT '层级',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_bom_view (bom_view_id),
    INDEX idx_parent (parent_id),
    INDEX idx_material (material_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='BOM关系表';

-- 属性Schema表
CREATE TABLE IF NOT EXISTS attribute_schemas (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    entity_type VARCHAR(50) NOT NULL COMMENT '实体类型(MATERIAL/DOCUMENT)',
    type_code VARCHAR(50) NOT NULL COMMENT '类型编码(物料类型或文档类型)',
    attr_type VARCHAR(50) NOT NULL COMMENT '属性类型(main/description/specification/custom)',
    schema_name VARCHAR(100) COMMENT 'Schema名称',
    schema_config JSON NOT NULL COMMENT 'Schema定义JSON',
    version VARCHAR(2) NOT NULL DEFAULT 'AA' COMMENT '版本号',
    status ENUM('DRAFT', 'RELEASED') DEFAULT 'DRAFT' COMMENT '状态',
    is_active BOOLEAN DEFAULT TRUE COMMENT '是否当前激活版本',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_by BIGINT NOT NULL DEFAULT 0,
    deleted_at TIMESTAMP NULL,
    INDEX idx_entity_type (entity_type),
    INDEX idx_type_code (type_code),
    INDEX idx_attr_type (attr_type),
    UNIQUE KEY uk_entity_type_version (entity_type, type_code, attr_type, version)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='属性Schema表';

-- 文档动态属性表
CREATE TABLE IF NOT EXISTS document_attributes (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    document_id BIGINT NOT NULL COMMENT '文档ID',
    attr_type VARCHAR(50) NOT NULL COMMENT '属性类型',
    attr_key VARCHAR(100) NOT NULL COMMENT '属性键',
    attr_value TEXT COMMENT '属性值',
    unit VARCHAR(20) COMMENT '单位',
    sort_order INT DEFAULT 0 COMMENT '排序',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_document (document_id),
    INDEX idx_document_type (document_id, attr_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文档动态属性表';

-- 添加属性配置管理权限
INSERT INTO permissions (name, code, module, description) VALUES
('属性配置管理', 'attribute:manage', 'admin', '管理物料/文档属性Schema');

-- 给ADMIN角色添加属性配置权限
INSERT INTO role_permissions (role_id, permission_id)
SELECT 1, id FROM permissions WHERE code = 'attribute:manage';
