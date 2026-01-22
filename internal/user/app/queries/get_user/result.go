package getuser

type GetUserByExternalIDResult struct {
	ID         string `json:"id"`
	CreatedAt  string `json:"created_at"`
	ExternalID string `json:"external_id"`
}
