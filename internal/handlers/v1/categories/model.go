package categories

import (
	"github.com/go-funcards/funapi/internal/gin/httputil"
	"github.com/go-funcards/funapi/proto/category_service/v1"
	"github.com/go-funcards/slice"
	"time"
)

type Category struct {
	CategoryID string    `json:"category_id"`
	OwnerID    string    `json:"owner_id"`
	BoardID    string    `json:"board_id"`
	Name       string    `json:"name"`
	Position   int32     `json:"position"`
	CreatedAt  time.Time `json:"created_at"`
}

type PageResponse struct {
	httputil.PageResponse
	Data []Category `json:"data"`
}

type PageRequest struct {
	httputil.PageRequest
	BoardID string `json:"-" form:"board_id" validate:"required,uuid4"`
}

func (dto PageRequest) toRead() *v1.CategoriesRequest {
	return &v1.CategoriesRequest{
		PageIndex: dto.Index,
		PageSize:  dto.Size,
		BoardIds:  []string{dto.BoardID},
	}
}

type CreateCategoryDTO struct {
	OwnerID  string `json:"-" ctx:"user_id" validate:"required,uuid4"`
	BoardID  string `json:"board_id" validate:"required,uuid4" format:"uuid"`
	Name     string `json:"name" validate:"omitempty,max=150"`
	Position int32  `json:"position"`
}

func (dto CreateCategoryDTO) toCreate(id string) *v1.CreateCategoryRequest {
	return &v1.CreateCategoryRequest{
		CategoryId: id,
		OwnerId:    dto.OwnerID,
		BoardId:    dto.BoardID,
		Name:       dto.Name,
		Position:   dto.Position,
	}
}

type UpdateCategoryDTO struct {
	CategoryID string `json:"category_id" uri:"category_id" validate:"required,uuid4" format:"uuid"`
	BoardID    string `json:"board_id" validate:"omitempty,uuid4" format:"uuid"`
	Name       string `json:"name,omitempty" validate:"omitempty,max=150"`
	Position   int32  `json:"position,omitempty"`
}

func (dto UpdateCategoryDTO) toUpdate() *v1.UpdateCategoryRequest {
	return &v1.UpdateCategoryRequest{
		CategoryId: dto.CategoryID,
		BoardId:    dto.BoardID,
		Name:       dto.Name,
		Position:   dto.Position,
	}
}

type UpdateManyCategoriesDTO struct {
	Data []UpdateCategoryDTO `json:"data" validate:"required,min=1,dive"`
}

func (dto UpdateManyCategoriesDTO) toUpdateMany() *v1.UpdateManyCategoriesRequest {
	return &v1.UpdateManyCategoriesRequest{
		Categories: slice.Map(dto.Data, func(item UpdateCategoryDTO) *v1.UpdateCategoryRequest {
			return item.toUpdate()
		}),
	}
}

type ReadCategoryDTO struct {
	CategoryID string `json:"-" uri:"category_id" validate:"required,uuid4"`
}

type DeleteCategoryDTO struct {
	CategoryID string `json:"-" uri:"category_id" validate:"required,uuid4"`
}

func (dto DeleteCategoryDTO) toDelete() *v1.DeleteCategoryRequest {
	return &v1.DeleteCategoryRequest{
		CategoryId: dto.CategoryID,
	}
}

func PageReq() PageRequest {
	return PageRequest{PageRequest: httputil.PageRequest{Size: 1}}
}

func PageResp(response *v1.CategoriesResponse, req PageRequest) PageResponse {
	return PageResponse{
		PageResponse: req.ToPageResponse(response.GetTotal()),
		Data: slice.Map(response.GetCategories(), func(b *v1.CategoriesResponse_Category) Category {
			return CreateCategory(b)
		}),
	}
}

func CreateCategory(response *v1.CategoriesResponse_Category) Category {
	return Category{
		CategoryID: response.GetCategoryId(),
		OwnerID:    response.GetOwnerId(),
		BoardID:    response.GetBoardId(),
		Name:       response.GetName(),
		Position:   response.GetPosition(),
		CreatedAt:  response.GetCreatedAt().AsTime(),
	}
}
