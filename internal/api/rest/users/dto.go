package users

import (
	"smtp-client/internal/api/rest/util"
	"smtp-client/internal/mailer/user"
)

type UpdateDto struct {
	Email   *string        `json:"email,omitempty"`
	Name    *string        `json:"name,omitempty"`
	Surname *string        `json:"surname,omitempty"`
	Meta    map[string]any `json:"meta,omitempty"`
}

type CreateDto struct {
	Email   string         `json:"email,omitempty"`
	Name    string         `json:"name,omitempty"`
	Surname string         `json:"surname,omitempty"`
	Meta    map[string]any `json:"meta,omitempty"`
}

type UserDto struct {
	ID      string         `json:"id,omitempty"`
	Email   string         `json:"email,omitempty"`
	Name    string         `json:"name,omitempty"`
	Surname string         `json:"surname,omitempty"`
	Meta    map[string]any `json:"meta,omitempty"`
}

func user2dto(user user.User) UserDto {
	return UserDto{
		ID:      user.ID,
		Email:   user.Email,
		Name:    user.Name,
		Surname: user.Surname,
		Meta:    util.ValidateEmptyMap(user.Meta),
	}
}
