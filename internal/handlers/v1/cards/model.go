package cards

import (
	"github.com/go-funcards/funapi/internal/gin/httputil"
	"github.com/go-funcards/funapi/proto/card_service/v1"
	"github.com/go-funcards/slice"
	"time"
)

type Attachment struct {
	AttachmentID string `json:"attachment_id"`
	Type         string `json:"type"`
}

type Card struct {
	CardID      string       `json:"card_id"`
	OwnerID     string       `json:"owner_id"`
	BoardID     string       `json:"board_id"`
	CategoryID  string       `json:"category_id"`
	Name        string       `json:"name"`
	Type        string       `json:"type"`
	Content     string       `json:"content"`
	Position    int32        `json:"position"`
	CreatedAt   time.Time    `json:"created_at"`
	Tags        []string     `json:"tags"`
	Attachments []Attachment `json:"attachments"`
}

type PageResponse struct {
	httputil.PageResponse
	Data []Card `json:"data"`
}

type PageRequest struct {
	httputil.PageRequest
	BoardID string `json:"-" form:"board_id" validate:"required,uuid4"`
}

func (dto PageRequest) toRead() *v1.CardsRequest {
	return &v1.CardsRequest{
		PageIndex: dto.Index,
		PageSize:  dto.Size,
		BoardIds:  []string{dto.BoardID},
	}
}

type CreateCardDTO struct {
	OwnerID    string   `json:"-" ctx:"user_id" validate:"required,uuid4"`
	BoardID    string   `json:"board_id" validate:"required,uuid4" format:"uuid"`
	CategoryID string   `json:"category_id" validate:"required,uuid4" format:"uuid"`
	Name       string   `json:"name" validate:"required,max=1000"`
	Type       string   `json:"type" validate:"required,oneof=UNK_CARD TEXT"`
	Content    string   `json:"content,omitempty" validate:"omitempty,max=10000"`
	Position   int32    `json:"position"`
	Tags       []string `json:"tags,omitempty" validate:"omitempty,dive,uuid4"`
}

func (dto CreateCardDTO) toCreate(id string) *v1.CreateCardRequest {
	return &v1.CreateCardRequest{
		CardId:     id,
		OwnerId:    dto.OwnerID,
		BoardId:    dto.BoardID,
		CategoryId: dto.CategoryID,
		Name:       dto.Name,
		Type:       v1.CardType(v1.CardType_value[dto.Type]),
		Content:    dto.Content,
		Position:   dto.Position,
		Tags:       dto.Tags,
	}
}

type UpdateCardDTO struct {
	CardID     string   `json:"card_id" uri:"card_id" validate:"required,uuid4" format:"uuid"`
	BoardID    string   `json:"board_id" validate:"omitempty,uuid4" format:"uuid"`
	CategoryID string   `json:"category_id,omitempty" validate:"omitempty,uuid4" format:"uuid"`
	Name       string   `json:"name,omitempty" validate:"omitempty,max=1000"`
	Content    string   `json:"content,omitempty" validate:"omitempty,max=10000"`
	Position   int32    `json:"position,omitempty"`
	Tags       []string `json:"tags,omitempty" validate:"omitempty,dive,uuid4"`
}

func (dto UpdateCardDTO) toUpdate() *v1.UpdateCardRequest {
	return &v1.UpdateCardRequest{
		CardId:     dto.CardID,
		BoardId:    dto.BoardID,
		CategoryId: dto.CategoryID,
		Name:       dto.Name,
		Content:    dto.Content,
		Position:   dto.Position,
		Tags:       dto.Tags,
	}
}

type UpdateManyCardsDTO struct {
	Data []UpdateCardDTO `json:"data" validate:"required,min=1,dive"`
}

func (dto UpdateManyCardsDTO) toUpdateMany() *v1.UpdateManyCardsRequest {
	return &v1.UpdateManyCardsRequest{
		Cards: slice.Map(dto.Data, func(item UpdateCardDTO) *v1.UpdateCardRequest {
			return item.toUpdate()
		}),
	}
}

type ReadCardDTO struct {
	CardID string `json:"-" uri:"card_id" validate:"required,uuid4"`
}

type DeleteCardDTO struct {
	CardID string `json:"-" uri:"card_id" validate:"required,uuid4"`
}

func (dto DeleteCardDTO) toDelete() *v1.DeleteCardRequest {
	return &v1.DeleteCardRequest{
		CardId: dto.CardID,
	}
}

func PageReq() PageRequest {
	return PageRequest{PageRequest: httputil.PageRequest{Size: 1}}
}

func PageResp(response *v1.CardsResponse, req PageRequest) PageResponse {
	return PageResponse{
		PageResponse: req.ToPageResponse(response.GetTotal()),
		Data: slice.Map(response.GetCards(), func(b *v1.CardsResponse_Card) Card {
			return CreateCard(b)
		}),
	}
}

func CreateCard(response *v1.CardsResponse_Card) Card {
	return Card{
		CardID:     response.GetCardId(),
		OwnerID:    response.GetOwnerId(),
		BoardID:    response.GetBoardId(),
		CategoryID: response.GetCategoryId(),
		Name:       response.GetName(),
		Type:       response.GetType().String(),
		Content:    response.GetContent(),
		Position:   response.GetPosition(),
		CreatedAt:  response.GetCreatedAt().AsTime(),
		Tags:       slice.Copy(response.GetTags()),
		Attachments: slice.Map(response.GetAttachments(), func(a *v1.CardsResponse_Card_Attachment) Attachment {
			return Attachment{
				AttachmentID: a.GetAttachmentId(),
				Type:         a.GetType().String(),
			}
		}),
	}
}
