// +build integration

package ligo_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/goller/ligo"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestReadFile(t *testing.T) {
	file := "L-L1_GWOSC_16KHZ_R1-1240215487-32.gwf"
	b, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatalf("unable to read file: %v", err)
	}

	// TODO(goller): this is a bit inefficient so we should switch to memory
	// mapping at some point.
	buf := bytes.NewReader(b)

	header := ligo.FileHeader{}
	binary.Read(buf, binary.LittleEndian, &header)
	fmt.Printf("%+#v\n", header)

	footer := ligo.FileFooter{}
	size := binary.Size(footer)
	fbuf := bytes.NewReader(b[len(b)-size:])
	binary.Read(fbuf, binary.LittleEndian, &footer)
	fmt.Printf("%+#v\n", footer)

	tbuf := bytes.NewReader(b[uint64(len(b))-footer.SeekTOC : len(b)-size])
	toc := ligo.TableOfContents{}
	binary.Read(tbuf, binary.LittleEndian, &toc.TOCHeader)

	toc.DataQuality = make([]uint32, toc.Frames)
	binary.Read(tbuf, binary.LittleEndian, &toc.DataQuality)
	toc.GTimeS = make([]uint32, toc.Frames)
	binary.Read(tbuf, binary.LittleEndian, &toc.GTimeS)
	toc.GTimeN = make([]uint32, toc.Frames)
	binary.Read(tbuf, binary.LittleEndian, &toc.GTimeN)
	toc.DT = make([]float64, toc.Frames)
	binary.Read(tbuf, binary.LittleEndian, &toc.DT)
	fmt.Printf("%+#v\n", toc)
}
