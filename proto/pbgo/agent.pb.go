// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: agent.proto

/*
	Package pbgo is a generated protocol buffer package.

	It is generated from these files:
		agent.proto
		game.proto

	It has these top-level messages:
		AgentConnect
		AgentConnected
		AgentForwardToSvr
		AgentForwardToCli
*/
package pbgo

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import actor "github.com/AsynkronIT/protoactor-go/actor"

import bytes "bytes"

import strings "strings"
import reflect "reflect"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// 节点连接
type AgentConnect struct {
	Sender     *actor.PID `protobuf:"bytes,1,opt,name=Sender" json:"Sender,omitempty"`
	ServerName string     `protobuf:"bytes,2,opt,name=ServerName,proto3" json:"ServerName,omitempty"`
}

func (m *AgentConnect) Reset()                    { *m = AgentConnect{} }
func (*AgentConnect) ProtoMessage()               {}
func (*AgentConnect) Descriptor() ([]byte, []int) { return fileDescriptorAgent, []int{0} }

func (m *AgentConnect) GetSender() *actor.PID {
	if m != nil {
		return m.Sender
	}
	return nil
}

func (m *AgentConnect) GetServerName() string {
	if m != nil {
		return m.ServerName
	}
	return ""
}

type AgentConnected struct {
	Message string     `protobuf:"bytes,1,opt,name=Message,proto3" json:"Message,omitempty"`
	Server  *actor.PID `protobuf:"bytes,2,opt,name=Server" json:"Server,omitempty"`
}

func (m *AgentConnected) Reset()                    { *m = AgentConnected{} }
func (*AgentConnected) ProtoMessage()               {}
func (*AgentConnected) Descriptor() ([]byte, []int) { return fileDescriptorAgent, []int{1} }

func (m *AgentConnected) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *AgentConnected) GetServer() *actor.PID {
	if m != nil {
		return m.Server
	}
	return nil
}

// 转发到服务器
type AgentForwardToSvr struct {
	ServerName string `protobuf:"bytes,1,opt,name=ServerName,proto3" json:"ServerName,omitempty"`
	Uid        string `protobuf:"bytes,2,opt,name=Uid,proto3" json:"Uid,omitempty"`
	Body       []byte `protobuf:"bytes,3,opt,name=Body,proto3" json:"Body,omitempty"`
}

func (m *AgentForwardToSvr) Reset()                    { *m = AgentForwardToSvr{} }
func (*AgentForwardToSvr) ProtoMessage()               {}
func (*AgentForwardToSvr) Descriptor() ([]byte, []int) { return fileDescriptorAgent, []int{2} }

func (m *AgentForwardToSvr) GetServerName() string {
	if m != nil {
		return m.ServerName
	}
	return ""
}

func (m *AgentForwardToSvr) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

func (m *AgentForwardToSvr) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

// 转发到客户端
type AgentForwardToCli struct {
	Uid  string `protobuf:"bytes,1,opt,name=Uid,proto3" json:"Uid,omitempty"`
	Body []byte `protobuf:"bytes,2,opt,name=Body,proto3" json:"Body,omitempty"`
}

func (m *AgentForwardToCli) Reset()                    { *m = AgentForwardToCli{} }
func (*AgentForwardToCli) ProtoMessage()               {}
func (*AgentForwardToCli) Descriptor() ([]byte, []int) { return fileDescriptorAgent, []int{3} }

func (m *AgentForwardToCli) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

func (m *AgentForwardToCli) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

func init() {
	proto.RegisterType((*AgentConnect)(nil), "pbgo.AgentConnect")
	proto.RegisterType((*AgentConnected)(nil), "pbgo.AgentConnected")
	proto.RegisterType((*AgentForwardToSvr)(nil), "pbgo.AgentForwardToSvr")
	proto.RegisterType((*AgentForwardToCli)(nil), "pbgo.AgentForwardToCli")
}
func (this *AgentConnect) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*AgentConnect)
	if !ok {
		that2, ok := that.(AgentConnect)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if !this.Sender.Equal(that1.Sender) {
		return false
	}
	if this.ServerName != that1.ServerName {
		return false
	}
	return true
}
func (this *AgentConnected) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*AgentConnected)
	if !ok {
		that2, ok := that.(AgentConnected)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if this.Message != that1.Message {
		return false
	}
	if !this.Server.Equal(that1.Server) {
		return false
	}
	return true
}
func (this *AgentForwardToSvr) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*AgentForwardToSvr)
	if !ok {
		that2, ok := that.(AgentForwardToSvr)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if this.ServerName != that1.ServerName {
		return false
	}
	if this.Uid != that1.Uid {
		return false
	}
	if !bytes.Equal(this.Body, that1.Body) {
		return false
	}
	return true
}
func (this *AgentForwardToCli) Equal(that interface{}) bool {
	if that == nil {
		if this == nil {
			return true
		}
		return false
	}

	that1, ok := that.(*AgentForwardToCli)
	if !ok {
		that2, ok := that.(AgentForwardToCli)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		if this == nil {
			return true
		}
		return false
	} else if this == nil {
		return false
	}
	if this.Uid != that1.Uid {
		return false
	}
	if !bytes.Equal(this.Body, that1.Body) {
		return false
	}
	return true
}
func (this *AgentConnect) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 6)
	s = append(s, "&pbgo.AgentConnect{")
	if this.Sender != nil {
		s = append(s, "Sender: "+fmt.Sprintf("%#v", this.Sender)+",\n")
	}
	s = append(s, "ServerName: "+fmt.Sprintf("%#v", this.ServerName)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *AgentConnected) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 6)
	s = append(s, "&pbgo.AgentConnected{")
	s = append(s, "Message: "+fmt.Sprintf("%#v", this.Message)+",\n")
	if this.Server != nil {
		s = append(s, "Server: "+fmt.Sprintf("%#v", this.Server)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *AgentForwardToSvr) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 7)
	s = append(s, "&pbgo.AgentForwardToSvr{")
	s = append(s, "ServerName: "+fmt.Sprintf("%#v", this.ServerName)+",\n")
	s = append(s, "Uid: "+fmt.Sprintf("%#v", this.Uid)+",\n")
	s = append(s, "Body: "+fmt.Sprintf("%#v", this.Body)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *AgentForwardToCli) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 6)
	s = append(s, "&pbgo.AgentForwardToCli{")
	s = append(s, "Uid: "+fmt.Sprintf("%#v", this.Uid)+",\n")
	s = append(s, "Body: "+fmt.Sprintf("%#v", this.Body)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringAgent(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *AgentConnect) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AgentConnect) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.Sender != nil {
		dAtA[i] = 0xa
		i++
		i = encodeVarintAgent(dAtA, i, uint64(m.Sender.Size()))
		n1, err := m.Sender.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n1
	}
	if len(m.ServerName) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintAgent(dAtA, i, uint64(len(m.ServerName)))
		i += copy(dAtA[i:], m.ServerName)
	}
	return i, nil
}

func (m *AgentConnected) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AgentConnected) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Message) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintAgent(dAtA, i, uint64(len(m.Message)))
		i += copy(dAtA[i:], m.Message)
	}
	if m.Server != nil {
		dAtA[i] = 0x12
		i++
		i = encodeVarintAgent(dAtA, i, uint64(m.Server.Size()))
		n2, err := m.Server.MarshalTo(dAtA[i:])
		if err != nil {
			return 0, err
		}
		i += n2
	}
	return i, nil
}

func (m *AgentForwardToSvr) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AgentForwardToSvr) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.ServerName) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintAgent(dAtA, i, uint64(len(m.ServerName)))
		i += copy(dAtA[i:], m.ServerName)
	}
	if len(m.Uid) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintAgent(dAtA, i, uint64(len(m.Uid)))
		i += copy(dAtA[i:], m.Uid)
	}
	if len(m.Body) > 0 {
		dAtA[i] = 0x1a
		i++
		i = encodeVarintAgent(dAtA, i, uint64(len(m.Body)))
		i += copy(dAtA[i:], m.Body)
	}
	return i, nil
}

func (m *AgentForwardToCli) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AgentForwardToCli) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if len(m.Uid) > 0 {
		dAtA[i] = 0xa
		i++
		i = encodeVarintAgent(dAtA, i, uint64(len(m.Uid)))
		i += copy(dAtA[i:], m.Uid)
	}
	if len(m.Body) > 0 {
		dAtA[i] = 0x12
		i++
		i = encodeVarintAgent(dAtA, i, uint64(len(m.Body)))
		i += copy(dAtA[i:], m.Body)
	}
	return i, nil
}

func encodeFixed64Agent(dAtA []byte, offset int, v uint64) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	dAtA[offset+4] = uint8(v >> 32)
	dAtA[offset+5] = uint8(v >> 40)
	dAtA[offset+6] = uint8(v >> 48)
	dAtA[offset+7] = uint8(v >> 56)
	return offset + 8
}
func encodeFixed32Agent(dAtA []byte, offset int, v uint32) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintAgent(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *AgentConnect) Size() (n int) {
	var l int
	_ = l
	if m.Sender != nil {
		l = m.Sender.Size()
		n += 1 + l + sovAgent(uint64(l))
	}
	l = len(m.ServerName)
	if l > 0 {
		n += 1 + l + sovAgent(uint64(l))
	}
	return n
}

func (m *AgentConnected) Size() (n int) {
	var l int
	_ = l
	l = len(m.Message)
	if l > 0 {
		n += 1 + l + sovAgent(uint64(l))
	}
	if m.Server != nil {
		l = m.Server.Size()
		n += 1 + l + sovAgent(uint64(l))
	}
	return n
}

func (m *AgentForwardToSvr) Size() (n int) {
	var l int
	_ = l
	l = len(m.ServerName)
	if l > 0 {
		n += 1 + l + sovAgent(uint64(l))
	}
	l = len(m.Uid)
	if l > 0 {
		n += 1 + l + sovAgent(uint64(l))
	}
	l = len(m.Body)
	if l > 0 {
		n += 1 + l + sovAgent(uint64(l))
	}
	return n
}

func (m *AgentForwardToCli) Size() (n int) {
	var l int
	_ = l
	l = len(m.Uid)
	if l > 0 {
		n += 1 + l + sovAgent(uint64(l))
	}
	l = len(m.Body)
	if l > 0 {
		n += 1 + l + sovAgent(uint64(l))
	}
	return n
}

func sovAgent(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozAgent(x uint64) (n int) {
	return sovAgent(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *AgentConnect) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&AgentConnect{`,
		`Sender:` + strings.Replace(fmt.Sprintf("%v", this.Sender), "PID", "actor.PID", 1) + `,`,
		`ServerName:` + fmt.Sprintf("%v", this.ServerName) + `,`,
		`}`,
	}, "")
	return s
}
func (this *AgentConnected) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&AgentConnected{`,
		`Message:` + fmt.Sprintf("%v", this.Message) + `,`,
		`Server:` + strings.Replace(fmt.Sprintf("%v", this.Server), "PID", "actor.PID", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *AgentForwardToSvr) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&AgentForwardToSvr{`,
		`ServerName:` + fmt.Sprintf("%v", this.ServerName) + `,`,
		`Uid:` + fmt.Sprintf("%v", this.Uid) + `,`,
		`Body:` + fmt.Sprintf("%v", this.Body) + `,`,
		`}`,
	}, "")
	return s
}
func (this *AgentForwardToCli) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&AgentForwardToCli{`,
		`Uid:` + fmt.Sprintf("%v", this.Uid) + `,`,
		`Body:` + fmt.Sprintf("%v", this.Body) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringAgent(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *AgentConnect) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAgent
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: AgentConnect: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AgentConnect: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sender", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAgent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthAgent
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Sender == nil {
				m.Sender = &actor.PID{}
			}
			if err := m.Sender.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ServerName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAgent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthAgent
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ServerName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAgent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthAgent
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *AgentConnected) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAgent
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: AgentConnected: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AgentConnected: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Message", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAgent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthAgent
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Message = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Server", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAgent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthAgent
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Server == nil {
				m.Server = &actor.PID{}
			}
			if err := m.Server.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAgent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthAgent
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *AgentForwardToSvr) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAgent
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: AgentForwardToSvr: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AgentForwardToSvr: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ServerName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAgent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthAgent
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ServerName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Uid", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAgent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthAgent
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Uid = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Body", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAgent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthAgent
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Body = append(m.Body[:0], dAtA[iNdEx:postIndex]...)
			if m.Body == nil {
				m.Body = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAgent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthAgent
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *AgentForwardToCli) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAgent
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: AgentForwardToCli: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AgentForwardToCli: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Uid", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAgent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthAgent
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Uid = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Body", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAgent
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthAgent
			}
			postIndex := iNdEx + byteLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Body = append(m.Body[:0], dAtA[iNdEx:postIndex]...)
			if m.Body == nil {
				m.Body = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAgent(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthAgent
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipAgent(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowAgent
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowAgent
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowAgent
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthAgent
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowAgent
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipAgent(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthAgent = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowAgent   = fmt.Errorf("proto: integer overflow")
)

func init() { proto.RegisterFile("agent.proto", fileDescriptorAgent) }

var fileDescriptorAgent = []byte{
	// 305 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x90, 0x3f, 0x4b, 0xc3, 0x40,
	0x18, 0xc6, 0xf3, 0xb6, 0xa5, 0xd2, 0x6b, 0x11, 0xbd, 0x29, 0x38, 0xbc, 0x94, 0x4c, 0x19, 0x34,
	0x01, 0x05, 0xc1, 0xb1, 0xad, 0x08, 0x1d, 0x2c, 0x92, 0xd4, 0xc1, 0x31, 0x7f, 0x8e, 0x18, 0xb4,
	0x77, 0xe5, 0x12, 0x2b, 0xdd, 0xfc, 0x08, 0x7e, 0x0c, 0x3f, 0x8a, 0x63, 0x47, 0x47, 0x73, 0x2e,
	0x8e, 0xfd, 0x08, 0x92, 0x4b, 0x8b, 0x25, 0xdd, 0x9e, 0xbc, 0x6f, 0x7e, 0xbf, 0xf7, 0xe1, 0x48,
	0x37, 0x48, 0x18, 0xcf, 0x9d, 0xb9, 0x14, 0xb9, 0xa0, 0xad, 0x79, 0x98, 0x88, 0x93, 0xcb, 0x24,
	0xcd, 0x1f, 0x5f, 0x42, 0x27, 0x12, 0x33, 0x77, 0x90, 0x2d, 0xf9, 0x93, 0x14, 0x7c, 0x3c, 0x75,
	0xf5, 0x2f, 0x41, 0x94, 0x0b, 0x79, 0x96, 0x08, 0x57, 0x87, 0x6a, 0x96, 0x55, 0xb4, 0xe5, 0x91,
	0xde, 0xa0, 0x94, 0x8d, 0x04, 0xe7, 0x2c, 0xca, 0xa9, 0x45, 0xda, 0x3e, 0xe3, 0x31, 0x93, 0x26,
	0xf4, 0xc1, 0xee, 0x9e, 0x13, 0x47, 0x43, 0xce, 0xdd, 0xf8, 0xda, 0xdb, 0x6c, 0x28, 0x12, 0xe2,
	0x33, 0xb9, 0x60, 0x72, 0x12, 0xcc, 0x98, 0xd9, 0xe8, 0x83, 0xdd, 0xf1, 0x76, 0x26, 0xd6, 0x84,
	0x1c, 0xee, 0x3a, 0x59, 0x4c, 0x4d, 0x72, 0x70, 0xcb, 0xb2, 0x2c, 0x48, 0x98, 0xd6, 0x76, 0xbc,
	0xed, 0x67, 0x75, 0xaf, 0x24, 0xb5, 0x67, 0xef, 0x5e, 0xb9, 0xb1, 0x1e, 0xc8, 0xb1, 0xf6, 0xdd,
	0x08, 0xf9, 0x1a, 0xc8, 0x78, 0x2a, 0xfc, 0x45, 0xbd, 0x04, 0xd4, 0x4b, 0xd0, 0x23, 0xd2, 0xbc,
	0x4f, 0xe3, 0x4d, 0xbb, 0x32, 0x52, 0x4a, 0x5a, 0x43, 0x11, 0x2f, 0xcd, 0x66, 0x1f, 0xec, 0x9e,
	0xa7, 0xb3, 0x75, 0x55, 0x57, 0x8f, 0x9e, 0xd3, 0x2d, 0x0a, 0xfb, 0x68, 0xe3, 0x1f, 0x1d, 0x9e,
	0xae, 0x0a, 0x34, 0xbe, 0x0a, 0x34, 0xd6, 0x05, 0xc2, 0x9b, 0x42, 0xf8, 0x50, 0x08, 0x9f, 0x0a,
	0x61, 0xa5, 0x10, 0xbe, 0x15, 0xc2, 0xaf, 0x42, 0x63, 0xad, 0x10, 0xde, 0x7f, 0xd0, 0x08, 0xdb,
	0xfa, 0xb9, 0x2f, 0xfe, 0x02, 0x00, 0x00, 0xff, 0xff, 0x6b, 0x4b, 0xad, 0x57, 0xbb, 0x01, 0x00,
	0x00,
}
