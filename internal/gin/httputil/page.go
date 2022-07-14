package httputil

type PageRequest struct {
	Index uint64 `json:"-" form:"page_index"`
	Size  uint32 `json:"-" form:"page_size" validate:"min=1,max=1000"`
}

type PageResponse struct {
	Index uint64 `json:"page_index"`
	Size  uint32 `json:"page_size"`
	Total uint64 `json:"total"`
}

func (p PageRequest) ToPageResponse(total uint64) PageResponse {
	return PageResponse{
		Index: p.Index,
		Size:  p.Size,
		Total: total,
	}
}
