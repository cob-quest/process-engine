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
	case "imageBuild":
		return "imageBuilder.toService.imageBuild", "Received trigger to start image build"
	case "challengeCreate":
		return "notification.toService.challengeCreate", "Received trigger to start challenge creation"
	case "imageBuilt":
		return "", "Image build completed"
	default:
		log.Fatalf("Unknown suffix: %s", suffix)
		return "nothing", "nothing"
	}
}
