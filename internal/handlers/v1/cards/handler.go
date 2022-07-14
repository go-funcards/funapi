package cards

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-funcards/funapi/internal/gin/binding"
	"github.com/go-funcards/funapi/internal/gin/httputil"
	"github.com/go-funcards/funapi/internal/handlers"
	"github.com/go-funcards/funapi/internal/handlers/v1/clientutil"
	v1Card "github.com/go-funcards/funapi/proto/card_service/v1"
	"github.com/go-funcards/slice"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
)

var _ handlers.Handler = (*Handler)(nil)

type Handler struct {
	*handlers.BaseBoard
	CardService v1Card.CardClient
	Log         *zap.Logger
}

func (h *Handler) Register(rg *gin.RouterGroup) {
	g := rg.Group("/cards")
	{
		g.GET("", h.list)
		g.POST("", h.create)
		g.PATCH("", h.updateMany)

		b := g.Group("/:card_id")
		{
			b.GET("", h.read)
			b.PATCH("", h.update)
			b.DELETE("", h.delete)
		}
	}
}

// @Summary Card List
// @Tags Cards
// @ModuleID listCard
// @Accept json
// @Produce json
// @Param board_id query string true "Board ID" format(uuid)
// @Param page_index query int false "Page Index" minimum(0)
// @Param page_size query int false "Page Size" minimum(1) maximum(1000)
// @Success 200 {object} cards.PageResponse
// @Failure 400,401,403,500 {object} httputil.APIError
// @Router /cards [get]
// @Security BearerAuth
func (h *Handler) list(c *gin.Context) {
	h.Log.Debug("card handler::list bind")
	req := PageReq()
	if !binding.BindQueryAndValidate(c, &req) {
		return
	}

	ctx := context.TODO()

	if !h.IsGranted(ctx, c, req.BoardID, "READ") {
		return
	}

	h.Log.Debug("card handler::list call gRPC /CardClient/GetCards")
	response, err := h.CardService.GetCards(ctx, req.toRead())
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, PageResp(response, req))
}

// @Summary Create Card
// @Tags Cards
// @ModuleID createCard
// @Accept json
// @Param payload body cards.CreateCardDTO true "Card data"
// @Success 201
// @Failure 400,401,403,422,500 {object} httputil.APIError
// @Header 201 {string} Location "/cards/{card_id}"
// @Router /cards [post]
// @Security BearerAuth
func (h *Handler) create(c *gin.Context) {
	h.Log.Debug("card handler::create bind")
	var dto CreateCardDTO
	if !binding.BindCtx(c, &dto) || !binding.BindBodyAndValidate(c, &dto) {
		return
	}

	ctx := context.TODO()

	if !h.IsGranted(ctx, c, dto.BoardID, "CREATE") {
		return
	}

	id := uuid.NewString()

	h.Log.Debug("card handler::create call gRPC /CardClient/CreateCard")
	if _, err := h.CardService.CreateCard(ctx, dto.toCreate(id)); err != nil {
		_ = c.Error(err)
		return
	}

	httputil.Created(c, id)
}

// @Summary Update Many Cards
// @Tags Cards
// @ModuleID updateManyCards
// @Accept json
// @Param payload body cards.UpdateManyCardsDTO true "Cards data"
// @Success 204
// @Failure 400,401,403,404,422,500 {object} httputil.APIError
// @Router /cards [patch]
// @Security BearerAuth
func (h *Handler) updateMany(c *gin.Context) {
	h.Log.Debug("card handler::updateMany bind")
	var dto UpdateManyCardsDTO
	if !binding.BindBodyAndValidate(c, &dto) {
		return
	}

	ctx := context.TODO()

	h.Log.Debug("card handler::updateMany call gRPC /CardClient/GetCards")
	data, err := h.CardService.GetCards(ctx, &v1Card.CardsRequest{
		PageIndex: 0,
		PageSize:  uint32(len(dto.Data)),
		CardIds: slice.Map(dto.Data, func(item UpdateCardDTO) string {
			return item.CardID
		}),
	})
	if err != nil {
		_ = c.Error(err)
		return
	}

	boards := make(map[string]bool)

	for _, card := range data.GetCards() {
		if _, ok := boards[card.GetBoardId()]; !ok {
			boards[card.GetBoardId()] = true
			if !h.IsGranted(ctx, c, card.GetBoardId(), "UPDATE") {
				return
			}
		}
	}

	for _, item := range dto.Data {
		if _, ok := boards[item.BoardID]; !ok && len(item.BoardID) > 0 {
			boards[item.BoardID] = true
			if !h.IsGranted(ctx, c, item.BoardID, "CREATE") {
				return
			}
		}
	}

	h.Log.Debug("card handler::updateMany call gRPC /CardClient/UpdateManyCards")
	if _, err = h.CardService.UpdateManyCards(ctx, dto.toUpdateMany()); err != nil {
		_ = c.Error(err)
		return
	}

	httputil.NoContent(c)
}

// @Summary Read Card
// @Tags Cards
// @ModuleID readCard
// @Accept json
// @Produce json
// @Param card_id path string true "Card ID" format(uuid)
// @Success 200 {object} cards.Card
// @Failure 400,401,403,404,500 {object} httputil.APIError
// @Router /cards/{card_id} [get]
// @Security BearerAuth
func (h *Handler) read(c *gin.Context) {
	h.Log.Debug("card handler::read bind")
	var dto ReadCardDTO
	if !binding.BindUriAndValidate(c, &dto) {
		return
	}

	ctx := context.TODO()

	h.Log.Debug("card handler::read call gRPC /CardClient/GetCard")
	card, err := clientutil.GetCard(ctx, h.CardService, dto.CardID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if !h.IsGranted(ctx, c, card.GetBoardId(), "READ") {
		return
	}

	c.JSON(http.StatusOK, CreateCard(card))
}

// @Summary Update Card
// @Tags Cards
// @ModuleID updateCard
// @Accept json
// @Param card_id path string true "Card ID" format(uuid)
// @Param payload body cards.UpdateCardDTO true "Card data"
// @Success 204
// @Failure 400,401,403,404,422,500 {object} httputil.APIError
// @Router /cards/{card_id} [patch]
// @Security BearerAuth
func (h *Handler) update(c *gin.Context) {
	h.Log.Debug("card handler::update bind")
	var dto UpdateCardDTO
	if !binding.BindBody(c, &dto) || !binding.BindUriAndValidate(c, &dto) {
		return
	}

	ctx := context.TODO()

	h.Log.Debug("card handler::update call gRPC /CardClient/GetCard")
	card, err := clientutil.GetCard(ctx, h.CardService, dto.CardID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if !h.IsGranted(ctx, c, card.GetBoardId(), "UPDATE") {
		return
	}

	if len(dto.BoardID) > 0 && dto.BoardID != card.GetBoardId() && !h.IsGranted(ctx, c, dto.BoardID, "CREATE") {
		return
	}

	h.Log.Debug("card handler::update call gRPC /CardClient/UpdateCard")
	if _, err = h.CardService.UpdateCard(ctx, dto.toUpdate()); err != nil {
		_ = c.Error(err)
		return
	}

	httputil.NoContent(c)
}

// @Summary Delete Card
// @Tags Cards
// @ModuleID deleteCard
// @Param card_id path string true "Card ID" format(uuid)
// @Success 204
// @Failure 400,401,403,404,500 {object} httputil.APIError
// @Router /cards/{card_id} [delete]
// @Security BearerAuth
func (h *Handler) delete(c *gin.Context) {
	h.Log.Debug("card handler::delete bind")
	var dto DeleteCardDTO
	if !binding.BindUriAndValidate(c, &dto) {
		return
	}

	ctx := context.TODO()

	h.Log.Debug("card handler::delete call gRPC /CardClient/GetCard")
	card, err := clientutil.GetCard(ctx, h.CardService, dto.CardID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	if !h.IsGranted(ctx, c, card.GetBoardId(), "DELETE") {
		return
	}

	h.Log.Debug("card handler::delete call gRPC /CardClient/DeleteCard")
	if _, err = h.CardService.DeleteCard(ctx, dto.toDelete()); err != nil {
		_ = c.Error(err)
		return
	}

	httputil.NoContent(c)
}
