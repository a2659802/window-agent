package message

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

var (
	NetworkEnding = binary.BigEndian
)

func Decode(r io.Reader) (*Message, error) {
	// 读header
	headBuf := make([]byte, HeaderSize)
	if nr, err := r.Read(headBuf); err != nil || nr != HeaderSize {
		return nil, fmt.Errorf("read message header error:%v", err)

	}
	ver := GetVersion(uint8(headBuf[0]))
	if ver == 0 || ver > CurrentProtocolVersion {
		return nil, fmt.Errorf("unknown message version:%v", ver)

	}

	// 读type
	typeBuf := make([]byte, TypeSize)
	if nr, err := r.Read(typeBuf); err != nil || nr != TypeSize {
		return nil, fmt.Errorf("read message type error:%v", err)

	}

	// 读len
	lenBuf := make([]byte, DataLenSize)
	if nr, err := io.ReadFull(r, lenBuf); err != nil || nr != DataLenSize {
		return nil, fmt.Errorf("read message len error:%v", err)

	}
	dataLen := NetworkEnding.Uint32(lenBuf)

	// 读body
	payloadBuf := make([]byte, dataLen)
	if nr, err := io.ReadFull(r, payloadBuf); err != nil || nr != int(dataLen) {
		return nil, fmt.Errorf("read message payload error:%v,get %v bytes", err, nr)

	}

	// 封装
	msg := &Message{
		Header:  uint8(headBuf[0]),
		Type:    MessageType(typeBuf[0]),
		DataLen: dataLen,
		RawData: payloadBuf,
	}

	return msg, nil
}

func Encode(w io.Writer, msg *Message) (err error) {
	if msg == nil {
		return fmt.Errorf("cannot encode nil message")
	}

	var payload []byte

	// 写入header
	ver := (CurrentProtocolVersion << (8 - ProtocolVersionBits))
	if _, err = w.Write([]byte{ver}); err != nil {
		return
	}
	msg.Header = ver

	// 写入Type
	if _, err = w.Write([]byte{byte(msg.Type)}); err != nil {
		return
	}

	if msg.Data != nil {
		// 序列化Data->RawData, 计算DataLen
		payload, err = json.Marshal(msg.Data)
		if err != nil {
			return
		}
	}

	// 写入Len
	lenBuf := make([]byte, DataLenSize)
	NetworkEnding.PutUint32(lenBuf, uint32(len(payload)))
	_, err = w.Write(lenBuf)
	msg.DataLen = uint32(len(payload))

	// 写入body
	if len(payload) > 0 {
		_, err = w.Write(payload)
	}

	return
}

// // 读header
// headBuf := make([]byte, message.HeaderSize)
// if nr, err := c.Read(headBuf); err != nil || nr != message.HeaderSize {
// 	logger.Errorf("read message header error:%v", err)
// 	return
// }
// ver := message.GetVersion(uint8(headBuf[0]))
// if ver == 0 || ver > message.CurrentProtocolVersion {
// 	logger.Error("unknown message version:%v", ver)
// 	return
// }

// // 读type
// typeBuf := make([]byte, message.TypeSize)
// if nr, err := c.Read(typeBuf); err != nil || nr != message.TypeSize {
// 	logger.Errorf("read message type error:%v", err)
// 	return
// }

// // 读len
// lenBuf := make([]byte, message.DataLenSize)
// if nr, err := io.ReadFull(c, lenBuf); err != nil || nr != message.DataLenSize {
// 	logger.Errorf("read message len error:%v", err)
// 	return
// }
// dataLen := NetworkEnding.Uint32(lenBuf)

// // 读body
// payloadBuf := make([]byte, dataLen)
// if nr, err := io.ReadFull(c, payloadBuf); err != nil || nr != int(dataLen) {
// 	logger.Errorf("read message payload error:%v,get %v bytes", err, nr)
// 	return
// }

// // 封装
// msg := message.Message{
// 	Header:  uint8(headBuf[0]),
// 	Type:    message.MessageType(typeBuf[0]),
// 	DataLen: dataLen,
// 	RawData: payloadBuf,
// }
