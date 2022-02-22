package output

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	log "github.com/sirupsen/logrus"
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
	filepath := filepath.Join(fo.outDir, "nddtest-"+strconv.Itoa(opi.Index)+".yaml")

	log.Infof("Writing to %s", filepath)

	f, err := os.Create(filepath)
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
