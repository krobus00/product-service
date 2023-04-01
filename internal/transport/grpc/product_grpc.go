package grpc

import (
	"context"
	"time"

	"github.com/krobus00/product-service/internal/model"
	"github.com/krobus00/product-service/internal/utils"
	pb "github.com/krobus00/product-service/pb/product"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (t *Delivery) Create(ctx context.Context, in *pb.CreateProductRequest) (*pb.Product, error) {
	defer func(tn time.Time) {
		_, _, fn := utils.Trace()
		utils.TimeTrack(tn, fn)
	}(time.Now())

	ctx = setUserIDCtx(ctx, in.GetUserId())

	payload := model.NewCreateProductPayloadFromProto(in)
	product, err := t.productUC.Create(ctx, payload)
	switch err {
	case nil:
	case model.ErrUnauthorizedAccess:
		return nil, status.Error(codes.Unauthenticated, err.Error())
	case model.ErrThumbnailNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	case model.ErrThumbnailTypeNotAllowed:
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	default:
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}

	return product.ToProto(), nil
}

func (t *Delivery) Update(ctx context.Context, in *pb.UpdateProductRequest) (*pb.Product, error) {
	defer func(tn time.Time) {
		_, _, fn := utils.Trace()
		utils.TimeTrack(tn, fn)
	}(time.Now())

	ctx = setUserIDCtx(ctx, in.GetUserId())

	payload := model.NewUpdateProductPayloadFromProto(in)
	product, err := t.productUC.Update(ctx, payload)
	switch err {
	case nil:
	case model.ErrUnauthorizedAccess:
		return nil, status.Error(codes.Unauthenticated, err.Error())
	case model.ErrProductNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	case model.ErrThumbnailNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	case model.ErrThumbnailTypeNotAllowed:
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	default:
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}

	return product.ToProto(), nil
}

func (t *Delivery) Delete(ctx context.Context, in *pb.DeleteProductRequest) (*pb.Empty, error) {
	defer func(tn time.Time) {
		_, _, fn := utils.Trace()
		utils.TimeTrack(tn, fn)
	}(time.Now())

	ctx = setUserIDCtx(ctx, in.GetUserId())

	err := t.productUC.Delete(ctx, in.GetId())
	switch err {
	case nil:
	case model.ErrUnauthorizedAccess:
		return nil, status.Error(codes.Unauthenticated, err.Error())
	case model.ErrProductNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	default:
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}

	return &pb.Empty{}, nil
}

func (t *Delivery) FindByID(ctx context.Context, in *pb.FindByIDRequest) (*pb.Product, error) {
	defer func(tn time.Time) {
		_, _, fn := utils.Trace()
		utils.TimeTrack(tn, fn)
	}(time.Now())

	ctx = setUserIDCtx(ctx, in.GetUserId())

	product, err := t.productUC.FindByID(ctx, in.GetId())
	switch err {
	case nil:
	case model.ErrUnauthorizedAccess:
		return nil, status.Error(codes.Unauthenticated, err.Error())
	case model.ErrProductNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	default:
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}

	return product.ToProto(), nil
}

func (t *Delivery) FindByIDs(ctx context.Context, in *pb.FindByIDsRequest) (*pb.FindByIDsResponse, error) {
	defer func(tn time.Time) {
		_, _, fn := utils.Trace()
		utils.TimeTrack(tn, fn)
	}(time.Now())

	ctx = setUserIDCtx(ctx, in.GetUserId())

	product, err := t.productUC.FindByIDs(ctx, in.GetIds())
	switch err {
	case nil:
	case model.ErrUnauthorizedAccess:
		return nil, status.Error(codes.Unauthenticated, err.Error())
	case model.ErrProductNotFound:
		return nil, status.Error(codes.NotFound, err.Error())
	default:
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}

	return &pb.FindByIDsResponse{
		Items: product.ToProto(),
	}, nil
}

func (t *Delivery) FindPaginatedIDs(ctx context.Context, in *pb.PaginationRequest) (*pb.PaginationResponse, error) {
	defer func(tn time.Time) {
		_, _, fn := utils.Trace()
		utils.TimeTrack(tn, fn)
	}(time.Now())

	ctx = setUserIDCtx(ctx, in.GetUserId())

	payload := model.NewPaginationPayloadFromProto(in)
	res, err := t.productUC.FindPaginatedIDs(ctx, payload)
	switch err {
	case nil:
	case model.ErrUnauthorizedAccess:
		return nil, status.Error(codes.Unauthenticated, err.Error())
	default:
		return nil, status.Error(codes.Internal, codes.Internal.String())
	}

	return res.ToProto(), nil
}
