package grpc

import (
	"context"
	uuid2 "github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
	"tournaments-core/internal/delivery/grpc/games_grpc"
	"tournaments-core/internal/domain/models"
	"tournaments-core/internal/domain/ports/repository"
	"tournaments-core/internal/domain/ports/usecase"
	usecase2 "tournaments-core/internal/usecase"
)

type games_server struct {
	games_grpc.UnimplementedGamesServiceServer
	usecase usecase.GamesUseCase
}

func NewGamesGrpcServer(gserver *grpc.Server, rep *repository.GamesRepository) {

	gamesServer := &games_server{
		usecase: usecase2.NewGamesUseCase(*rep, 10*time.Second),
	}

	games_grpc.RegisterGamesServiceServer(gserver, gamesServer)
}

func (s games_server) FetchById(ctx context.Context, request *games_grpc.IdGameRequest) (*games_grpc.GameResponse, error) {
	uuid, err := uuid2.Parse(request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	r, err := s.usecase.FetchById(ctx, uuid)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	gameStartProto := timestamppb.New(r.GameStart)
	if err := gameStartProto.CheckValid(); err != nil {
		return nil, status.Errorf(codes.Internal, "invalid time: %v", err)
	}

	return &games_grpc.GameResponse{
		Id:         uuid.String(),
		GameStart:  gameStartProto,
		GameTypeId: r.GameTypeID.String(),
	}, nil
}

func (s games_server) DeleteById(ctx context.Context, request *games_grpc.IdGameRequest) (*emptypb.Empty, error) {
	uuid, err := uuid2.Parse(request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	err = s.usecase.DeleteById(ctx, uuid)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s games_server) Update(ctx context.Context, request *games_grpc.GameRequest) (*emptypb.Empty, error) {
	uuid, err := uuid2.Parse(request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	gameTypeUuid, err := uuid2.Parse(request.GetGameTypeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	var game *models.Game
	game = &models.Game{
		GameID:     uuid,
		GameStart:  request.GameStart.AsTime(),
		GameTypeID: gameTypeUuid,
	}

	err = s.usecase.Update(ctx, game)
	if err != nil {
		return nil, status.Errorf(codes.Canceled, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s games_server) Create(ctx context.Context, request *games_grpc.GameCreateRequest) (*emptypb.Empty, error) {
	var game *models.Game

	gameTypeUuid, err := uuid2.Parse(request.GetGameTypeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	game = &models.Game{
		GameID:     uuid2.New(),
		GameStart:  request.GameStart.AsTime(),
		GameTypeID: gameTypeUuid,
	}
	err = s.usecase.Create(ctx, game)
	if err != nil {
		return nil, status.Errorf(codes.Canceled, err.Error())
	}
	return &emptypb.Empty{}, nil
}
