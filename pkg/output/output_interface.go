package output

import (
	"bytes"
)

type OutputPlugin interface {
	Commit(data *bytes.Buffer, opi *OutputPluginInfo) error
	Delete(data *bytes.Buffer, opi *OutputPluginInfo) error
}

type OutputPluginInfo struct {
	Index   int
	Outname string
}
