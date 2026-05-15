package services

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gorm.io/datatypes"
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

func marshalStrings(arr []string) (datatypes.JSON, error) {
	if arr == nil {
		arr = []string{}
	}
	b, err := json.Marshal(arr)
	if err != nil {
		return nil, err
	}
	return datatypes.JSON(b), nil
}
