package framing

import (
	"encoding/binary"
	"fmt"

	"github.com/rsocket/rsocket-go/internal/common"
)

// FrameResumeOK represents a frame of ResumeOK.
type FrameResumeOK struct {
	*BaseFrame
}

func (p *FrameResumeOK) String() string {
	return fmt.Sprintf("FrameResumeOK{%s,lastReceivedClientPosition=%d}", p.header, p.LastReceivedClientPosition())
}

// Validate validate current frame.
func (p *FrameResumeOK) Validate() (err error) {
	return
}

// LastReceivedClientPosition returns last received client position.
func (p *FrameResumeOK) LastReceivedClientPosition() uint64 {
	raw := p.body.Bytes()
	return binary.BigEndian.Uint64(raw)
}

// NewResumeOK creates a new frame of ResumeOK.
func NewResumeOK(position uint64) *FrameResumeOK {
	var b8 [8]byte
	binary.BigEndian.PutUint64(b8[:], position)
	bf := common.BorrowByteBuffer()
	_, err := bf.Write(b8[:])
	if err != nil {
		common.ReturnByteBuffer(bf)
		panic(err)
	}
	return &FrameResumeOK{
		&BaseFrame{
			header: NewFrameHeader(0, FrameTypeResumeOK),
			body:   bf,
		},
	}
}
