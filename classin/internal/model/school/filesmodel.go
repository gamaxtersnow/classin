package school

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FilesModel = (*customFilesModel)(nil)

const FileStatusNormal = 1
const FileStatusDeleted = 2

type (
	// FilesModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFilesModel.
	FilesModel interface {
		filesModel
	}

	customFilesModel struct {
		*defaultFilesModel
	}
)

// NewFilesModel returns a model for the database table.
func NewFilesModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FilesModel {
	return &customFilesModel{
		defaultFilesModel: newFilesModel(conn, c, opts...),
	}
}
