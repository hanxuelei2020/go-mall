package logic

import (
	"context"
	"encoding/json"
	"go-mall/service/user/rpc/userclient"

	"go-mall/service/user/api/internal/svc"
	"go-mall/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserInfoLogic) UserInfo() (resp *types.UserInfoResponse, err error) {
	// 查询数据
	uid, err := l.ctx.Value("uid").(json.Number).Int64()
	if err != nil {
		return nil, err
	}
	userInfo, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &userclient.UserInfoRequest{
		Id: uid,
	})
	if err != nil {
		return nil, err
	}

	return &types.UserInfoResponse{
		Id:     userInfo.Id,
		Name:   userInfo.Name,
		Gender: userInfo.Gender,
		Mobile: userInfo.Mobile,
	}, nil
}
