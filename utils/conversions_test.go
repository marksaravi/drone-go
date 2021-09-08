package utils

import (
	"testing"
)

func compareByteArrays(ba1, ba2 []byte) bool {
	l1 := len(ba1)
	if l1 != len(ba2) {
		return false
	}
	for i := 0; i < l1; i++ {
		if ba1[i] != ba2[i] {
			return false
		}
	}
	return true
}

func compareBoolArrays(ba1, ba2 []bool) bool {
	l1 := len(ba1)
	if l1 != len(ba2) {
		return false
	}
	for i := 0; i < l1; i++ {
		if ba1[i] != ba2[i] {
			return false
		}
	}
	return true
}

func TestPositiveFloat32ToBytes(t *testing.T) {
	got := Float32ToBytes(1324.456)
	want := []byte{152, 142, 165, 68}
	if !compareByteArrays(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestBytesToPositiveFloat32(t *testing.T) {
	got := Float32FromBytes([]byte{152, 142, 165, 68})
	var want float32 = 1324.456
	if want != got {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestNegativeFloat32ToBytes(t *testing.T) {
	got := Float32ToBytes(-360.742)
	want := []byte{250, 94, 180, 195}
	if !compareByteArrays(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestBytesToNetagiveFloat32(t *testing.T) {
	got := Float32FromBytes([]byte{250, 94, 180, 195})
	var want float32 = -360.742
	if want != got {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestBoolArrayToByte(t *testing.T) {
	got := BoolArrayToByte([8]bool{true, true, false, false, true, true, false, true})
	var want byte = 179
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestBoolArrayFromByte(t *testing.T) {
	got := BoolArrayFromByte(179)
	var want [8]bool = [8]bool{true, true, false, false, true, true, false, true}
	if !compareBoolArrays(got[0:], want[0:]) {
		t.Errorf("got %v, want %v", got, want)
	}
}
