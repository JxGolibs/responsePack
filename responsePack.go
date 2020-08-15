package responsePack

import (
	"reflect"
	"time"
"github.com/kataras/iris/v12"
)

//400 接受到请求，但处理失败
//401 未登陆，或者登陆失效
//403 权限不足 //404 不存在

// 40000  数据解析失败
// 40001  数据不能为空
// 40002  数据已经存在
// 40003  短信验证失败
// 40004  用户不存在
// 40005  token已经过期

type (
	Response struct {
		Timestamp int64  `json:"timestamp"` //生成时间
		Code    int         `json:"code"`
		Status  int         `json:"status"`
		Message string      `json:"message"`
		Data    *Data        `json:"data,omitempty"`
		ExtDat     interface{} `json:"ext,omitempty"`
	}

	Data struct {
		Content interface{} `json:"content,omitempty"`
		Page    *Page        `json:"page,omitempty"`
	}

	Page struct {
		TotalRecord      int `json:"total"`      //总记录数量
		TotalPages int `json:"totalPages"` //总页数
		PageNo     int `json:"pageNo"`     //当前页号
		PageSize   int `json:"pageSize"`   //每页数据量
	}
)

func ParsePage(ctx iris.Context) *Page {
	if ctx == nil || !ctx.URLParamExists("page_no") {
		//无分页处理
		return nil
	}
	return &Page{
		TotalRecord: 0,
		TotalPages: 0,
		PageNo: ctx.URLParamIntDefault("page_no", 1),
		PageSize:  ctx.URLParamIntDefault("page_size", 5),
	}
}

func IsNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return i==nil
}


func New(code int, status int, msg string, content interface{}, page ...*Page) *Response {
	r := &Response{
		Code:    code,
		Status:  status,
		Message: msg,
	}
	if !IsNil(content)   {
		r.Data = &Data{
			Content: content,
		}
		if len(page) > 0 && page[0]!=nil {
			r.Data.Page = page[0]
		}
	}

	return r
}

func Success(msg string, content interface{}, page...*Page) *Response {
	r := New(200, 200, msg, content, page...)
	return r
}

func NotFound(msg string, content interface{}, page ...*Page) *Response {
	r := New(404, 404, msg, content, page...)
	return r
}

func Fail(msg string, content interface{}, page ...*Page) *Response {
	r := New(500, 500, msg, content, page...)
	return r
}

func (r *Response) Ext(ext interface{}) *Response {
	r.ExtDat = ext
	return r
}


func (r *Response) SetCode(code int) *Response {
	r.Code = code
	return r
}

func (r *Response) Statu(status int) *Response {
	r.Status = status
	return r
}

func (r *Response) SetMsg(msg string) *Response {
	r.Message = msg
	return r
}

func (r *Response) SetPage(page *Page) *Response {
	r.Data.Page = page
	return r
}

func (r *Response) JSON(ctx iris.Context) {
	r.Timestamp =  time.Now().Unix()
	ctx.StatusCode(r.Code)
	ctx.JSON(r)
}
