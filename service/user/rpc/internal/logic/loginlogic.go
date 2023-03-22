package logic

import (
	"context"
	"go-mall/common/cryptx"
	"go-mall/service/user/model"
	"google.golang.org/grpc/status"

	"go-mall/service/user/rpc/internal/svc"
	"go-mall/service/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginRequest) (*user.LoginResponse, error) {
	// 查询用户是否存在
	userDB, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, in.Mobile)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Error(100, "用户不存在")
		}
		return nil, status.Error(500, err.Error())
	}

	// 判断密码是否正确
	password := cryptx.PasswordEncrypt(l.svcCtx.Config.Salt, in.Password)
	if password != userDB.Password {
		return nil, status.Error(100, "用户名密码不正确")
	}

	return &user.LoginResponse{
		Id:     userDB.Id,
		Name:   userDB.Name,
		Gender: userDB.Gender,
		Mobile: userDB.Mobile,
	}, nil
}
