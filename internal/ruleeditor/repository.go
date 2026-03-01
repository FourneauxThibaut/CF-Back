package ruleeditor

import (
	"context"
)

// Repository defines data access for rule systems (Phase 2 can swap implementation).
type Repository interface {
	// Systems
	GetSystemsByUser(ctx context.Context, userID string) ([]RuleSystem, error)
	GetSystemByID(ctx context.Context, systemID string) (*RuleSystem, error)
	CreateSystem(ctx context.Context, userID string, name, description string, blockDefs []BlockDefinition) (*RuleSystem, error)
	UpdateSystem(ctx context.Context, systemID, userID string, name, description *string) (*RuleSystem, error)
	DeleteSystem(ctx context.Context, systemID, userID string) error

	// Rules
	AddRule(ctx context.Context, systemID, userID string, name, description, icon string, order int) (*Rule, error)
	UpdateRule(ctx context.Context, systemID, ruleID, userID string, name, description, icon *string, isActive *bool, order *int) (*Rule, error)
	DeleteRule(ctx context.Context, systemID, ruleID, userID string) error
	ReorderRules(ctx context.Context, systemID, userID string, orderedIDs []string) error

	// Blocks
	AddBlock(ctx context.Context, systemID, ruleID, userID string, block RuleBlock) (*RuleBlock, error)
	UpdateBlock(ctx context.Context, systemID, ruleID, blockID, userID string, segments []Segment, order *int) (*RuleBlock, error)
	DeleteBlock(ctx context.Context, systemID, ruleID, blockID, userID string) error
	ReorderBlocks(ctx context.Context, systemID, ruleID, userID string, orderedIDs []string) error

	// Block definitions
	GetBlockDefinitions(ctx context.Context, systemID, userID string) ([]BlockDefinition, error)
	AddBlockDefinition(ctx context.Context, systemID, userID string, def BlockDefinition) (*BlockDefinition, error)
	UpdateBlockDefinition(ctx context.Context, systemID, defID, userID string, def BlockDefinition) (*BlockDefinition, error)
	DeleteBlockDefinition(ctx context.Context, systemID, defID, userID string) error
}
