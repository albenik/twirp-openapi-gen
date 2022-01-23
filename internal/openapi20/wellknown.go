package openapi20

import (
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	TypeInteger = "integer"
	TypeNumber  = "number"
	TypeString  = "string"
	TypeBoolean = "boolean"
	TypeArray   = "array"
	TypeObject  = "object"

	FormatNone     = ""
	FormatInt32    = "int32"
	FormatInt64    = "int64"
	FormatFloat    = "float"
	FormatDouble   = "double"
	FormatByte     = "byte"
	FormatBinary   = "binary"
	FormatDate     = "date"
	FormatDateTime = "date-time"
	FormatPassword = "password"
)

// Datatypes from https://swagger.io/specification/v2/
var (
	IntegerType  = DataTypeDesc{Type: TypeInteger, Format: FormatInt32}
	LongType     = DataTypeDesc{Type: TypeInteger, Format: FormatInt64}
	Uint64Type   = DataTypeDesc{Type: TypeString, Format: "uint64"} // non-standart type
	FloatType    = DataTypeDesc{Type: TypeNumber, Format: FormatFloat}
	DoubleType   = DataTypeDesc{Type: TypeNumber, Format: FormatDouble}
	StringType   = DataTypeDesc{Type: TypeString, Format: FormatNone}
	ByteType     = DataTypeDesc{Type: TypeString, Format: FormatByte}
	BinaryType   = DataTypeDesc{Type: TypeString, Format: FormatBinary}
	BooleanType  = DataTypeDesc{Type: TypeBoolean, Format: FormatNone}
	DateType     = DataTypeDesc{Type: TypeString, Format: FormatDate}
	DateTimeType = DataTypeDesc{Type: TypeString, Format: FormatDateTime}
	PasswordType = DataTypeDesc{Type: TypeString, Format: FormatPassword}

	protoKindTypes = map[protoreflect.Kind]DataTypeDesc{
		protoreflect.Int32Kind:  IntegerType,
		protoreflect.Sint32Kind: IntegerType,
		protoreflect.Uint32Kind: LongType,
		protoreflect.Int64Kind:  LongType,
		protoreflect.Sint64Kind: LongType,
		protoreflect.Uint64Kind: Uint64Type,
		protoreflect.FloatKind:  FloatType,
		protoreflect.DoubleKind: DoubleType,
		protoreflect.BoolKind:   BooleanType,
		protoreflect.StringKind: StringType,
		protoreflect.BytesKind:  ByteType,
	}

	protoWellknownTypes = map[protoreflect.FullName]DataTypeDesc{
		"google.protobuf.StringValue": StringType,
		"google.protobuf.BytesValue":  ByteType,
		"google.protobuf.BoolValue":   BooleanType,
		"google.protobuf.Int32Value":  IntegerType,
		"google.protobuf.Int64Value":  LongType,
		"google.protobuf.UInt32Value": LongType,
		"google.protobuf.UInt64Value": Uint64Type,
		"google.protobuf.FloatValue":  FloatType,
		"google.protobuf.DoubleValue": DoubleType,
		"google.protobuf.Timestamp":   DateTimeType,
		"google.protobuf.Duration":    StringType,
	}
)

type DataTypeDesc struct {
	Type   string
	Format string
}

func (d DataTypeDesc) Schema() *Schema {
	return &Schema{Def: &SchemaDef{
		Type:   StringOrArray{d.Type},
		Format: d.Format,
	}}
}
