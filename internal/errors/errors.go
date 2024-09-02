package errors

import "errors"

var (
	// ErrInvalidRefreshToken - invalid refresh token
	ErrInvalidRefreshToken = errors.New("invalid refresh token")

	// ErrInvalidAccessToken - invalid access token
	ErrInvalidAccessToken = errors.New("invalid access token")

	// ErrWrongPassword - wrong password
	ErrWrongPassword = errors.New("wrong password")

	// ErrGenerateToken - failed to generate token
	ErrGenerateToken = errors.New("failed to generate token")

	// ErrMetadataNotProvided - metadata not provided
	ErrMetadataNotProvided = errors.New("metadata is not provided")

	// ErrAuthorizationHeader - authorization header is not provided
	ErrAuthorizationHeader = errors.New("authorization header is not provided")

	// ErrInvalidHeaderFormat - invalid authorization header format
	ErrInvalidHeaderFormat = errors.New("invalid authorization header format")

	// ErrGetAccessibleRoles  - failed to get accessible roles
	ErrGetAccessibleRoles = errors.New("failed to get accessible roles")

	// ErrAccessDenied        - access denied
	ErrAccessDenied = errors.New("access denied")
)
