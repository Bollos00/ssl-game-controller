// Code generated by protoc-gen-go. DO NOT EDIT.
// source: ssl_gc_engine.proto

package engine

import (
	fmt "fmt"
	geom "github.com/RoboCup-SSL/ssl-game-controller/internal/app/geom"
	state "github.com/RoboCup-SSL/ssl-game-controller/internal/app/state"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// The GC state contains settings and state independent of the match state
type GcState struct {
	// The state of each team
	TeamState map[string]*GcStateTeam `protobuf:"bytes,1,rep,name=team_state,json=teamState" json:"team_state,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	// the states of the auto referees
	AutoRefState map[string]*GcStateAutoRef `protobuf:"bytes,2,rep,name=auto_ref_state,json=autoRefState" json:"auto_ref_state,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	// the states of the attached trackers
	TrackerState map[string]*GcStateTracker `protobuf:"bytes,3,rep,name=tracker_state,json=trackerState" json:"tracker_state,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	// the state of the currently selected tracker
	TrackerStateGc *GcStateTracker `protobuf:"bytes,4,opt,name=tracker_state_gc,json=trackerStateGc" json:"tracker_state_gc,omitempty"`
	// can the match be continued right now?
	ReadyToContinue      *bool    `protobuf:"varint,5,opt,name=ready_to_continue,json=readyToContinue" json:"ready_to_continue,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GcState) Reset()         { *m = GcState{} }
func (m *GcState) String() string { return proto.CompactTextString(m) }
func (*GcState) ProtoMessage()    {}
func (*GcState) Descriptor() ([]byte, []int) {
	return fileDescriptor_16bd66f36a99a3d1, []int{0}
}

func (m *GcState) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GcState.Unmarshal(m, b)
}
func (m *GcState) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GcState.Marshal(b, m, deterministic)
}
func (m *GcState) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GcState.Merge(m, src)
}
func (m *GcState) XXX_Size() int {
	return xxx_messageInfo_GcState.Size(m)
}
func (m *GcState) XXX_DiscardUnknown() {
	xxx_messageInfo_GcState.DiscardUnknown(m)
}

var xxx_messageInfo_GcState proto.InternalMessageInfo

func (m *GcState) GetTeamState() map[string]*GcStateTeam {
	if m != nil {
		return m.TeamState
	}
	return nil
}

func (m *GcState) GetAutoRefState() map[string]*GcStateAutoRef {
	if m != nil {
		return m.AutoRefState
	}
	return nil
}

func (m *GcState) GetTrackerState() map[string]*GcStateTracker {
	if m != nil {
		return m.TrackerState
	}
	return nil
}

func (m *GcState) GetTrackerStateGc() *GcStateTracker {
	if m != nil {
		return m.TrackerStateGc
	}
	return nil
}

func (m *GcState) GetReadyToContinue() bool {
	if m != nil && m.ReadyToContinue != nil {
		return *m.ReadyToContinue
	}
	return false
}

// The GC state for a singl eteam
type GcStateTeam struct {
	// true: The team is connected
	Connected *bool `protobuf:"varint,1,opt,name=connected" json:"connected,omitempty"`
	// true: The team connected via TLS with a verified certificate
	ConnectionVerified *bool `protobuf:"varint,2,opt,name=connection_verified,json=connectionVerified" json:"connection_verified,omitempty"`
	// true: The remote control for the team is connected
	RemoteControlConnected *bool `protobuf:"varint,3,opt,name=remote_control_connected,json=remoteControlConnected" json:"remote_control_connected,omitempty"`
	// true: The remote control for the team connected via TLS with a verified certificate
	RemoteControlConnectionVerified *bool    `protobuf:"varint,4,opt,name=remote_control_connection_verified,json=remoteControlConnectionVerified" json:"remote_control_connection_verified,omitempty"`
	XXX_NoUnkeyedLiteral            struct{} `json:"-"`
	XXX_unrecognized                []byte   `json:"-"`
	XXX_sizecache                   int32    `json:"-"`
}

func (m *GcStateTeam) Reset()         { *m = GcStateTeam{} }
func (m *GcStateTeam) String() string { return proto.CompactTextString(m) }
func (*GcStateTeam) ProtoMessage()    {}
func (*GcStateTeam) Descriptor() ([]byte, []int) {
	return fileDescriptor_16bd66f36a99a3d1, []int{1}
}

func (m *GcStateTeam) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GcStateTeam.Unmarshal(m, b)
}
func (m *GcStateTeam) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GcStateTeam.Marshal(b, m, deterministic)
}
func (m *GcStateTeam) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GcStateTeam.Merge(m, src)
}
func (m *GcStateTeam) XXX_Size() int {
	return xxx_messageInfo_GcStateTeam.Size(m)
}
func (m *GcStateTeam) XXX_DiscardUnknown() {
	xxx_messageInfo_GcStateTeam.DiscardUnknown(m)
}

var xxx_messageInfo_GcStateTeam proto.InternalMessageInfo

func (m *GcStateTeam) GetConnected() bool {
	if m != nil && m.Connected != nil {
		return *m.Connected
	}
	return false
}

func (m *GcStateTeam) GetConnectionVerified() bool {
	if m != nil && m.ConnectionVerified != nil {
		return *m.ConnectionVerified
	}
	return false
}

func (m *GcStateTeam) GetRemoteControlConnected() bool {
	if m != nil && m.RemoteControlConnected != nil {
		return *m.RemoteControlConnected
	}
	return false
}

func (m *GcStateTeam) GetRemoteControlConnectionVerified() bool {
	if m != nil && m.RemoteControlConnectionVerified != nil {
		return *m.RemoteControlConnectionVerified
	}
	return false
}

// The GC state of an auto referee
type GcStateAutoRef struct {
	// true: The autoRef connected via TLS with a verified certificate
	ConnectionVerified   *bool    `protobuf:"varint,1,opt,name=connection_verified,json=connectionVerified" json:"connection_verified,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GcStateAutoRef) Reset()         { *m = GcStateAutoRef{} }
func (m *GcStateAutoRef) String() string { return proto.CompactTextString(m) }
func (*GcStateAutoRef) ProtoMessage()    {}
func (*GcStateAutoRef) Descriptor() ([]byte, []int) {
	return fileDescriptor_16bd66f36a99a3d1, []int{2}
}

func (m *GcStateAutoRef) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GcStateAutoRef.Unmarshal(m, b)
}
func (m *GcStateAutoRef) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GcStateAutoRef.Marshal(b, m, deterministic)
}
func (m *GcStateAutoRef) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GcStateAutoRef.Merge(m, src)
}
func (m *GcStateAutoRef) XXX_Size() int {
	return xxx_messageInfo_GcStateAutoRef.Size(m)
}
func (m *GcStateAutoRef) XXX_DiscardUnknown() {
	xxx_messageInfo_GcStateAutoRef.DiscardUnknown(m)
}

var xxx_messageInfo_GcStateAutoRef proto.InternalMessageInfo

func (m *GcStateAutoRef) GetConnectionVerified() bool {
	if m != nil && m.ConnectionVerified != nil {
		return *m.ConnectionVerified
	}
	return false
}

// GC state of a tracker
type GcStateTracker struct {
	// Name of the source
	SourceName *string `protobuf:"bytes,1,opt,name=source_name,json=sourceName" json:"source_name,omitempty"`
	// UUID of the source
	Uuid *string `protobuf:"bytes,4,opt,name=uuid" json:"uuid,omitempty"`
	// Current ball
	Ball *Ball `protobuf:"bytes,2,opt,name=ball" json:"ball,omitempty"`
	// Current robots
	Robots               []*Robot `protobuf:"bytes,3,rep,name=robots" json:"robots,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GcStateTracker) Reset()         { *m = GcStateTracker{} }
func (m *GcStateTracker) String() string { return proto.CompactTextString(m) }
func (*GcStateTracker) ProtoMessage()    {}
func (*GcStateTracker) Descriptor() ([]byte, []int) {
	return fileDescriptor_16bd66f36a99a3d1, []int{3}
}

func (m *GcStateTracker) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GcStateTracker.Unmarshal(m, b)
}
func (m *GcStateTracker) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GcStateTracker.Marshal(b, m, deterministic)
}
func (m *GcStateTracker) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GcStateTracker.Merge(m, src)
}
func (m *GcStateTracker) XXX_Size() int {
	return xxx_messageInfo_GcStateTracker.Size(m)
}
func (m *GcStateTracker) XXX_DiscardUnknown() {
	xxx_messageInfo_GcStateTracker.DiscardUnknown(m)
}

var xxx_messageInfo_GcStateTracker proto.InternalMessageInfo

func (m *GcStateTracker) GetSourceName() string {
	if m != nil && m.SourceName != nil {
		return *m.SourceName
	}
	return ""
}

func (m *GcStateTracker) GetUuid() string {
	if m != nil && m.Uuid != nil {
		return *m.Uuid
	}
	return ""
}

func (m *GcStateTracker) GetBall() *Ball {
	if m != nil {
		return m.Ball
	}
	return nil
}

func (m *GcStateTracker) GetRobots() []*Robot {
	if m != nil {
		return m.Robots
	}
	return nil
}

// The ball state
type Ball struct {
	// ball position [m]
	Pos *geom.Vector3 `protobuf:"bytes,1,opt,name=pos" json:"pos,omitempty"`
	// ball velocity [m/s]
	Vel                  *geom.Vector3 `protobuf:"bytes,2,opt,name=vel" json:"vel,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *Ball) Reset()         { *m = Ball{} }
func (m *Ball) String() string { return proto.CompactTextString(m) }
func (*Ball) ProtoMessage()    {}
func (*Ball) Descriptor() ([]byte, []int) {
	return fileDescriptor_16bd66f36a99a3d1, []int{4}
}

func (m *Ball) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Ball.Unmarshal(m, b)
}
func (m *Ball) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Ball.Marshal(b, m, deterministic)
}
func (m *Ball) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Ball.Merge(m, src)
}
func (m *Ball) XXX_Size() int {
	return xxx_messageInfo_Ball.Size(m)
}
func (m *Ball) XXX_DiscardUnknown() {
	xxx_messageInfo_Ball.DiscardUnknown(m)
}

var xxx_messageInfo_Ball proto.InternalMessageInfo

func (m *Ball) GetPos() *geom.Vector3 {
	if m != nil {
		return m.Pos
	}
	return nil
}

func (m *Ball) GetVel() *geom.Vector3 {
	if m != nil {
		return m.Vel
	}
	return nil
}

// The robot state
type Robot struct {
	// robot id and team
	Id *state.RobotId `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	// robot position [m]
	Pos                  *geom.Vector2 `protobuf:"bytes,2,opt,name=pos" json:"pos,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *Robot) Reset()         { *m = Robot{} }
func (m *Robot) String() string { return proto.CompactTextString(m) }
func (*Robot) ProtoMessage()    {}
func (*Robot) Descriptor() ([]byte, []int) {
	return fileDescriptor_16bd66f36a99a3d1, []int{5}
}

func (m *Robot) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Robot.Unmarshal(m, b)
}
func (m *Robot) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Robot.Marshal(b, m, deterministic)
}
func (m *Robot) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Robot.Merge(m, src)
}
func (m *Robot) XXX_Size() int {
	return xxx_messageInfo_Robot.Size(m)
}
func (m *Robot) XXX_DiscardUnknown() {
	xxx_messageInfo_Robot.DiscardUnknown(m)
}

var xxx_messageInfo_Robot proto.InternalMessageInfo

func (m *Robot) GetId() *state.RobotId {
	if m != nil {
		return m.Id
	}
	return nil
}

func (m *Robot) GetPos() *geom.Vector2 {
	if m != nil {
		return m.Pos
	}
	return nil
}

func init() {
	proto.RegisterType((*GcState)(nil), "GcState")
	proto.RegisterMapType((map[string]*GcStateAutoRef)(nil), "GcState.AutoRefStateEntry")
	proto.RegisterMapType((map[string]*GcStateTeam)(nil), "GcState.TeamStateEntry")
	proto.RegisterMapType((map[string]*GcStateTracker)(nil), "GcState.TrackerStateEntry")
	proto.RegisterType((*GcStateTeam)(nil), "GcStateTeam")
	proto.RegisterType((*GcStateAutoRef)(nil), "GcStateAutoRef")
	proto.RegisterType((*GcStateTracker)(nil), "GcStateTracker")
	proto.RegisterType((*Ball)(nil), "Ball")
	proto.RegisterType((*Robot)(nil), "Robot")
}

func init() {
	proto.RegisterFile("ssl_gc_engine.proto", fileDescriptor_16bd66f36a99a3d1)
}

var fileDescriptor_16bd66f36a99a3d1 = []byte{
	// 577 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x94, 0xdf, 0x6a, 0xdb, 0x30,
	0x18, 0xc5, 0x71, 0xfe, 0x74, 0xc9, 0x97, 0x2e, 0x6d, 0x55, 0xb6, 0x65, 0x61, 0xac, 0xc1, 0x30,
	0x08, 0x83, 0x3a, 0x90, 0xc1, 0xe8, 0x06, 0xeb, 0xda, 0x86, 0x51, 0xf6, 0x87, 0x31, 0xd4, 0xd2,
	0x8b, 0xdd, 0x18, 0x55, 0xf9, 0xea, 0x99, 0xda, 0x52, 0x90, 0xe5, 0x40, 0xee, 0xf6, 0x1e, 0x7b,
	0xbc, 0xbd, 0xc8, 0xb0, 0x2c, 0x3b, 0x76, 0x9b, 0x5e, 0xec, 0x4e, 0x3a, 0xc7, 0xe7, 0xf7, 0x49,
	0x47, 0x24, 0xb0, 0x9f, 0x24, 0x91, 0x1f, 0x70, 0x1f, 0x45, 0x10, 0x0a, 0xf4, 0x16, 0x4a, 0x6a,
	0x39, 0x7c, 0x62, 0xc5, 0x00, 0x65, 0x8c, 0x5a, 0xad, 0xac, 0x5c, 0x7c, 0xcb, 0x65, 0x1c, 0x4b,
	0x91, 0x8b, 0xee, 0x9f, 0x16, 0x3c, 0x3a, 0xe7, 0x17, 0x9a, 0x69, 0x24, 0x6f, 0x01, 0x34, 0xb2,
	0xd8, 0x4f, 0xb2, 0xdd, 0xc0, 0x19, 0x35, 0xc7, 0xbd, 0xe9, 0x33, 0xcf, 0xba, 0xde, 0x25, 0xb2,
	0xd8, 0xac, 0x3e, 0x09, 0xad, 0x56, 0xb4, 0xab, 0x8b, 0x3d, 0x39, 0x81, 0x3e, 0x4b, 0xb5, 0xf4,
	0x15, 0xde, 0xd8, 0x6c, 0xc3, 0x64, 0x87, 0x65, 0xf6, 0x34, 0xd5, 0x92, 0xe2, 0x4d, 0x25, 0xbe,
	0xcd, 0x2a, 0x12, 0xf9, 0x08, 0x8f, 0xb5, 0x62, 0xfc, 0x16, 0x95, 0x05, 0x34, 0xef, 0x00, 0x2e,
	0x73, 0xb7, 0x0a, 0xd0, 0x15, 0x89, 0xbc, 0x83, 0xdd, 0x1a, 0xc0, 0x0f, 0xf8, 0xa0, 0x35, 0x72,
	0xc6, 0xbd, 0xe9, 0x4e, 0xc1, 0xb0, 0x08, 0xda, 0xaf, 0x06, 0xcf, 0x39, 0x79, 0x0d, 0x7b, 0x0a,
	0xd9, 0x7c, 0xe5, 0x6b, 0xe9, 0x73, 0x29, 0x74, 0x28, 0x52, 0x1c, 0xb4, 0x47, 0xce, 0xb8, 0x43,
	0x77, 0x8c, 0x71, 0x29, 0x67, 0x56, 0x1e, 0x7e, 0x81, 0x7e, 0xbd, 0x06, 0xb2, 0x0b, 0xcd, 0x5b,
	0x5c, 0x0d, 0x9c, 0x91, 0x33, 0xee, 0xd2, 0x6c, 0x49, 0x5c, 0x68, 0x2f, 0x59, 0x94, 0x66, 0x25,
	0x64, 0xf3, 0xb7, 0xcb, 0xf9, 0xc8, 0x62, 0x9a, 0x5b, 0xef, 0x1b, 0x47, 0xce, 0xf0, 0x07, 0xec,
	0xdd, 0xab, 0x65, 0x03, 0xee, 0x55, 0x1d, 0x57, 0x5e, 0xc7, 0x66, 0xef, 0x10, 0xef, 0xf5, 0xf4,
	0x1f, 0xc4, 0xa2, 0xa0, 0x35, 0xd1, 0xfd, 0xeb, 0x40, 0xaf, 0x72, 0x7c, 0xf2, 0x02, 0xba, 0x5c,
	0x0a, 0x81, 0x5c, 0xe3, 0xdc, 0x20, 0x3b, 0x74, 0x2d, 0x90, 0x09, 0xec, 0xdb, 0x4d, 0x28, 0x85,
	0xbf, 0x44, 0x15, 0xde, 0x84, 0x38, 0x37, 0x63, 0x3a, 0x94, 0xac, 0xad, 0x2b, 0xeb, 0x90, 0x23,
	0x18, 0x28, 0x8c, 0xa5, 0x46, 0x53, 0xbc, 0x92, 0x91, 0xbf, 0xa6, 0x37, 0x4d, 0xea, 0x69, 0xee,
	0xcf, 0x72, 0x7b, 0x56, 0x8e, 0xfa, 0x0a, 0xee, 0xe6, 0x64, 0x6d, 0x72, 0xcb, 0x30, 0x0e, 0x36,
	0x31, 0x2a, 0xc7, 0x70, 0x4f, 0xa1, 0x5f, 0x2f, 0xf5, 0xa1, 0x9b, 0x38, 0x0f, 0xdd, 0xc4, 0xfd,
	0xed, 0x94, 0x0c, 0x5b, 0x23, 0x39, 0x80, 0x5e, 0x22, 0x53, 0xc5, 0xd1, 0x17, 0x2c, 0x46, 0xfb,
	0x00, 0x90, 0x4b, 0xdf, 0x59, 0x8c, 0x84, 0x40, 0x2b, 0x4d, 0xc3, 0xfc, 0x94, 0x5d, 0x6a, 0xd6,
	0xe4, 0x39, 0xb4, 0xae, 0x59, 0x14, 0xd9, 0xa7, 0x69, 0x7b, 0x67, 0x2c, 0x8a, 0xa8, 0x91, 0xc8,
	0x4b, 0xd8, 0x52, 0xf2, 0x5a, 0xea, 0xc4, 0xfe, 0x38, 0xb6, 0x3c, 0x9a, 0x6d, 0xa9, 0x55, 0xdd,
	0x63, 0x68, 0x65, 0x5f, 0x93, 0x21, 0x34, 0x17, 0x32, 0x31, 0xf3, 0x7a, 0xd3, 0x8e, 0x77, 0x85,
	0x5c, 0x4b, 0xf5, 0x86, 0x66, 0x62, 0xe6, 0x2d, 0xb1, 0xa0, 0x57, 0xbc, 0x25, 0x46, 0xee, 0x07,
	0x68, 0x1b, 0x20, 0x19, 0x40, 0x23, 0x9c, 0x97, 0x79, 0xa3, 0x7d, 0x9e, 0xd3, 0x46, 0x38, 0x2f,
	0xd0, 0xf5, 0xf8, 0xd4, 0xa0, 0xcf, 0x4e, 0x7e, 0x1e, 0x07, 0xa1, 0xfe, 0x95, 0x5e, 0x7b, 0x5c,
	0xc6, 0x93, 0x2c, 0x35, 0x4b, 0x17, 0x87, 0x17, 0x17, 0xdf, 0x26, 0x49, 0x12, 0x1d, 0x06, 0x2c,
	0xc6, 0x43, 0xfb, 0x54, 0x11, 0xaa, 0x49, 0x28, 0x34, 0x2a, 0xc1, 0xa2, 0x09, 0x5b, 0x2c, 0x26,
	0xf9, 0x9f, 0xd7, 0xbf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xb4, 0x90, 0x94, 0x34, 0xcc, 0x04, 0x00,
	0x00,
}
