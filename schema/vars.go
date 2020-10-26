// Package schema contains the variables
package schema

// MediaType is media types
type MediaType string

const (
	ProtoV2     MediaType = "application/vnd.protobuf.proto.v2"
	ProtoV3     MediaType = "application/vnd.protobuf.proto.v3"
	Content     MediaType = "application/vnd.protobuf.content.v1.bin"
	ContentGzip MediaType = "application/vnd.protobuf.content.v1.bin+gzip"
	CodeGenGo   MediaType = "application/vnd.protobuf.codegen.go.v1.tar+gzip"
)

func (mt MediaType) String() string {
	return string(mt)
}