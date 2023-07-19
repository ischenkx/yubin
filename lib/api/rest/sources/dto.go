package sources

import (
	yubin "yubin/src"
	"yubin/src/mail"
)

type SourceDto struct {
	Name     string `json:"name"`
	Address  string `json:"address,omitempty"`
	Password string `json:"password,omitempty"`
	Host     string `json:"host,omitempty"`
	Port     int    `json:"port,omitempty"`
}

func source2dto(source yubin.NamedSource) SourceDto {
	return SourceDto{
		Name:     source.Name,
		Address:  source.Address,
		Password: source.Password,
		Host:     source.Host,
		Port:     source.Port,
	}
}

func dto2source(source SourceDto) yubin.NamedSource {
	return yubin.NamedSource{
		Name: source.Name,
		Source: mail.Source{
			Address:  source.Address,
			Password: source.Password,
			Host:     source.Host,
			Port:     source.Port,
		},
	}
}

type UpdateSourceDto struct {
	Name     string  `json:"name"`
	Address  *string `json:"address,omitempty"`
	Password *string `json:"password,omitempty"`
	Host     *string `json:"host,omitempty"`
	Port     *int    `json:"port,omitempty"`
}
