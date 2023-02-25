package swagger

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var swagger = Swagger{
	OpenAPI: "3.0.3",
	Paths:   map[string]map[string]Operation{},
	Components: Component{
		Schemas:         map[string]ComponentSchema{},
		RequestBodies:   map[string]RequestBody{},
		SecuritySchemes: map[string]SecurityScheme{},
	},
}

func Init(a Application) error {
	swagger.Info = Info{
		Title: a.Name,
		Contact: Contact{
			Email: "test@email.com",
		},
		Version: "3.0",
	}
	swagger.Servers = []Server{{"http://localhost:8080"}}

	for _, c := range a.Controllers {
		swagger.AddTag(c.Name, c.Description)
		for _, r := range c.Routes {
			swagger.SetPaths(r, c.Name)
		}
	}

	doc, err := yaml.Marshal(swagger)
	if err != nil {
		return errors.Wrap(err, "swagger marshal")
	}

	if err := os.WriteFile("./framework/swagger/dist/swagger.yaml", doc, 0o644); err != nil {
		return errors.Wrap(err, "swagger yaml write")
	}
	return nil
}
