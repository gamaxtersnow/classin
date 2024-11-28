package classin

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"meishiedu.com/classin/internal/types"
	"net/http"
	"net/url"
)

const classVodListUrl = "https://dynamic.eeo.cn/saasajax/school.ajax.php?action=getClassVodList"

type (
	LessonModel struct {
		client *http.Client
		ctx    context.Context
	}

	PlaySet struct {
		Url string `json:"Url"`
	}

	File struct {
		FileId   string    `json:"FileId"`
		FileName string    `json:"FileName"`
		FileSize string    `json:"Size"`
		PlaySet  []PlaySet `json:"Playset"`
	}

	VodInfo struct {
		ClassId  int64  `json:"Id"`
		AllCount int64  `json:"AllCount"`
		FileList []File `json:"FileList"`
	}

	ClassResponse struct {
		ErrorInfo types.ErrorInfo           `json:"error_info"`
		Data      struct{ VodInfo VodInfo } `json:"data"`
	}
)

func NewLessonModel(client *http.Client) *LessonModel {
	return &LessonModel{
		client: client,
	}
}
func (c *LessonModel) GetLessonClipList(ctx context.Context, cookie string, params url.Values) (lessonVideos *ClassResponse, err error) {
	request := &http.Request{}
	request, err = http.NewRequest(http.MethodPost, classVodListUrl, bytes.NewBufferString(params.Encode()))
	if err != nil {
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
	lessonVideos = &ClassResponse{}
	err = json.Unmarshal(body, lessonVideos)
	if err != nil {
		return nil, err
	}
	return lessonVideos, nil
}
