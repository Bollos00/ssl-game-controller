// Code generated by protoc-gen-go. DO NOT EDIT.
// source: ssl_vision_geometry.proto

package vision

import (
	fmt "fmt"
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

// A 2D float vector.
type Vector2F struct {
	X                    *float32 `protobuf:"fixed32,1,req,name=x" json:"x,omitempty"`
	Y                    *float32 `protobuf:"fixed32,2,req,name=y" json:"y,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Vector2F) Reset()         { *m = Vector2F{} }
func (m *Vector2F) String() string { return proto.CompactTextString(m) }
func (*Vector2F) ProtoMessage()    {}
func (*Vector2F) Descriptor() ([]byte, []int) {
	return fileDescriptor_3d55d410315aeff8, []int{0}
}

func (m *Vector2F) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Vector2F.Unmarshal(m, b)
}
func (m *Vector2F) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Vector2F.Marshal(b, m, deterministic)
}
func (m *Vector2F) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Vector2F.Merge(m, src)
}
func (m *Vector2F) XXX_Size() int {
	return xxx_messageInfo_Vector2F.Size(m)
}
func (m *Vector2F) XXX_DiscardUnknown() {
	xxx_messageInfo_Vector2F.DiscardUnknown(m)
}

var xxx_messageInfo_Vector2F proto.InternalMessageInfo

func (m *Vector2F) GetX() float32 {
	if m != nil && m.X != nil {
		return *m.X
	}
	return 0
}

func (m *Vector2F) GetY() float32 {
	if m != nil && m.Y != nil {
		return *m.Y
	}
	return 0
}

// Represents a field marking as a line segment represented by a start point p1,
// and end point p2, and a line thickness. The start and end points are along
// the center of the line, so the thickness of the line extends by thickness / 2
// on either side of the line.
type SSL_FieldLineSegment struct {
	// Name of this field marking.
	Name *string `protobuf:"bytes,1,req,name=name" json:"name,omitempty"`
	// Start point of the line segment.
	P1 *Vector2F `protobuf:"bytes,2,req,name=p1" json:"p1,omitempty"`
	// End point of the line segment.
	P2 *Vector2F `protobuf:"bytes,3,req,name=p2" json:"p2,omitempty"`
	// Thickness of the line segment.
	Thickness            *float32 `protobuf:"fixed32,4,req,name=thickness" json:"thickness,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SSL_FieldLineSegment) Reset()         { *m = SSL_FieldLineSegment{} }
func (m *SSL_FieldLineSegment) String() string { return proto.CompactTextString(m) }
func (*SSL_FieldLineSegment) ProtoMessage()    {}
func (*SSL_FieldLineSegment) Descriptor() ([]byte, []int) {
	return fileDescriptor_3d55d410315aeff8, []int{1}
}

func (m *SSL_FieldLineSegment) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SSL_FieldLineSegment.Unmarshal(m, b)
}
func (m *SSL_FieldLineSegment) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SSL_FieldLineSegment.Marshal(b, m, deterministic)
}
func (m *SSL_FieldLineSegment) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SSL_FieldLineSegment.Merge(m, src)
}
func (m *SSL_FieldLineSegment) XXX_Size() int {
	return xxx_messageInfo_SSL_FieldLineSegment.Size(m)
}
func (m *SSL_FieldLineSegment) XXX_DiscardUnknown() {
	xxx_messageInfo_SSL_FieldLineSegment.DiscardUnknown(m)
}

var xxx_messageInfo_SSL_FieldLineSegment proto.InternalMessageInfo

func (m *SSL_FieldLineSegment) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *SSL_FieldLineSegment) GetP1() *Vector2F {
	if m != nil {
		return m.P1
	}
	return nil
}

func (m *SSL_FieldLineSegment) GetP2() *Vector2F {
	if m != nil {
		return m.P2
	}
	return nil
}

func (m *SSL_FieldLineSegment) GetThickness() float32 {
	if m != nil && m.Thickness != nil {
		return *m.Thickness
	}
	return 0
}

// Represents a field marking as a circular arc segment represented by center point, a
// start angle, an end angle, and an arc thickness.
type SSL_FieldCicularArc struct {
	// Name of this field marking.
	Name *string `protobuf:"bytes,1,req,name=name" json:"name,omitempty"`
	// Center point of the circular arc.
	Center *Vector2F `protobuf:"bytes,2,req,name=center" json:"center,omitempty"`
	// Radius of the arc.
	Radius *float32 `protobuf:"fixed32,3,req,name=radius" json:"radius,omitempty"`
	// Start angle in counter-clockwise order.
	A1 *float32 `protobuf:"fixed32,4,req,name=a1" json:"a1,omitempty"`
	// End angle in counter-clockwise order.
	A2 *float32 `protobuf:"fixed32,5,req,name=a2" json:"a2,omitempty"`
	// Thickness of the arc.
	Thickness            *float32 `protobuf:"fixed32,6,req,name=thickness" json:"thickness,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SSL_FieldCicularArc) Reset()         { *m = SSL_FieldCicularArc{} }
func (m *SSL_FieldCicularArc) String() string { return proto.CompactTextString(m) }
func (*SSL_FieldCicularArc) ProtoMessage()    {}
func (*SSL_FieldCicularArc) Descriptor() ([]byte, []int) {
	return fileDescriptor_3d55d410315aeff8, []int{2}
}

func (m *SSL_FieldCicularArc) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SSL_FieldCicularArc.Unmarshal(m, b)
}
func (m *SSL_FieldCicularArc) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SSL_FieldCicularArc.Marshal(b, m, deterministic)
}
func (m *SSL_FieldCicularArc) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SSL_FieldCicularArc.Merge(m, src)
}
func (m *SSL_FieldCicularArc) XXX_Size() int {
	return xxx_messageInfo_SSL_FieldCicularArc.Size(m)
}
func (m *SSL_FieldCicularArc) XXX_DiscardUnknown() {
	xxx_messageInfo_SSL_FieldCicularArc.DiscardUnknown(m)
}

var xxx_messageInfo_SSL_FieldCicularArc proto.InternalMessageInfo

func (m *SSL_FieldCicularArc) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *SSL_FieldCicularArc) GetCenter() *Vector2F {
	if m != nil {
		return m.Center
	}
	return nil
}

func (m *SSL_FieldCicularArc) GetRadius() float32 {
	if m != nil && m.Radius != nil {
		return *m.Radius
	}
	return 0
}

func (m *SSL_FieldCicularArc) GetA1() float32 {
	if m != nil && m.A1 != nil {
		return *m.A1
	}
	return 0
}

func (m *SSL_FieldCicularArc) GetA2() float32 {
	if m != nil && m.A2 != nil {
		return *m.A2
	}
	return 0
}

func (m *SSL_FieldCicularArc) GetThickness() float32 {
	if m != nil && m.Thickness != nil {
		return *m.Thickness
	}
	return 0
}

type SSL_GeometryFieldSize struct {
	FieldLength          *int32                  `protobuf:"varint,1,req,name=field_length,json=fieldLength" json:"field_length,omitempty"`
	FieldWidth           *int32                  `protobuf:"varint,2,req,name=field_width,json=fieldWidth" json:"field_width,omitempty"`
	GoalWidth            *int32                  `protobuf:"varint,3,req,name=goal_width,json=goalWidth" json:"goal_width,omitempty"`
	GoalDepth            *int32                  `protobuf:"varint,4,req,name=goal_depth,json=goalDepth" json:"goal_depth,omitempty"`
	BoundaryWidth        *int32                  `protobuf:"varint,5,req,name=boundary_width,json=boundaryWidth" json:"boundary_width,omitempty"`
	FieldLines           []*SSL_FieldLineSegment `protobuf:"bytes,6,rep,name=field_lines,json=fieldLines" json:"field_lines,omitempty"`
	FieldArcs            []*SSL_FieldCicularArc  `protobuf:"bytes,7,rep,name=field_arcs,json=fieldArcs" json:"field_arcs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *SSL_GeometryFieldSize) Reset()         { *m = SSL_GeometryFieldSize{} }
func (m *SSL_GeometryFieldSize) String() string { return proto.CompactTextString(m) }
func (*SSL_GeometryFieldSize) ProtoMessage()    {}
func (*SSL_GeometryFieldSize) Descriptor() ([]byte, []int) {
	return fileDescriptor_3d55d410315aeff8, []int{3}
}

func (m *SSL_GeometryFieldSize) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SSL_GeometryFieldSize.Unmarshal(m, b)
}
func (m *SSL_GeometryFieldSize) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SSL_GeometryFieldSize.Marshal(b, m, deterministic)
}
func (m *SSL_GeometryFieldSize) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SSL_GeometryFieldSize.Merge(m, src)
}
func (m *SSL_GeometryFieldSize) XXX_Size() int {
	return xxx_messageInfo_SSL_GeometryFieldSize.Size(m)
}
func (m *SSL_GeometryFieldSize) XXX_DiscardUnknown() {
	xxx_messageInfo_SSL_GeometryFieldSize.DiscardUnknown(m)
}

var xxx_messageInfo_SSL_GeometryFieldSize proto.InternalMessageInfo

func (m *SSL_GeometryFieldSize) GetFieldLength() int32 {
	if m != nil && m.FieldLength != nil {
		return *m.FieldLength
	}
	return 0
}

func (m *SSL_GeometryFieldSize) GetFieldWidth() int32 {
	if m != nil && m.FieldWidth != nil {
		return *m.FieldWidth
	}
	return 0
}

func (m *SSL_GeometryFieldSize) GetGoalWidth() int32 {
	if m != nil && m.GoalWidth != nil {
		return *m.GoalWidth
	}
	return 0
}

func (m *SSL_GeometryFieldSize) GetGoalDepth() int32 {
	if m != nil && m.GoalDepth != nil {
		return *m.GoalDepth
	}
	return 0
}

func (m *SSL_GeometryFieldSize) GetBoundaryWidth() int32 {
	if m != nil && m.BoundaryWidth != nil {
		return *m.BoundaryWidth
	}
	return 0
}

func (m *SSL_GeometryFieldSize) GetFieldLines() []*SSL_FieldLineSegment {
	if m != nil {
		return m.FieldLines
	}
	return nil
}

func (m *SSL_GeometryFieldSize) GetFieldArcs() []*SSL_FieldCicularArc {
	if m != nil {
		return m.FieldArcs
	}
	return nil
}

type SSL_GeometryCameraCalibration struct {
	CameraId             *uint32  `protobuf:"varint,1,req,name=camera_id,json=cameraId" json:"camera_id,omitempty"`
	FocalLength          *float32 `protobuf:"fixed32,2,req,name=focal_length,json=focalLength" json:"focal_length,omitempty"`
	PrincipalPointX      *float32 `protobuf:"fixed32,3,req,name=principal_point_x,json=principalPointX" json:"principal_point_x,omitempty"`
	PrincipalPointY      *float32 `protobuf:"fixed32,4,req,name=principal_point_y,json=principalPointY" json:"principal_point_y,omitempty"`
	Distortion           *float32 `protobuf:"fixed32,5,req,name=distortion" json:"distortion,omitempty"`
	Q0                   *float32 `protobuf:"fixed32,6,req,name=q0" json:"q0,omitempty"`
	Q1                   *float32 `protobuf:"fixed32,7,req,name=q1" json:"q1,omitempty"`
	Q2                   *float32 `protobuf:"fixed32,8,req,name=q2" json:"q2,omitempty"`
	Q3                   *float32 `protobuf:"fixed32,9,req,name=q3" json:"q3,omitempty"`
	Tx                   *float32 `protobuf:"fixed32,10,req,name=tx" json:"tx,omitempty"`
	Ty                   *float32 `protobuf:"fixed32,11,req,name=ty" json:"ty,omitempty"`
	Tz                   *float32 `protobuf:"fixed32,12,req,name=tz" json:"tz,omitempty"`
	DerivedCameraWorldTx *float32 `protobuf:"fixed32,13,opt,name=derived_camera_world_tx,json=derivedCameraWorldTx" json:"derived_camera_world_tx,omitempty"`
	DerivedCameraWorldTy *float32 `protobuf:"fixed32,14,opt,name=derived_camera_world_ty,json=derivedCameraWorldTy" json:"derived_camera_world_ty,omitempty"`
	DerivedCameraWorldTz *float32 `protobuf:"fixed32,15,opt,name=derived_camera_world_tz,json=derivedCameraWorldTz" json:"derived_camera_world_tz,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SSL_GeometryCameraCalibration) Reset()         { *m = SSL_GeometryCameraCalibration{} }
func (m *SSL_GeometryCameraCalibration) String() string { return proto.CompactTextString(m) }
func (*SSL_GeometryCameraCalibration) ProtoMessage()    {}
func (*SSL_GeometryCameraCalibration) Descriptor() ([]byte, []int) {
	return fileDescriptor_3d55d410315aeff8, []int{4}
}

func (m *SSL_GeometryCameraCalibration) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SSL_GeometryCameraCalibration.Unmarshal(m, b)
}
func (m *SSL_GeometryCameraCalibration) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SSL_GeometryCameraCalibration.Marshal(b, m, deterministic)
}
func (m *SSL_GeometryCameraCalibration) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SSL_GeometryCameraCalibration.Merge(m, src)
}
func (m *SSL_GeometryCameraCalibration) XXX_Size() int {
	return xxx_messageInfo_SSL_GeometryCameraCalibration.Size(m)
}
func (m *SSL_GeometryCameraCalibration) XXX_DiscardUnknown() {
	xxx_messageInfo_SSL_GeometryCameraCalibration.DiscardUnknown(m)
}

var xxx_messageInfo_SSL_GeometryCameraCalibration proto.InternalMessageInfo

func (m *SSL_GeometryCameraCalibration) GetCameraId() uint32 {
	if m != nil && m.CameraId != nil {
		return *m.CameraId
	}
	return 0
}

func (m *SSL_GeometryCameraCalibration) GetFocalLength() float32 {
	if m != nil && m.FocalLength != nil {
		return *m.FocalLength
	}
	return 0
}

func (m *SSL_GeometryCameraCalibration) GetPrincipalPointX() float32 {
	if m != nil && m.PrincipalPointX != nil {
		return *m.PrincipalPointX
	}
	return 0
}

func (m *SSL_GeometryCameraCalibration) GetPrincipalPointY() float32 {
	if m != nil && m.PrincipalPointY != nil {
		return *m.PrincipalPointY
	}
	return 0
}

func (m *SSL_GeometryCameraCalibration) GetDistortion() float32 {
	if m != nil && m.Distortion != nil {
		return *m.Distortion
	}
	return 0
}

func (m *SSL_GeometryCameraCalibration) GetQ0() float32 {
	if m != nil && m.Q0 != nil {
		return *m.Q0
	}
	return 0
}

func (m *SSL_GeometryCameraCalibration) GetQ1() float32 {
	if m != nil && m.Q1 != nil {
		return *m.Q1
	}
	return 0
}

func (m *SSL_GeometryCameraCalibration) GetQ2() float32 {
	if m != nil && m.Q2 != nil {
		return *m.Q2
	}
	return 0
}

func (m *SSL_GeometryCameraCalibration) GetQ3() float32 {
	if m != nil && m.Q3 != nil {
		return *m.Q3
	}
	return 0
}

func (m *SSL_GeometryCameraCalibration) GetTx() float32 {
	if m != nil && m.Tx != nil {
		return *m.Tx
	}
	return 0
}

func (m *SSL_GeometryCameraCalibration) GetTy() float32 {
	if m != nil && m.Ty != nil {
		return *m.Ty
	}
	return 0
}

func (m *SSL_GeometryCameraCalibration) GetTz() float32 {
	if m != nil && m.Tz != nil {
		return *m.Tz
	}
	return 0
}

func (m *SSL_GeometryCameraCalibration) GetDerivedCameraWorldTx() float32 {
	if m != nil && m.DerivedCameraWorldTx != nil {
		return *m.DerivedCameraWorldTx
	}
	return 0
}

func (m *SSL_GeometryCameraCalibration) GetDerivedCameraWorldTy() float32 {
	if m != nil && m.DerivedCameraWorldTy != nil {
		return *m.DerivedCameraWorldTy
	}
	return 0
}

func (m *SSL_GeometryCameraCalibration) GetDerivedCameraWorldTz() float32 {
	if m != nil && m.DerivedCameraWorldTz != nil {
		return *m.DerivedCameraWorldTz
	}
	return 0
}

type SSL_GeometryData struct {
	Field                *SSL_GeometryFieldSize           `protobuf:"bytes,1,req,name=field" json:"field,omitempty"`
	Calib                []*SSL_GeometryCameraCalibration `protobuf:"bytes,2,rep,name=calib" json:"calib,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                         `json:"-"`
	XXX_unrecognized     []byte                           `json:"-"`
	XXX_sizecache        int32                            `json:"-"`
}

func (m *SSL_GeometryData) Reset()         { *m = SSL_GeometryData{} }
func (m *SSL_GeometryData) String() string { return proto.CompactTextString(m) }
func (*SSL_GeometryData) ProtoMessage()    {}
func (*SSL_GeometryData) Descriptor() ([]byte, []int) {
	return fileDescriptor_3d55d410315aeff8, []int{5}
}

func (m *SSL_GeometryData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SSL_GeometryData.Unmarshal(m, b)
}
func (m *SSL_GeometryData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SSL_GeometryData.Marshal(b, m, deterministic)
}
func (m *SSL_GeometryData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SSL_GeometryData.Merge(m, src)
}
func (m *SSL_GeometryData) XXX_Size() int {
	return xxx_messageInfo_SSL_GeometryData.Size(m)
}
func (m *SSL_GeometryData) XXX_DiscardUnknown() {
	xxx_messageInfo_SSL_GeometryData.DiscardUnknown(m)
}

var xxx_messageInfo_SSL_GeometryData proto.InternalMessageInfo

func (m *SSL_GeometryData) GetField() *SSL_GeometryFieldSize {
	if m != nil {
		return m.Field
	}
	return nil
}

func (m *SSL_GeometryData) GetCalib() []*SSL_GeometryCameraCalibration {
	if m != nil {
		return m.Calib
	}
	return nil
}

func init() {
	proto.RegisterType((*Vector2F)(nil), "Vector2f")
	proto.RegisterType((*SSL_FieldLineSegment)(nil), "SSL_FieldLineSegment")
	proto.RegisterType((*SSL_FieldCicularArc)(nil), "SSL_FieldCicularArc")
	proto.RegisterType((*SSL_GeometryFieldSize)(nil), "SSL_GeometryFieldSize")
	proto.RegisterType((*SSL_GeometryCameraCalibration)(nil), "SSL_GeometryCameraCalibration")
	proto.RegisterType((*SSL_GeometryData)(nil), "SSL_GeometryData")
}

func init() {
	proto.RegisterFile("ssl_vision_geometry.proto", fileDescriptor_3d55d410315aeff8)
}

var fileDescriptor_3d55d410315aeff8 = []byte{
	// 660 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x54, 0x5d, 0x4f, 0x13, 0x41,
	0x14, 0x4d, 0xb7, 0x14, 0xe8, 0x2d, 0x1f, 0x3a, 0x02, 0x0e, 0x51, 0x10, 0x9b, 0x68, 0x88, 0x91,
	0x96, 0x2e, 0xea, 0xa3, 0x11, 0x4b, 0x34, 0x26, 0x7d, 0x30, 0x5b, 0x23, 0xea, 0xcb, 0x66, 0xba,
	0x3b, 0xb4, 0x13, 0xb7, 0x33, 0xcb, 0xec, 0x14, 0xba, 0x7d, 0xf0, 0xc7, 0xf8, 0xee, 0x3f, 0xf3,
	0x47, 0x98, 0xb9, 0x3b, 0x45, 0x84, 0xd6, 0xb7, 0x9e, 0x73, 0xee, 0x9d, 0x39, 0x7b, 0xef, 0xe9,
	0xc0, 0x76, 0x96, 0x25, 0xe1, 0x85, 0xc8, 0x84, 0x92, 0x61, 0x9f, 0xab, 0x21, 0x37, 0x3a, 0x6f,
	0xa4, 0x5a, 0x19, 0x55, 0x7f, 0x0a, 0xcb, 0x9f, 0x79, 0x64, 0x94, 0xf6, 0xcf, 0xc8, 0x0a, 0x94,
	0xc6, 0xb4, 0xb4, 0xe7, 0xed, 0x7b, 0x41, 0x69, 0x6c, 0x51, 0x4e, 0xbd, 0x02, 0xe5, 0xf5, 0x1f,
	0xb0, 0xd1, 0xed, 0x76, 0xc2, 0x77, 0x82, 0x27, 0x71, 0x47, 0x48, 0xde, 0xe5, 0xfd, 0x21, 0x97,
	0x86, 0x10, 0x58, 0x90, 0x6c, 0xc8, 0xb1, 0xad, 0x1a, 0xe0, 0x6f, 0xb2, 0x0d, 0x5e, 0xda, 0xc2,
	0xd6, 0x9a, 0x5f, 0x6d, 0x4c, 0x8f, 0x0f, 0xbc, 0xb4, 0x85, 0x92, 0x4f, 0xcb, 0xb7, 0x25, 0x9f,
	0x3c, 0x84, 0xaa, 0x19, 0x88, 0xe8, 0xbb, 0xe4, 0x59, 0x46, 0x17, 0xf0, 0xde, 0xbf, 0x44, 0xfd,
	0x67, 0x09, 0xee, 0x5d, 0x19, 0x68, 0x8b, 0x68, 0x94, 0x30, 0x7d, 0xac, 0xa3, 0x99, 0xf7, 0x3f,
	0x86, 0xc5, 0x88, 0x4b, 0xc3, 0xf5, 0x6d, 0x0f, 0x4e, 0x20, 0x5b, 0xb0, 0xa8, 0x59, 0x2c, 0x46,
	0x19, 0x7a, 0xf1, 0x02, 0x87, 0xc8, 0x1a, 0x78, 0xac, 0xe5, 0x6e, 0xf7, 0x58, 0x0b, 0xb1, 0x4f,
	0x2b, 0x0e, 0xdf, 0x30, 0xb9, 0x78, 0xd3, 0xe4, 0x2f, 0x0f, 0x36, 0xad, 0xc9, 0xf7, 0x6e, 0xc6,
	0x68, 0xb6, 0x2b, 0x26, 0xd6, 0xd2, 0xca, 0x99, 0x05, 0x61, 0xc2, 0x65, 0xdf, 0x0c, 0xd0, 0x6e,
	0x25, 0xa8, 0x21, 0xd7, 0x41, 0x8a, 0x3c, 0x82, 0x02, 0x86, 0x97, 0x22, 0x36, 0x03, 0xb4, 0x5e,
	0x09, 0x00, 0xa9, 0x53, 0xcb, 0x90, 0x1d, 0x80, 0xbe, 0x62, 0x89, 0xd3, 0xcb, 0xa8, 0x57, 0x2d,
	0xf3, 0xaf, 0x1c, 0xf3, 0xd4, 0x0c, 0xf0, 0x13, 0x9c, 0x7c, 0x62, 0x09, 0xf2, 0x04, 0xd6, 0x7a,
	0x6a, 0x24, 0x63, 0xa6, 0x73, 0x77, 0x42, 0x05, 0x4b, 0x56, 0xa7, 0x6c, 0x71, 0xca, 0xab, 0xa9,
	0x8b, 0x44, 0x48, 0x6e, 0x3f, 0xb1, 0xbc, 0x5f, 0xf3, 0x37, 0x1b, 0xb3, 0x76, 0xef, 0xcc, 0x59,
	0x26, 0x23, 0x47, 0x50, 0xa0, 0x90, 0xe9, 0x28, 0xa3, 0x4b, 0xd8, 0xb6, 0xd1, 0x98, 0xb1, 0xb1,
	0xa0, 0x8a, 0x75, 0xc7, 0x3a, 0xca, 0xea, 0xbf, 0xcb, 0xb0, 0x73, 0x7d, 0x5e, 0x6d, 0x36, 0xe4,
	0x9a, 0xb5, 0x59, 0x22, 0x7a, 0x9a, 0x19, 0xa1, 0x24, 0x79, 0x00, 0xd5, 0x08, 0xc9, 0x50, 0xc4,
	0x38, 0xb4, 0xd5, 0x60, 0xb9, 0x20, 0x3e, 0xc4, 0x38, 0x54, 0x15, 0xb1, 0x64, 0x3a, 0xd4, 0x22,
	0xac, 0x35, 0xe4, 0xdc, 0x50, 0x9f, 0xc1, 0xdd, 0x54, 0x0b, 0x19, 0x89, 0x94, 0x25, 0x61, 0xaa,
	0x84, 0x34, 0xe1, 0xd8, 0xad, 0x7c, 0xfd, 0x4a, 0xf8, 0x68, 0xf9, 0x2f, 0xb3, 0x6a, 0x73, 0x17,
	0x85, 0x1b, 0xb5, 0x5f, 0xc9, 0x2e, 0x40, 0x2c, 0x32, 0xa3, 0xb4, 0x75, 0xe9, 0xf2, 0x71, 0x8d,
	0xb1, 0xb9, 0x39, 0x3f, 0x74, 0x01, 0xf1, 0xce, 0x0f, 0x11, 0xb7, 0xe8, 0x92, 0xc3, 0x98, 0xab,
	0x73, 0x9f, 0x2e, 0x3b, 0xec, 0x23, 0x3e, 0xa2, 0x55, 0x87, 0x8f, 0x2c, 0x36, 0x63, 0x0a, 0x05,
	0x36, 0x63, 0xc4, 0x39, 0xad, 0x39, 0x9c, 0x23, 0x9e, 0xd0, 0x15, 0x87, 0x27, 0xe4, 0x25, 0xdc,
	0x8f, 0xb9, 0x16, 0x17, 0x3c, 0x0e, 0xdd, 0xbc, 0x2e, 0x95, 0x4e, 0xe2, 0xd0, 0x8c, 0xe9, 0xea,
	0x5e, 0x69, 0xdf, 0x0b, 0x36, 0x9c, 0x5c, 0x8c, 0xf8, 0xd4, 0x8a, 0x9f, 0xc6, 0xf3, 0xdb, 0x72,
	0xba, 0x36, 0xb7, 0x2d, 0x9f, 0xdf, 0x36, 0xa1, 0xeb, 0x73, 0xdb, 0x26, 0xf5, 0x0b, 0xb8, 0x73,
	0x7d, 0xdb, 0x27, 0xcc, 0x30, 0xf2, 0x1c, 0x2a, 0x98, 0x07, 0x5c, 0x6e, 0xcd, 0xdf, 0x6a, 0xcc,
	0xfc, 0xff, 0x04, 0x45, 0x11, 0x79, 0x01, 0x95, 0xc8, 0xa6, 0x83, 0x7a, 0x18, 0xb0, 0xdd, 0xc6,
	0x7f, 0xd3, 0x13, 0x14, 0xc5, 0x6f, 0xdf, 0x7c, 0x7b, 0xdd, 0x17, 0x66, 0x30, 0xea, 0x35, 0x22,
	0x35, 0x6c, 0x06, 0xaa, 0xa7, 0xda, 0xa3, 0xf4, 0xa0, 0xdb, 0xed, 0x34, 0xb3, 0x2c, 0x39, 0xe8,
	0xb3, 0x21, 0x3f, 0x88, 0x94, 0x34, 0x5a, 0x25, 0x09, 0xd7, 0x4d, 0x61, 0x1f, 0x08, 0xc9, 0x92,
	0x26, 0x4b, 0xd3, 0x66, 0xf1, 0x68, 0xfe, 0x09, 0x00, 0x00, 0xff, 0xff, 0x87, 0x42, 0x9c, 0x14,
	0x41, 0x05, 0x00, 0x00,
}
