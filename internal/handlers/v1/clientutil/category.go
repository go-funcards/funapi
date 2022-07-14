package clientutil

import (
	"context"
	"github.com/go-funcards/funapi/internal/gin/httputil"
	v1Category "github.com/go-funcards/funapi/proto/category_service/v1"
)

func CategoriesRequest(id string) *v1Category.CategoriesRequest {
	return &v1Category.CategoriesRequest{
		PageIndex:   0,
		PageSize:    1,
		CategoryIds: []string{id},
	}
}

func GetCategory(ctx context.Context, client v1Category.CategoryClient, id string) (*v1Category.CategoriesResponse_Category, error) {
	response, err := client.GetCategories(ctx, CategoriesRequest(id))
	if err != nil {
		return nil, err
	}
	if len(response.GetCategories()) != 1 {
		return nil, httputil.ErrNotFound
	}
	return response.GetCategories()[0], nil
}
