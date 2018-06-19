package types

import (
	"encoding/xml"
	"testing"
)

func TestRecording_String(t *testing.T) {
	type fields struct {
		XMLName         xml.Name
		DiskID          string
		EventID         string
		EventTrigger    string
		Locked          string
		RecordingID     string
		RecordingStatus string
		RecordingType   string
		Source          string
		StartTime       string
		StartTimeLocal  string
		StopTime        string
		StopTimeLocal   string
		Video           *Video
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Simple recording",
			fields: fields{
				StartTimeLocal: "2018-06-08T16:38:45.761378-7:00",
				StopTimeLocal:  "2018-06-10T16:38:45.761378-7:00",
			},
			want: "2018-06-08T16:38:45.761378-7:00 - 2018-06-10T16:38:45.761378-7:00",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Recording{
				XMLName:         tt.fields.XMLName,
				DiskID:          tt.fields.DiskID,
				EventID:         tt.fields.EventID,
				EventTrigger:    tt.fields.EventTrigger,
				Locked:          tt.fields.Locked,
				RecordingID:     tt.fields.RecordingID,
				RecordingStatus: tt.fields.RecordingStatus,
				RecordingType:   tt.fields.RecordingType,
				Source:          tt.fields.Source,
				StartTime:       tt.fields.StartTime,
				StartTimeLocal:  tt.fields.StartTimeLocal,
				StopTime:        tt.fields.StopTime,
				StopTimeLocal:   tt.fields.StopTimeLocal,
				Video:           tt.fields.Video,
			}
			if got := r.String(); got != tt.want {
				t.Errorf("Recording.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecordings_String(t *testing.T) {
	type fields struct {
		XMLName                 xml.Name
		NumberOfRecordings      string
		TotalNumberOfRecordings string
		Recording               []*Recording
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recordings := &Recordings{
				XMLName:                 tt.fields.XMLName,
				NumberOfRecordings:      tt.fields.NumberOfRecordings,
				TotalNumberOfRecordings: tt.fields.TotalNumberOfRecordings,
				Recording:               tt.fields.Recording,
			}
			if got := recordings.String(); got != tt.want {
				t.Errorf("Recordings.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
