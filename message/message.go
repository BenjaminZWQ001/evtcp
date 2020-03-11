package message

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
)

const HEADER_LEN = 4

type MessagePack struct {
	ConnPoint *net.TCPConn
	Protocol
}

type Protocol struct {
	Length  uint32
	Content []byte
}

func (p *Protocol) ParseContent() (map[string]interface{}, error) {
	var object map[string]interface{}
	unmarshal := json.Unmarshal(p.Content, &object)
	if unmarshal != nil {
		return object, unmarshal
	}
	return object, nil
}

type Content struct {
	ServiceId string `json:"serviceId"`
	Data      string `json:"data"`
}

func Packet(content string) []byte {
	buffer := make([]byte, HEADER_LEN+len(content))
	// 将buffer前面四个字节设置为包长度，大端序
	binary.BigEndian.PutUint32(buffer[0:4], uint32(len(content)))
	copy(buffer[4:], content)
	return buffer
}

func UnPacket(c *net.TCPConn) (*Protocol, error) {
	var (
		p      = &Protocol{}
		header = make([]byte, HEADER_LEN)
	)
	conn := bufio.NewReader(c)
	_, err := io.ReadFull(conn, header)
	if err != nil {
		return p, err
	}
	p.Length = binary.BigEndian.Uint32(header) //转换成10进制的数字
	contentByte := make([]byte, p.Length)
	_, e := io.ReadFull(conn, contentByte) //读取内容
	if e != nil {
		return p, e
	}
	p.Content = contentByte
	return p, nil
}
