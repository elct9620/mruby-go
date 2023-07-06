package rite_test

import (
	"bytes"
	"testing"

	"github.com/elct9620/mruby-go/rite"
	"github.com/google/go-cmp/cmp"
)

func Test_Load(t *testing.T) {
	buffer := bytes.NewBuffer([]byte{
		0x52, 0x49, 0x54, 0x45, 0x30, 0x33, 0x30, 0x30, 0x00, 0x00, 0x00, 0x41, 0x4d, 0x41, 0x54, 0x5a,
		0x30, 0x30, 0x30, 0x30, 0x49, 0x52, 0x45, 0x50, 0x00, 0x00, 0x00, 0x25, 0x30, 0x33, 0x30, 0x30,
		0x00, 0x00, 0x00, 0x19, 0x00, 0x01, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05,
		0x08, 0x01, 0x38, 0x01, 0x69, 0x00, 0x00, 0x00, 0x00, 0x45, 0x4e, 0x44, 0x00, 0x00, 0x00, 0x00,
		0x08,
	})

	entity, err := rite.Load(buffer)
	if err != nil {
		t.Fatal(err)
	}

	expectedHeader := `#<RITE id="RITE" version="03.00" size="65" compiler="MATZ#0000">`
	header := entity.Header().String()
	if !cmp.Equal(expectedHeader, header) {
		t.Fatal("RITE header mismatch", cmp.Diff(expectedHeader, header))
	}

	expectedSectionIdent := []rite.SectionType{
		rite.TypeIREP,
	}
	expectedSectionSize := len(expectedSectionIdent)
	sections := entity.Sections()
	sectionSize := len(sections)

	if !cmp.Equal(expectedSectionSize, sectionSize) {
		t.Fatal("RITE Section size mismatch", cmp.Diff(expectedSectionSize, sectionSize))
	}

	for idx, section := range sections {
		sectionType := section.Type()
		if !cmp.Equal(expectedSectionIdent[idx], sectionType) {
			t.Fatal("RITE Section Identity mismatch", cmp.Diff(expectedSectionIdent[idx], sectionType))
		}
	}
}
