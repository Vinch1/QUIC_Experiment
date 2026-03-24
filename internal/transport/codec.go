package transport

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

func WriteFrame(writer io.Writer, value any) error {
	body, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal frame: %w", err)
	}

	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(body)))

	if _, err := writer.Write(header); err != nil {
		return fmt.Errorf("write frame header: %w", err)
	}
	if _, err := writer.Write(body); err != nil {
		return fmt.Errorf("write frame body: %w", err)
	}
	return nil
}

func ReadFrame(reader io.Reader, value any) error {
	header := make([]byte, 4)
	if _, err := io.ReadFull(reader, header); err != nil {
		return fmt.Errorf("read frame header: %w", err)
	}

	size := binary.BigEndian.Uint32(header)
	body := make([]byte, size)
	if _, err := io.ReadFull(reader, body); err != nil {
		return fmt.Errorf("read frame body: %w", err)
	}

	if err := json.Unmarshal(body, value); err != nil {
		return fmt.Errorf("unmarshal frame: %w", err)
	}
	return nil
}
