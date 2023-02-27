package swagger

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var S = Swagger{
	OpenAPI: "3.0.3",
	Paths:   map[string]Path{},
	Components: Component{
		Schemas:         map[string]ComponentSchema{},
		RequestBodies:   map[string]RequestBody{},
		SecuritySchemes: map[string]SecurityScheme{},
	},
}

func Init(name string) error {
	S.Info = Info{
		Title: name,
		Contact: Contact{
			Email: "test@email.com",
		},
		Version: "3.0",
	}
	S.Servers = []Server{{"http://localhost:8080"}}
	S.SetSecuritySchemes()

	doc, err := yaml.Marshal(S)
	if err != nil {
		return errors.Wrap(err, "swagger marshal")
	}

	if err := os.WriteFile("./framework/swagger/dist/swagger.yaml", doc, 0o644); err != nil {
		return errors.Wrap(err, "swagger yaml write")
	}
	return nil
}
