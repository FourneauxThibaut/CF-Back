package ruleeditor

// DefaultBlockDefinitions returns the default block templates for a new system.
func DefaultBlockDefinitions() []BlockDefinition {
	return []BlockDefinition{
		{
			ID:   string(BlockTypeTrigger),
			Type: BlockTypeTrigger,
			Label: "QUAND",
			Color: "#e8a838",
			Icon:  "⚡",
			TemplateSegments: []Segment{
				{Type: "text", Value: "Le joueur"},
				{Type: "dropdown", ID: "trigger_type", Value: "fait un test", Options: []string{"fait un test", "attaque", "utilise un objet", "lance un sort", "commence son tour"}},
				{Type: "text", Value: "de"},
				{Type: "dropdown", ID: "trigger_target", Value: "Compétence", Options: []string{"Compétence", "Attribut", "Sauvegarde", "Résistance"}},
			},
		},
		{
			ID:   string(BlockTypeRoll),
			Type: BlockTypeRoll,
			Label: "LANCER",
			Color: "#4a9eff",
			Icon:  "🎲",
			TemplateSegments: []Segment{
				{Type: "text", Value: "Lancer"},
				{Type: "dropdown", ID: "dice_amount", Value: "selon niveau", Options: []string{"selon niveau", "1D", "2D", "3D", "4D", "5D"}},
				{Type: "text", Value: "de type"},
				{Type: "dropdown", ID: "dice_type", Value: "selon Caractéristique", Options: []string{"selon Caractéristique", "D4", "D6", "D8", "D10", "D12", "D20"}},
			},
		},
		{
			ID:   string(BlockTypeCondition),
			Type: BlockTypeCondition,
			Label: "SI",
			Color: "#aa6fff",
			Icon:  "◆",
			TemplateSegments: []Segment{
				{Type: "text", Value: "Si"},
				{Type: "dropdown", ID: "cond_subject", Value: "résultat", Options: []string{"résultat", "nombre de 1", "dés identiques", "valeur max", "valeur min"}},
				{Type: "dropdown", ID: "cond_operator", Value: "≥", Options: []string{"≥", ">", "=", "≤", "<", "≠"}},
				{Type: "input", ID: "cond_value", Value: "", Placeholder: "valeur"},
			},
		},
		{
			ID:   string(BlockTypeAction),
			Type: BlockTypeAction,
			Label: "ALORS",
			Color: "#44ddaa",
			Icon:  "→",
			TemplateSegments: []Segment{
				{Type: "dropdown", ID: "action_type", Value: "Réussite", Options: []string{"Réussite", "Échec", "Échec Critique", "Réussite Critique", "Appliquer dégâts", "Ajouter état", "Consommer ressource"}},
			},
		},
		{
			ID:   string(BlockTypeModifier),
			Type: BlockTypeModifier,
			Label: "AVEC",
			Color: "#ff6b8a",
			Icon:  "✦",
			TemplateSegments: []Segment{
				{Type: "text", Value: "Appliquer"},
				{Type: "dropdown", ID: "mod_type", Value: "bonus", Options: []string{"bonus", "malus", "avantage", "désavantage"}},
				{Type: "text", Value: "de"},
				{Type: "input", ID: "mod_value", Value: "", Placeholder: "+/-"},
				{Type: "text", Value: "depuis"},
				{Type: "dropdown", ID: "mod_source", Value: "objet équipé", Options: []string{"objet équipé", "compétence passive", "état actif", "environnement"}},
			},
		},
		{
			ID:   string(BlockTypeOption),
			Type: BlockTypeOption,
			Label: "OPTION",
			Color: "#8bc4ff",
			Icon:  "?",
			TemplateSegments: []Segment{
				{Type: "text", Value: "Le joueur peut"},
				{Type: "dropdown", ID: "option_action", Value: "Relancer tous les dés", Options: []string{"Relancer tous les dés", "Relancer un dé", "Garder le résultat", "Annuler l'action", "Dépenser un point"}},
			},
		},
	}
}
