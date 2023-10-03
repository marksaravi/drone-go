package utils_test

import (
	"testing"
	"time"

	"github.com/marksaravi/drone-go/utils"
)

func TestSerializeFloat64(t *testing.T) {
	testCases := []struct{ f, want float64 }{
		{
			f:    96.237,
			want: 96.2,
		},
		{
			f:    17.87,
			want: 17.9,
		},
		{
			f:    -16.79,
			want: -16.8,
		},
		{
			f:    -345.245,
			want: -345.2,
		},
	}
	for i, tc := range testCases {
		data := utils.SerializeFloat64(tc.f)
		got := utils.DeSerializeFloat64(data)

		if got != tc.want {
			t.Errorf("#%2d: testing %f, wanted %f, got %f", i, tc.f, tc.want, got)
		}
	}
}

func TestSerializeInt(t *testing.T) {
	testCases := []struct{ n, want int16 }{
		{
			n:    96,
			want: 96,
		},
		{
			n:    -17,
			want: -17,
		},
		{
			n:    160,
			want: 160,
		},
		{
			n:    -345,
			want: -345,
		},
	}
	for i, tc := range testCases {
		data := utils.SerializeInt(tc.n)
		got := utils.DeSerializeInt(data)

		if got != tc.want {
			t.Errorf("#%2d: testing %d, wanted %d, got %d", i, tc.n, tc.want, got)
		}
	}
}

func TestSerializeDuration(t *testing.T) {
	testCases := []struct{ dur, want time.Duration }{
		{
			dur:  time.Duration(time.Microsecond * 8774623),
			want: time.Duration(time.Microsecond * 8774600),
		},
		{
			dur:  time.Duration(time.Microsecond * 4439038),
			want: time.Duration(time.Microsecond * 4439000),
		},
	}
	for i, tc := range testCases {
		data := utils.SerializeDuration(tc.dur)
		got := utils.DeSerializeDuration(data)

		if got != tc.want {
			t.Errorf("#%2d: testing %v, wanted %v, got %v", i, tc.dur, tc.want, got)
		}
	}
}
