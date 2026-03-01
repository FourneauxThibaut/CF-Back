package ruleeditor

import "time"

// BlockType is the kind of rule block.
type BlockType string

const (
	BlockTypeTrigger   BlockType = "trigger"
	BlockTypeRoll      BlockType = "roll"
	BlockTypeCondition BlockType = "condition"
	BlockTypeAction    BlockType = "action"
	BlockTypeModifier  BlockType = "modifier"
	BlockTypeOption    BlockType = "option"
)

// RefType is the type of a reference segment.
type RefType string

const (
	RefTypeSkill   RefType = "skill"
	RefTypeItem    RefType = "item"
	RefTypeEntity  RefType = "entity"
	RefTypeStat    RefType = "stat"
)

// RuleSystem is the root document (one per game system).
type RuleSystem struct {
	ID               string            `json:"id"`
	UserID           string            `json:"userId"`
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	Rules            []Rule            `json:"rules"`
	BlockDefinitions []BlockDefinition `json:"blockDefinitions"`
	CreatedAt        time.Time         `json:"createdAt"`
	UpdatedAt        time.Time         `json:"updatedAt"`
}

// Rule is a single named rule (e.g. "Test de Compétence") with ordered blocks.
type Rule struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Icon        string      `json:"icon"`
	Blocks      []RuleBlock `json:"blocks"`
	IsActive    bool        `json:"isActive"`
	Order       int         `json:"order"`
}

// RuleBlock is one block in a rule (trigger, roll, condition, etc.).
type RuleBlock struct {
	ID       string    `json:"id"`
	Type     BlockType `json:"type"`
	Segments []Segment `json:"segments"`
	Order    int       `json:"order"`
}

// Segment is the content inside a block (text, dropdown, input, reference).
type Segment struct {
	Type  string `json:"type"` // "text" | "dropdown" | "input" | "reference"
	Value string `json:"value,omitempty"`

	// dropdown / input / reference
	ID          string   `json:"id,omitempty"`
	Options     []string `json:"options,omitempty"`
	Placeholder string   `json:"placeholder,omitempty"`
	RefType     RefType  `json:"refType,omitempty"`
	RefID       *string  `json:"refId,omitempty"`
}

// BlockDefinition is a reusable template for creating blocks.
type BlockDefinition struct {
	ID               string    `json:"id,omitempty"` // optional; default templates use type as id
	Type             BlockType `json:"type"`
	Label            string    `json:"label"`
	Color            string    `json:"color"`
	Icon             string    `json:"icon"`
	TemplateSegments []Segment `json:"templateSegments"`
}
