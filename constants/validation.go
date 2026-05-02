package constants

import "errors"

// ==================== AUTHENTICATION & AUTHORIZATION ====================
const (
	ErrEmailAlreadyExists   = "email already registered"
	ErrInvalidCredentials   = "invalid email or password"
	ErrAccountDeactivated   = "account is deactivated"
	ErrUserNotFound         = "user not found"
	ErrInvalidToken         = "invalid or expired token"
	ErrMissingAuthHeader    = "authorization header is missing"
	ErrInvalidAuthFormat    = "authorization header must be: Bearer {token}"
	ErrUnauthorized         = "unauthorized access"
	ErrForbidden            = "access denied – insufficient permissions"
	ErrTokenExpired         = "token has expired"
	ErrTokenInvalid         = "token is malformed or invalid"
	ErrPasswordResetFailed  = "password reset failed"
	ErrOldPasswordIncorrect = "current password is incorrect"
	ErrWeakPassword         = "password is too weak – use at least 8 chars with mix"
)

// ==================== USER PROFILE ====================
const (
	ErrFetchProfile       = "failed to fetch user profile"
	ErrUpdateProfile      = "failed to update user profile"
	ErrInvalidPhoneFormat = "phone number must be in E.164 format (e.g., +1234567890)"
	ErrEmailNotVerified   = "email not verified – please verify your email first"
	ErrVerificationFailed = "email verification failed"
)

// ==================== PRODUCTS ====================
const (
	ErrProductNotFound     = "product not found"
	ErrProductOutOfStock   = "product is out of stock"
	ErrProductInactive     = "product is no longer available"
	ErrInvalidProductPrice = "invalid product price (must be positive)"
	ErrInvalidProductStock = "invalid stock quantity (must be >= 0)"
	ErrDuplicateProductSKU = "product with this SKU already exists"
	ErrCategoryNotFound    = "category not found"
	ErrInvalidProductData  = "invalid product data"
)

// ==================== ORDERS ====================
const (
	ErrOrderNotFound          = "order not found"
	ErrInvalidOrderStatus     = "invalid order status transition"
	ErrOrderAlreadyPaid       = "order already paid – cannot modify"
	ErrOrderCancelled         = "order is cancelled"
	ErrOrderExpired           = "order has expired"
	ErrEmptyCart              = "cannot place order with empty cart"
	ErrInvalidShippingAddress = "shipping address is invalid"
)

// ==================== CART ====================
const (
	ErrCartItemNotFound     = "item not found in cart"
	ErrCartEmpty            = "cart is empty"
	ErrInvalidQuantity      = "invalid quantity (must be at least 1)"
	ErrQuantityExceedsStock = "requested quantity exceeds available stock"
)

// ==================== PAYMENTS ====================
const (
	ErrPaymentFailed           = "payment processing failed"
	ErrPaymentMethodInvalid    = "invalid payment method"
	ErrPaymentAlreadyProcessed = "payment already processed for this order"
	ErrInsufficientFunds       = "insufficient funds"
	ErrPaymentTimeout          = "payment gateway timeout – please try again"
)

// ==================== SHIPPING ====================
const (
	ErrShippingNotAvailable  = "shipping not available for this address"
	ErrInvalidTrackingNumber = "invalid tracking number"
	ErrShippingDelay         = "shipping carrier reported delay"
)

// ==================== RATE LIMITING / REQUEST ====================
const (
	ErrTooManyRequests   = "too many requests – please slow down"
	ErrRateLimitExceeded = "rate limit exceeded – try again later"
)

// ==================== GENERAL / DATABASE ====================
const (
	ErrValidationFailed      = "validation failed – check your input"
	ErrInvalidRequestFormat  = "request body is malformed"
	ErrConflict              = "resource conflict – duplicate or state mismatch"
	ErrNotFound              = "resource not found"
	ErrBadRequest            = "bad request – missing or invalid parameters"
	ErrDatabaseOperation     = "database operation failed"
	ErrDependencyUnavailable = "external service is temporarily unavailable"
	MsgRegistrationFailed    = "registration failed"
	MsgLoginFailed           = "login failed"
	ErrorMissingAuthHeader   = "authorization header is missing"
	ErrorInvalidAuthFormat   = "authorization header must be in the format: Bearer {token}"
	ErrorUnauthorized        = "You don't have a required permission"
	ErrorForbidden           = "you don't have a required role"
	ErrorInvalidToken        = "invalid token error"
)

// ==================== SUCCESS MESSAGES (remaining as they are) ====================
const (
	MsgLogoutSuccess = "logout successful"
	// Add others if needed (already defined in earlier step)

)

var (
	// ErrTaskIsNotAFuncError is the error panicked when a task passed to `Job.Do`
	ErrTaskIsNotAFuncError = errors.New("the `task` your a scheduling must be of type func")

	// ErrMissmatchedTaskParams is the error panicked when someone passes too many or too few params to `Job.Do`
	ErrMissmatchedTaskParams = errors.New("the `task` your a scheduling must be of type func")

	// ErrJobIsNotInitialized is the error panicked when a job is scheduled that was not initialized
	ErrJobIsNotInitialized = errors.New("this job was not intialized")

	// ErrIncorrectTimeFormat is the error panicked when `At` is passed an incorrect time
	ErrIncorrectTimeFormat = errors.New("the time format is incorrect")

	// ErrIntervalNotValid error panicked when the interval is not valid
	ErrIntervalNotValid = errors.New("the interval must be greater than 0")
)
