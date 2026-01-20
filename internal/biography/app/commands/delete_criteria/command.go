package delete_criteria

type DeleteCriteriaCommand struct {
	UserID      string   `json:"user_id"`
	CriteriaIDs []string `json:"criteria_ids"`
}
