package output

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type FileOutput struct {
	outDir string
}

func NewFileOutput(outDir string) *FileOutput {
	return &FileOutput{
		outDir: outDir,
	}
}

func (fo *FileOutput) Commit(data *bytes.Buffer, opi *OutputPluginInfo) error {
	f, err := os.Create(filepath.Join(fo.outDir, "nddtest-"+strconv.Itoa(opi.Index)+".yaml"))
	if err != nil {
		return err
	}

	_, err = data.WriteTo(f)
	if err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

func (fo *FileOutput) Delete(data *bytes.Buffer, opi *OutputPluginInfo) error {
	return fmt.Errorf("Not implemented")
}
