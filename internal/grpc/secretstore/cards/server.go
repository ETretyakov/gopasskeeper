package cards

import (
	"context"
	"errors"
	"gopasskeeper/internal/domain/models"
	"gopasskeeper/internal/grpc/interceptors"
	cardsv1 "gopasskeeper/internal/grpc/secretstore/cards/gen/cards"
	"gopasskeeper/internal/lib/validators"
	"gopasskeeper/internal/logger"
	"gopasskeeper/internal/services/secretstore/cards"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverAPI struct {
	cardsv1.UnimplementedCardsServer
	cards Cards
}

// Cards is an interface for cards service.
type Cards interface {
	Add(
		ctx context.Context,
		uid string,
		name string,
		number string,
		mask string,
		month int32,
		year int32,
		cvc string,
		pin string,
	) (*models.Message, error)
	GetSecret(
		ctx context.Context,
		uid string,
		cardID string,
	) (*models.CardSecret, error)
	Search(
		ctx context.Context,
		uid string,
		schema *models.CardSearchRequest,
	) (*models.CardSearchResponse, error)
	Remove(
		ctx context.Context,
		uid string,
		cardID string,
	) (*models.Message, error)
}

// RegisterCards is a method to register Cards service within grpc server.
func RegisterCards(gRPCServer *grpc.Server, cards Cards) {
	cardsv1.RegisterCardsServer(gRPCServer, &serverAPI{cards: cards})
}

// Add is a method to add a card record.
func (s *serverAPI) Add(
	ctx context.Context,
	in *cardsv1.CardAddRequest,
) (*cardsv1.CardAddResponse, error) {
	creditCard := validators.NewCreditCard(
		in.GetNumber(),
		in.GetMonth(),
		in.GetYear(),
		in.GetCvc(),
		in.GetPin(),
	)

	if err := creditCard.Validate(); err != nil {
		logger.Error("failed to validate card", err)
		return nil, status.Error(codes.InvalidArgument, "credit card is invalid")
	}

	uid, err := interceptors.ExtractUID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed extract uid")
	}

	msg, err := s.cards.Add(
		ctx,
		uid,
		in.GetName(),
		creditCard.Number,
		creditCard.Mask(),
		creditCard.Month,
		creditCard.Year,
		creditCard.CVC,
		creditCard.PIN,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed add card")
	}

	return &cardsv1.CardAddResponse{Status: msg.Status, Msg: msg.Msg}, nil
}

// GetSecret is a method to get a card secret.
func (s *serverAPI) GetSecret(
	ctx context.Context,
	in *cardsv1.CardSecretRequest,
) (*cardsv1.CardSecretResponse, error) {
	if in.GetSecretId() == "" {
		return nil, status.Error(codes.InvalidArgument, "secret_id is required")
	}

	uid, err := interceptors.ExtractUID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed extract uid")
	}

	cardSecret, err := s.cards.GetSecret(ctx, uid, in.GetSecretId())
	if err != nil {
		if errors.Is(err, cards.ErrCardNotFound) {
			return nil, status.Error(codes.InvalidArgument, "failed to find card")
		}

		return nil, status.Error(codes.Internal, "failed to get secret")
	}

	return &cardsv1.CardSecretResponse{
		Name:   cardSecret.Name,
		Number: cardSecret.Number,
		Month:  cardSecret.Month,
		Year:   cardSecret.Year,
		Cvc:    cardSecret.CVC,
		Pin:    cardSecret.PIN,
	}, nil
}

// Search is a method to search among cards records.
func (s *serverAPI) Search(
	ctx context.Context,
	in *cardsv1.CardSearchRequest,
) (*cardsv1.CardSearchResponse, error) {
	if in.GetLimit() == 0 {
		return nil, status.Error(codes.InvalidArgument, "limit can't be 0")
	}

	uid, err := interceptors.ExtractUID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed extract uid")
	}

	searchResponse, err := s.cards.Search(
		ctx,
		uid,
		&models.CardSearchRequest{
			Substring: in.GetSubstring(),
			Offset:    uint64(in.GetOffset()),
			Limit:     uint32(in.GetLimit()),
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to search cards")
	}

	items := []*cardsv1.CardSearchItem{}
	for _, i := range searchResponse.Items {
		items = append(items, &cardsv1.CardSearchItem{
			Id:   i.ID,
			Name: i.Name,
			Mask: i.Mask,
		})
	}

	return &cardsv1.CardSearchResponse{
		Count: int64(searchResponse.Count),
		Items: items,
	}, nil
}

// Remove is a method to remove a card record.
func (s *serverAPI) Remove(
	ctx context.Context,
	in *cardsv1.CardRemoveRequest,
) (*cardsv1.CardRemoveResponse, error) {
	if in.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "secret_id is required")
	}

	uid, err := interceptors.ExtractUID(ctx)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "failed extract uid")
	}

	msg, err := s.cards.Remove(ctx, uid, in.GetId())
	if err != nil {
		if errors.Is(err, cards.ErrCardNotFound) {
			return nil, status.Error(codes.InvalidArgument, "failed to find card")
		}
		return nil, status.Error(codes.Internal, "failed to get card")
	}

	return &cardsv1.CardRemoveResponse{
		Status: msg.Status,
		Msg:    msg.Msg,
	}, nil
}
