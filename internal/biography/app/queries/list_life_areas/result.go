package list_live_area

type ListLifeAreaResult struct {
	Items []*LifeArea `json:"items"`
}

type LifeArea struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	Goal      string  `json:"goal"`
	ParentID  *string `json:"parent_id"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
