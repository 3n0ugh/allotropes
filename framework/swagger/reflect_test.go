package swagger

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/3n0ugh/allotropes/internal/pagination"
	"github.com/stretchr/testify/assert"
)

type TestRoute struct {
	Name        string
	Description string
	Method      string
	Path        string
	Headers     map[string]string
	Middlewares []func(http.Handler) http.Handler
	Request     any
	Response    any
}

func (r TestRoute) GetRequestModel() any                              { return r.Request }
func (r TestRoute) GetResponseModel() any                             { return r.Response }
func (r TestRoute) GetMethod() string                                 { return r.Method }
func (r TestRoute) GetDescription() string                            { return r.Description }
func (r TestRoute) GetMiddlewares() []func(http.Handler) http.Handler { return r.Middlewares }
func (r TestRoute) GetHeaders() map[string]string                     { return r.Headers }
func (r TestRoute) GetPath() string                                   { return r.Path }

type Req struct {
	Elma   []string `query:"elma"`
	UserID int      `query:"userId"`
}
type Res struct {
	Domates
	Patates string `json:"patates"`
}

type Domates struct {
	Adet int `json:"adet"`
}

func auth(_ http.Handler) http.Handler { return nil }

func TestSwagger_AddTag(t *testing.T) {
	expected := Swagger{
		Tags: []Tag{
			{Name: "test", Description: "abcdef"},
		},
	}

	actual := Swagger{}
	actual.AddTag("test", "abcdef")

	assert.Equal(t, expected, actual)
}

func TestSwagger_SetSecuritySchemes(t *testing.T) {
	expected := Swagger{
		Components: Component{
			SecuritySchemes: map[string]SecurityScheme{
				"bearerAuth": {
					Type:         "http",
					Scheme:       "bearer",
					BearerFormat: "JWT",
				},
			},
		},
	}

	actual := Swagger{Components: Component{SecuritySchemes: map[string]SecurityScheme{}}}
	actual.SetSecuritySchemes()

	assert.Equal(t, expected, actual)
}

func TestSwagger_SetPaths(t *testing.T) {
	expected := Swagger{
		OpenAPI: "3.0.3",
		Paths: map[string]Path{
			"/v1/test": {
				Get: Operation{
					Tags:        []string{"testTag"},
					Description: "testDesc",
					RequestBody: RequestBody{
						Content:  map[string]Content(nil),
						Required: false,
					},
					Responses: map[string]Response{
						"200": {
							Description: "Success",
							Headers: map[string]Header{
								"Authorization": {
									Description: "",
									Schema: HeaderSchema{
										Type:   "string",
										Format: "string",
									},
								},
							}, Content: map[string]Content{
								"application/json": {
									Schema: ContentSchema{Ref: "#/components/schemas/Res"},
								},
							},
						}, "500": {
							Description: "Error",
							Headers:     map[string]Header(nil),
							Content: map[string]Content{
								"application/json": {Schema: ContentSchema{Ref: "#/components/schemas/Error"}},
							},
						},
					},
					Security: []map[string][]string{
						{"bearerAuth": {}},
					},
					Parameters: []Parameter{
						{
							Name:    "elma",
							In:      "query",
							Explode: true,
							Schema: HeaderSchema{
								Type: "string",
							},
						},
						{
							Name: "userId",
							In:   "query",
							Schema: HeaderSchema{
								Type:   "integer",
								Format: "int32",
							},
						},
					},
				},
			},
		},
		Components: Component{
			Schemas:         map[string]ComponentSchema{},
			RequestBodies:   map[string]RequestBody{},
			SecuritySchemes: map[string]SecurityScheme{},
		},
	}

	r := TestRoute{
		Name:        "testRoute",
		Description: "testDesc",
		Method:      "GET",
		Path:        "/v1/test",
		Headers: map[string]string{
			"Authorization": "abcdef",
		},
		Middlewares: []func(http.Handler) http.Handler{auth},
		Request:     Req{},
		Response:    Res{},
	}

	actual := Swagger{
		OpenAPI: "3.0.3",
		Paths:   map[string]Path{},
		Components: Component{
			Schemas:         map[string]ComponentSchema{},
			RequestBodies:   map[string]RequestBody{},
			SecuritySchemes: map[string]SecurityScheme{},
		},
	}

	actual.SetPaths(r, "testTag")

	assert.Equal(t, expected, actual)
}

func TestOperation_SetSecurity(t *testing.T) {
	expected := Operation{
		Security: []map[string][]string{
			{"bearerAuth": {}},
		},
	}

	actual := Operation{
		Security: []map[string][]string{},
	}

	actual.SetSecurity([]func(http.Handler) http.Handler{auth})

	assert.Equal(t, expected, actual)
}

func TestOperation_SetRequestBody(t *testing.T) {
	type A struct {
		B string `query:"b"`
	}

	type C struct {
		D int `path:"d"`
	}

	type X struct {
		A
		Name    string   `query:"name" description:"abcdef" required:"true" deprecated:"true"`
		Planes  []string `query:"planes"`
		ID      float64  `path:"id"`
		Test    int      `json:"test"`
		boolean bool     `query:"boolean"`
		C       C        `json:"c"`
	}

	expected := Operation{
		RequestBody: RequestBody{
			Content: map[string]Content{
				"application/json": {
					Schema: ContentSchema{
						Ref: ReferencePrefix + "X",
					},
				},
			},
		},
	}

	actual := Operation{}
	actual.SetRequestBody(reflect.TypeOf(X{}))

	assert.Equal(t, expected, actual)
}

func TestOperation_SetParameters(t *testing.T) {
	type X struct {
		Name    string   `query:"name" description:"abcdef" required:"true" deprecated:"true"`
		Planes  []string `query:"planes"`
		ID      float64  `path:"id"`
		Test    int      `json:"test"`
		boolean bool     `query:"boolean"`
	}

	expected := Operation{
		Parameters: []Parameter{
			{
				Name:        "name",
				In:          "query",
				Description: "abcdef",
				Required:    true,
				Deprecated:  true,
				Schema: HeaderSchema{
					Type: "string",
				},
			},
			{
				Name:    "planes",
				In:      "query",
				Explode: true,
				Schema: HeaderSchema{
					Type: "string",
				},
			},
			{
				Name: "id",
				In:   "path",
				Schema: HeaderSchema{
					Type:   "number",
					Format: "double",
				},
			},
		},
	}

	actual := Operation{}
	actual.SetParameters(reflect.TypeOf(X{}))

	assert.Equal(t, expected, actual)
}

func TestOperation_SetResponses(t *testing.T) {
	type X struct{}

	expected := Operation{
		Responses: map[string]Response{
			"200": {
				Description: "Success",
				Headers: map[string]Header{
					"Authorization": {Schema: HeaderSchema{"string", "string"}},
				},
				Content: map[string]Content{
					"application/json": {
						Schema: ContentSchema{
							Ref: ReferencePrefix + "X",
						},
					},
				},
			},
			"500": {
				Description: "Error",
				Content: map[string]Content{
					"application/json": {
						Schema: ContentSchema{
							Ref: ReferencePrefix + "Error",
						},
					},
				},
			},
		},
	}

	actual := Operation{
		Responses: map[string]Response{},
	}
	actual.SetResponses(reflect.TypeOf(X{}), map[string]string{"Authorization": ""})

	assert.Equal(t, expected, actual)
}

func TestSwagger_SetSchema(t *testing.T) {
	type X struct {
		B bool `json:"b"`
		s string
	}

	type Given struct {
		Name       string           `query:"name"`
		s          string           `header:"s"`
		B          bool             `path:"b"`
		R          rune             `header:"r"`
		I          int              `json:"i"`
		U          uint             `json:"u"`
		F          float64          `json:"f"`
		T          time.Time        `query:"t"`
		X          []X              `json:"x"`
		Pagination pagination.Model `json:"pagination"`
	}

	testCases := map[string]struct {
		Given    any
		Expected Swagger
	}{
		"should add component schemas": {
			Given: Given{},
			Expected: Swagger{
				Components: Component{
					Schemas: map[string]ComponentSchema{
						"Given": {
							XSwaggerRouterModel: RouterModelPrefix + "Given",
							Properties: map[string]Property{
								"i": {
									Type:   "integer",
									Format: "int32",
								},
								"u": {
									Type:   "integer",
									Format: "int32",
								},
								"f": {
									Type:   "number",
									Format: "double",
								},
								"x": {
									Type: "array",
									Items: Items{
										Ref: ReferencePrefix + "X",
									},
								},
								"pagination": {
									Ref: ReferencePrefix + "Model",
								},
							},
							Type: "object",
						},
						"X": {
							XSwaggerRouterModel: RouterModelPrefix + "X",
							Properties: map[string]Property{
								"b": {
									Type:   "boolean",
									Format: "",
								},
							},
							Type: "object",
						},
						"Model": {
							XSwaggerRouterModel: RouterModelPrefix + "Model",
							Properties: map[string]Property{
								"rel": {
									Type: "string",
								},
								"next": {
									Type: "string",
								},
								"prev": {
									Type: "string",
								},
							},
							Type: "object",
						},
					},
				},
			},
		},
		"should not add component schemas when given has no exported field": {
			Given: struct {
				name string `query:"name"`
			}{},
			Expected: Swagger{
				Components: Component{
					Schemas: map[string]ComponentSchema{
						"": {
							XSwaggerRouterModel: RouterModelPrefix,
							Properties:          map[string]Property{},
							Type:                "object",
						},
					},
				},
			},
		},
		"should not add component schemas when given has no valid tag": {
			Given: struct {
				Name string `test:"name"`
			}{},
			Expected: Swagger{
				Components: Component{
					Schemas: map[string]ComponentSchema{
						"": {
							XSwaggerRouterModel: RouterModelPrefix,
							Properties:          map[string]Property{},
							Type:                "object",
						},
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		s := Swagger{
			Components: Component{
				Schemas: map[string]ComponentSchema{},
			},
		}
		t.Run(name, func(t *testing.T) {
			s.SetSchema(reflect.TypeOf(tc.Given))
			assert.Equal(t, tc.Expected, s)
		})
	}
}

func Test_StringSliceContains(t *testing.T) {
	type Given struct {
		Slice []string
		S     string
	}

	testCases := map[string]struct {
		Given    Given
		Expected bool
	}{
		"should return true when slice contains the given string": {
			Given: Given{
				Slice: []string{"car", "tren"},
				S:     "tren",
			},
			Expected: true,
		},
		"should return false when slice doesn't contain the given string": {
			Given: Given{
				Slice: []string{"car", "tren"},
				S:     "plane",
			},
			Expected: false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			actual := StringSliceContains(tc.Given.Slice, tc.Given.S)
			assert.Equal(t, tc.Expected, actual)
		})
	}
}

func Test_IsValidTag(t *testing.T) {
	type Given struct {
		Slice []string
		Tag   reflect.StructTag
	}

	type Expected struct {
		Key   string
		Value string
	}

	testCases := map[string]struct {
		Given    Given
		Expected Expected
	}{
		"should return tag key and value when tag is valid": {
			Given: Given{
				Slice: []string{"json", "query", "path"},
				Tag:   "query:\"test\"",
			},
			Expected: Expected{
				Key:   "test",
				Value: "query",
			},
		},
		"should return empty tag key and value when tag is not valid": {
			Given: Given{
				Slice: []string{"json", "query", "path"},
				Tag:   "header:\"test\"",
			},
			Expected: Expected{
				Key:   "",
				Value: "",
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			actualKey, actualVal := isValidTag(tc.Given.Slice, tc.Given.Tag)

			assert.Equal(t, tc.Expected.Key, actualKey)
			assert.Equal(t, tc.Expected.Value, actualVal)
		})
	}
}
