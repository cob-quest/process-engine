package util

import (
	"log"
	"strings"
)

func getSuffix(routingKey string) string {
	parts := strings.Split(routingKey, ".")
	log.Println("Suffix is: ", parts[len(parts)-1])
	return parts[len(parts)-1]
}

func DetermineNewRoutingKeyAndEventName(routingKey string) (string, string) {
	suffix := getSuffix(routingKey)
	switch suffix {
	// Trigger - Image Builder
	case "imageBuild":
		return "imageBuilder.toService.imageBuild", "Received trigger to start image build"
	// Image Builder - Process Engine
	case "imageBuilt":
		return "", "Image build completed"
	// Trigger - Assignment
	case "challengeCreate":
		return "assignment.toService.challengeCreate", "Received trigger to start challenge creation"
	// Assignment - Notification
	case "challengeCreated":
		return "notification.toService.emailSend", "Assignment creation completed"
	// Notification - Process Engine
	case "emailSent":
		return "", "Notification sending completed"
	default:
		log.Fatalf("Unknown suffix: %s", suffix)
		return "nothing", "nothing"
	}
}

func DetermineBuildStatus(m map[string]interface{}) bool {
	buildStatus, buildStatusExists := m["buildStatus"]; 
	if buildStatusExists {
		buildStatus = strings.ToLower(buildStatus.(string))
		switch buildStatus {
		case "success":
			return true
		case "failure":
			return false
		default:
			log.Fatalf("Invalid build status: %s", buildStatus)
			return false
		}
	}
	log.Fatalf("Build status does not exist")
	return false
}