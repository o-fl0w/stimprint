package ff

import "testing"

func TestFrequencyHints_Set(t *testing.T) {
	var fhs FrequencyHints
	sfhs := "10:#ffffff;20:#000000"
	err := fhs.Set(sfhs)
	if err != nil {
		t.Errorf("set: %v", err)
	}
	if fhs[0].Hz != 10 {
		t.Errorf("hz=%d, want=10", fhs[0].Hz)
	}
	if fhs[0].Color != "#ffffff" {
		t.Errorf("color=%s, want=#ffffff", fhs[0].Color)
	}
	if fhs[1].Hz != 20 {
		t.Errorf("hz=%d, want=20", fhs[1].Hz)
	}
	if fhs[1].Color != "#000000" {
		t.Errorf("color=%s, want=#000000", fhs[1].Color)
	}
	s := fhs.String()
	if s != sfhs {
		t.Errorf("string()=%s, want=%s", s, sfhs)
	}
}
