package openapi20

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var versionRx = regexp.MustCompile(`(?s)\.(v\d+[^.]*)\.[^.]+$`)

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) GenerateSchema(host, preifx string, service *protogen.Service) (*Swagger, error) {
	ver := "undefined"
	match := versionRx.FindStringSubmatch(string(service.Desc.FullName()))
	if match != nil {
		ver = match[1]
	}

	swagger := &Swagger{
		Swagger: Version20{},
		Info: Info{
			Title:          string(service.Desc.Name()),
			Version:        ver,
			Description:    string(service.Comments.Leading),
			TermsOfService: "",
			Contact:        nil,
			License:        nil,
		},
		Host:        host,
		BasePath:    "/" + strings.TrimPrefix(path.Join(preifx, string(service.Desc.FullName())), "/"),
		Schemes:     []string{"https"},
		Consumes:    []string{"application/json"},
		Produces:    []string{"application/json"},
		Paths:       make(Paths),
		Definitions: make(Definitions),
	}

	for _, method := range service.Methods {
		summary, description := splitComment(string(method.Comments.Leading))

		swagger.Paths["/"+string(method.Desc.Name())] = &Path{
			POST: &Operation{
				OperationID: string(method.Desc.FullName()),
				Summary:     summary,
				Description: description,
				Parameters: []Parameter{{
					Name:        "body",
					In:          "body",
					Required:    true,
					Description: "",
					Schema:      ref(method.Input.Desc.Name()),
				}},
				Responses: Responses{
					200: &Response{
						Description: "",
						Headers:     nil,
						Schema:      ref(method.Output.Desc.Name()),
					},
				},
			},
		}

		if err := g.collectMessageDefinitions(swagger.Definitions, method.Input); err != nil {
			return nil, err
		}

		if err := g.collectMessageDefinitions(swagger.Definitions, method.Output); err != nil {
			return nil, err
		}
	}

	return swagger, nil
}

func (g *Generator) collectMessageDefinitions(defs Definitions, message *protogen.Message) error {
	title := string(message.Desc.Name())
	summary, description := splitComment(string(message.Comments.Leading))
	if summary != "" {
		title += " — " + strings.TrimRight(strings.TrimSpace(summary), ".,;!?")
	}

	def := &SchemaDef{
		Type:        StringOrArray{"object"},
		Title:       title,
		Description: description,
		Properties:  make(Properties),
	}

	for _, field := range message.Fields {
		s, err := g.fieldSchema(defs, field)
		if err != nil {
			return err
		}

		def.Properties[string(field.Desc.Name())] = s
	}

	defs[string(message.Desc.Name())] = &Schema{Def: def}
	return nil
}

func (g *Generator) fieldSchema(defs Definitions, field *protogen.Field) (*Schema, error) {
	//nolint:exhaustive
	switch field.Desc.Kind() {
	case protoreflect.EnumKind:
		schema := StringType.Schema()
		schema.Def.Enum = make([]string, 0, len(field.Enum.Values))

		for _, v := range field.Enum.Values {
			schema.Def.Enum = append(schema.Def.Enum, string(v.Desc.Name()))
		}

		return schema, nil

	case protoreflect.MessageKind:
		if t, ok := protoWellknownTypes[field.Message.Desc.FullName()]; ok {
			return t.Schema(), nil
		}

		if field.Desc.IsMap() {
			return g.mapFieldSchema(defs, field)
		}

		if err := g.collectMessageDefinitions(defs, field.Message); err != nil {
			return nil, err
		}

		fieldRef := ref(field.Message.Desc.Name())

		if field.Desc.IsList() {
			return &Schema{Def: &SchemaDef{
				Type:  StringOrArray{TypeArray},
				Items: fieldRef,
			}}, nil
		}

		return fieldRef, nil

	default:
		if t, ok := protoKindTypes[field.Desc.Kind()]; ok {
			fieldSchema := t.Schema()

			if field.Desc.IsList() {
				return &Schema{Def: &SchemaDef{
					Type:  StringOrArray{TypeArray},
					Items: fieldSchema,
				}}, nil
			}

			return fieldSchema, nil
		}
		return nil, fmt.Errorf("field %s: unknown kind %s", field.Desc.FullName(), field.Desc.Kind())
	}
}

func (g *Generator) mapFieldSchema(defs Definitions, field *protogen.Field) (*Schema, error) {
	k := field.Desc.MapKey().Kind()
	if k != protoreflect.StringKind {
		return nil, fmt.Errorf("field %s: unsupported map key kind %s", field.Desc.FullName(), k)
	}

	var vs *Schema

	val := field.Desc.MapValue()
	for _, ff := range field.Message.Fields {
		if ff.Desc == val {
			var err error
			if vs, err = g.fieldSchema(defs, ff); err != nil {
				return nil, err
			}
			break
		}
	}

	if vs == nil {
		return nil, fmt.Errorf("field %s: map value type %s not found", field.Desc.FullName(), val.FullName())
	}

	return &Schema{Def: &SchemaDef{Type: StringOrArray{TypeObject}, AdditionalProperties: vs}}, nil
}

func ref(name protoreflect.Name) *Schema {
	return &Schema{Ref: fmt.Sprintf("#/definitions/%s", name)}
}

func splitComment(s string) (sum, desc string) {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	sum = lines[0]
	if len(lines) > 1 {
		desc = strings.Join(lines[1:], "\n")
	}
	return
}
