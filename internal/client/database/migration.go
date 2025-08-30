package database

import (
	"fmt"
	"log"
	"time"

	"github.com/FatWang1/open_ticket_system/internal/models"
	"gorm.io/gorm"
)

// Migration 数据库迁移结构
type Migration struct {
	ID        uint   `gorm:"primaryKey"`
	Version   string `gorm:"type:varchar(100);not null;uniqueIndex"`
	AppliedAt int64  `gorm:"not null"`
}

// Migrator 数据库迁移器
type Migrator struct {
	db *gorm.DB
}

// NewMigrator 创建迁移器实例
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{db: db}
}

// RunMigrations 运行所有迁移
func (m *Migrator) RunMigrations() error {
	// 创建迁移表
	if err := m.db.AutoMigrate(&Migration{}); err != nil {
		return fmt.Errorf("failed to create migration table: %w", err)
	}

	// 定义迁移版本
	migrations := []struct {
		version string
		up      func(*gorm.DB) error
	}{
		{
			version: "001_initial_schema",
			up:      m.migration001InitialSchema,
		},
		{
			version: "002_add_template_fields",
			up:      m.migration002AddTemplateFields,
		},
		{
			version: "003_add_indexes",
			up:      m.migration003AddIndexes,
		},
	}

	// 执行迁移
	for _, migration := range migrations {
		if err := m.runMigration(migration.version, migration.up); err != nil {
			return fmt.Errorf("migration %s failed: %w", migration.version, err)
		}
	}

	log.Println("All migrations completed successfully")
	return nil
}

// runMigration 运行单个迁移
func (m *Migrator) runMigration(version string, up func(*gorm.DB) error) error {
	// 检查是否已经应用
	var count int64
	if err := m.db.Model(&Migration{}).Where("version = ?", version).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check migration status: %w", err)
	}

	if count > 0 {
		log.Printf("Migration %s already applied, skipping", version)
		return nil
	}

	// 开始事务
	tx := m.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 执行迁移
	if err := up(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute migration %s: %w", version, err)
	}

	// 记录迁移
	migration := &Migration{
		Version:   version,
		AppliedAt: time.Now().Unix(),
	}
	if err := tx.Create(migration).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to record migration %s: %w", version, err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit migration %s: %w", version, err)
	}

	log.Printf("Migration %s applied successfully", version)
	return nil
}

// migration001InitialSchema 初始数据库结构
func (m *Migrator) migration001InitialSchema(tx *gorm.DB) error {
	// 创建基础表结构，按照依赖关系顺序创建
	// 1. 先创建没有外键依赖的表
	if err := tx.AutoMigrate(&models.User{}); err != nil {
		return fmt.Errorf("failed to migrate User: %w", err)
	}

	if err := tx.AutoMigrate(&models.TicketTemplate{}); err != nil {
		return fmt.Errorf("failed to migrate TicketTemplate: %w", err)
	}

	if err := tx.AutoMigrate(&models.Ticket{}); err != nil {
		return fmt.Errorf("failed to migrate Ticket: %w", err)
	}

	// 2. 创建有外键依赖的表
	if err := tx.AutoMigrate(&models.TicketOperator{}); err != nil {
		return fmt.Errorf("failed to migrate TicketOperator: %w", err)
	}

	if err := tx.AutoMigrate(&models.TicketOperatedUser{}); err != nil {
		return fmt.Errorf("failed to migrate TicketOperatedUser: %w", err)
	}

	if err := tx.AutoMigrate(&models.TemplateEndStep{}); err != nil {
		return fmt.Errorf("failed to migrate TemplateEndStep: %w", err)
	}

	if err := tx.AutoMigrate(&models.StepConfigDB{}); err != nil {
		return fmt.Errorf("failed to migrate StepConfigDB: %w", err)
	}

	if err := tx.AutoMigrate(&models.StepOperator{}); err != nil {
		return fmt.Errorf("failed to migrate StepOperator: %w", err)
	}

	if err := tx.AutoMigrate(&models.NextStep{}); err != nil {
		return fmt.Errorf("failed to migrate NextStep: %w", err)
	}

	return nil
}

// migration002AddTemplateFields 添加模板字段
func (m *Migrator) migration002AddTemplateFields(tx *gorm.DB) error {
	// 检查字段是否已存在
	if tx.Migrator().HasColumn(&models.TicketTemplate{}, "name") {
		return nil
	}

	// 添加缺失的字段
	if err := tx.Exec("ALTER TABLE ticket_templates ADD COLUMN name VARCHAR(100) NOT NULL DEFAULT ''").Error; err != nil {
		return err
	}
	if err := tx.Exec("ALTER TABLE ticket_templates ADD COLUMN memo TEXT").Error; err != nil {
		return err
	}
	if err := tx.Exec("ALTER TABLE ticket_templates ADD COLUMN version VARCHAR(50) NOT NULL DEFAULT '1.0'").Error; err != nil {
		return err
	}
	if err := tx.Exec("ALTER TABLE ticket_templates ADD COLUMN creator VARCHAR(50) NOT NULL DEFAULT ''").Error; err != nil {
		return err
	}

	return nil
}

// migration003AddIndexes 添加索引
func (m *Migrator) migration003AddIndexes(tx *gorm.DB) error {
	// 添加复合索引
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_tickets_status_created_at ON tickets(status, created_at)",
		"CREATE INDEX IF NOT EXISTS idx_tickets_template_status ON tickets(template_id, status)",
		"CREATE INDEX IF NOT EXISTS idx_ticket_operators_ticket_operator ON ticket_operators(ticket_id, operator)",
		"CREATE INDEX IF NOT EXISTS idx_step_configs_template_step ON step_configs(ticket_template_id, step)",
	}

	for _, index := range indexes {
		if err := tx.Exec(index).Error; err != nil {
			return err
		}
	}

	return nil
}

// GetMigrationStatus 获取迁移状态
func (m *Migrator) GetMigrationStatus() ([]Migration, error) {
	var migrations []Migration
	if err := m.db.Order("version").Find(&migrations).Error; err != nil {
		return nil, fmt.Errorf("failed to get migration status: %w", err)
	}
	return migrations, nil
}

// RollbackMigration 回滚迁移
func (m *Migrator) RollbackMigration(version string) error {
	// 检查迁移是否存在
	var migration Migration
	if err := m.db.Where("version = ?", version).First(&migration).Error; err != nil {
		return fmt.Errorf("migration %s not found: %w", version, err)
	}

	// 开始事务
	tx := m.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 执行回滚逻辑（这里需要根据具体迁移实现）
	if err := m.executeRollbackLogic(tx, version); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to execute rollback logic: %w", err)
	}

	// 删除迁移记录
	if err := tx.Delete(&migration).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete migration record: %w", err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit rollback: %w", err)
	}

	log.Printf("Migration %s rolled back successfully", version)
	return nil
}

// executeRollbackLogic 执行具体的回滚逻辑
func (m *Migrator) executeRollbackLogic(tx *gorm.DB, version string) error {
	switch version {
	case "001":
		// 回滚001：删除所有表
		tables := []string{
			"next_steps",
			"step_operators",
			"step_configs",
			"template_end_steps",
			"ticket_templates",
			"ticket_operated_users",
			"ticket_operators",
			"tickets",
			"users",
		}

		for _, table := range tables {
			if err := tx.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table)).Error; err != nil {
				return fmt.Errorf("failed to drop table %s: %w", table, err)
			}
		}

	case "002":
		// 回滚002：删除添加的字段
		if err := tx.Exec("ALTER TABLE ticket_templates DROP COLUMN IF EXISTS name").Error; err != nil {
			return fmt.Errorf("failed to drop column name: %w", err)
		}
		if err := tx.Exec("ALTER TABLE ticket_templates DROP COLUMN IF EXISTS memo").Error; err != nil {
			return fmt.Errorf("failed to drop column memo: %w", err)
		}
		if err := tx.Exec("ALTER TABLE ticket_templates DROP COLUMN IF EXISTS version").Error; err != nil {
			return fmt.Errorf("failed to drop column version: %w", err)
		}
		if err := tx.Exec("ALTER TABLE ticket_templates DROP COLUMN IF EXISTS creator").Error; err != nil {
			return fmt.Errorf("failed to drop column creator: %w", err)
		}

	case "003":
		// 回滚003：删除添加的索引
		indexes := []string{
			"idx_tickets_status_created_at",
			"idx_tickets_template_status",
			"idx_ticket_operators_ticket_operator",
			"idx_step_configs_template_step",
		}

		for _, index := range indexes {
			if err := tx.Exec(fmt.Sprintf("DROP INDEX IF EXISTS %s", index)).Error; err != nil {
				return fmt.Errorf("failed to drop index %s: %w", index, err)
			}
		}

	default:
		return fmt.Errorf("unsupported migration version for rollback: %s", version)
	}

	return nil
}
