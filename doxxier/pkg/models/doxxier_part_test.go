package models

import (
	"reflect"
	"testing"
	"time"
)

func TestWithId(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
		want DoxxierPartOption
	}{
		{
			name: "Set ID to a specific value",
			args: args{id: "12345"},
			want: func(dp *DoxxierPart) {
				dp.Id = "12345"
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dp := &DoxxierPart{}
			got := DoxxierPartWithId(tt.args.id)
			got(dp)
			wantDp := &DoxxierPart{}
			tt.want(wantDp)
			if !reflect.DeepEqual(dp, wantDp) {
				t.Errorf("WithId() = %v, want %v", dp, wantDp)
			}
		})
	}
}
func TestWithContent(t *testing.T) {
	type args struct {
		content []byte
	}
	tests := []struct {
		name string
		args args
		want DoxxierPartOption
	}{
		{
			name: "Set Content to a specific value",
			args: args{content: []byte("test content")},
			want: func(dp *DoxxierPart) {
				dp.Content = []byte("test content")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dp := &DoxxierPart{}
			got := DoxxierPartWithContent(tt.args.content)
			got(dp)
			wantDp := &DoxxierPart{}
			tt.want(wantDp)
			if !reflect.DeepEqual(dp, wantDp) {
				t.Errorf("WithContent() = %v, want %v", dp, wantDp)
			}
		})
	}
}
func TestNewDoxxierPart(t *testing.T) {
	tests := []struct {
		name   string
		params []DoxxierPartOption
		want   *DoxxierPart
	}{
		{
			name:   "No options provided",
			params: nil,
			want: &DoxxierPart{
				Id: "", // Id will be a new UUID, so we can't predict it here
			},
		},
		{
			name: "WithId option provided",
			params: []DoxxierPartOption{
				DoxxierPartWithId("12345"),
			},
			want: &DoxxierPart{
				Id: "12345",
			},
		},
		{
			name: "WithContent option provided",
			params: []DoxxierPartOption{
				DoxxierPartWithContent([]byte("test content")),
			},
			want: &DoxxierPart{
				Content: []byte("test content"),
			},
		},
		{
			name: "Multiple options provided",
			params: []DoxxierPartOption{
				DoxxierPartWithId("12345"),
				DoxxierPartWithContent([]byte("test content")),
			},
			want: &DoxxierPart{
				Id:      "12345",
				Content: []byte("test content"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDoxxierPart(tt.params...)
			if tt.name == "No options provided" {
				if got.Id == "" {
					t.Errorf("NewDoxxierPart() Id is empty, expected a new UUID")
				}
				got.Id = "" // Reset Id for comparison
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDoxxierPart() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestToJson(t *testing.T) {

	parsedTime := time.Date(0001, 1, 1, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name    string
		dp      *DoxxierPart
		want    string
		wantErr bool
	}{
		{
			name: "Valid DoxxierPart",
			dp: &DoxxierPart{
				Id:              "12345",
				Content:         []byte("test content"),
				Context:         "test context",
				Descriptors:     map[string]string{"key1": "value1"},
				Priority:        1,
				PrivacyLevel:    2,
				Transformations: []string{"transformation1"},
				Metadata: Metadata{
					Gps: GpsInfo{
						Latitude:     0,
						Longitude:    0,
						Date:         parsedTime,
						Time:         0,
						Altitude:     0,
						AltitudeRef:  false,
						LatitudeRef:  false,
						LongitudeRef: false,
					},
					OriginalDateTime: parsedTime,
					ModifiedDateTime: parsedTime,
					CreationDateTime: parsedTime,
				},
			},
			want:    `{"size":12,"id":"12345","context":"test context","descriptors":{"key1":"value1"},"priority":1,"privacy_level":2,"transformations":["transformation1"],"metadata":{"Gps":{"Latitude":0,"Longitude":0,"Date":"0001-01-01T00:00:00Z","Time":0,"Altitude":0,"AltitudeRef":false,"LatitudeRef":false,"LongitudeRef":false},"OriginalDateTime":"0001-01-01T00:00:00Z","ModifiedDateTime":"0001-01-01T00:00:00Z","CreationDateTime":"0001-01-01T00:00:00Z"}}`,
			wantErr: false,
		},
		{
			name:    "Empty DoxxierPart",
			dp:      &DoxxierPart{},
			want:    `{"size":0,"id":"","context":"","descriptors":null,"priority":0,"privacy_level":0,"transformations":null,"metadata":{"Gps":{"Latitude":0,"Longitude":0,"Date":"0001-01-01T00:00:00Z","Time":0,"Altitude":0,"AltitudeRef":false,"LatitudeRef":false,"LongitudeRef":false},"OriginalDateTime":"0001-01-01T00:00:00Z","ModifiedDateTime":"0001-01-01T00:00:00Z","CreationDateTime":"0001-01-01T00:00:00Z"}}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.dp.ToJson()
			if (err != nil) != tt.wantErr {
				t.Errorf("ToJson() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ToJson() = %v, want %v", got, tt.want)
			}
		})
	}
}
