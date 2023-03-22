package logic

import (
	"context"
	"go-mall/common/cryptx"
	"go-mall/service/user/model"
	"go-mall/service/user/rpc/internal/svc"
	"go-mall/service/user/rpc/user"
	"google.golang.org/grpc/status"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterRequest) (*user.RegisterResponse, error) {
	// 首先判断是否已经注册了
	_, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, in.Mobile)
	// 如果用户已经存在了
	if err == nil {
		l.Logger.Infof("手机号为%v的用户已经存在", in.Mobile)
		return nil, status.Error(100, "该用户已存在")
	}
	// 如果用户没有找到
	if err == model.ErrNotFound {
		newUser := model.User{
			Name:     in.Name,
			Gender:   in.Gender,
			Mobile:   in.Mobile,
			Password: cryptx.PasswordEncrypt(l.svcCtx.Config.Salt, in.Password),
		}

		// 注册用户
		res, err := l.svcCtx.UserModel.Insert(l.ctx, &newUser)
		if err != nil {
			l.Logger.Error("数据库插入失败")
			return nil, status.Error(500, "数据库插入失败")
		}

		newUser.Id, err = res.LastInsertId()
		if err != nil {
			l.Logger.Error("内部服务器错误")
			return nil, status.Error(500, err.Error())
		}

		// 返回数据
		return &user.RegisterResponse{
			Id:     newUser.Id,
			Name:   newUser.Name,
			Gender: newUser.Gender,
			Mobile: newUser.Mobile,
		}, nil
	}
	return nil, status.Error(500, err.Error())
}
