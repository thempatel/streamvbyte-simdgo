package util

import (
	"encoding/binary"
	"reflect"
	"testing"
)

func TestVarintRoundTrip(t *testing.T) {
	count := 8
	nums := GenUint32(count)
	out := make([]byte, count*binary.MaxVarintLen32)
	written := PutVarint(nums, out)
	out = out[:written]

	actual := make([]uint32, count)
	read := GetVarint(out, actual)

	if written != read {
		t.Fatalf("expected to read %d, got %d", written, read)
	}

	if !reflect.DeepEqual(nums, actual) {
		t.Fatalf("expected %+v, got %+v", nums, actual)
	}
}

func TestVarintDeltaRoundTrip(t *testing.T) {
	count := 8
	nums := GenUint32(count)
	SortUint32(nums)
	out := make([]byte, count*binary.MaxVarintLen32)
	written := PutDeltaVarint(nums, out, 0)
	out = out[:written]

	actual := make([]uint32, count)
	read := GetDeltaVarint(out, actual, 0)

	if written != read {
		t.Fatalf("expected to read %d, got %d", written, read)
	}

	if !reflect.DeepEqual(nums, actual) {
		t.Fatalf("expected %+v, got %+v", nums, actual)
	}
}
