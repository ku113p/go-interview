package create_criteria

type CreateCriterion struct {
	Description string `json:"description"`
}

type CreateCriteriaCommand struct {
	NodeID       string   `json:"node_id"`
	Descriptions []string `json:"descriptions"`
}
