package database

import (
	"reflect"
	"testing"
)

func TestStoreAndGetURL(t *testing.T) {
	Connect()

	tests := []struct {
		name    string
		url     string
		id      string
		want    string
		wantErr bool
	}{
		{
			name:    "Store and retrieve a URL",
			url:     "https://example.com",
			id:      "test123",
			want:    "https://example.com",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Store the URL
			StoreURL(tt.url, tt.id)

			// Retrieve the URL
			got, err := GetURL(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllURLs(t *testing.T) {
	Connect()
	FlushAll()

	tests := []struct {
		name    string
		want    map[string]URL
		wantErr bool
		prep    func()
	}{
		{
			name: "Get all URLs",
			want: map[string]URL{
				"test123": {URL: "https://example.com", Counter: 0},
			},
			wantErr: false,
			prep: func() {
				FlushAll()
				StoreURL("https://example.com", "test123")
			},
		},
		{
			name:    "Get all URLs with no URLs",
			want:    map[string]URL{},
			wantErr: false,
			prep: func() {
				FlushAll()
			},
		},
		{
			name:    "Get all URLs with an error",
			want:    nil,
			wantErr: true,
			prep: func() {
				FlushAll()
				_ = client.Close()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.prep()

			got, err := GetAllURLs()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllURLs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllURLs() got = %v, want %v", got, tt.want)
			}
		})
	}
}
