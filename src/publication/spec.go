package publication

type Spec struct {
	Destination Destination
	Sender      string
	Template    string
	Source      string
	Properties  map[string]any
}
