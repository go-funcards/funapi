package clientutil

import (
	"context"
	"github.com/go-funcards/funapi/internal/gin/httputil"
	v1Tag "github.com/go-funcards/funapi/proto/tag_service/v1"
)

func TagsRequest(id string) *v1Tag.TagsRequest {
	return &v1Tag.TagsRequest{
		PageIndex: 0,
		PageSize:  1,
		TagIds:    []string{id},
	}
}

func GetTag(ctx context.Context, client v1Tag.TagClient, id string) (*v1Tag.TagsResponse_Tag, error) {
	response, err := client.GetTags(ctx, TagsRequest(id))
	if err != nil {
		return nil, err
	}
	if len(response.GetTags()) != 1 {
		return nil, httputil.ErrNotFound
	}
	return response.GetTags()[0], nil
}
