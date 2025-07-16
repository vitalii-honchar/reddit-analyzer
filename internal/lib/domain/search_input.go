package domain

type SearchInput struct {
	Query string `json:"query" jsonschema_description:"Search query"`
}
