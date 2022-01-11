// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1-devel
// 	protoc        v3.18.0
// source: api/query.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type InvoiceStatus int32

const (
	InvoiceStatus_AVAILABLE InvoiceStatus = 0
	InvoiceStatus_FINANCED  InvoiceStatus = 1
	InvoiceStatus_APPROVED  InvoiceStatus = 2
	InvoiceStatus_REVERSED  InvoiceStatus = 3
)

// Enum value maps for InvoiceStatus.
var (
	InvoiceStatus_name = map[int32]string{
		0: "AVAILABLE",
		1: "FINANCED",
		2: "APPROVED",
		3: "REVERSED",
	}
	InvoiceStatus_value = map[string]int32{
		"AVAILABLE": 0,
		"FINANCED":  1,
		"APPROVED":  2,
		"REVERSED":  3,
	}
)

func (x InvoiceStatus) Enum() *InvoiceStatus {
	p := new(InvoiceStatus)
	*p = x
	return p
}

func (x InvoiceStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (InvoiceStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_api_query_proto_enumTypes[0].Descriptor()
}

func (InvoiceStatus) Type() protoreflect.EnumType {
	return &file_api_query_proto_enumTypes[0]
}

func (x InvoiceStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use InvoiceStatus.Descriptor instead.
func (InvoiceStatus) EnumDescriptor() ([]byte, []int) {
	return file_api_query_proto_rawDescGZIP(), []int{0}
}

type AllInvestorsReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Investors []*Investor `protobuf:"bytes,1,rep,name=investors,proto3" json:"investors,omitempty"`
}

func (x *AllInvestorsReply) Reset() {
	*x = AllInvestorsReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_query_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AllInvestorsReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AllInvestorsReply) ProtoMessage() {}

func (x *AllInvestorsReply) ProtoReflect() protoreflect.Message {
	mi := &file_api_query_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AllInvestorsReply.ProtoReflect.Descriptor instead.
func (*AllInvestorsReply) Descriptor() ([]byte, []int) {
	return file_api_query_proto_rawDescGZIP(), []int{0}
}

func (x *AllInvestorsReply) GetInvestors() []*Investor {
	if x != nil {
		return x.Investors
	}
	return nil
}

type Investor struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       *UUID  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Balance  *Money `protobuf:"bytes,2,opt,name=balance,proto3" json:"balance,omitempty"`
	Reserved *Money `protobuf:"bytes,3,opt,name=reserved,proto3" json:"reserved,omitempty"`
}

func (x *Investor) Reset() {
	*x = Investor{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_query_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Investor) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Investor) ProtoMessage() {}

func (x *Investor) ProtoReflect() protoreflect.Message {
	mi := &file_api_query_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Investor.ProtoReflect.Descriptor instead.
func (*Investor) Descriptor() ([]byte, []int) {
	return file_api_query_proto_rawDescGZIP(), []int{1}
}

func (x *Investor) GetId() *UUID {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *Investor) GetBalance() *Money {
	if x != nil {
		return x.Balance
	}
	return nil
}

func (x *Investor) GetReserved() *Money {
	if x != nil {
		return x.Reserved
	}
	return nil
}

type AllInvoicesReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Invoices []*Invoice `protobuf:"bytes,1,rep,name=invoices,proto3" json:"invoices,omitempty"`
}

func (x *AllInvoicesReply) Reset() {
	*x = AllInvoicesReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_query_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AllInvoicesReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AllInvoicesReply) ProtoMessage() {}

func (x *AllInvoicesReply) ProtoReflect() protoreflect.Message {
	mi := &file_api_query_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AllInvoicesReply.ProtoReflect.Descriptor instead.
func (*AllInvoicesReply) Descriptor() ([]byte, []int) {
	return file_api_query_proto_rawDescGZIP(), []int{2}
}

func (x *AllInvoicesReply) GetInvoices() []*Invoice {
	if x != nil {
		return x.Invoices
	}
	return nil
}

type Invoice struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          *UUID         `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	IssuerId    *UUID         `protobuf:"bytes,2,opt,name=issuer_id,json=issuerId,proto3" json:"issuer_id,omitempty"`
	AskingPrice *Money        `protobuf:"bytes,3,opt,name=asking_price,json=askingPrice,proto3" json:"asking_price,omitempty"`
	Status      InvoiceStatus `protobuf:"varint,4,opt,name=status,proto3,enum=InvoiceStatus" json:"status,omitempty"`
	WinningBid  *Bid          `protobuf:"bytes,5,opt,name=winning_bid,json=winningBid,proto3,oneof" json:"winning_bid,omitempty"`
}

func (x *Invoice) Reset() {
	*x = Invoice{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_query_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Invoice) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Invoice) ProtoMessage() {}

func (x *Invoice) ProtoReflect() protoreflect.Message {
	mi := &file_api_query_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Invoice.ProtoReflect.Descriptor instead.
func (*Invoice) Descriptor() ([]byte, []int) {
	return file_api_query_proto_rawDescGZIP(), []int{3}
}

func (x *Invoice) GetId() *UUID {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *Invoice) GetIssuerId() *UUID {
	if x != nil {
		return x.IssuerId
	}
	return nil
}

func (x *Invoice) GetAskingPrice() *Money {
	if x != nil {
		return x.AskingPrice
	}
	return nil
}

func (x *Invoice) GetStatus() InvoiceStatus {
	if x != nil {
		return x.Status
	}
	return InvoiceStatus_AVAILABLE
}

func (x *Invoice) GetWinningBid() *Bid {
	if x != nil {
		return x.WinningBid
	}
	return nil
}

type Bid struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	InvestorId *UUID  `protobuf:"bytes,1,opt,name=investor_id,json=investorId,proto3" json:"investor_id,omitempty"`
	Amount     *Money `protobuf:"bytes,2,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (x *Bid) Reset() {
	*x = Bid{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_query_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Bid) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Bid) ProtoMessage() {}

func (x *Bid) ProtoReflect() protoreflect.Message {
	mi := &file_api_query_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Bid.ProtoReflect.Descriptor instead.
func (*Bid) Descriptor() ([]byte, []int) {
	return file_api_query_proto_rawDescGZIP(), []int{4}
}

func (x *Bid) GetInvestorId() *UUID {
	if x != nil {
		return x.InvestorId
	}
	return nil
}

func (x *Bid) GetAmount() *Money {
	if x != nil {
		return x.Amount
	}
	return nil
}

var File_api_query_proto protoreflect.FileDescriptor

var file_api_query_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x61, 0x70, 0x69, 0x2f, 0x71, 0x75, 0x65, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x13, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x69, 0x6e, 0x61, 0x6e, 0x63, 0x69, 0x6e, 0x67,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x3c, 0x0a, 0x11, 0x41, 0x6c, 0x6c, 0x49, 0x6e, 0x76, 0x65, 0x73, 0x74,
	0x6f, 0x72, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x27, 0x0a, 0x09, 0x69, 0x6e, 0x76, 0x65,
	0x73, 0x74, 0x6f, 0x72, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x09, 0x2e, 0x49, 0x6e,
	0x76, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x52, 0x09, 0x69, 0x6e, 0x76, 0x65, 0x73, 0x74, 0x6f, 0x72,
	0x73, 0x22, 0x67, 0x0a, 0x08, 0x49, 0x6e, 0x76, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x12, 0x15, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x55, 0x55, 0x49, 0x44,
	0x52, 0x02, 0x69, 0x64, 0x12, 0x20, 0x0a, 0x07, 0x62, 0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x4d, 0x6f, 0x6e, 0x65, 0x79, 0x52, 0x07, 0x62,
	0x61, 0x6c, 0x61, 0x6e, 0x63, 0x65, 0x12, 0x22, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x4d, 0x6f, 0x6e, 0x65, 0x79,
	0x52, 0x08, 0x72, 0x65, 0x73, 0x65, 0x72, 0x76, 0x65, 0x64, 0x22, 0x38, 0x0a, 0x10, 0x41, 0x6c,
	0x6c, 0x49, 0x6e, 0x76, 0x6f, 0x69, 0x63, 0x65, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x24,
	0x0a, 0x08, 0x69, 0x6e, 0x76, 0x6f, 0x69, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x08, 0x2e, 0x49, 0x6e, 0x76, 0x6f, 0x69, 0x63, 0x65, 0x52, 0x08, 0x69, 0x6e, 0x76, 0x6f,
	0x69, 0x63, 0x65, 0x73, 0x22, 0xd3, 0x01, 0x0a, 0x07, 0x49, 0x6e, 0x76, 0x6f, 0x69, 0x63, 0x65,
	0x12, 0x15, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x55,
	0x55, 0x49, 0x44, 0x52, 0x02, 0x69, 0x64, 0x12, 0x22, 0x0a, 0x09, 0x69, 0x73, 0x73, 0x75, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x55, 0x55, 0x49,
	0x44, 0x52, 0x08, 0x69, 0x73, 0x73, 0x75, 0x65, 0x72, 0x49, 0x64, 0x12, 0x29, 0x0a, 0x0c, 0x61,
	0x73, 0x6b, 0x69, 0x6e, 0x67, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x06, 0x2e, 0x4d, 0x6f, 0x6e, 0x65, 0x79, 0x52, 0x0b, 0x61, 0x73, 0x6b, 0x69, 0x6e,
	0x67, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x26, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0e, 0x2e, 0x49, 0x6e, 0x76, 0x6f, 0x69, 0x63, 0x65,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x2a,
	0x0a, 0x0b, 0x77, 0x69, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x5f, 0x62, 0x69, 0x64, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x04, 0x2e, 0x42, 0x69, 0x64, 0x48, 0x00, 0x52, 0x0a, 0x77, 0x69, 0x6e,
	0x6e, 0x69, 0x6e, 0x67, 0x42, 0x69, 0x64, 0x88, 0x01, 0x01, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x77,
	0x69, 0x6e, 0x6e, 0x69, 0x6e, 0x67, 0x5f, 0x62, 0x69, 0x64, 0x22, 0x4d, 0x0a, 0x03, 0x42, 0x69,
	0x64, 0x12, 0x26, 0x0a, 0x0b, 0x69, 0x6e, 0x76, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x55, 0x55, 0x49, 0x44, 0x52, 0x0a, 0x69,
	0x6e, 0x76, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x49, 0x64, 0x12, 0x1e, 0x0a, 0x06, 0x61, 0x6d, 0x6f,
	0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x4d, 0x6f, 0x6e, 0x65,
	0x79, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x2a, 0x48, 0x0a, 0x0d, 0x49, 0x6e, 0x76,
	0x6f, 0x69, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0d, 0x0a, 0x09, 0x41, 0x56,
	0x41, 0x49, 0x4c, 0x41, 0x42, 0x4c, 0x45, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x46, 0x49, 0x4e,
	0x41, 0x4e, 0x43, 0x45, 0x44, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x41, 0x50, 0x50, 0x52, 0x4f,
	0x56, 0x45, 0x44, 0x10, 0x02, 0x12, 0x0c, 0x0a, 0x08, 0x52, 0x45, 0x56, 0x45, 0x52, 0x53, 0x45,
	0x44, 0x10, 0x03, 0x32, 0x7f, 0x0a, 0x07, 0x51, 0x75, 0x65, 0x72, 0x69, 0x65, 0x73, 0x12, 0x3a,
	0x0a, 0x0c, 0x41, 0x6c, 0x6c, 0x49, 0x6e, 0x76, 0x65, 0x73, 0x74, 0x6f, 0x72, 0x73, 0x12, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x12, 0x2e, 0x41, 0x6c, 0x6c, 0x49, 0x6e, 0x76, 0x65,
	0x73, 0x74, 0x6f, 0x72, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x38, 0x0a, 0x0b, 0x41, 0x6c,
	0x6c, 0x49, 0x6e, 0x76, 0x6f, 0x69, 0x63, 0x65, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x1a, 0x11, 0x2e, 0x41, 0x6c, 0x6c, 0x49, 0x6e, 0x76, 0x6f, 0x69, 0x63, 0x65, 0x73, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x42, 0x10, 0x5a, 0x0e, 0x2e, 0x2e, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x67,
	0x72, 0x70, 0x63, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_query_proto_rawDescOnce sync.Once
	file_api_query_proto_rawDescData = file_api_query_proto_rawDesc
)

func file_api_query_proto_rawDescGZIP() []byte {
	file_api_query_proto_rawDescOnce.Do(func() {
		file_api_query_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_query_proto_rawDescData)
	})
	return file_api_query_proto_rawDescData
}

var file_api_query_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_query_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_api_query_proto_goTypes = []interface{}{
	(InvoiceStatus)(0),        // 0: InvoiceStatus
	(*AllInvestorsReply)(nil), // 1: AllInvestorsReply
	(*Investor)(nil),          // 2: Investor
	(*AllInvoicesReply)(nil),  // 3: AllInvoicesReply
	(*Invoice)(nil),           // 4: Invoice
	(*Bid)(nil),               // 5: Bid
	(*UUID)(nil),              // 6: UUID
	(*Money)(nil),             // 7: Money
	(*emptypb.Empty)(nil),     // 8: google.protobuf.Empty
}
var file_api_query_proto_depIdxs = []int32{
	2,  // 0: AllInvestorsReply.investors:type_name -> Investor
	6,  // 1: Investor.id:type_name -> UUID
	7,  // 2: Investor.balance:type_name -> Money
	7,  // 3: Investor.reserved:type_name -> Money
	4,  // 4: AllInvoicesReply.invoices:type_name -> Invoice
	6,  // 5: Invoice.id:type_name -> UUID
	6,  // 6: Invoice.issuer_id:type_name -> UUID
	7,  // 7: Invoice.asking_price:type_name -> Money
	0,  // 8: Invoice.status:type_name -> InvoiceStatus
	5,  // 9: Invoice.winning_bid:type_name -> Bid
	6,  // 10: Bid.investor_id:type_name -> UUID
	7,  // 11: Bid.amount:type_name -> Money
	8,  // 12: Queries.AllInvestors:input_type -> google.protobuf.Empty
	8,  // 13: Queries.AllInvoices:input_type -> google.protobuf.Empty
	1,  // 14: Queries.AllInvestors:output_type -> AllInvestorsReply
	3,  // 15: Queries.AllInvoices:output_type -> AllInvoicesReply
	14, // [14:16] is the sub-list for method output_type
	12, // [12:14] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_api_query_proto_init() }
func file_api_query_proto_init() {
	if File_api_query_proto != nil {
		return
	}
	file_api_financing_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_api_query_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AllInvestorsReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_query_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Investor); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_query_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AllInvoicesReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_query_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Invoice); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_api_query_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Bid); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_api_query_proto_msgTypes[3].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_query_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_query_proto_goTypes,
		DependencyIndexes: file_api_query_proto_depIdxs,
		EnumInfos:         file_api_query_proto_enumTypes,
		MessageInfos:      file_api_query_proto_msgTypes,
	}.Build()
	File_api_query_proto = out.File
	file_api_query_proto_rawDesc = nil
	file_api_query_proto_goTypes = nil
	file_api_query_proto_depIdxs = nil
}