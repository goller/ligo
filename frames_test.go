package ligo

import (
	"fmt"
	"testing"
)

func TestDataQuality(t *testing.T) {
	tests := []struct {
		dq   DataQuality
		want uint32
	}{
		{TAMA300, 1 << 0},
		{Virgo, 1 << 2},
		{GEO600, 1 << 4},
		{LIGOHanford2km, 1 << 6},
		{LIGOHanford4km, 1 << 8},
		{LIGOLivingston4km, 1 << 10},
		{LIGOCaltech, 1 << 12},
		{ALLEGRO, 1 << 14},
		{AURIGA, 1 << 16},
		{EXPLORER, 1 << 18},
		{NIOBE, 1 << 20},
		{NAUTILUS, 1 << 22},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			if uint32(tt.dq) != tt.want {
				t.Errorf("DataQuality iota incorrect %d %d", tt.dq, tt.want)
			}
		})
	}
}
