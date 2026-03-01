package ruleeditor

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// JSONBRepository implements Repository with PostgreSQL JSONB storage.
type JSONBRepository struct {
	pool *pgxpool.Pool
}

// NewJSONBRepository returns a repository that stores rule systems in JSONB.
func NewJSONBRepository(pool *pgxpool.Pool) *JSONBRepository {
	return &JSONBRepository{pool: pool}
}

func (r *JSONBRepository) GetSystemsByUser(ctx context.Context, userID string) ([]RuleSystem, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, name, description, rules, block_definitions, created_at, updated_at
		 FROM rule_systems WHERE user_id = $1 ORDER BY updated_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []RuleSystem
	for rows.Next() {
		var sys RuleSystem
		var rulesJSON, defsJSON []byte
		if err := rows.Scan(&sys.ID, &sys.UserID, &sys.Name, &sys.Description, &rulesJSON, &defsJSON, &sys.CreatedAt, &sys.UpdatedAt); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(rulesJSON, &sys.Rules); err != nil {
			return nil, err
		}
		if err := json.Unmarshal(defsJSON, &sys.BlockDefinitions); err != nil {
			return nil, err
		}
		list = append(list, sys)
	}
	return list, rows.Err()
}

func (r *JSONBRepository) GetSystemByID(ctx context.Context, systemID string) (*RuleSystem, error) {
	var sys RuleSystem
	var rulesJSON, defsJSON []byte
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, name, description, rules, block_definitions, created_at, updated_at
		 FROM rule_systems WHERE id = $1`,
		systemID,
	).Scan(&sys.ID, &sys.UserID, &sys.Name, &sys.Description, &rulesJSON, &defsJSON, &sys.CreatedAt, &sys.UpdatedAt)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(rulesJSON, &sys.Rules); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(defsJSON, &sys.BlockDefinitions); err != nil {
		return nil, err
	}
	return &sys, nil
}

func (r *JSONBRepository) CreateSystem(ctx context.Context, userID string, name, description string, blockDefs []BlockDefinition) (*RuleSystem, error) {
	id := uuid.New().String()
	rulesJSON, _ := json.Marshal([]Rule{})
	defsJSON, err := json.Marshal(blockDefs)
	if err != nil {
		return nil, err
	}
	_, err = r.pool.Exec(ctx,
		`INSERT INTO rule_systems (id, user_id, name, description, rules, block_definitions)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		id, userID, name, description, rulesJSON, defsJSON,
	)
	if err != nil {
		return nil, err
	}
	return r.GetSystemByID(ctx, id)
}

func (r *JSONBRepository) UpdateSystem(ctx context.Context, systemID, userID string, name, description *string) (*RuleSystem, error) {
	sys, err := r.getSystemForUpdate(ctx, systemID, userID)
	if err != nil {
		return nil, err
	}
	if name != nil {
		sys.Name = *name
	}
	if description != nil {
		sys.Description = *description
	}
	return r.saveSystem(ctx, sys)
}

func (r *JSONBRepository) DeleteSystem(ctx context.Context, systemID, userID string) error {
	res, err := r.pool.Exec(ctx, `DELETE FROM rule_systems WHERE id = $1 AND user_id = $2`, systemID, userID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("rule system not found")
	}
	return nil
}

func (r *JSONBRepository) AddRule(ctx context.Context, systemID, userID string, name, description, icon string, order int) (*Rule, error) {
	sys, err := r.getSystemForUpdate(ctx, systemID, userID)
	if err != nil {
		return nil, err
	}
	rule := Rule{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Icon:        icon,
		Blocks:      []RuleBlock{},
		IsActive:    true,
		Order:       order,
	}
	sys.Rules = append(sys.Rules, rule)
	_, err = r.saveSystem(ctx, sys)
	if err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *JSONBRepository) UpdateRule(ctx context.Context, systemID, ruleID, userID string, name, description, icon *string, isActive *bool, order *int) (*Rule, error) {
	sys, err := r.getSystemForUpdate(ctx, systemID, userID)
	if err != nil {
		return nil, err
	}
	for i := range sys.Rules {
		if sys.Rules[i].ID == ruleID {
			if name != nil {
				sys.Rules[i].Name = *name
			}
			if description != nil {
				sys.Rules[i].Description = *description
			}
			if icon != nil {
				sys.Rules[i].Icon = *icon
			}
			if isActive != nil {
				sys.Rules[i].IsActive = *isActive
			}
			if order != nil {
				sys.Rules[i].Order = *order
			}
			_, err = r.saveSystem(ctx, sys)
			if err != nil {
				return nil, err
			}
			return &sys.Rules[i], nil
		}
	}
	return nil, fmt.Errorf("rule not found")
}

func (r *JSONBRepository) DeleteRule(ctx context.Context, systemID, ruleID, userID string) error {
	sys, err := r.getSystemForUpdate(ctx, systemID, userID)
	if err != nil {
		return err
	}
	for i, rule := range sys.Rules {
		if rule.ID == ruleID {
			sys.Rules = append(sys.Rules[:i], sys.Rules[i+1:]...)
			_, err = r.saveSystem(ctx, sys)
			return err
		}
	}
	return fmt.Errorf("rule not found")
}

func (r *JSONBRepository) ReorderRules(ctx context.Context, systemID, userID string, orderedIDs []string) error {
	sys, err := r.getSystemForUpdate(ctx, systemID, userID)
	if err != nil {
		return err
	}
	byID := make(map[string]Rule)
	for _, rule := range sys.Rules {
		byID[rule.ID] = rule
	}
	sys.Rules = make([]Rule, 0, len(orderedIDs))
	for i, id := range orderedIDs {
		if rule, ok := byID[id]; ok {
			rule.Order = i
			sys.Rules = append(sys.Rules, rule)
		}
	}
	_, err = r.saveSystem(ctx, sys)
	return err
}

func (r *JSONBRepository) AddBlock(ctx context.Context, systemID, ruleID, userID string, block RuleBlock) (*RuleBlock, error) {
	sys, err := r.getSystemForUpdate(ctx, systemID, userID)
	if err != nil {
		return nil, err
	}
	for i := range sys.Rules {
		if sys.Rules[i].ID == ruleID {
			if block.ID == "" {
				block.ID = uuid.New().String()
			}
			block.Order = len(sys.Rules[i].Blocks)
			sys.Rules[i].Blocks = append(sys.Rules[i].Blocks, block)
			_, err = r.saveSystem(ctx, sys)
			if err != nil {
				return nil, err
			}
			return &sys.Rules[i].Blocks[len(sys.Rules[i].Blocks)-1], nil
		}
	}
	return nil, fmt.Errorf("rule not found")
}

func (r *JSONBRepository) UpdateBlock(ctx context.Context, systemID, ruleID, blockID, userID string, segments []Segment, order *int) (*RuleBlock, error) {
	sys, err := r.getSystemForUpdate(ctx, systemID, userID)
	if err != nil {
		return nil, err
	}
	for i := range sys.Rules {
		if sys.Rules[i].ID == ruleID {
			for j := range sys.Rules[i].Blocks {
				if sys.Rules[i].Blocks[j].ID == blockID {
					if segments != nil {
						sys.Rules[i].Blocks[j].Segments = segments
					}
					if order != nil {
						sys.Rules[i].Blocks[j].Order = *order
					}
					_, err = r.saveSystem(ctx, sys)
					if err != nil {
						return nil, err
					}
					return &sys.Rules[i].Blocks[j], nil
				}
			}
			return nil, fmt.Errorf("block not found")
		}
	}
	return nil, fmt.Errorf("rule not found")
}

func (r *JSONBRepository) DeleteBlock(ctx context.Context, systemID, ruleID, blockID, userID string) error {
	sys, err := r.getSystemForUpdate(ctx, systemID, userID)
	if err != nil {
		return err
	}
	for i := range sys.Rules {
		if sys.Rules[i].ID == ruleID {
			for j, b := range sys.Rules[i].Blocks {
				if b.ID == blockID {
					sys.Rules[i].Blocks = append(sys.Rules[i].Blocks[:j], sys.Rules[i].Blocks[j+1:]...)
					_, err = r.saveSystem(ctx, sys)
					return err
				}
			}
			return fmt.Errorf("block not found")
		}
	}
	return fmt.Errorf("rule not found")
}

func (r *JSONBRepository) ReorderBlocks(ctx context.Context, systemID, ruleID, userID string, orderedIDs []string) error {
	sys, err := r.getSystemForUpdate(ctx, systemID, userID)
	if err != nil {
		return err
	}
	for i := range sys.Rules {
		if sys.Rules[i].ID == ruleID {
			byID := make(map[string]RuleBlock)
			for _, b := range sys.Rules[i].Blocks {
				byID[b.ID] = b
			}
			sys.Rules[i].Blocks = make([]RuleBlock, 0, len(orderedIDs))
			for j, id := range orderedIDs {
				if b, ok := byID[id]; ok {
					b.Order = j
					sys.Rules[i].Blocks = append(sys.Rules[i].Blocks, b)
				}
			}
			_, err = r.saveSystem(ctx, sys)
			return err
		}
	}
	return fmt.Errorf("rule not found")
}

func (r *JSONBRepository) GetBlockDefinitions(ctx context.Context, systemID, userID string) ([]BlockDefinition, error) {
	sys, err := r.GetSystemByID(ctx, systemID)
	if err != nil {
		return nil, err
	}
	if sys.UserID != userID {
		return nil, fmt.Errorf("forbidden")
	}
	return sys.BlockDefinitions, nil
}

func (r *JSONBRepository) AddBlockDefinition(ctx context.Context, systemID, userID string, def BlockDefinition) (*BlockDefinition, error) {
	sys, err := r.getSystemForUpdate(ctx, systemID, userID)
	if err != nil {
		return nil, err
	}
	if def.ID == "" {
		def.ID = uuid.New().String()
	}
	sys.BlockDefinitions = append(sys.BlockDefinitions, def)
	_, err = r.saveSystem(ctx, sys)
	if err != nil {
		return nil, err
	}
	return &sys.BlockDefinitions[len(sys.BlockDefinitions)-1], nil
}

func (r *JSONBRepository) UpdateBlockDefinition(ctx context.Context, systemID, defID, userID string, def BlockDefinition) (*BlockDefinition, error) {
	sys, err := r.getSystemForUpdate(ctx, systemID, userID)
	if err != nil {
		return nil, err
	}
	for i := range sys.BlockDefinitions {
		if sys.BlockDefinitions[i].ID == defID || string(sys.BlockDefinitions[i].Type) == defID || sys.BlockDefinitions[i].Label == defID {
			sys.BlockDefinitions[i] = def
			if sys.BlockDefinitions[i].ID == "" {
				sys.BlockDefinitions[i].ID = defID
			}
			_, err = r.saveSystem(ctx, sys)
			if err != nil {
				return nil, err
			}
			return &sys.BlockDefinitions[i], nil
		}
	}
	return nil, fmt.Errorf("block definition not found")
}

func (r *JSONBRepository) DeleteBlockDefinition(ctx context.Context, systemID, defID, userID string) error {
	sys, err := r.getSystemForUpdate(ctx, systemID, userID)
	if err != nil {
		return err
	}
	for i, d := range sys.BlockDefinitions {
		if d.ID == defID || string(d.Type) == defID || d.Label == defID {
			sys.BlockDefinitions = append(sys.BlockDefinitions[:i], sys.BlockDefinitions[i+1:]...)
			_, err = r.saveSystem(ctx, sys)
			return err
		}
	}
	return fmt.Errorf("block definition not found")
}

func (r *JSONBRepository) getSystemForUpdate(ctx context.Context, systemID, userID string) (*RuleSystem, error) {
	sys, err := r.GetSystemByID(ctx, systemID)
	if err != nil {
		return nil, err
	}
	if sys.UserID != userID {
		return nil, fmt.Errorf("forbidden")
	}
	return sys, nil
}

func (r *JSONBRepository) saveSystem(ctx context.Context, sys *RuleSystem) (*RuleSystem, error) {
	rulesJSON, err := json.Marshal(sys.Rules)
	if err != nil {
		return nil, err
	}
	defsJSON, err := json.Marshal(sys.BlockDefinitions)
	if err != nil {
		return nil, err
	}
	_, err = r.pool.Exec(ctx,
		`UPDATE rule_systems SET name = $1, description = $2, rules = $3, block_definitions = $4, updated_at = NOW() WHERE id = $5 AND user_id = $6`,
		sys.Name, sys.Description, rulesJSON, defsJSON, sys.ID, sys.UserID,
	)
	if err != nil {
		return nil, err
	}
	return r.GetSystemByID(ctx, sys.ID)
}
