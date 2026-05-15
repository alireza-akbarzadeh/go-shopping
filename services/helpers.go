package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
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

func marshalStrings(arr []string) datatypes.JSON {
	if arr == nil {
		arr = []string{}
	}
	b, _ := json.Marshal(arr)
	return datatypes.JSON(b)
}

func isDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
