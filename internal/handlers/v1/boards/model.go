package boards

import (
	"github.com/go-funcards/funapi/internal/gin/httputil"
	v1Authz "github.com/go-funcards/funapi/proto/authz_service/v1"
	"github.com/go-funcards/funapi/proto/board_service/v1"
	"github.com/go-funcards/slice"
	"time"
)

type Member struct {
	MemberID string   `json:"member_id"`
	Roles    []string `json:"roles"`
}

type Board struct {
	BoardID     string    `json:"board_id"`
	OwnerID     string    `json:"owner_id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Data        string    `json:"data"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Members     []Member  `json:"members"`
}

type PageRequest struct {
	httputil.PageRequest
	UserID string `json:"-" ctx:"user_id" validate:"required,uuid4"`
}

func (dto PageRequest) toRead() *v1.BoardsRequest {
	return &v1.BoardsRequest{
		PageIndex: dto.Index,
		PageSize:  dto.Size,
		OwnerIds:  []string{dto.UserID},
		MemberIds: []string{dto.UserID},
	}
}

type PageResponse struct {
	httputil.PageResponse
	Data []Board `json:"data"`
}

type CreateBoardDTO struct {
	OwnerID     string `json:"-" ctx:"user_id" validate:"required,uuid4"`
	Name        string `json:"name" validate:"required,max=150"`
	Type        string `json:"type" validate:"required,oneof=UNK_BOARD CARDS"`
	Data        string `json:"data" validate:"omitempty,max=10000"`
	Description string `json:"description,omitempty" validate:"omitempty,max=1000"`
}

func (dto CreateBoardDTO) toCreate(id string) *v1.CreateBoardRequest {
	return &v1.CreateBoardRequest{
		BoardId:     id,
		OwnerId:     dto.OwnerID,
		Name:        dto.Name,
		Data:        dto.Data,
		Description: dto.Description,
		Type:        v1.BoardType(v1.BoardType_value[dto.Type]),
	}
}

type UpdateBoardDTO struct {
	BoardID     string `json:"-" uri:"board_id" validate:"required,uuid4"`
	Name        string `json:"name,omitempty" validate:"omitempty,max=150"`
	Data        string `json:"data" validate:"omitempty,max=10000"`
	Description string `json:"description,omitempty" validate:"omitempty,max=1000"`
}

func (dto UpdateBoardDTO) toUpdate() *v1.UpdateBoardRequest {
	return &v1.UpdateBoardRequest{
		BoardId:     dto.BoardID,
		Name:        dto.Name,
		Data:        dto.Data,
		Description: dto.Description,
	}
}

type ReadBoardDTO struct {
	BoardID string `json:"-" uri:"board_id" validate:"required,uuid4"`
}

type DeleteBoardDTO struct {
	BoardID string `json:"-" uri:"board_id" validate:"required,uuid4"`
}

func (dto DeleteBoardDTO) toDelete() *v1.DeleteBoardRequest {
	return &v1.DeleteBoardRequest{
		BoardId: dto.BoardID,
	}
}

func (dto DeleteBoardDTO) toDeleteRef() *v1Authz.DeleteRefRequest {
	return &v1Authz.DeleteRefRequest{
		RefId: dto.BoardID,
	}
}

func PageReq() PageRequest {
	return PageRequest{PageRequest: httputil.PageRequest{Size: 1}}
}

func PageResp(response *v1.BoardsResponse, req PageRequest) PageResponse {
	return PageResponse{
		PageResponse: req.ToPageResponse(response.GetTotal()),
		Data: slice.Map(response.GetBoards(), func(b *v1.BoardsResponse_Board) Board {
			return CreateBoard(b)
		}),
	}
}

func CreateBoard(response *v1.BoardsResponse_Board) Board {
	return Board{
		BoardID:     response.GetBoardId(),
		OwnerID:     response.GetOwnerId(),
		Name:        response.GetName(),
		Type:        response.GetType().String(),
		Data:        response.GetData(),
		Description: response.GetDescription(),
		CreatedAt:   response.GetCreatedAt().AsTime(),
		Members: slice.Map(response.GetMembers(), func(m *v1.BoardsResponse_Board_Member) Member {
			return Member{
				MemberID: m.GetMemberId(),
				Roles:    slice.Copy(m.GetRoles()),
			}
		}),
	}
}
