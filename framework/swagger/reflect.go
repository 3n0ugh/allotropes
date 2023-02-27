package swagger

import (
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

const (
	RouterModelPrefix = "io.swagger.model."
	ReferencePrefix   = "#/components/schemas/"
)

var (
	validParameterTags = []string{"header", "query", "path"}
	validBodyTags      = []string{"json"}
	validMethods       = []string{"get", "put", "post", "head", "delete", "trace", "options", "path"}
	authMiddlewares    = []string{"auth"}
)

var Types = map[string]HeaderSchema{
	"string": {Type: "string", Format: ""},

	"bool": {Type: "boolean", Format: ""},

	"rune": {Type: "integer", Format: ""},

	"int":   {Type: "integer", Format: "int32"},
	"int8":  {Type: "integer", Format: ""},
	"int16": {Type: "integer", Format: ""},
	"int32": {Type: "integer", Format: "int32"},
	"int64": {Type: "integer", Format: "int64"},

	"uint":   {Type: "integer", Format: "int32"},
	"uint8":  {Type: "integer", Format: ""},
	"uint16": {Type: "integer", Format: ""},
	"uint32": {Type: "integer", Format: "int32"},
	"uint64": {Type: "integer", Format: "int64"},

	"float32": {Type: "number", Format: "float"},
	"float64": {Type: "number", Format: "double"},

	"time.Time": {Type: "string", Format: "date"},
}

type Route interface {
	GetMethod() string
	GetDescription() string
	GetMiddlewares() []func(http.Handler) http.Handler
	GetHeaders() map[string]string
	GetPath() string
	GetRequestModel() any
	GetResponseModel() any
}

func (s *Swagger) AddTag(name, desc string) {
	s.Tags = append(s.Tags, Tag{
		Name:        name,
		Description: desc,
	})
}

func (s *Swagger) SetSecuritySchemes() {
	s.Components.SecuritySchemes["bearerAuth"] = SecurityScheme{
		Type:         "http",
		Scheme:       "bearer",
		BearerFormat: "JWT",
	}
}

func (s *Swagger) SetPaths(r Route, tag string) {
	if !StringSliceContains(validMethods, r.GetMethod()) {
		return
	}

	operation := Operation{
		Tags:        []string{tag},
		Description: r.GetDescription(),
		RequestBody: RequestBody{},
		Responses:   map[string]Response{},
		Security:    []map[string][]string{},
		Parameters:  []Parameter{},
	}

	operation.SetSecurity(r.GetMiddlewares())
	operation.SetRequestBody(reflect.TypeOf(r.GetRequestModel()))
	operation.SetParameters(reflect.TypeOf(r.GetRequestModel()))
	operation.SetResponses(reflect.TypeOf(r.GetResponseModel()), r.GetHeaders())

	p := s.Paths[r.GetPath()]
	switch r.GetMethod() {
	case http.MethodGet:
		p.Get = operation
	case http.MethodPost:
		p.Post = operation
	case http.MethodPut:
		p.Put = operation
	case http.MethodDelete:
		p.Delete = operation
	case http.MethodOptions:
		p.Options = operation
	case http.MethodHead:
		p.Head = operation
	case http.MethodPatch:
		p.Patch = operation
	case http.MethodTrace:
		p.Trace = operation
	}
	s.Paths[r.GetPath()] = p
}

func (o *Operation) SetSecurity(middlewares []func(http.Handler) http.Handler) {
	for _, m := range middlewares {
		n := runtime.FuncForPC(reflect.ValueOf(m).Pointer()).Name()
		name := strings.Split(n, ".")
		if StringSliceContains(authMiddlewares, name[len(name)-1]) {
			o.Security = append(o.Security, map[string][]string{"bearerAuth": {}})
		}
	}
}

func (o *Operation) SetRequestBody(val reflect.Type) {
	obj := val
	if obj.Kind() == reflect.Ptr {
		obj = obj.Elem()
	}

	var isBodyTagFound bool
	for _, field := range reflect.VisibleFields(obj) {
		if s, _ := isValidTag(validBodyTags, field.Tag); field.Tag != "" && s != "" {
			isBodyTagFound = true
		}
	}

	if !isBodyTagFound {
		return
	}

	o.RequestBody = RequestBody{
		Content: map[string]Content{
			"application/json": {
				Schema: ContentSchema{
					Ref: ReferencePrefix + val.Name(),
				},
			},
		},
	}

	S.SetSchema(val)
}

func (o *Operation) SetParameters(val reflect.Type) {
	obj := val
	if obj.Kind() == reflect.Ptr {
		obj = obj.Elem()
	}

	for _, field := range reflect.VisibleFields(obj) {
		tag, tagKey := isValidTag(validParameterTags, field.Tag)
		if field.Anonymous || tag == "" || !field.IsExported() {
			continue
		}

		required, _ := strconv.ParseBool(field.Tag.Get("required"))
		deprecated, _ := strconv.ParseBool(field.Tag.Get("deprecated"))

		p := Parameter{
			Name:        tag,
			In:          tagKey,
			Description: field.Tag.Get("description"),
			Required:    required,
			Deprecated:  deprecated,
			Schema:      Types[field.Type.String()],
		}

		switch field.Type.Kind() {
		case reflect.Struct, reflect.Map:
		case reflect.Slice, reflect.Array:
			p.Explode = true
			p.Schema = Types[field.Type.Elem().String()]
			fallthrough
		default:
			o.Parameters = append(o.Parameters, p)
		}
	}
}

func (o *Operation) SetResponses(val reflect.Type, headers map[string]string) {
	r := Response{
		Description: "Success",
		Headers:     map[string]Header{},
		Content: map[string]Content{
			"application/json": {
				Schema: ContentSchema{
					Ref: ReferencePrefix + val.Name(),
				},
			},
		},
	}

	for h := range headers {
		r.Headers[h] = Header{Schema: HeaderSchema{"string", "string"}}
	}
	o.Responses["200"] = r

	o.Responses["500"] = Response{
		Description: "Error",
		Content: map[string]Content{
			"application/json": {
				Schema: ContentSchema{
					Ref: ReferencePrefix + "Error",
				},
			},
		},
	}

	S.SetSchema(val)
}

func (s *Swagger) SetSchema(val reflect.Type) {
	o := val
	if o.Kind() == reflect.Ptr {
		o = o.Elem()
	}

	properties := map[string]Property{}
	for _, field := range reflect.VisibleFields(o) {
		tag, _ := isValidTag(validBodyTags, field.Tag)
		if field.Anonymous || tag == "" || !field.IsExported() {
			continue
		}

		if field.Type.Kind() == reflect.Struct && field.Type.String() != "time.Time" {
			s.SetSchema(field.Type)

			properties[tag] = Property{
				Ref: ReferencePrefix + field.Name,
			}
		} else if field.Type.Kind() == reflect.Slice || field.Type.Kind() == reflect.Array {
			if field.Type.Elem().Kind() == reflect.Struct {
				properties[tag] = Property{
					Type:  "array",
					Items: Items{Ref: ReferencePrefix + field.Type.Elem().Name()},
				}
				s.SetSchema(field.Type.Elem())
			} else {
				properties[tag] = Property{
					Type: "array",
					Items: Items{
						Type: Types[field.Type.Elem().Kind().String()].Type,
					},
				}
			}
		} else if field.Type.Kind() != reflect.Map {
			properties[tag] = Property{
				Type:   Types[field.Type.String()].Type,
				Format: Types[field.Type.String()].Format,
			}
		}
	}

	if len(reflect.VisibleFields(o)) > 0 {
		s.Components.Schemas[o.Name()] = ComponentSchema{
			XSwaggerRouterModel: RouterModelPrefix + o.Name(),
			Properties:          properties,
			Type:                "object",
		}
	}
}

func StringSliceContains(slice []string, s string) bool {
	for _, ss := range slice {
		if strings.EqualFold(ss, s) {
			return true
		}
	}
	return false
}

func isValidTag(tags []string, tag reflect.StructTag) (string, string) {
	for _, validTag := range tags {
		if t := tag.Get(validTag); t != "" {
			return t, validTag
		}
	}
	return "", ""
}
