package database

import (
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
