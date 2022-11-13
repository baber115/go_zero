package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"go_zero_demo/service/user/rpc/types/user"

	"go_zero_demo/service/search/api/internal/svc"
	"go_zero_demo/service/search/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSearchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchLogic {
	return &SearchLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchLogic) Search(req *types.SearchReq) (resp *types.SearchResp, err error) {
	logx.Infof("userId: %v", l.ctx.Value("userId"))     // 这里的key和生成jwt token时传入的key一致
	logx.Infof("username: %v", l.ctx.Value("username")) // 这里的key和生成jwt token时传入的key一致
	userIdNumber := json.Number(fmt.Sprintf("%v", l.ctx.Value("userId")))
	logx.Infof("userId: %s", userIdNumber)
	userId, err := userIdNumber.Int64()
	if err != nil {
		return nil, err
	}

	// 使用user rpc
	_, err = l.svcCtx.UserRpc.GetUser(l.ctx, &user.IdReq{
		Id: userId,
	})
	if err != nil {
		return nil, err
	}

	return &types.SearchResp{
		Name:  req.Name,
		Count: 100,
	}, nil
}
