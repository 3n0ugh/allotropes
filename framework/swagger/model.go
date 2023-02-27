package swagger

type Swagger struct {
	OpenAPI      string          `yaml:"openapi,omitempty"`
	Info         Info            `yaml:"info,omitempty"`
	ExternalDocs ExternalDoc     `yaml:"externalDocs,omitempty"`
	Servers      []Server        `yaml:"servers,omitempty"`
	Tags         []Tag           `yaml:"tags,omitempty"`
	Paths        map[string]Path `yaml:"paths,omitempty"`
	Components   Component       `yaml:"components,omitempty"`
}

type Info struct {
	Title          string  `yaml:"title,omitempty"`
	Description    string  `yaml:"description,omitempty"`
	TermsOfService string  `yaml:"termsOfService,omitempty"`
	Contact        Contact `yaml:"contact,omitempty"`
	License        License `yaml:"license,omitempty"`
	Version        string  `yaml:"version,omitempty"`
}

type Contact struct {
	Name  string `yaml:"name,omitempty"`
	URL   string `yaml:"url,omitempty"`
	Email string `yaml:"email,omitempty"`
}

type License struct {
	Name       string `yaml:"name,omitempty"`
	Identifier string `yaml:"identifier,omitempty"`
	URL        string `yaml:"url,omitempty"`
}

type Server struct {
	URL string `yaml:"url,omitempty"`
}

type ExternalDoc struct {
	URL         string `yaml:"url,omitempty"`
	Description string `yaml:"description,omitempty"`
}

type Tag struct {
	Name         string      `yaml:"name,omitempty"`
	Description  string      `yaml:"description,omitempty"`
	ExternalDocs ExternalDoc `yaml:"externalDocs,omitempty"`
}

type Path struct {
	Get     Operation `yaml:"get,omitempty"`
	Put     Operation `yaml:"put,omitempty"`
	Post    Operation `yaml:"post,omitempty"`
	Delete  Operation `yaml:"delete,omitempty"`
	Options Operation `yaml:"options,omitempty"`
	Head    Operation `yaml:"head,omitempty"`
	Patch   Operation `yaml:"patch,omitempty"`
	Trace   Operation `yaml:"trace,omitempty"`
}

type Operation struct {
	Tags        []string              `yaml:"tags,omitempty"`
	Summary     string                `yaml:"summary,omitempty"`
	Description string                `yaml:"description,omitempty"`
	OperationID string                `yaml:"operationId,omitempty"`
	RequestBody RequestBody           `yaml:"requestBody,omitempty"`
	Responses   map[string]Response   `yaml:"responses,omitempty"`
	Security    []map[string][]string `yaml:"security,omitempty"`
	Parameters  []Parameter           `yaml:"parameters,omitempty"`
}

type RequestBody struct {
	Description string             `yaml:"description,omitempty"`
	Content     map[string]Content `yaml:"content,omitempty"`
	Required    bool               `yaml:"required"`
}

type Content struct {
	Schema ContentSchema `yaml:"schema,omitempty"`
}

type ContentSchema struct {
	Ref string `yaml:"$ref,omitempty"`
}

type Response struct {
	Description string             `yaml:"description,omitempty"`
	Headers     map[string]Header  `yaml:"headers,omitempty"`
	Content     map[string]Content `yaml:"content,omitempty"`
}

type Header struct {
	Description string       `yaml:"description,omitempty"`
	Schema      HeaderSchema `yaml:"schema,omitempty"`
}

type HeaderSchema struct {
	Type   string `yaml:"type,omitempty"`
	Format string `yaml:"format,omitempty"`
}

type Parameter struct {
	Name            string       `yaml:"name,omitempty"`
	In              string       `yaml:"in,omitempty"`
	Description     string       `yaml:"description,omitempty"`
	Required        bool         `yaml:"required"`
	Deprecated      bool         `yaml:"deprecated,omitempty"`
	AllowEmptyValue bool         `yaml:"allowEmptyValue,omitempty"`
	Style           string       `yaml:"style,omitempty"`
	Explode         bool         `yaml:"explode,omitempty"`
	AllowReserved   bool         `yaml:"allowReserved,omitempty"`
	Schema          HeaderSchema `yaml:"schema,omitempty"`
}

type Component struct {
	Schemas         map[string]ComponentSchema `yaml:"schemas,omitempty"`
	RequestBodies   map[string]RequestBody     `yaml:"requestBodies,omitempty"`
	SecuritySchemes map[string]SecurityScheme  `yaml:"securitySchemes,omitempty"`
}

type ComponentSchema struct {
	XSwaggerRouterModel string              `yaml:"x-swagger-router-model"`
	Properties          map[string]Property `yaml:"properties,omitempty"`
	Type                string              `yaml:"type,omitempty"`
}

type Property struct {
	Type        string `yaml:"type,omitempty"`
	Description string `yaml:"description,omitempty"`
	Format      string `yaml:"format,omitempty"`
	Example     any    `yaml:"example,omitempty"`
	Enum        []any  `yaml:"enum,omitempty"`
	Ref         string `yaml:"$ref,omitempty"`
	Items       Items  `yaml:"items,omitempty"`
}

type Items struct {
	Type string `yaml:"type,omitempty"`
	Ref  string `yaml:"$ref,omitempty"`
}

type SecurityScheme struct {
	Type         string `yaml:"type,omitempty"`
	Scheme       string `yaml:"scheme,omitempty"`
	BearerFormat string `yaml:"bearerFormat,omitempty"`
}
