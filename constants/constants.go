// Package constants holds shared immutable configuration values, error messages,
// and HTTP status codes used across the shopping platform.
package constants

import (
	"errors"
	"time"
)

type ContextKey string

// ==================== Defaults ====================
const (
	DefaultProtectedAPIPort       int           = 8080
	DefaultPublicAPIPort          int           = 8081
	DefaultHiddenAPIPort          int           = 8079
	DefaultHost                   string        = "0.0.0.0"
	DefaultDevHost                string        = "127.0.0.1"
	DefaultLogLevel               string        = "warn"
	DefaultDevLogLevel            string        = "debug"
	DefaultCORSAllowOrigins       string        = "*"
	DefaultDBPlatform             string        = DBPlatformSQLite
	DefaultDBTimezone             string        = DBTimezoneUTC
	DefaultDBSSLMode              string        = DBSSLModeDisabled
	DefaultSQLiteDBName           string        = "sqlite.db"
	DefaultLoggerTimestampFormat  string        = "2006-01-02 15:04:05.00000"
	DefaultRequestTimeoutDuration time.Duration = 60 * time.Second
	DefaultWatcherSleepInterval   time.Duration = 5 * time.Second
	DefaultGzipLevel              int           = 5
	DefaultLimit                  int           = 20
	MaxLimit                      int           = 100
	MinLimit                      int           = 1
	MinOffset                     int           = 0
)

// ==================== Feature flags ====================
const (
	FeatureService   = "service"
	FeatureOryKratos = "ory_kratos"
	FeatureOryKeto   = "ory_keto"
	FeatureDatabase  = "database"
	FeatureCORS      = "cors"
	FeatureGzip      = "gzip"
	FeatureRedis     = "redis"
)

// ==================== Generic words ====================
const (
	WordDatabase       = "database"
	WordDatabaseServer = "database_server"
	WordInternalCode   = "internalCode"
	WordServiceCode    = "serviceCode"
)

// ==================== App‑specific names ====================
const (
	NameHealthPath      = "/alive"
	NameHealthReadyPath = "/ready"
	NameCORSConfig      = "CORSAllowOrigins"
	NameTimeoutDuration = "RequestTimeoutDuration"
	NameCmdDBMigrate    = "migrate"
	NameCmdDBRollback   = "rollback"
	NameCmdDBSeed       = "seed"
)

// ==================== Database platforms & settings ====================
const (
	DBPlatformPostgres  = "postgres"
	DBPlatformMySQL     = "mysql"
	DBPlatformSQLite    = "sqlite"
	DBSSLModeEnabled    = "require"
	DBSSLModeDisabled   = "disable"
	DBTimezoneUTC       = "Etc/GMT"
	DBTimezoneMelbourne = "Australia/Melbourne"
)

// ==================== Headers & context keys ====================
const (
	HeaderContentType     = "Content-Type"
	HeaderContentTypeJSON = "application/json; charset=utf-8"
	HeaderAuthorization   = "Authorization"
	HeaderAuthBearerWord  = "Bearer"
	HeaderKratosCookie    = "ory_kratos_session"
)

var (
	RequestIDKey ContextKey = "request_id"
	UserIDKey    ContextKey = "uid"
)

// ==================== Output messages ====================
const (
	MsgServerShuttingDown       = "server is shutting down"
	MsgNotAcceptable            = "not acceptable"
	MsgMissingAcceptHeader      = "unknown accept format"
	MsgMissingContentTypeHeader = "unknown content format"
	MsgSuccess                  = "success"
	MsgError                    = "error"
	MsgValidationError          = "validation error"
	MsgRouteNotFound            = "route not found"
	MsgRecordNotFound           = "record not found"
	MsgDependencyNotFound       = "dependency not found"
	MsgSessionNotFound          = "session not found"
	MsgAccessIDsNotFound        = "access ids not found or not readable"
	MsgNotAuthorized            = "not authorized"
	MsgIDNotReadable            = "ID not found or not readable"
	MsgUnableToBindBody         = "error binding body"
	MsgForbidden                = "forbidden"
	MsgUnknownDBPlatform        = "unknown database platform"
	MsgInternalServer           = "internal server error"
	MsgShutdownServerCompleted  = "Graceful shutdown complete."
	MsgUserRegisterSuccess      = "User registered successfully"
	ErrShipmentNotFound         = "shipment not found"
)

// ==================== Log levels ====================
var (
	LogLevels = []string{"debug", "info", "warn", "error", "fatal", "panic"}
)

// ==================== Predefined errors ====================
var (
	ErrNotAuthorized     = errors.New(MsgNotAuthorized)
	ErrSessionNotFound   = errors.New(MsgSessionNotFound)
	ErrIDNotFound        = errors.New(MsgIDNotReadable)
	ErrAccessIDsNotFound = errors.New(MsgAccessIDsNotFound)
	ErrBindingBody       = errors.New(MsgUnableToBindBody)
	ErrUnknownDBPlatform = errors.New(MsgUnknownDBPlatform)
	ErrInternalServer    = errors.New(MsgInternalServer)
)
