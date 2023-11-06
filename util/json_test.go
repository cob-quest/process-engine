package util

import (
	"reflect"
	"testing"

	"gitlab.com/cs302-2023/g3-team8/project/process-engine/models"
)

func TestUnmarshalJson(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name:    "Valid JSON",
			input:   []byte(`{"key": "value"}`),
			want:    map[string]interface{}{"key": "value"},
			wantErr: false,
		},
		{
			name: "Valid JSON 2",
			input: []byte(`{
				"corId": "a123b456-7890-1abc-2def-3ghi456jklm",
				"creatorName": "cs302",
				"challengeName": "challenge_0",
				"imageName": "image_0",
				"duration": 60,
				"participants": [
				  "challenger85@google.com",
				  "challenger74@google.com",
				  "challenger100@google.com"
				]
			  }`),
			want: map[string]interface{}{
				"corId":         "a123b456-7890-1abc-2def-3ghi456jklm",
				"creatorName":   "cs302",
				"challengeName": "challenge_0",
				"imageName":     "image_0",
				"duration":      float64(60),
				"participants": []interface{}{
					"challenger85@google.com",
					"challenger74@google.com",
					"challenger100@google.com",
				},
			},
			wantErr: false,
		},
		{
			name:    "Invalid JSON",
			input:   []byte(`{"key": "value"`),
			want:    nil,
			wantErr: true,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && !tt.wantErr {
					t.Errorf("UnmarshalJson() caused a panic for input %s", tt.input)
				}
			}()

			got := UnmarshalJson(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnmarshalJson() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapJsonToProcessEngine(t *testing.T) {
	tests := []struct {
		name    string
		input   map[string]interface{}
		want    *models.ProcessEngine
		wantErr bool
	}{
		{
			name: "Valid mapping",
			input: map[string]interface{}{
				"corId":         "a123b456-7890-1abc-2def-3ghi456jklm",
				"event":         "event1",
				"eventStatus":   "status1",
				"creatorName":   "cs302",
				"challengeName": "challenge_0",
				"imageName":     "image_0",
			},
			want: &models.ProcessEngine{
				CorId:         ptr("a123b456-7890-1abc-2def-3ghi456jklm"),
				Event:         ptr("event1"),
				EventStatus:   ptr("status1"),
				CreatorName:   ptr("cs302"),
				ChallengeName: ptr("challenge_0"),
				ImageName:     ptr("image_0"),
			},
			wantErr: false,
		},
		{
			name: "Invalid mapping",
			input: map[string]interface{}{
				"corId":       "a123b456-7890-1abc-2def-3ghi456jklm",
				"event":       "event1",
				"eventStatus": "status1",
				"creatorName": "cs302",
			},
			want: &models.ProcessEngine{
				CorId:       ptr("a123b456-7890-1abc-2def-3ghi456jklm"),
				Event:       ptr("event2"),
				EventStatus: ptr("status1"),
				CreatorName: ptr("cs302"),
			},
			wantErr: true,
		},
		{
			name: "Error mapping",
			input: map[string]interface{}{
				"corId":       10,
			},
			want: nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if !tt.wantErr {
						t.Errorf("MapJsonToProcessEngine() caused a panic for input %v", tt.input)
					}
				}
			}()

			got := MapJsonToProcessEngine(tt.input)
			if !reflect.DeepEqual(got, tt.want) && !tt.wantErr {
				t.Errorf("MapJsonToProcessEngine() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// ptr is a helper function to easily create pointers to strings for the tests
func ptr(s string) *string {
	return &s
}
