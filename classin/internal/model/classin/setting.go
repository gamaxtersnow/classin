package classin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"meishiedu.com/classin/internal/types"
	"net/http"
	"net/url"
	"strconv"
)

const classSetNumUrl = "https://dynamic.eeo.cn/saasajax/course.ajax.php?action=editCourseClass"

type (
	SettingModel struct {
		client *http.Client
		ctx    context.Context
	}
	SettingResponse struct {
		ErrorInfo types.ErrorInfo `json:"error_info"`
		Data      []SetNumInfo    `json:"data"`
	}
	SetNumInfo struct {
		ClassId    int64  `json:"classId"`
		ClassName  string `json:"className"`
		ClassBtime string `json:"classBtime"`
		Errno      int64  `json:"errno"`
		Error      string `json:"error"`
	}
)

func NewSettingModel(client *http.Client) *SettingModel {
	return &SettingModel{
		client: client,
	}
}
func (c *SettingModel) ClassSetNum(ctx context.Context, cookie string, courseId int64, info []types.ClassJsonInfo) (setNumResponse *SettingResponse, err error) {
	// 将结构体转换为 JSON 字符串
	jsonData, err := json.Marshal(info)
	if err != nil {
		// 处理 JSON 序列化错误
		return nil, err
	}
	request := &http.Request{}
	params := url.Values{}
	params.Set("courseId", strconv.FormatInt(courseId, 10))
	params.Set("classJson", string(jsonData))
	logx.Infof("ClassSetNum=====params: %+v", params)
	request, err = http.NewRequest(http.MethodPost, classSetNumUrl, bytes.NewBufferString(params.Encode()))
	if err != nil {
		// 处理 JSON 序列化错误
		fmt.Println("JSON marshal error:", err)
		return nil, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("User-Agent", ctx.Value("User-Agent").(string))
	request.Header.Set("Cookie", cookie)
	response := &http.Response{}
	response, err = c.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = response.Body.Close()
	}()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))
	setNumResponse = &SettingResponse{}
	err = json.Unmarshal(body, setNumResponse)
	if err != nil {
		return nil, err
	}
	return setNumResponse, nil
}
