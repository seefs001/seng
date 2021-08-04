package seng

type Route struct {
	// Data for routing
	pos  uint32 // Position in stack -> important for the sort of the matched routes
	use  bool   // USE matches path prefixes
	star bool   // Path equals '*'
	root bool   // Path equals '/'
	path string
	//routeParser routeParser

	Method   string        `json:"method"`
	Path     string        `json:"path"`
	Params   []string      `json:"params"`
	Handlers []HandlerFunc `json:"-"`
}

func NewRoute() *Route {
	return &Route{}
}
