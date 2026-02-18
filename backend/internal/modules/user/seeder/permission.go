package seeder

import (
	"log"

	"gorm.io/gorm"
	"plm/internal/modules/user/model"
)

// DefaultPermissions 定义系统所有权限
// 每次迭代添加新功能时，需要在这里添加对应的权限
var DefaultPermissions = []model.Permission{
	// 用户管理
	{Name: "查看用户", Code: "user:view", Module: "user", Description: "用户列表和详情"},
	{Name: "创建用户", Code: "user:create", Module: "user", Description: "新建用户"},
	{Name: "编辑用户", Code: "user:edit", Module: "user", Description: "修改用户信息"},
	{Name: "删除用户", Code: "user:delete", Module: "user", Description: "删除用户"},
	// 角色管理
	{Name: "查看角色", Code: "role:view", Module: "role", Description: "角色列表和详情"},
	{Name: "管理角色", Code: "role:manage", Module: "role", Description: "创建/编辑/删除角色"},
	// 物料管理
	{Name: "查看物料", Code: "material:view", Module: "material", Description: "物料列表和详情"},
	{Name: "创建物料", Code: "material:create", Module: "material", Description: "新建物料"},
	{Name: "编辑物料", Code: "material:edit", Module: "material", Description: "修改物料信息"},
	{Name: "删除物料", Code: "material:delete", Module: "material", Description: "删除物料"},
	{Name: "发布物料", Code: "material:release", Module: "material", Description: "发布物料"},
	// 文档管理
	{Name: "查看文档", Code: "document:view", Module: "document", Description: "文档列表和详情"},
	{Name: "上传文档", Code: "document:upload", Module: "document", Description: "上传新文档"},
	{Name: "编辑文档", Code: "document:edit", Module: "document", Description: "修改文档信息"},
	{Name: "删除文档", Code: "document:delete", Module: "document", Description: "删除文档"},
	{Name: "下载文档", Code: "document:download", Module: "document", Description: "下载文件"},
	// BOM管理
	{Name: "查看BOM", Code: "bom:view", Module: "bom", Description: "BOM列表和详情"},
	{Name: "创建BOM", Code: "bom:create", Module: "bom", Description: "新建BOM视图"},
	{Name: "编辑BOM", Code: "bom:edit", Module: "bom", Description: "修改BOM结构"},
	{Name: "删除BOM", Code: "bom:delete", Module: "bom", Description: "删除BOM视图"},
	{Name: "导出BOM", Code: "bom:export", Module: "bom", Description: "导出Excel"},
	// 流程管理
	{Name: "查看流程", Code: "workflow:view", Module: "workflow", Description: "流程列表和详情"},
	{Name: "发起流程", Code: "workflow:initiate", Module: "workflow", Description: "发起审批"},
	{Name: "审批流程", Code: "workflow:approve", Module: "workflow", Description: "同意/驳回/转交"},
	// 系统管理
	{Name: "系统管理", Code: "admin:all", Module: "admin", Description: "全部权限"},
}

// SyncPermissions 同步权限到数据库
// 只添加不存在的权限，不会删除或修改已有权限
func SyncPermissions(db *gorm.DB) error {
	log.Println("[PermissionSeeder] Starting permission sync...")

	var synced, skipped int

	for _, perm := range DefaultPermissions {
		var existing model.Permission
		result := db.Where("code = ?", perm.Code).First(&existing)

		if result.Error == gorm.ErrRecordNotFound {
			// 权限不存在，创建它
			if err := db.Create(&perm).Error; err != nil {
				log.Printf("[PermissionSeeder] Failed to create permission %s: %v", perm.Code, err)
				return err
			}
			synced++
			log.Printf("[PermissionSeeder] Created permission: %s (%s)", perm.Code, perm.Name)
		} else if result.Error != nil {
			return result.Error
		} else {
			skipped++
		}
	}

	log.Printf("[PermissionSeeder] Sync complete. Created: %d, Skipped: %d", synced, skipped)

	// 如果有新权限创建，同步给ADMIN角色
	if synced > 0 {
		if err := syncAdminPermissions(db); err != nil {
			return err
		}
	}

	return nil
}

// syncAdminPermissions 将所有权限分配给ADMIN角色
func syncAdminPermissions(db *gorm.DB) error {
	var adminRole model.Role
	if err := db.Where("code = ?", "ADMIN").First(&adminRole).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("[PermissionSeeder] ADMIN role not found, skipping permission assignment")
			return nil
		}
		return err
	}

	// 获取所有权限
	var allPermissions []model.Permission
	if err := db.Find(&allPermissions).Error; err != nil {
		return err
	}

	// 获取当前ADMIN角色已有的权限
	var currentPermissions []model.Permission
	db.Model(&adminRole).Association("Permissions").Find(&currentPermissions)

	// 找出需要添加的权限
	currentMap := make(map[uint]bool)
	for _, p := range currentPermissions {
		currentMap[p.ID] = true
	}

	var newPermissions []model.Permission
	for _, p := range allPermissions {
		if !currentMap[p.ID] {
			newPermissions = append(newPermissions, p)
		}
	}

	if len(newPermissions) > 0 {
		if err := db.Model(&adminRole).Association("Permissions").Append(newPermissions); err != nil {
			return err
		}
		log.Printf("[PermissionSeeder] Assigned %d new permissions to ADMIN role", len(newPermissions))
	}

	return nil
}
