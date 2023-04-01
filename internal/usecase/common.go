package usecase

import (
	"context"
	"fmt"
	"strconv"

	authPB "github.com/krobus00/auth-service/pb/auth"
	"github.com/krobus00/product-service/internal/constant"
	"github.com/krobus00/product-service/internal/model"
)

func getUserIDFromCtx(ctx context.Context) string {
	ctxUserID := ctx.Value(constant.KeyUserIDCtx)

	userID := fmt.Sprintf("%v", ctxUserID)
	if userID == "" {
		return constant.GuestID
	}
	return userID
}

func getDataSource(ctx context.Context) int {
	ctxData := ctx.Value(constant.KeyDataSource)

	data, err := strconv.Atoi(fmt.Sprintf("%v", ctxData))
	if err != nil {
		return constant.SourceOS
	}
	return data
}

func hasAccess(ctx context.Context, authClient authPB.AuthServiceClient, permissions []string) error {
	userID := getUserIDFromCtx(ctx)

	permissions = append(permissions, constant.PermissionFullAccess)
	res, err := authClient.HasAccess(ctx, &authPB.HasAccessRequest{
		UserId:      userID,
		Permissions: permissions,
	})

	if err != nil {
		return model.ErrUnauthorizedAccess
	}
	if res == nil {
		return model.ErrUnauthorizedAccess
	}
	if !res.Value {
		return model.ErrUnauthorizedAccess
	}
	return nil
}
