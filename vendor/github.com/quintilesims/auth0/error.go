package auth0

import "fmt"

type Error struct {
	Type        string `json:"error"`
	Description string `json:"error_description"`
	URI         string `json:"error_uri"`
}

func (e Error) Error() string {
	if e.URI == "" {
		return fmt.Sprintf("%s: %s", e.Type, e.Description)
	}

	return fmt.Sprintf("%s: %s (%s)", e.Type, e.Description, e.URI)
}
