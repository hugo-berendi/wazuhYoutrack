package api

type Issue struct {
	ID           string                   `json:"id"`
	Type         string                   `json:"$type"`
	Summary      string                   `json:"summary"`
	Description  string                   `json:"description"`
	CustomFields []map[string]interface{} `json:"customFields"`
}
