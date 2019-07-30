// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/cloud/automl/v1beta1/translation.proto

package automl

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

// Dataset metadata that is specific to translation.
type TranslationDatasetMetadata struct {
	// Required. The BCP-47 language code of the source language.
	SourceLanguageCode string `protobuf:"bytes,1,opt,name=source_language_code,json=sourceLanguageCode,proto3" json:"source_language_code,omitempty"`
	// Required. The BCP-47 language code of the target language.
	TargetLanguageCode   string   `protobuf:"bytes,2,opt,name=target_language_code,json=targetLanguageCode,proto3" json:"target_language_code,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TranslationDatasetMetadata) Reset()         { *m = TranslationDatasetMetadata{} }
func (m *TranslationDatasetMetadata) String() string { return proto.CompactTextString(m) }
func (*TranslationDatasetMetadata) ProtoMessage()    {}
func (*TranslationDatasetMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_74f6484316c15700, []int{0}
}

func (m *TranslationDatasetMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TranslationDatasetMetadata.Unmarshal(m, b)
}
func (m *TranslationDatasetMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TranslationDatasetMetadata.Marshal(b, m, deterministic)
}
func (m *TranslationDatasetMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TranslationDatasetMetadata.Merge(m, src)
}
func (m *TranslationDatasetMetadata) XXX_Size() int {
	return xxx_messageInfo_TranslationDatasetMetadata.Size(m)
}
func (m *TranslationDatasetMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_TranslationDatasetMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_TranslationDatasetMetadata proto.InternalMessageInfo

func (m *TranslationDatasetMetadata) GetSourceLanguageCode() string {
	if m != nil {
		return m.SourceLanguageCode
	}
	return ""
}

func (m *TranslationDatasetMetadata) GetTargetLanguageCode() string {
	if m != nil {
		return m.TargetLanguageCode
	}
	return ""
}

// Evaluation metrics for the dataset.
type TranslationEvaluationMetrics struct {
	// Output only. BLEU score.
	BleuScore float64 `protobuf:"fixed64,1,opt,name=bleu_score,json=bleuScore,proto3" json:"bleu_score,omitempty"`
	// Output only. BLEU score for base model.
	BaseBleuScore        float64  `protobuf:"fixed64,2,opt,name=base_bleu_score,json=baseBleuScore,proto3" json:"base_bleu_score,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TranslationEvaluationMetrics) Reset()         { *m = TranslationEvaluationMetrics{} }
func (m *TranslationEvaluationMetrics) String() string { return proto.CompactTextString(m) }
func (*TranslationEvaluationMetrics) ProtoMessage()    {}
func (*TranslationEvaluationMetrics) Descriptor() ([]byte, []int) {
	return fileDescriptor_74f6484316c15700, []int{1}
}

func (m *TranslationEvaluationMetrics) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TranslationEvaluationMetrics.Unmarshal(m, b)
}
func (m *TranslationEvaluationMetrics) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TranslationEvaluationMetrics.Marshal(b, m, deterministic)
}
func (m *TranslationEvaluationMetrics) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TranslationEvaluationMetrics.Merge(m, src)
}
func (m *TranslationEvaluationMetrics) XXX_Size() int {
	return xxx_messageInfo_TranslationEvaluationMetrics.Size(m)
}
func (m *TranslationEvaluationMetrics) XXX_DiscardUnknown() {
	xxx_messageInfo_TranslationEvaluationMetrics.DiscardUnknown(m)
}

var xxx_messageInfo_TranslationEvaluationMetrics proto.InternalMessageInfo

func (m *TranslationEvaluationMetrics) GetBleuScore() float64 {
	if m != nil {
		return m.BleuScore
	}
	return 0
}

func (m *TranslationEvaluationMetrics) GetBaseBleuScore() float64 {
	if m != nil {
		return m.BaseBleuScore
	}
	return 0
}

// Model metadata that is specific to translation.
type TranslationModelMetadata struct {
	// The resource name of the model to use as a baseline to train the custom
	// model. If unset, we use the default base model provided by Google
	// Translate. Format:
	// `projects/{project_id}/locations/{location_id}/models/{model_id}`
	BaseModel string `protobuf:"bytes,1,opt,name=base_model,json=baseModel,proto3" json:"base_model,omitempty"`
	// Output only. Inferred from the dataset.
	// The source languge (The BCP-47 language code) that is used for training.
	SourceLanguageCode string `protobuf:"bytes,2,opt,name=source_language_code,json=sourceLanguageCode,proto3" json:"source_language_code,omitempty"`
	// Output only. The target languge (The BCP-47 language code) that is used for
	// training.
	TargetLanguageCode   string   `protobuf:"bytes,3,opt,name=target_language_code,json=targetLanguageCode,proto3" json:"target_language_code,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TranslationModelMetadata) Reset()         { *m = TranslationModelMetadata{} }
func (m *TranslationModelMetadata) String() string { return proto.CompactTextString(m) }
func (*TranslationModelMetadata) ProtoMessage()    {}
func (*TranslationModelMetadata) Descriptor() ([]byte, []int) {
	return fileDescriptor_74f6484316c15700, []int{2}
}

func (m *TranslationModelMetadata) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TranslationModelMetadata.Unmarshal(m, b)
}
func (m *TranslationModelMetadata) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TranslationModelMetadata.Marshal(b, m, deterministic)
}
func (m *TranslationModelMetadata) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TranslationModelMetadata.Merge(m, src)
}
func (m *TranslationModelMetadata) XXX_Size() int {
	return xxx_messageInfo_TranslationModelMetadata.Size(m)
}
func (m *TranslationModelMetadata) XXX_DiscardUnknown() {
	xxx_messageInfo_TranslationModelMetadata.DiscardUnknown(m)
}

var xxx_messageInfo_TranslationModelMetadata proto.InternalMessageInfo

func (m *TranslationModelMetadata) GetBaseModel() string {
	if m != nil {
		return m.BaseModel
	}
	return ""
}

func (m *TranslationModelMetadata) GetSourceLanguageCode() string {
	if m != nil {
		return m.SourceLanguageCode
	}
	return ""
}

func (m *TranslationModelMetadata) GetTargetLanguageCode() string {
	if m != nil {
		return m.TargetLanguageCode
	}
	return ""
}

// Annotation details specific to translation.
type TranslationAnnotation struct {
	// Output only . The translated content.
	TranslatedContent    *TextSnippet `protobuf:"bytes,1,opt,name=translated_content,json=translatedContent,proto3" json:"translated_content,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *TranslationAnnotation) Reset()         { *m = TranslationAnnotation{} }
func (m *TranslationAnnotation) String() string { return proto.CompactTextString(m) }
func (*TranslationAnnotation) ProtoMessage()    {}
func (*TranslationAnnotation) Descriptor() ([]byte, []int) {
	return fileDescriptor_74f6484316c15700, []int{3}
}

func (m *TranslationAnnotation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TranslationAnnotation.Unmarshal(m, b)
}
func (m *TranslationAnnotation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TranslationAnnotation.Marshal(b, m, deterministic)
}
func (m *TranslationAnnotation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TranslationAnnotation.Merge(m, src)
}
func (m *TranslationAnnotation) XXX_Size() int {
	return xxx_messageInfo_TranslationAnnotation.Size(m)
}
func (m *TranslationAnnotation) XXX_DiscardUnknown() {
	xxx_messageInfo_TranslationAnnotation.DiscardUnknown(m)
}

var xxx_messageInfo_TranslationAnnotation proto.InternalMessageInfo

func (m *TranslationAnnotation) GetTranslatedContent() *TextSnippet {
	if m != nil {
		return m.TranslatedContent
	}
	return nil
}

func init() {
	proto.RegisterType((*TranslationDatasetMetadata)(nil), "google.cloud.automl.v1beta1.TranslationDatasetMetadata")
	proto.RegisterType((*TranslationEvaluationMetrics)(nil), "google.cloud.automl.v1beta1.TranslationEvaluationMetrics")
	proto.RegisterType((*TranslationModelMetadata)(nil), "google.cloud.automl.v1beta1.TranslationModelMetadata")
	proto.RegisterType((*TranslationAnnotation)(nil), "google.cloud.automl.v1beta1.TranslationAnnotation")
}

func init() {
	proto.RegisterFile("google/cloud/automl/v1beta1/translation.proto", fileDescriptor_74f6484316c15700)
}

var fileDescriptor_74f6484316c15700 = []byte{
	// 419 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x92, 0xdf, 0xaa, 0xd3, 0x40,
	0x10, 0xc6, 0x49, 0x04, 0xa1, 0x2b, 0xa2, 0x06, 0x85, 0x43, 0xcf, 0xf1, 0x0f, 0xbd, 0x90, 0x73,
	0xa1, 0x89, 0x47, 0xef, 0xe2, 0x55, 0x5b, 0xc5, 0x9b, 0x13, 0x38, 0xf4, 0x14, 0x05, 0x29, 0x84,
	0x49, 0x32, 0x2c, 0x81, 0xcd, 0x4e, 0xc8, 0x4e, 0x8a, 0x97, 0xbe, 0x88, 0xcf, 0xe1, 0x3b, 0xf8,
	0x28, 0x3e, 0x85, 0xec, 0x6e, 0xda, 0x06, 0xb1, 0x05, 0xef, 0xda, 0xf9, 0x7e, 0xf3, 0x7d, 0xb3,
	0x93, 0x11, 0xaf, 0x25, 0x91, 0x54, 0x98, 0x94, 0x8a, 0xfa, 0x2a, 0x81, 0x9e, 0xa9, 0x51, 0xc9,
	0xf6, 0xaa, 0x40, 0x86, 0xab, 0x84, 0x3b, 0xd0, 0x46, 0x01, 0xd7, 0xa4, 0xe3, 0xb6, 0x23, 0xa6,
	0xe8, 0xdc, 0xe3, 0xb1, 0xc3, 0x63, 0x8f, 0xc7, 0x03, 0x3e, 0x7d, 0x75, 0xca, 0xab, 0x02, 0x86,
	0xbc, 0x66, 0x6c, 0x8c, 0xb7, 0x9a, 0x5e, 0x0c, 0x34, 0xb4, 0x75, 0x02, 0x5a, 0x13, 0xbb, 0x9c,
	0x41, 0x9d, 0x7d, 0x0f, 0xc4, 0x74, 0x7d, 0x88, 0xff, 0x00, 0x0c, 0x06, 0x39, 0x43, 0x06, 0x6b,
	0x14, 0xbd, 0x11, 0x8f, 0x0d, 0xf5, 0x5d, 0x89, 0xb9, 0x02, 0x2d, 0x7b, 0x90, 0x98, 0x97, 0x54,
	0xe1, 0x59, 0xf0, 0x22, 0xb8, 0x9c, 0xac, 0x22, 0xaf, 0x5d, 0x0f, 0xd2, 0x92, 0x2a, 0xb4, 0x1d,
	0x0c, 0x9d, 0x44, 0xfe, 0xab, 0x23, 0xf4, 0x1d, 0x5e, 0x1b, 0x77, 0xcc, 0x50, 0x5c, 0x8c, 0x26,
	0xf8, 0xb8, 0x05, 0xd5, 0xbb, 0x5f, 0x19, 0x72, 0x57, 0x97, 0x26, 0x7a, 0x2a, 0x44, 0xa1, 0xb0,
	0xcf, 0x4d, 0x49, 0x9d, 0x4f, 0x0e, 0x56, 0x13, 0x5b, 0xb9, 0xb5, 0x85, 0xe8, 0xa5, 0x78, 0x50,
	0x80, 0xc1, 0x7c, 0xc4, 0x84, 0x8e, 0xb9, 0x6f, 0xcb, 0x8b, 0x1d, 0x37, 0xfb, 0x11, 0x88, 0xb3,
	0x51, 0x4e, 0x46, 0x15, 0xaa, 0xfd, 0x3b, 0x6d, 0x86, 0x35, 0x69, 0x6c, 0x75, 0x78, 0xdd, 0xc4,
	0x56, 0x1c, 0x76, 0x74, 0x0d, 0xe1, 0x7f, 0xaf, 0xe1, 0xce, 0xd1, 0x35, 0xb4, 0xe2, 0xc9, 0x68,
	0xbc, 0xf9, 0xfe, 0x4b, 0x45, 0x5f, 0x44, 0xb4, 0x3b, 0x10, 0xac, 0xf2, 0x92, 0x34, 0xa3, 0x66,
	0x37, 0xe3, 0xbd, 0xb7, 0x97, 0xf1, 0x89, 0x43, 0x89, 0xd7, 0xf8, 0x8d, 0x6f, 0x75, 0xdd, 0xb6,
	0xc8, 0xab, 0x47, 0x07, 0x8f, 0xa5, 0xb7, 0x58, 0xfc, 0x0c, 0xc4, 0xf3, 0x92, 0x9a, 0x53, 0x16,
	0x8b, 0x87, 0xa3, 0x99, 0x6e, 0xec, 0xc5, 0xdc, 0x04, 0x5f, 0xe7, 0x43, 0x83, 0x24, 0xfb, 0xb6,
	0x98, 0x3a, 0x99, 0x48, 0xd4, 0xee, 0x9e, 0x12, 0x2f, 0x41, 0x5b, 0x9b, 0x7f, 0x9e, 0xe7, 0x7b,
	0xff, 0xf7, 0x57, 0x78, 0xfe, 0xc9, 0x81, 0x9b, 0xa5, 0x85, 0x36, 0xf3, 0x9e, 0x29, 0x53, 0x9b,
	0xcf, 0x1e, 0xfa, 0x1d, 0x3e, 0xf3, 0x6a, 0x9a, 0x3a, 0x39, 0x4d, 0x9d, 0x7e, 0x9d, 0xa6, 0x03,
	0x50, 0xdc, 0x75, 0x61, 0xef, 0xfe, 0x04, 0x00, 0x00, 0xff, 0xff, 0x63, 0x05, 0x99, 0xec, 0x56,
	0x03, 0x00, 0x00,
}