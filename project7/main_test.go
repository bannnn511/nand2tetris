package main

import (
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"strings"
	"testing"
)

type mockFile struct {
	CloseFunc  func() error
	ReadFunc   func(p []byte) (int, error)
	ReadAtFunc func(p []byte, off int64) (int, error)
	SeekFunc   func(offset int64, whence int) (int64, error)
	StatFunc   func() (os.FileInfo, error)
}

var _ io.Reader = (*mockFile)(nil)

func (m *mockFile) Close() error {
	return m.CloseFunc()
}

func (m *mockFile) Read(p []byte) (int, error) {
	return m.ReadFunc(p)
}

func (m *mockFile) ReadAt(p []byte, off int64) (int, error) {
	return m.ReadAtFunc(p, off)
}

func (m *mockFile) Seek(offset int64, whence int) (int64, error) {
	return m.SeekFunc(offset, whence)
}

func (m *mockFile) Stat() (os.FileInfo, error) {
	return m.StatFunc()
}

func TestParser_advance(t *testing.T) {

	type fields struct {
		command string
	}

	tests := []struct {
		name          string
		fields        fields
		expectCmdType CommandType
		expectArg0    string
		expectArg1    string
		expectArg2    string
	}{
		{
			"push constant 17",
			fields{
				command: "push constant 17",
			},
			CPUSH,
			"push",
			"constant",
			"17",
		},
		{
			"pop temp 5",
			fields{
				command: "pop temp 5",
			},
			CPOP,
			"pop",
			"temp",
			"5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := io.NopCloser(strings.NewReader(tt.fields.command))
			file := &mockFile{
				ReadFunc: func(p []byte) (int, error) {
					return reader.Read(p)
				},
				CloseFunc: func() error {
					return nil
				},
			}
			p := NewParser(file)
			if p.hasMoreCommand() {
				p.advance()
			}
			assert.Equal(t, tt.expectCmdType, p.CommandType())
			assert.Equal(t, tt.expectArg0, p.arg0)
			assert.Equal(t, tt.expectArg1, p.arg1)
			assert.Equal(t, tt.expectArg2, p.arg2)
		})
	}
}
