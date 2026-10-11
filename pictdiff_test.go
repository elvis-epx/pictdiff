package main

import "testing"

func TestAbs(t *testing.T) {
	cases := []struct {
		name string
		in   int
		want int
	}{
		{"zero", 0, 0},
		{"positive", 42, 42},
		{"negative", -42, 42},
		{"min int", -1 << 62, 1 << 62},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Abs(tc.in)
			if got != tc.want {
				t.Fatalf("Abs(%d) = %d, want %d", tc.in, got, tc.want)
			}
		})
	}
}

func TestPreMultAlpha(t *testing.T) {
	cases := []struct {
		name    string
		pixel   []uint8
		channel int
		want    int
	}{
		{"full opacity red", []uint8{255, 0, 0, 255}, 0, 255},
		{"half opacity red", []uint8{255, 0, 0, 128}, 0, 128},
		{"zero alpha", []uint8{255, 0, 0, 0}, 0, 0},
		{"green channel half opacity", []uint8{0, 200, 0, 128}, 1, 100},
		{"blue channel full opacity", []uint8{0, 0, 100, 255}, 2, 100},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := PreMultAlpha(tc.pixel, tc.channel)
			if got != tc.want {
				t.Fatalf("PreMultAlpha(%v, %d) = %d, want %d", tc.pixel, tc.channel, got, tc.want)
			}
		})
	}
}
