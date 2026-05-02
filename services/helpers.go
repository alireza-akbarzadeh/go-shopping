package services

import (
	"fmt"
	"strings"
	"time"
)

func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	return slug
}

func generateOrderNumber(userID uint) string {
	return fmt.Sprintf("ORD-%d-%d", userID, time.Now().UnixNano())
}
