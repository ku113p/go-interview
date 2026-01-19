package get

type Result struct {
	ID        string       `json:"id"`
	Title     string       `json:"title"`
	Goal      string       `json:"goal"`
	ParentID  *string      `json:"parent_id"`
	Criteria  []*Criterion `json:"criteria"`
	CreatedAt string       `json:"created_at"`
	UpdatedAt string       `json:"updated_at"`
}

type Criterion struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
}
