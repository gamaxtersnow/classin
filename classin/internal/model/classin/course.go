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

const courseListUrl = "https://dynamic.eeo.cn/saasajax/course.ajax.php?action=getClassListNew"

type (
	CourseModel struct {
		client *http.Client
	}

	CourseResponse struct {
		ErrorInfo types.ErrorInfo `json:"error_info"`
		Data      ResponseData    `json:"data"`
	}
	ResponseData struct {
		TotalClassNum int64        `json:"totalClassNum"`
		ClassList     []CourseInfo `json:"classList"`
	}
	CourseInfo struct {
		ClassId        int64       `json:"id"`
		CourseId       int64       `json:"courseId"`
		ClassName      string      `json:"className"`
		CourseName     string      `json:"courseName"`
		ClassStartTime int64       `json:"classBtime"`
		ClassEndTime   int64       `json:"classEtime"`
		StudentNum     int64       `json:"studentNum"`
		SeatNum        int64       `json:"seatNum"`
		IsHd           int64       `json:"isHd"`
		IsDc           int64       `json:"isDc"`
		TeacherInfo    TeacherInfo `json:"teacherInfo"`
	}
	TeacherInfo struct {
		Mobile string `json:"mobile"`
		Name   string `json:"name"`
	}
)

func NewCourseModel(client *http.Client) *CourseModel {
	return &CourseModel{
		client: client,
	}
}
func (c *CourseModel) GetCourseList(ctx context.Context, cookie string, params url.Values) (courses *CourseResponse, err error) {
	request := &http.Request{}
	request, err = http.NewRequest(http.MethodPost, courseListUrl, bytes.NewBufferString(params.Encode()))
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
	courses = &CourseResponse{}
	err = json.Unmarshal(body, courses)
	if err != nil {
		return nil, err
	}
	return courses, nil
}
