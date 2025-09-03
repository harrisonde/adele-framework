package render

import (
	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"github.com/petaki/inertia-go"
)

type Render struct {
	Renderer       string
	RootPath       string
	Directory      string
	Secure         bool
	Port           string
	ServerName     string
	JetViews       *jet.Set
	Session        *scs.SessionManager
	InertiaManager *inertia.Inertia
}

type TemplateData struct {
	IsAuthenticated bool
	IntMap          map[string]int
	StringMap       map[string]string
	FloatMap        map[string]float32
	Data            map[string]interface{} // Use interface data can be anything.
	CSRFToken       string
	Port            string
	ServerName      string
	Secure          bool
	Error           string
	Flash           string
}
