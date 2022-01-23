package openapi20

import (
	"bytes"
	"encoding/json"
	"sort"
	"strconv"
)

type Version20 struct{}

func (Version20) MarshalJSON() ([]byte, error) {
	return []byte(`"2.0"`), nil
}

type Swagger struct {
	Swagger     Version20   `json:"swagger"`
	Info        Info        `json:"info"`
	Host        string      `json:"host,omitempty"`
	BasePath    string      `json:"basePath"`
	Schemes     []string    `json:"schemes,omitempty"`
	Consumes    []string    `json:"consumes,omitempty"`
	Produces    []string    `json:"produces,omitempty"`
	Paths       Paths       `json:"paths"`
	Definitions Definitions `json:"definitions,omitempty"`
}

type Info struct {
	Title          string   `json:"title"`
	Version        string   `json:"version"`
	Description    string   `json:"description,omitempty"`
	TermsOfService string   `json:"termsOfService,omitempty"`
	Contact        *Contact `json:"contact,omitempty"`
	License        *License `json:"license,omitempty"`
}

type Schema struct {
	Ref string
	Def *SchemaDef
}

func (s *Schema) MarshalJSON() ([]byte, error) {
	if s.Ref != "" {
		ref := struct {
			Ref string `json:"$ref"`
		}{Ref: s.Ref}
		return json.Marshal(&ref)
	}
	return json.Marshal(s.Def)
}

type SchemaDef struct {
	Title                string        `json:"title,omitempty"`
	Description          string        `json:"description,omitempty"`
	Type                 StringOrArray `json:"type"`
	Format               string        `json:"format,omitempty"`
	Required             []string      `json:"required,omitempty"`
	Default              interface{}   `json:"default,omitempty"`
	Minimum              *float64      `json:"minimum,omitempty"`
	Maximum              *float64      `json:"maximum,omitempty"`
	MultipleOf           *float64      `json:"multipleOf,omitempty"`
	ExclusiveMinimum     bool          `json:"exclusiveMinimum,omitempty"`
	ExclusiveMaximum     bool          `json:"exclusiveMaximum,omitempty"`
	MaxItems             *int64        `json:"maxItems,omitempty"`
	MinItems             *int64        `json:"minItems,omitempty"`
	Properties           Properties    `json:"properties,omitempty"`
	AdditionalProperties *Schema       `json:"additionalProperties,omitempty"`
	Items                *Schema       `json:"items,omitempty"`
	Enum                 []string      `json:"enum,omitempty"`
	AllOf                []*Schema     `json:"allOf,omitempty"`
}

type Properties map[string]*Schema

func (p Properties) MarshalJSON() ([]byte, error) {
	keys := make([]string, 0, len(p))
	for k := range p {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	buf := bytes.NewBuffer(nil)
	buf.WriteRune('{')

	//nolint:varnamelen
	for i, k := range keys {
		b, err := json.Marshal(p[k])
		if err != nil {
			return nil, err
		}

		if i > 0 {
			buf.WriteRune(',')
		}
		buf.WriteString(strconv.Quote(k))
		buf.WriteRune(':')
		buf.Write(b)
	}

	buf.WriteRune('}')
	return buf.Bytes(), nil
}

type Definitions map[string]*Schema

func (d Definitions) MarshalJSON() ([]byte, error) {
	keys := make([]string, 0, len(d))
	for k := range d {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	buf := bytes.NewBuffer(nil)
	buf.WriteRune('{')

	//nolint:varnamelen
	for i, k := range keys {
		b, err := json.Marshal(d[k])
		if err != nil {
			return nil, err
		}

		if i > 0 {
			buf.WriteRune(',')
		}
		buf.WriteString(strconv.Quote(k))
		buf.WriteRune(':')
		buf.Write(b)
	}

	buf.WriteRune('}')
	return buf.Bytes(), nil
}

type Contact struct {
	Name  string `json:"name,omitempty"`
	URL   string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
}

type License struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

type Paths map[string]*Path

func (p Paths) MarshalJSON() ([]byte, error) {
	keys := make([]string, 0, len(p))
	for k := range p {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	buf := bytes.NewBuffer(nil)
	buf.WriteRune('{')

	//nolint:varnamelen
	for i, k := range keys {
		b, err := json.Marshal(p[k])
		if err != nil {
			return nil, err
		}

		if i > 0 {
			buf.WriteRune(',')
		}
		buf.WriteString(strconv.Quote(k))
		buf.WriteRune(':')
		buf.Write(b)
	}

	buf.WriteRune('}')
	return buf.Bytes(), nil
}

type Path struct {
	GET  *Operation `json:"get,omitempty"`
	POST *Operation `json:"post,omitempty"`
}

type Operation struct {
	OperationID  string        `json:"operationId,omitempty"`
	Summary      string        `json:"summary,omitempty"`
	Description  string        `json:"description,omitempty"`
	Tags         []string      `json:"tags,omitempty"`
	Consumes     []string      `json:"consumes,omitempty"`
	Produces     []string      `json:"produces,omitempty"`
	Parameters   []Parameter   `json:"parameters,omitempty"`
	Responses    Responses     `json:"responses"`
	ExternalDocs *ExternalDocs `json:"externalDocs,omitempty"`
}

type ExternalDocs struct {
	URL         string `json:"url"`
	Description string `json:"description,omitempty"`
}

type Headers map[string]*Header

func (h Headers) MarshalJSON() ([]byte, error) {
	keys := make([]string, 0, len(h))
	for k := range h {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	buf := bytes.NewBuffer(nil)
	buf.WriteRune('{')

	//nolint:varnamelen
	for i, k := range keys {
		b, err := json.Marshal(h[k])
		if err != nil {
			return nil, err
		}

		if i > 0 {
			buf.WriteRune(',')
		}
		buf.WriteString(strconv.Quote(k))
		buf.WriteRune(':')
		buf.Write(b)
	}

	buf.WriteRune('}')
	return buf.Bytes(), nil
}

type Header struct {
	Description string `json:"description,omitempty"`
	Type        string `json:"type"`
	Format      string `json:"format,omitempty"`
}

type Parameter struct {
	Name        string  `json:"name"`
	In          string  `json:"in"`
	Required    bool    `json:"required"`
	Description string  `json:"description,omitempty"`
	Schema      *Schema `json:"schema,omitempty"`
}

type Responses map[int]*Response

func (r Responses) MarshalJSON() ([]byte, error) {
	keys := make([]int, 0, len(r))
	for k := range r {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	buf := bytes.NewBuffer(nil)
	buf.WriteRune('{')

	//nolint:varnamelen
	for i, k := range keys {
		b, err := json.Marshal(r[k])
		if err != nil {
			return nil, err
		}

		if i > 0 {
			buf.WriteRune(',')
		}
		buf.WriteRune('"')
		buf.WriteString(strconv.FormatInt(int64(k), 10)) //nolint:gomnd
		buf.WriteString(`":`)
		buf.Write(b)
	}

	buf.WriteRune('}')
	return buf.Bytes(), nil
}

type Response struct {
	Description string  `json:"description"`
	Headers     Headers `json:"headers,omitempty"`
	Schema      *Schema `json:"schema,omitempty"`
}

type StringOrArray []string

func (s StringOrArray) MarshalJSON() ([]byte, error) {
	if len(s) == 1 {
		return json.Marshal([]string(s)[0])
	}
	return json.Marshal([]string(s))
}
