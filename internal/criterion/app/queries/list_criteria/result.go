package listcriteria

type ListCriteriaResult struct {
	Items []Criterion `json:"items"`
}

type Criterion struct {
	ID          string `json:"id"`
	NodeID      string `json:"node_id"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
}
