package logic

import (
	"context"
	"go-mall/service/product/model"
	"google.golang.org/grpc/status"

	"go-mall/service/product/rpc/internal/svc"
	"go-mall/service/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type DetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DetailLogic) Detail(in *product.DetailRequest) (*product.DetailResponse, error) {
	// 查询产品是否存在
	res, err := l.svcCtx.ProductModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, status.Errorf(100, "产品 id 为 %v 的产品不存在", in.Id)
		}
		// 如果不是没有找到的问题, 那么即使服务器的问题
		return nil, status.Error(500, err.Error())
	}

	return &product.DetailResponse{
		Id:     res.Id,
		Name:   res.Name,
		Desc:   res.Desc,
		Stock:  res.Stock,
		Amount: res.Amount,
		Status: res.Status,
	}, nil
}
