package file

import (
	"codeup.aliyun.com/61b84a04fa282c88e1039838/utils"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"meishiedu.com/classin/internal/model"
	"time"

	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileViewByUUIDLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetFileViewByUUIDLogic 根据uuid获取文件预览信息
func NewGetFileViewByUUIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileViewByUUIDLogic {
	return &GetFileViewByUUIDLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFileViewByUUIDLogic) GetFileViewByUUID(req *types.FileViewByUUIDReq) (resp *types.FileViewResp, err error) {
	file, err := l.svcCtx.SFilesModel.FindOneByUuid(l.ctx, req.Uuid)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errors.New("文件不存在")
		}
		return nil, err
	}
	url, _ := l.svcCtx.MsOssModel.GetSignUrl(l.ctx, file.ObjectKey, 3600)
	fileKey := file.Uuid
	expireTime := time.Now().Add(time.Minute * 30).Unix()
	if req.ExpireTime > 0 {
		expireTime = time.Now().Add(time.Minute * time.Duration(req.ExpireTime)).Unix()
	}
	payloadMap := utils.CreateFileViewPayload(url, fileKey, file.Name, file.FileType, expireTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(payloadMap))
	tokenString, err := token.SignedString([]byte(l.svcCtx.Config.OfficeAuth.AccessSecret))
	if err != nil {
		return nil, err
	}
	resp = &types.FileViewResp{}
	resp.ErrorInfo.ErrorCode = 1
	resp.ErrorInfo.ErrorMsg = "获取成功"
	resp.FileViewInfo.Token = tokenString
	resp.FileViewInfo.Name = file.Name
	resp.FileViewInfo.Key = fileKey
	resp.FileViewInfo.Address = url
	resp.FileViewInfo.Type = file.FileType
	return resp, nil
}
