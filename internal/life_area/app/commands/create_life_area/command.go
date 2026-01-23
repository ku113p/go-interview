package create_life_area

type CreateLifeAreaCommand struct {
	UserID   string  `json:"user_id"`
	ParentID *string `json:"parent_id"`
	Title    string  `json:"title"`
	Goal     string  `json:"goal"`
}
