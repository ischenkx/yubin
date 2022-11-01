package user

import "smtp-client/pkg/data/crud"

type Repo interface {
	crud.CRUD[string, User]
}
