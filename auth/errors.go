package auth

import (
	core "github.com/jacob-ebey/graphql-core"
)

var NotAuthenticatedError = &core.WrappedError{
	Message: "Not authenticated.",
	Code:    "NOT_AUTHENTICATED",
}

var NotAuthorizedError = &core.WrappedError{
	Message: "Not authorized.",
	Code:    "NOT_AUTHORIZED",
}
