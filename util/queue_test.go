// +build unit

package util

import (
	"testing"
)

func TestGetSuffix(t *testing.T) {
	tests := []struct {
		name       string
		routingKey string
		want       string
	}{
		{"Simple case", "test.imageCreate", "imageCreate"},
		{"Multiple dots", "test.multiple.dots.imageCreated", "imageCreated"},
		{"No dots", "nodots", "nodots"},
		{"Empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if suffix := getSuffix(tt.routingKey); suffix != tt.want {
				t.Errorf("GetSuffix() = %v, want %v", suffix, tt.want)
			}
		})
	}
}

func TestDetermineNewRoutingKeyAndEventName(t *testing.T) {
	tests := []struct {
		name       string
		routingKey string
		queueName  string
		wantQueue  string
		wantEvent  string
	}{
		{"Image Create", "platform.fromService.imageCreate", "", "imageBuilder.toService.imageCreate", "imageCreate"},
		{"Image Created", "imageBuilder.fromService.imageCreated", "", "", "imageCreate"},
		{"Challenge Create", "platform.fromService.challengeCreate", "", "challenge.toService.challengeCreate", "challengeCreate"},
		{"Challenge Created", "challenge.fromService.challengeCreated", "", "notification.toService.emailSend", "challengeCreate"},
		{"Email Sent", "notification.fromService.emailSent", "", "", "emailSend"},
		{"Challenge Start", "platform.fromService.challengeStart", "", "challenge.toService.challengeStart", "challengeStart"},
		{"Challenge Started", "challenge.fromService.challengeStarted", "", "", "challengeStart"},
		{"Unknown Suffix", "unknown.suffix", "", "error", "error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotQueue, gotEvent := DetermineNewRoutingKeyAndEventName(tt.routingKey, tt.queueName)
			if gotQueue != tt.wantQueue || gotEvent != tt.wantEvent {
				t.Errorf("DetermineNewRoutingKeyAndEventName() gotQueue = %v, want %v; gotEvent = %v, want %v",
					gotQueue, tt.wantQueue, gotEvent, tt.wantEvent)
			}
		})
	}
}