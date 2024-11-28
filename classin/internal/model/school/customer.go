package school

import (
	"codeup.aliyun.com/61b84a04fa282c88e1039838/xiaoxiaosdk"
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"meishiedu.com/classin/internal/model/oldcrm"
	"meishiedu.com/classin/internal/types"
)

var _ CustomerModel = (*customCustomerModel)(nil)

type (
	CustomerModel interface {
		GetCustomerBasicData(ctx context.Context, ids []int64) ([]types.CustomerBasicData, error)
	}
	customCustomerModel struct {
		oldcrmCustomerModel oldcrm.MsCustomerModel
		cache               cache.Cache
	}
)

func NewCustomerModel(client *xiaoxiaosdk.HttpClient, conn sqlx.SqlConn, cache cache.Cache) CustomerModel {
	return &customCustomerModel{
		oldcrmCustomerModel: oldcrm.NewMsCustomerModel(conn),
		cache:               cache,
	}
}
func (t *customCustomerModel) GetCustomerBasicData(ctx context.Context, ids []int64) ([]types.CustomerBasicData, error) {
	//查询 ms_customer 的数据
	customerList, err := t.oldcrmCustomerModel.FindListByIds(ctx, ids)
	fmt.Println("GetCustomerBasicData/customerList: ", customerList)
	if err != nil {
		fmt.Println("GetCustomerBasicData/customerList/error: ", err)
		return nil, err
	}

	//组合数据结果
	var respData []types.CustomerBasicData
	for _, customerItem := range customerList {
		respData = append(respData, types.CustomerBasicData{
			Id:            customerItem.Id,
			Name:          customerItem.Name.String,
			Type:          customerItem.Type.Int64,
			CustomerStage: customerItem.CustomerStage.Int64,
		})
	}
	return respData, nil
}
