package logic

import (
	"context"
	"go-mall/common/jwtx"
	"go-mall/service/user/rpc/userclient"
	"time"

	"go-mall/service/user/api/internal/svc"
	"go-mall/service/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	// 调用远程服务
	loginInfo, err := l.svcCtx.UserRpc.Login(l.ctx, &userclient.LoginRequest{
		Mobile:   req.Mobile,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	// 生成 Token
	now := time.Now().Unix()
	token, err := jwtx.GetToken(l.svcCtx.Config.Auth.AccessSecret,
		now, accessExpire, loginInfo.Id)

	// 处理 token
	if err != nil {
		return nil, err
	}
	return &types.LoginResponse{
		AccessToken:  token,
		AccessExpire: now + accessExpire,
	}, nil
}
