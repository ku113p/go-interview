package change_life_area_goal

type ChangeLifeAreaGoalCommand struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Goal   string `json:"goal"`
}
