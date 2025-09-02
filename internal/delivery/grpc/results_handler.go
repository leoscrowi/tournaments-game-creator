package grpc

import (
	"context"
	uuid2 "github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
	"tournaments-core/internal/delivery/grpc/results_grpc"
	"tournaments-core/internal/domain/models"
	"tournaments-core/internal/domain/ports/repository"
	"tournaments-core/internal/domain/ports/usecase"
	usecase2 "tournaments-core/internal/usecase"
)

type res_server struct {
	results_grpc.UnimplementedResultsServiceServer
	usecase usecase.ResultsUseCase
}

func NewResultsGrpcServer(gserver *grpc.Server, rep *repository.ResultsRepository) {

	resultsServer := &res_server{
		usecase: usecase2.NewResultsUseCase(*rep, 10*time.Second),
	}

	results_grpc.RegisterResultsServiceServer(gserver, resultsServer)
}

func (s res_server) FetchById(ctx context.Context, request *results_grpc.IdResultRequest) (*results_grpc.ResultResponse, error) {
	uuid, err := uuid2.Parse(request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	r, err := s.usecase.FetchById(ctx, uuid)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}

	return &results_grpc.ResultResponse{
		Id:       uuid.String(),
		GameId:   r.GameID.String(),
		WinnerId: r.WinnerID.String(),
		Comment:  r.Comment,
	}, nil
}

func (s res_server) DeleteById(ctx context.Context, request *results_grpc.IdResultRequest) (*emptypb.Empty, error) {
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

func (s res_server) Create(ctx context.Context, request *results_grpc.ResultCreateRequest) (*emptypb.Empty, error) {
	gameId, err := uuid2.Parse(request.GameId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	winnerId, err := uuid2.Parse(request.WinnerId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	var result *models.Result
	result = &models.Result{
		ResultID: uuid2.New(),
		GameID:   gameId,
		WinnerID: winnerId,
		Comment:  request.Comment,
	}

	err = s.usecase.Create(ctx, result)
	if err != nil {
		return nil, status.Errorf(codes.Canceled, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s res_server) Update(ctx context.Context, request *results_grpc.ResultRequest) (*emptypb.Empty, error) {
	uuid, err := uuid2.Parse(request.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	gameId, err := uuid2.Parse(request.GetGameId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	winnerId, err := uuid2.Parse(request.GetWinnerId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	var result *models.Result
	result = &models.Result{
		ResultID: uuid,
		GameID:   gameId,
		WinnerID: winnerId,
		Comment:  request.Comment,
	}

	err = s.usecase.Update(ctx, result)
	if err != nil {
		return nil, status.Errorf(codes.Canceled, err.Error())
	}

	return &emptypb.Empty{}, nil
}
