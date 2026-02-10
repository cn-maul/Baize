package provider

import (
	"io"
	"strings"
)

// Error 错误结构
type Error struct {
	Message string `json:"message"`
}

// LineScanner 行扫描器
type LineScanner struct {
	scanner *strings.Builder
	reader  io.Reader
}

// NewLineScanner 创建一个新的行扫描器
func NewLineScanner(reader io.Reader) *LineScanner {
	return &LineScanner{
		scanner: &strings.Builder{},
		reader:  reader,
	}
}

// Scan 扫描下一行
func (s *LineScanner) Scan() bool {
	s.scanner.Reset()
	buffer := make([]byte, 1)
	for {
		n, err := s.reader.Read(buffer)
		if n > 0 {
			if buffer[0] == '\n' {
				return true
			}
			s.scanner.Write(buffer[:n])
		}
		if err != nil {
			return s.scanner.Len() > 0
		}
	}
}

// Text 返回当前行的文本
func (s *LineScanner) Text() string {
	return s.scanner.String()
}

// Err 返回扫描过程中的错误
func (s *LineScanner) Err() error {
	return nil
}
