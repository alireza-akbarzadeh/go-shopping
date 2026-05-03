# Junior Developer Upgrade Guide

## What to look for first

### 1. Duplicate patterns
- Many controllers repeat the same validation and error-handling blocks.
- Look for repeated `ShouldBindJSON`, `ctrl.validate.Struct(req)`, and `errors.As(err, &appErr)`.
- Extract these into reusable helpers or middleware.

### 2. Hard-coded values
- Status values and role strings are used directly in many places.
- Use constants such as `const CartStatusActive = "active"` and `const OrderStatusPending = "pending"`.

### 3. Error handling style
- The app uses `AppError`, but not consistently.
- Always return a typed error from the service layer and map it once in controller or middleware.

### 4. Controller responsibilities
- Controllers should only handle HTTP details, not business rules.
- Business logic belongs in services; controllers should only parse requests and return responses.

### 5. Input validation
- Right now validation is repeated in each endpoint.
- Use request DTOs with struct tags and a validation helper so that validation error responses are consistent.

## Practical improvements

### Improve request validation
- Add a helper like `BindAndValidate(c, &req)`.
- Let it return a standard validation response when parsing fails.

### Centralize API responses
- Use a single response helper package for success, created, error, and validation responses.
- Avoid manually writing `c.JSON` in every controller.

### Clean up filtering and pagination
- Create a helper for pagination defaults.
- Create a filter builder for query parameters so controller logic is small.

### Use constants for domain status
- Example:
  - `const CartStatusActive = "active"`
  - `const CartStatusConverted = "converted"`
  - `const OrderStatusPending = "pending"`
- This reduces copy/paste bugs and makes future refactor easier.

### Add tests early
- Start with service unit tests where business rules are isolated.
- Then add simple endpoint tests for the `cart` and `orders` APIs.

## Recommended next PRs

1. **Refactor validation and request binding**
   - Create a package or helper in `utils`/`controllers`.
   - Use it in at least two controllers.

2. **Standardize app errors**
   - Add an error mapper function.
   - Remove repeated `errors.As` logic from controllers.

3. **Add one unit test suite**
   - Focus on `services/cart_service.go` or `services/orders_service.go`.
   - Use a mock DB or in-memory test database.

4. **Document the architecture**
   - Add a simple `docs/architecture.md` describing how controllers, services, models, and routes connect.

## Why this matters
- Less duplicated code means fewer bugs.
- Centralized validation and error handling makes the API easier to maintain.
- Proper tests give confidence when changing checkout, cart, or order logic.
- Cleaner code helps the project grow from a junior prototype to a production backend.
