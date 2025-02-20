package framing

import (
	"encoding/binary"
	"fmt"

	"github.com/rsocket/rsocket-go/internal/common"
)

const (
	initReqLen                = 4
	minRequestChannelFrameLen = initReqLen
)

// FrameRequestChannel is frame for RequestChannel.
type FrameRequestChannel struct {
	*BaseFrame
}

// Validate returns error if frame is invalid.
func (p *FrameRequestChannel) Validate() (err error) {
	if p.body.Len() < minRequestChannelFrameLen {
		err = errIncompleteFrame
	}
	return
}

func (p *FrameRequestChannel) String() string {
	m, _ := p.MetadataUTF8()
	return fmt.Sprintf("FrameRequestChannel{%s,data=%s,metadata=%s,initialRequestN=%d}",
		p.header, p.DataUTF8(), m, p.InitialRequestN())
}

// InitialRequestN returns initial N.
func (p *FrameRequestChannel) InitialRequestN() uint32 {
	return binary.BigEndian.Uint32(p.body.Bytes())
}

// Metadata returns metadata bytes.
func (p *FrameRequestChannel) Metadata() ([]byte, bool) {
	return p.trySliceMetadata(initReqLen)
}

// Data returns data bytes.
func (p *FrameRequestChannel) Data() []byte {
	return p.trySliceData(initReqLen)
}

// MetadataUTF8 returns metadata as UTF8 string.
func (p *FrameRequestChannel) MetadataUTF8() (metadata string, ok bool) {
	raw, ok := p.Metadata()
	if ok {
		metadata = common.Bytes2str(raw)
	}
	return
}

// DataUTF8 returns data as UTF8 string.
func (p *FrameRequestChannel) DataUTF8() string {
	return common.Bytes2str(p.Data())
}

// NewFrameRequestChannel returns a new RequestChannel frame.
func NewFrameRequestChannel(sid uint32, n uint32, data, metadata []byte, flags ...FrameFlag) *FrameRequestChannel {
	fg := newFlags(flags...)
	bf := common.BorrowByteBuffer()
	var b4 [4]byte
	binary.BigEndian.PutUint32(b4[:], n)
	if _, err := bf.Write(b4[:]); err != nil {
		common.ReturnByteBuffer(bf)
		panic(err)
	}
	if len(metadata) > 0 {
		fg |= FlagMetadata
		if err := bf.WriteUint24(len(metadata)); err != nil {
			common.ReturnByteBuffer(bf)
			panic(err)
		}
		if _, err := bf.Write(metadata); err != nil {
			common.ReturnByteBuffer(bf)
			panic(err)
		}
	}
	if len(data) > 0 {
		if _, err := bf.Write(data); err != nil {
			common.ReturnByteBuffer(bf)
			panic(err)
		}
	}
	return &FrameRequestChannel{
		&BaseFrame{
			header: NewFrameHeader(sid, FrameTypeRequestChannel, fg),
			body:   bf,
		},
	}
}
