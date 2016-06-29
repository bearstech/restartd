// Code generated by protoc-gen-go.
// source: protocol/command.proto
// DO NOT EDIT!

/*
Package protocol is a generated protocol buffer package.

It is generated from these files:
	protocol/command.proto

It has these top-level messages:
	Message
	Response
*/
package protocol

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

type Message_Commands int32

const (
	Message_status  Message_Commands = 0
	Message_start   Message_Commands = 1
	Message_stop    Message_Commands = 2
	Message_restart Message_Commands = 3
	Message_reload  Message_Commands = 4
)

var Message_Commands_name = map[int32]string{
	0: "status",
	1: "start",
	2: "stop",
	3: "restart",
	4: "reload",
}
var Message_Commands_value = map[string]int32{
	"status":  0,
	"start":   1,
	"stop":    2,
	"restart": 3,
	"reload":  4,
}

func (x Message_Commands) Enum() *Message_Commands {
	p := new(Message_Commands)
	*p = x
	return p
}
func (x Message_Commands) String() string {
	return proto.EnumName(Message_Commands_name, int32(x))
}
func (x *Message_Commands) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Message_Commands_value, data, "Message_Commands")
	if err != nil {
		return err
	}
	*x = Message_Commands(value)
	return nil
}
func (Message_Commands) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type Response_Codes int32

const (
	Response_success     Response_Codes = 0
	Response_suc_reload  Response_Codes = 1
	Response_suc_status  Response_Codes = 2
	Response_suc_start   Response_Codes = 3
	Response_suc_stop    Response_Codes = 4
	Response_suc_restart Response_Codes = 5
	Response_error       Response_Codes = 100
	Response_err_missing Response_Codes = 101
	Response_err_parsing Response_Codes = 102
	Response_err_reading Response_Codes = 103
	Response_err_cmd     Response_Codes = 104
	Response_err_status  Response_Codes = 105
	Response_err_start   Response_Codes = 106
	Response_err_stop    Response_Codes = 107
	Response_err_restart Response_Codes = 108
	Response_err_reload  Response_Codes = 109
)

var Response_Codes_name = map[int32]string{
	0:   "success",
	1:   "suc_reload",
	2:   "suc_status",
	3:   "suc_start",
	4:   "suc_stop",
	5:   "suc_restart",
	100: "error",
	101: "err_missing",
	102: "err_parsing",
	103: "err_reading",
	104: "err_cmd",
	105: "err_status",
	106: "err_start",
	107: "err_stop",
	108: "err_restart",
	109: "err_reload",
}
var Response_Codes_value = map[string]int32{
	"success":     0,
	"suc_reload":  1,
	"suc_status":  2,
	"suc_start":   3,
	"suc_stop":    4,
	"suc_restart": 5,
	"error":       100,
	"err_missing": 101,
	"err_parsing": 102,
	"err_reading": 103,
	"err_cmd":     104,
	"err_status":  105,
	"err_start":   106,
	"err_stop":    107,
	"err_restart": 108,
	"err_reload":  109,
}

func (x Response_Codes) Enum() *Response_Codes {
	p := new(Response_Codes)
	*p = x
	return p
}
func (x Response_Codes) String() string {
	return proto.EnumName(Response_Codes_name, int32(x))
}
func (x *Response_Codes) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(Response_Codes_value, data, "Response_Codes")
	if err != nil {
		return err
	}
	*x = Response_Codes(value)
	return nil
}
func (Response_Codes) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1, 0} }

type Message struct {
	Service          *string           `protobuf:"bytes,1,req,name=Service" json:"Service,omitempty"`
	Command          *Message_Commands `protobuf:"varint,2,req,name=Command,enum=protocol.Message_Commands" json:"Command,omitempty"`
	XXX_unrecognized []byte            `json:"-"`
}

func (m *Message) Reset()                    { *m = Message{} }
func (m *Message) String() string            { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()               {}
func (*Message) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Message) GetService() string {
	if m != nil && m.Service != nil {
		return *m.Service
	}
	return ""
}

func (m *Message) GetCommand() Message_Commands {
	if m != nil && m.Command != nil {
		return *m.Command
	}
	return Message_status
}

type Response struct {
	Code             *Response_Codes `protobuf:"varint,1,req,name=Code,enum=protocol.Response_Codes" json:"Code,omitempty"`
	Message          *string         `protobuf:"bytes,2,opt,name=Message" json:"Message,omitempty"`
	XXX_unrecognized []byte          `json:"-"`
}

func (m *Response) Reset()                    { *m = Response{} }
func (m *Response) String() string            { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()               {}
func (*Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Response) GetCode() Response_Codes {
	if m != nil && m.Code != nil {
		return *m.Code
	}
	return Response_success
}

func (m *Response) GetMessage() string {
	if m != nil && m.Message != nil {
		return *m.Message
	}
	return ""
}

func init() {
	proto.RegisterType((*Message)(nil), "protocol.Message")
	proto.RegisterType((*Response)(nil), "protocol.Response")
	proto.RegisterEnum("protocol.Message_Commands", Message_Commands_name, Message_Commands_value)
	proto.RegisterEnum("protocol.Response_Codes", Response_Codes_name, Response_Codes_value)
}

func init() { proto.RegisterFile("protocol/command.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 306 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x44, 0x8f, 0x41, 0x4e, 0xf3, 0x30,
	0x10, 0x85, 0x7f, 0xe7, 0x4f, 0x49, 0x32, 0x85, 0x62, 0x79, 0x81, 0x2a, 0x56, 0x28, 0x0b, 0x54,
	0x09, 0x29, 0x48, 0xbd, 0x42, 0xd9, 0xb2, 0x81, 0x03, 0x54, 0x96, 0x6d, 0x4a, 0xa0, 0xa9, 0x23,
	0x4f, 0xca, 0x2d, 0x38, 0x14, 0xe7, 0x62, 0xc3, 0x4c, 0x6d, 0xab, 0xbb, 0x79, 0xef, 0xcd, 0xf8,
	0x7d, 0x86, 0x9b, 0x31, 0xf8, 0xc9, 0x1b, 0xbf, 0x7f, 0x34, 0x7e, 0x18, 0xf4, 0xc1, 0x76, 0x27,
	0x43, 0xd5, 0xd9, 0x6f, 0xbf, 0x05, 0x54, 0xcf, 0x0e, 0x51, 0xef, 0x9c, 0xba, 0x86, 0xea, 0xd5,
	0x85, 0xaf, 0xde, 0xb8, 0xa5, 0xb8, 0x2b, 0x56, 0x8d, 0x7a, 0x80, 0x6a, 0x13, 0xef, 0x96, 0x05,
	0x19, 0x8b, 0xf5, 0x6d, 0x97, 0x0f, 0xbb, 0x74, 0xd4, 0xa5, 0x05, 0x6c, 0x9f, 0xa0, 0xce, 0xb3,
	0x02, 0xb8, 0xc0, 0x49, 0x4f, 0x47, 0x94, 0xff, 0x54, 0x03, 0x33, 0x9a, 0xc3, 0x24, 0x85, 0xaa,
	0xa1, 0xc4, 0xc9, 0x8f, 0xb2, 0x50, 0x73, 0xa8, 0x82, 0x8b, 0xf6, 0x7f, 0xde, 0x0e, 0x6e, 0xef,
	0xb5, 0x95, 0x65, 0xfb, 0x53, 0x40, 0xfd, 0xe2, 0x70, 0xf4, 0x07, 0x74, 0xea, 0x1e, 0xca, 0x8d,
	0xb7, 0x91, 0x66, 0xb1, 0x5e, 0x9e, 0xcb, 0xf3, 0x46, 0xc7, 0x31, 0x32, 0x78, 0xc2, 0x21, 0x4e,
	0xb1, 0x6a, 0xda, 0x5f, 0x01, 0xb3, 0x18, 0x51, 0x11, 0x1e, 0x8d, 0xa1, 0x94, 0x50, 0x16, 0x00,
	0x24, 0xb6, 0xa9, 0x4c, 0x64, 0x9d, 0x50, 0x0b, 0x75, 0x05, 0x4d, 0xd2, 0x27, 0xae, 0x4b, 0xa8,
	0xa3, 0x24, 0xe4, 0x92, 0x4a, 0xe6, 0xf1, 0x38, 0xc6, 0x33, 0xfe, 0x98, 0x0b, 0xc1, 0x07, 0x69,
	0x39, 0xa3, 0x71, 0x3b, 0xf4, 0x88, 0xfd, 0x61, 0x27, 0x5d, 0x36, 0x46, 0x1d, 0x4e, 0xc6, 0x5b,
	0x36, 0x82, 0xd3, 0x96, 0x8d, 0x1d, 0x83, 0xb1, 0x61, 0x06, 0x2b, 0xdf, 0x19, 0x84, 0x45, 0x02,
	0xe9, 0x19, 0x24, 0x69, 0x6a, 0xfa, 0x60, 0x90, 0x28, 0x09, 0xe4, 0xf3, 0xfc, 0x54, 0x8c, 0xf7,
	0xf9, 0x3a, 0x7d, 0x6b, 0xf8, 0x0b, 0x00, 0x00, 0xff, 0xff, 0xf3, 0xbe, 0x7e, 0xd1, 0xf6, 0x01,
	0x00, 0x00,
}
