/*
@Date: 2021/1/12 下午2:24
@Author: max.liu
@File : code
@Desc:
*/

package common

type RetCode int

const (
	SUCCESS   RetCode = 0
	FORBIDDEN RetCode = 4030
	FAILED    RetCode = 6000 + iota
	ErrorDatabaseRead
	ErrorDatabaseWrite
	ErrInvalidParams
	ErrInvalidJSONParams
	ErrorPrivilege
	ErrorResourceNotExist
	ErrorCallOtherSrv
	ErrGetList
	ErrGetDetail
	ErrDelete
	ErrCreate
)

var codeMsg = map[RetCode]string{
	SUCCESS:               "成功",
	FAILED:                "失败",
	FORBIDDEN:             "无权限",
	ErrorDatabaseRead:     "查询错误",
	ErrorDatabaseWrite:    "保存失败",
	ErrInvalidParams:      "参数错误",
	ErrInvalidJSONParams:  "参数不是合法的JSON",
	ErrorPrivilege:        "权限错误",
	ErrorResourceNotExist: "资源不存在",
	ErrorCallOtherSrv:     "调用第三方服务异常",
	ErrGetList:            "列表数据查询失败",
	ErrGetDetail:          "详情数据查询失败",
	ErrDelete:             "删除失败",
}

func GetMsg(code RetCode) string {
	return codeMsg[code]
}

func NewErrorCode(code RetCode, err error) CodeWithErr {
	return CodeWithErr{RetCode: code, ErrInfo: err}
}

type CodeWithErr struct {
	RetCode RetCode
	ErrInfo error
}

func (c CodeWithErr) Error() string {
	if c.ErrInfo == nil {
		return ""
	}
	return c.ErrInfo.Error()
}
