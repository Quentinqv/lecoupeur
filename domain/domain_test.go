package domain

import (
	"testing"
)

func TestCheckURL(t *testing.T) {
	type args struct {
		urlInput string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "CheckURL with valid URL",
			args: args{
				urlInput: "https://example.com",
			},
			want: true,
		},
		{
			name: "CheckURL with invalid URL",
			args: args{
				urlInput: "example.com",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckURL(tt.args.urlInput); got != tt.want {
				t.Errorf("CheckURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateUniqueID(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "GenerateUniqueID with length 6",
			want: "123456",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateUniqueID(); len(got) != len(tt.want) {
				t.Errorf("GenerateUniqueID() = %v, want %v", len(got), len(tt.want))
			}
		})
	}
}

func TestShortenURL(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ShortenURL with valid URL",
			args: args{
				url: "https://example.com",
			},
			want: "http://localhost:8080/abcdef",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShortenURL(tt.args.url); len(got) != len(tt.want) {
				t.Errorf("ShortenURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
