package main

import (
	"testing"
)

func TestInitialModelGameFilter(t *testing.T) {
	tests := []struct {
		name       string
		gameFilter string
		wantFilter string
	}{
		{
			name:       "empty filter loads all",
			gameFilter: "",
			wantFilter: "",
		},
		{
			name:       "t5 filter",
			gameFilter: "t5",
			wantFilter: "t5",
		},
		{
			name:       "t6 filter",
			gameFilter: "t6",
			wantFilter: "t6",
		},
		{
			name:       "all filter loads all",
			gameFilter: "all",
			wantFilter: "all",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := initialModel(tt.gameFilter)
			if m.gameFilter != tt.wantFilter {
				t.Errorf("initialModel(%q).gameFilter = %q, want %q",
					tt.gameFilter, m.gameFilter, tt.wantFilter)
			}
		})
	}
}

func TestGetGamesToLoad(t *testing.T) {
	tests := []struct {
		name     string
		gameFlag string
		want     []string
	}{
		{
			name:     "empty flag returns all games",
			gameFlag: "",
			want:     []string{"t5", "t6"},
		},
		{
			name:     "all flag returns all games",
			gameFlag: "all",
			want:     []string{"t5", "t6"},
		},
		{
			name:     "t5 flag returns only t5",
			gameFlag: "t5",
			want:     []string{"t5"},
		},
		{
			name:     "t6 flag returns only t6",
			gameFlag: "t6",
			want:     []string{"t6"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set global gameFlag for test
			oldGameFlag := gameFlag
			gameFlag = tt.gameFlag
			defer func() { gameFlag = oldGameFlag }()

			got := getGamesToLoad()
			if len(got) != len(tt.want) {
				t.Errorf("getGamesToLoad() = %v, want %v", got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("getGamesToLoad()[%d] = %v, want %v", i, got, tt.want)
				}
			}
		})
	}
}
