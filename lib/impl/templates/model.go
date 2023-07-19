package templates

type Model struct {
	Name         string
	Data         string
	Meta         map[string]any
	SubTemplates map[string]Model
}
