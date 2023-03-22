package logic

import (
	"context"
	"go-mall/service/order/model"
	"go-mall/service/product/rpc/product"
	"go-mall/service/user/rpc/user"
	"google.golang.org/grpc/status"

	"go-mall/service/order/rpc/internal/svc"
	"go-mall/service/order/rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateLogic) Create(in *order.CreateRequest) (*order.CreateResponse, error) {
	// 查询用户是否存在
	_, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &user.UserInfoRequest{
		Id: in.Uid,
	})
	if err != nil {
		return nil, err
	}

	// 查询产品是否存在
	productRes, err := l.svcCtx.ProductRpc.Detail(l.ctx, &product.DetailRequest{
		Id: in.Pid,
	})
	if err != nil {
		return nil, err
	}
	// 判断产品库存是否充足
	if productRes.Stock <= 0 {
		return nil, status.Error(500, "产品库存不足")
	}

	newOrder := model.Order{
		Uid:    in.Uid,
		Pid:    in.Pid,
		Amount: in.Amount,
		Status: 0,
	}
	// 创建订单
	res, err := l.svcCtx.OrderModel.Insert(l.ctx, &newOrder)
	if err != nil {
		return nil, status.Error(500, err.Error())
	}

	newOrder.Id, err = res.LastInsertId()
	if err != nil {
		return nil, status.Error(500, err.Error())
	}
	// 更新产品库存, 后期可以增加锁定库存的功能, 使用 rabbitmq 防止订单构建失败
	_, err = l.svcCtx.ProductRpc.Update(l.ctx, &product.UpdateRequest{
		Id:     productRes.Id,
		Name:   productRes.Name,
		Desc:   productRes.Desc,
		Stock:  productRes.Stock - 1,
		Amount: productRes.Amount,
		Status: productRes.Status,
	})
	if err != nil {
		return nil, err
	}

	return &order.CreateResponse{
		Id: newOrder.Id,
	}, nil
}
