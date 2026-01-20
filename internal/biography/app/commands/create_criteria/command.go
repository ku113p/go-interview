package create_criteria

type CreateCriterion struct {
	Description string `json:"description"`
}

type CreateCriteriaCommand struct {
	LifeAreaID string            `json:"life_area_id"`
	UserID     string            `json:"user_id"`
	Criteria   []CreateCriterion `json:"criteria"`
}
