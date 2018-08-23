// Code generated by protoc-gen-go. DO NOT EDIT.
// source: test.proto

package example

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type FOO int32

const (
	FOO_x FOO = 17
)

var FOO_name = map[int32]string{
	17: "x",
}
var FOO_value = map[string]int32{
	"x": 17,
}

func (x FOO) Enum() *FOO {
	p := new(FOO)
	*p = x
	return p
}
func (x FOO) String() string {
	return proto.EnumName(FOO_name, int32(x))
}
func (x *FOO) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(FOO_value, data, "FOO")
	if err != nil {
		return err
	}
	*x = FOO(value)
	return nil
}
func (FOO) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_test_d1cdc29eb1e2b687, []int{0}
}

type Test struct {
	Label                *string             `protobuf:"bytes,1,req,name=label" json:"label,omitempty"`
	Type                 *int32              `protobuf:"varint,2,opt,name=type,def=77" json:"type,omitempty"`
	Reps                 []int64             `protobuf:"varint,3,rep,name=reps" json:"reps,omitempty"`
	Optionalgroup        *Test_OptionalGroup `protobuf:"group,4,opt,name=OptionalGroup" json:"optionalgroup,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *Test) Reset()         { *m = Test{} }
func (m *Test) String() string { return proto.CompactTextString(m) }
func (*Test) ProtoMessage()    {}
func (*Test) Descriptor() ([]byte, []int) {
	return fileDescriptor_test_d1cdc29eb1e2b687, []int{0}
}
func (m *Test) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Test.Unmarshal(m, b)
}
func (m *Test) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Test.Marshal(b, m, deterministic)
}
func (dst *Test) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Test.Merge(dst, src)
}
func (m *Test) XXX_Size() int {
	return xxx_messageInfo_Test.Size(m)
}
func (m *Test) XXX_DiscardUnknown() {
	xxx_messageInfo_Test.DiscardUnknown(m)
}

var xxx_messageInfo_Test proto.InternalMessageInfo

const Default_Test_Type int32 = 77

func (m *Test) GetLabel() string {
	if m != nil && m.Label != nil {
		return *m.Label
	}
	return ""
}

func (m *Test) GetType() int32 {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return Default_Test_Type
}

func (m *Test) GetReps() []int64 {
	if m != nil {
		return m.Reps
	}
	return nil
}

func (m *Test) GetOptionalgroup() *Test_OptionalGroup {
	if m != nil {
		return m.Optionalgroup
	}
	return nil
}

type Test_OptionalGroup struct {
	RequiredField        *string  `protobuf:"bytes,5,req,name=RequiredField" json:"RequiredField,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Test_OptionalGroup) Reset()         { *m = Test_OptionalGroup{} }
func (m *Test_OptionalGroup) String() string { return proto.CompactTextString(m) }
func (*Test_OptionalGroup) ProtoMessage()    {}
func (*Test_OptionalGroup) Descriptor() ([]byte, []int) {
	return fileDescriptor_test_d1cdc29eb1e2b687, []int{0, 0}
}
func (m *Test_OptionalGroup) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Test_OptionalGroup.Unmarshal(m, b)
}
func (m *Test_OptionalGroup) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Test_OptionalGroup.Marshal(b, m, deterministic)
}
func (dst *Test_OptionalGroup) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Test_OptionalGroup.Merge(dst, src)
}
func (m *Test_OptionalGroup) XXX_Size() int {
	return xxx_messageInfo_Test_OptionalGroup.Size(m)
}
func (m *Test_OptionalGroup) XXX_DiscardUnknown() {
	xxx_messageInfo_Test_OptionalGroup.DiscardUnknown(m)
}

var xxx_messageInfo_Test_OptionalGroup proto.InternalMessageInfo

func (m *Test_OptionalGroup) GetRequiredField() string {
	if m != nil && m.RequiredField != nil {
		return *m.RequiredField
	}
	return ""
}

func init() {
	proto.RegisterType((*Test)(nil), "example.Test")
	proto.RegisterType((*Test_OptionalGroup)(nil), "example.Test.OptionalGroup")
	proto.RegisterEnum("example.FOO", FOO_name, FOO_value)
}

func init() { proto.RegisterFile("test.proto", fileDescriptor_test_d1cdc29eb1e2b687) }

var fileDescriptor_test_d1cdc29eb1e2b687 = []byte{
	// 179 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0xcd, 0x31, 0x0a, 0xc2, 0x30,
	0x14, 0xc6, 0x71, 0xd2, 0x24, 0x88, 0x8f, 0x06, 0x6a, 0x40, 0x08, 0xba, 0x04, 0x07, 0x09, 0x0e,
	0x19, 0xba, 0x14, 0x3c, 0x40, 0x1d, 0x0b, 0xe2, 0x05, 0x2a, 0x7d, 0x48, 0x21, 0x9a, 0x98, 0xa6,
	0x50, 0x8f, 0xe1, 0x8d, 0xa5, 0xd0, 0xc5, 0xf5, 0xe3, 0xcf, 0xef, 0x03, 0x48, 0x38, 0x24, 0x1b,
	0xa2, 0x4f, 0x5e, 0xae, 0x70, 0x6a, 0x9f, 0xc1, 0xe1, 0xe1, 0x4b, 0x80, 0xdd, 0x70, 0x48, 0x52,
	0x00, 0x77, 0xed, 0x1d, 0x9d, 0x22, 0x3a, 0x33, 0x6b, 0x59, 0x00, 0x4b, 0x9f, 0x80, 0x2a, 0xd3,
	0xc4, 0xf0, 0x73, 0x56, 0x55, 0x32, 0x07, 0x16, 0x31, 0x0c, 0x8a, 0x6a, 0x6a, 0xa8, 0x2c, 0x41,
	0xf8, 0x90, 0x7a, 0xff, 0x6a, 0xdd, 0x23, 0xfa, 0x31, 0x28, 0xa6, 0x89, 0x81, 0x72, 0x6f, 0x17,
	0xd8, 0xce, 0xa8, 0x6d, 0x96, 0xe4, 0x32, 0x27, 0xbb, 0x23, 0x88, 0xbf, 0x41, 0x6e, 0x41, 0x5c,
	0xf1, 0x3d, 0xf6, 0x11, 0xbb, 0xba, 0x47, 0xd7, 0x29, 0x3e, 0x7f, 0x9f, 0x72, 0xa0, 0x75, 0xd3,
	0x48, 0x0e, 0x64, 0x2a, 0x36, 0xbf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xc4, 0xf7, 0x55, 0x52, 0xb7,
	0x00, 0x00, 0x00,
}
