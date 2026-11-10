package ibt

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/OJPARKINSON/ibt/headers"
)

func TestInit(t *testing.T) {
	reader, err := NewIbtReader(".testing/valid_test_file.ibt")
	if err != nil {
		t.Fatalf("failed to open test file: %v", err)
	}
	defer func() { _ = reader.Close() }()

	header, err := headers.ParseHeaders(reader)
	if err != nil {
		t.Fatalf("failed to parse headers: %v", err)
	}

	test := NewDirectStructParser(reader, header)

	fmt.Printf("current: %+v\n", test.current)
	dh, _ := json.MarshalIndent(test.header.DiskHeader, "", "  ")
	fmt.Println(string(dh))
	bs, _ := json.MarshalIndent(test.header.TelemetryHeader, "", "  ")
	fmt.Println(string(bs))
	fmt.Printf("maxOffset: %+v\n", test.maxOffset)
	// fmt.Println(test.mmapData)
	fmt.Printf("sessionStartTime: %+v\n", test.sessionStartTime)

	if !reflect.DeepEqual(test.current, 0) {
		t.Fatalf("Does not start at 0")
	}
}
