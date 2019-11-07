package nux

import (
	"os"
	"testing"
)

func TestRoot(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Root(); got != tt.want {
				t.Errorf("Root() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRootWithValue(t *testing.T) {
	os.Setenv(nuxRootFs, "/var")
	tests := []struct {
		name string
		want string
	}{
		{
			want: "/var",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Root(); got != tt.want {
				t.Errorf("Root() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRootWithErrorFormatValue(t *testing.T) {
	os.Setenv(nuxRootFs, "/var:/tmp")
	tests := []struct {
		name string
		want string
	}{
		{
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Root(); got != tt.want {
				t.Errorf("Root() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRootWithNotExistsValue(t *testing.T) {
	os.Setenv(nuxRootFs, "/var_NotExists")
	tests := []struct {
		name string
		want string
	}{
		{
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Root(); got != tt.want {
				t.Errorf("Root() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRootWithRelativeValue(t *testing.T) {
	os.Setenv(nuxRootFs, "var_NotExists")
	tests := []struct {
		name string
		want string
	}{
		{
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Root(); got != tt.want {
				t.Errorf("Root() = %v, want %v", got, tt.want)
			}
		})
	}
}
