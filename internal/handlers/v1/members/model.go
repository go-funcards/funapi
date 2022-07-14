package members

import (
	v1Authz "github.com/go-funcards/funapi/proto/authz_service/v1"
	v1Board "github.com/go-funcards/funapi/proto/board_service/v1"
)

type SaveMemberDTO struct {
	BoardID  string   `json:"-" uri:"board_id" validate:"required,uuid4"`
	MemberID string   `json:"-" uri:"member_id" validate:"required,uuid4"`
	Roles    []string `json:"roles" validate:"required,dive,min=1,max=50"`
}

func (dto SaveMemberDTO) toUpdate() *v1Board.UpdateBoardRequest {
	return &v1Board.UpdateBoardRequest{
		BoardId: dto.BoardID,
		Members: []*v1Board.UpdateBoardRequest_Member{
			{MemberId: dto.MemberID, Roles: dto.Roles, Delete: false},
		},
	}
}

func (dto SaveMemberDTO) toSaveSub() *v1Authz.SaveSubRequest {
	return &v1Authz.SaveSubRequest{
		SubId: dto.MemberID,
		Refs: []*v1Authz.SaveSubRequest_Ref{
			{RefId: dto.BoardID, Roles: dto.Roles, Delete: false},
		},
	}
}

type DeleteMemberDTO struct {
	BoardID  string `json:"-" uri:"board_id" validate:"required,uuid4"`
	MemberID string `json:"-" uri:"member_id" validate:"required,uuid4"`
}

func (dto DeleteMemberDTO) toUpdate() *v1Board.UpdateBoardRequest {
	return &v1Board.UpdateBoardRequest{
		BoardId: dto.BoardID,
		Members: []*v1Board.UpdateBoardRequest_Member{
			{MemberId: dto.MemberID, Delete: true},
		},
	}
}

func (dto DeleteMemberDTO) toSaveSub() *v1Authz.SaveSubRequest {
	return &v1Authz.SaveSubRequest{
		SubId: dto.MemberID,
		Refs: []*v1Authz.SaveSubRequest_Ref{
			{RefId: dto.BoardID, Delete: true},
		},
	}
}
