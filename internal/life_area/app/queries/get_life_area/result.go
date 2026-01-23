package get_life_area

type GetLifeAreaResult struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	Goal      string  `json:"goal"`
	ParentID  *string `json:"parent_id"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
