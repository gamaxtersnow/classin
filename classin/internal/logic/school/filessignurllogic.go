package school

import (
	"context"
	"errors"
	"fmt"

	"meishiedu.com/classin/internal/svc"
	"meishiedu.com/classin/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FilesSignUrlLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewFilesSignUrlLogic 查看文件链接
func NewFilesSignUrlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FilesSignUrlLogic {
	return &FilesSignUrlLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FilesSignUrlLogic) FilesSignUrl(req *types.FileSignUrlRequest) (resp *types.FileSignUrlResponse, err error) {
	// 参数验证
	if req == nil {
		return nil, errors.New("请求参数不能为空")
	}
	if req.FileKey == "" {
		return nil, errors.New("文件key不能为空")
	}

	// 记录请求日志
	logx.Infof("获取文件签名URL, FileKey: %s", req.FileKey)

	// 获取签名URL，默认1小时有效期
	signUrl, err := l.svcCtx.MsOssModel.GetSignUrl(l.ctx, req.FileKey, 3600)
	if err != nil {
		logx.Errorf("获取签名URL失败: %v", err)
		return nil, fmt.Errorf("获取文件访问链接失败: %v", err)
	}
	return &types.FileSignUrlResponse{
		ErrorInfo: types.ErrorInfo{
			ErrorCode: 1,
			ErrorMsg:  "success",
		},
		Url: signUrl,
	}, nil
}
