package constant

// HTTP Methods
const (
	HTTP_GET     = "GET"
	HTTP_POST    = "POST"
	HTTP_PUT     = "PUT"
	HTTP_DELETE  = "DELETE"
	HTTP_PATCH   = "PATCH"
	HTTP_OPTIONS = "OPTIONS"
)

// Common HTTP Headers
const (
	HEADER_AUTHORIZATION      = "Authorization"
	HEADER_CONTENT_TYPE      = "Content-Type"
	HEADER_ACCEPT            = "Accept"
	HEADER_USER_AGENT        = "User-Agent"
	HEADER_X_REQUEST_ID      = "X-Request-ID"
	HEADER_X_CORRELATION_ID  = "X-Correlation-ID"
)

// Content Types
const (
	CONTENT_TYPE_JSON           = "application/json"
	CONTENT_TYPE_FORM          = "application/x-www-form-urlencoded"
	CONTENT_TYPE_MULTIPART     = "multipart/form-data"
	CONTENT_TYPE_TEXT          = "text/plain"
	CONTENT_TYPE_HTML          = "text/html"
	CONTENT_TYPE_OCTET_STREAM  = "application/octet-stream"
)
