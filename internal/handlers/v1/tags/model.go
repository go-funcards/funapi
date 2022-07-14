package tags

import (
	"github.com/go-funcards/funapi/internal/gin/httputil"
	"github.com/go-funcards/funapi/proto/tag_service/v1"
	"github.com/go-funcards/slice"
	"time"
)

type Tag struct {
	TagID     string    `json:"tag_id"`
	OwnerID   string    `json:"owner_id"`
	BoardID   string    `json:"board_id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}

type PageResponse struct {
	httputil.PageResponse
	Data []Tag `json:"data"`
}

type PageRequest struct {
	httputil.PageRequest
	BoardID string `json:"-" form:"board_id" validate:"required,uuid4"`
}

func (dto PageRequest) toRead() *v1.TagsRequest {
	return &v1.TagsRequest{
		PageIndex: dto.Index,
		PageSize:  dto.Size,
		BoardIds:  []string{dto.BoardID},
	}
}

type CreateTagDTO struct {
	OwnerID string `json:"-" ctx:"user_id" validate:"required,uuid4"`
	BoardID string `json:"board_id" validate:"required,uuid4" format:"uuid"`
	Name    string `json:"name" validate:"omitempty,max=100"`
	Color   string `json:"color" validate:"required,max=50"`
}

func (dto CreateTagDTO) toCreate(id string) *v1.CreateTagRequest {
	return &v1.CreateTagRequest{
		TagId:   id,
		OwnerId: dto.OwnerID,
		BoardId: dto.BoardID,
		Name:    dto.Name,
		Color:   dto.Color,
	}
}

type UpdateTagDTO struct {
	TagID   string `json:"-" uri:"tag_id" validate:"required,uuid4"`
	BoardID string `json:"board_id" validate:"omitempty,uuid4" format:"uuid"`
	Name    string `json:"name,omitempty" validate:"omitempty,max=100"`
	Color   string `json:"color,omitempty" validate:"omitempty,max=50"`
}

func (dto UpdateTagDTO) toUpdate() *v1.UpdateTagRequest {
	return &v1.UpdateTagRequest{
		TagId:   dto.TagID,
		BoardId: dto.BoardID,
		Name:    dto.Name,
		Color:   dto.Color,
	}
}

type ReadTagDTO struct {
	TagID string `json:"-" uri:"tag_id" validate:"required,uuid4"`
}

type DeleteTagDTO struct {
	TagID string `json:"-" uri:"tag_id" validate:"required,uuid4"`
}

func (dto DeleteTagDTO) toDelete() *v1.DeleteTagRequest {
	return &v1.DeleteTagRequest{
		TagId: dto.TagID,
	}
}

func PageReq() PageRequest {
	return PageRequest{PageRequest: httputil.PageRequest{Size: 1}}
}

func PageResp(response *v1.TagsResponse, req PageRequest) PageResponse {
	return PageResponse{
		PageResponse: req.ToPageResponse(response.GetTotal()),
		Data: slice.Map(response.GetTags(), func(b *v1.TagsResponse_Tag) Tag {
			return CreateTag(b)
		}),
	}
}

func CreateTag(response *v1.TagsResponse_Tag) Tag {
	return Tag{
		TagID:     response.GetTagId(),
		OwnerID:   response.GetOwnerId(),
		BoardID:   response.GetBoardId(),
		Name:      response.GetName(),
		Color:     response.GetColor(),
		CreatedAt: response.GetCreatedAt().AsTime(),
	}
}
