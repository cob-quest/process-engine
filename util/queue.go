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

func DetermineNewRoutingKeyAndEventName(routingKey string, queueName string) (newQueueName string, eventName string) {
	suffix := getSuffix(routingKey)
	switch suffix {
	// Platform - Image Builder
	case "imageCreate":
		return "imageBuilder.toService.imageCreate", "imageCreate"
	// Image Builder - Process Engine
	case "imageCreated":
		return "", "imageCreate"
	// Platform - Challenge
	case "challengeCreate":
		return "challenge.toService.challengeCreate", "challengeCreate"
	// Challenge - Notification
	case "challengeCreated":
		return "notification.toService.emailSend", "challengeCreate"
	// Notification - Process Engine
	case "emailSent":
		return "", "emailSend"
	// Platform - Challenge
	case "challengeStart":
		return "challenge.toService.challengeStart", "challengeStart"
	case "challengeStarted":
		return "", "challengeStart"
	default:
		log.Printf("Unknown suffix: %s", suffix)
		return "error", "error"
	}
}

func DetermineBuildStatus(m map[string]interface{}) string {
	buildStatus, buildStatusExists := m["buildStatus"]; 
	if buildStatusExists {
		buildStatus = strings.ToLower(buildStatus.(string))
		switch buildStatus {
		case "success":
			return ""
		case "failure":
			return ""
		default:
			log.Fatalf("Invalid build status: %s", buildStatus)
			return ""
		}
	}
	log.Fatalf("Build status does not exist")
	return ""
}