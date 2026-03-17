package toollearn

import "strings"

type ErrorClass string

const (
AuthFailure    ErrorClass = "auth_failure"
RateLimit      ErrorClass = "rate_limit"
NotFound       ErrorClass = "not_found"
Timeout        ErrorClass = "timeout"
ParseError     ErrorClass = "parse_error"
NetworkError   ErrorClass = "network_error"
PermissionErr  ErrorClass = "permission_error"
UnknownError   ErrorClass = "unknown"
)

func ClassifyError(errMsg string) ErrorClass {
lower := strings.ToLower(errMsg)
switch {
case strings.Contains(lower, "401") || strings.Contains(lower, "unauthorized") || strings.Contains(lower, "auth"):
return AuthFailure
case strings.Contains(lower, "429") || strings.Contains(lower, "rate limit"):
return RateLimit
case strings.Contains(lower, "404") || strings.Contains(lower, "not found"):
return NotFound
case strings.Contains(lower, "timeout") || strings.Contains(lower, "deadline"):
return Timeout
case strings.Contains(lower, "parse") || strings.Contains(lower, "unmarshal") || strings.Contains(lower, "json"):
return ParseError
case strings.Contains(lower, "network") || strings.Contains(lower, "connection"):
return NetworkError
case strings.Contains(lower, "permission") || strings.Contains(lower, "forbidden") || strings.Contains(lower, "403"):
return PermissionErr
}
return UnknownError
}
