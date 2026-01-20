package change_life_area_parent

type ChangeLifeAreaParentCommand struct {
	ID       string  `json:"id"`
	UserID   string  `json:"user_id"`
	ParentID *string `json:"parent_id"`
}
