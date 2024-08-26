package access

import (
	"context"
	"strings"

	errorAccess "github.com/BelyaevEI/microservices_auth/internal/errors"
	"github.com/BelyaevEI/microservices_auth/internal/model"
	"github.com/BelyaevEI/microservices_auth/internal/utils"

	"google.golang.org/grpc/metadata"
)

func (s *serv) Check(ctx context.Context, endpointAddress string) error {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errorAccess.ErrMetadataNotProvided
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return errorAccess.ErrAuthorizationHeader
	}

	if !strings.HasPrefix(authHeader[0], s.authPrefix) {
		return errorAccess.ErrInvalidHeaderFormat
	}

	accessToken := strings.TrimPrefix(authHeader[0], s.authPrefix)

	claims, err := utils.VerifyToken(accessToken, []byte(s.accessTokenSecretKey))
	if err != nil {
		return errorAccess.ErrInvalidAccessToken
	}

	accessibleMap, err := s.accessibleRoles(ctx)
	if err != nil {
		return errorAccess.ErrGetAccessibleRoles
	}

	roles, ok := accessibleMap[endpointAddress]
	if !ok {
		return nil
	}

	if _, ok = roles[claims.Role]; ok {
		return nil
	}

	return errorAccess.ErrAccessDenied
}

// Returns a map with the endpoint address and the role that has access to it
func (s *serv) accessibleRoles(_ context.Context) (map[string]map[model.Role]struct{}, error) {
	return accessibleRoles, nil
}
