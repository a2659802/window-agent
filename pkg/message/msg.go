package message

type MessageType uint8

const (
	MsgTypePing = MessageType(iota)
	MsgTypePong
	MsgTypePwdRequest
	MsgTypeResponse
)

// 0000----  前4位是版本号，后面4位保留
const (
	CurrentProtocolVersion = uint8(1)
	ProtocolVersionBits    = 4
)

const (
	HeaderSize  = 1
	TypeSize    = 1
	DataLenSize = 4
)

// 传输时的格式HTLV
/*
| 1B    | 4B    | ??B     |
| TYPE  | LEN   | PAYLOAD |
*/
type Message struct {
	Header  uint8
	Type    MessageType
	DataLen uint32
	RawData []byte
	Data    interface{} // 简单起见，用JSON传输数据
}

type StatusCode int

const (
	StatusOK         = StatusCode(0b0001) // 操作成功
	StatusBadRequest = StatusCode(0b0010) // Type 不支持等
	StatusFailed     = StatusCode(0b0100) // 操作失败
)

type ResponseMessage struct {
	Code   StatusCode `json:"code"`
	ID     string     `json:"id"`
	Reason string     `json:"reason"`
}

func GetVersion(header uint8) uint8 {
	return header >> ProtocolVersionBits
}
