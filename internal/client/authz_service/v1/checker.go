package v1

import (
	"context"
	"encoding/json"
	"github.com/go-funcards/funapi/proto/authz_service/v1"
	"github.com/go-funcards/grpc-pool"
	"google.golang.org/grpc"
)

var _ v1.AuthorizationCheckerClient = (*CheckerService)(nil)

type CheckerService struct {
	Pool *grpcpool.Pool
}

func (s *CheckerService) IsGranted(ctx context.Context, in *v1.IsGrantedRequest, opts ...grpc.CallOption) (*v1.Granted, error) {
	conn, err := s.Pool.Get(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return v1.NewAuthorizationCheckerClient(conn).IsGranted(ctx, in, opts...)
}

type paramObject struct {
	Name  string `json:"name"`
	Owner string `json:"owner,omitempty"`
	Ref   string `json:"ref,omitempty"`
}

func NewObject(name, owner, ref string) string {
	data, err := json.Marshal(paramObject{
		Name:  name,
		Owner: owner,
		Ref:   ref,
	})
	if err != nil {
		return ""
	}
	return string(data)
}

func Is(sub, obj, act string) *v1.IsGrantedRequest {
	return &v1.IsGrantedRequest{
		Params: []string{sub, obj, act},
	}
}

func IsGranted(ctx context.Context, client v1.AuthorizationCheckerClient, sub, obj, act string) bool {
	result, err := client.IsGranted(ctx, Is(sub, obj, act))
	if err != nil {
		return false
	}
	return result.GetYes()
}
