package random

import "testing"

// making tests for "NewRandomString(length int) string" function in random.go. Length of returned string must be equal to length argument.
func TestNewRandomString(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{
			name: "size = 1",
			size: 1,
		},
		{
			name: "size = 5",
			size: 5,
		},
		{
			name: "size = 10",
			size: 10,
		},
		{
			name: "size = 20",
			size: 20,
		},
		{
			name: "size = 50",
			size: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRandomString(tt.size)
			if len(got) != tt.size {
				t.Errorf("NewRandomString() = %v, want %v", len(got), tt.size)
			}
		})
	}
}
