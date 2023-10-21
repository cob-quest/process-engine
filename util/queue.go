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

func DetermineNewRoutingKey(routingKey string) string {
	suffix := getSuffix(routingKey)
	switch suffix {
	case "imageBuild":
		return "imageBuilder.toService.imageBuild"
	case "challengeCreate":
		return "notification.toService.challengeCreate"
	default:
		log.Fatalf("Unknown suffix: %s", suffix)
		return "nothing"
	}
}
