package ruleeditor

import (
	"context"
	"fmt"
)

// Service holds business logic for the rule editor (repository-agnostic).
type Service struct {
	repo Repository
}

// NewService returns a new rule editor service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// ListSystems returns all rule systems for the user.
func (s *Service) ListSystems(ctx context.Context, userID string) ([]RuleSystem, error) {
	return s.repo.GetSystemsByUser(ctx, userID)
}

// GetSystem returns a single rule system by ID (must belong to user).
func (s *Service) GetSystem(ctx context.Context, systemID, userID string) (*RuleSystem, error) {
	sys, err := s.repo.GetSystemByID(ctx, systemID)
	if err != nil {
		return nil, err
	}
	if sys.UserID != userID {
		return nil, fmt.Errorf("forbidden")
	}
	return sys, nil
}

// CreateSystem creates a new rule system with default block definitions.
func (s *Service) CreateSystem(ctx context.Context, userID string, name, description string) (*RuleSystem, error) {
	if name == "" {
		name = "Nouveau système"
	}
	return s.repo.CreateSystem(ctx, userID, name, description, DefaultBlockDefinitions())
}

// UpdateSystem updates system name/description.
func (s *Service) UpdateSystem(ctx context.Context, systemID, userID string, name, description *string) (*RuleSystem, error) {
	return s.repo.UpdateSystem(ctx, systemID, userID, name, description)
}

// DeleteSystem deletes a rule system.
func (s *Service) DeleteSystem(ctx context.Context, systemID, userID string) error {
	return s.repo.DeleteSystem(ctx, systemID, userID)
}

// AddRule adds a rule to a system.
func (s *Service) AddRule(ctx context.Context, systemID, userID string, name, description, icon string, order int) (*Rule, error) {
	return s.repo.AddRule(ctx, systemID, userID, name, description, icon, order)
}

// UpdateRule updates a rule.
func (s *Service) UpdateRule(ctx context.Context, systemID, ruleID, userID string, name, description, icon *string, isActive *bool, order *int) (*Rule, error) {
	return s.repo.UpdateRule(ctx, systemID, ruleID, userID, name, description, icon, isActive, order)
}

// DeleteRule deletes a rule.
func (s *Service) DeleteRule(ctx context.Context, systemID, ruleID, userID string) error {
	return s.repo.DeleteRule(ctx, systemID, ruleID, userID)
}

// ReorderRules reorders rules by ordered IDs.
func (s *Service) ReorderRules(ctx context.Context, systemID, userID string, orderedIDs []string) error {
	return s.repo.ReorderRules(ctx, systemID, userID, orderedIDs)
}

// AddBlock adds a block to a rule.
func (s *Service) AddBlock(ctx context.Context, systemID, ruleID, userID string, block RuleBlock) (*RuleBlock, error) {
	return s.repo.AddBlock(ctx, systemID, ruleID, userID, block)
}

// UpdateBlock updates a block (segments and/or order).
func (s *Service) UpdateBlock(ctx context.Context, systemID, ruleID, blockID, userID string, segments []Segment, order *int) (*RuleBlock, error) {
	return s.repo.UpdateBlock(ctx, systemID, ruleID, blockID, userID, segments, order)
}

// DeleteBlock deletes a block.
func (s *Service) DeleteBlock(ctx context.Context, systemID, ruleID, blockID, userID string) error {
	return s.repo.DeleteBlock(ctx, systemID, ruleID, blockID, userID)
}

// ReorderBlocks reorders blocks in a rule.
func (s *Service) ReorderBlocks(ctx context.Context, systemID, ruleID, userID string, orderedIDs []string) error {
	return s.repo.ReorderBlocks(ctx, systemID, ruleID, userID, orderedIDs)
}

// GetBlockDefinitions returns block definitions for a system.
func (s *Service) GetBlockDefinitions(ctx context.Context, systemID, userID string) ([]BlockDefinition, error) {
	return s.repo.GetBlockDefinitions(ctx, systemID, userID)
}

// AddBlockDefinition adds a custom block definition.
func (s *Service) AddBlockDefinition(ctx context.Context, systemID, userID string, def BlockDefinition) (*BlockDefinition, error) {
	return s.repo.AddBlockDefinition(ctx, systemID, userID, def)
}

// UpdateBlockDefinition updates a block definition.
func (s *Service) UpdateBlockDefinition(ctx context.Context, systemID, defID, userID string, def BlockDefinition) (*BlockDefinition, error) {
	return s.repo.UpdateBlockDefinition(ctx, systemID, defID, userID, def)
}

// DeleteBlockDefinition deletes a block definition.
func (s *Service) DeleteBlockDefinition(ctx context.Context, systemID, defID, userID string) error {
	return s.repo.DeleteBlockDefinition(ctx, systemID, defID, userID)
}
