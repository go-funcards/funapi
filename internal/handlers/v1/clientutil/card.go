package clientutil

import (
	"context"
	"github.com/go-funcards/funapi/internal/gin/httputil"
	v1Card "github.com/go-funcards/funapi/proto/card_service/v1"
)

func CardsRequest(id string) *v1Card.CardsRequest {
	return &v1Card.CardsRequest{
		PageIndex: 0,
		PageSize:  1,
		CardIds:   []string{id},
	}
}

func GetCard(ctx context.Context, client v1Card.CardClient, id string) (*v1Card.CardsResponse_Card, error) {
	response, err := client.GetCards(ctx, CardsRequest(id))
	if err != nil {
		return nil, err
	}
	if len(response.GetCards()) != 1 {
		return nil, httputil.ErrNotFound
	}
	return response.GetCards()[0], nil
}
