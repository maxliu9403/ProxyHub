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
	ErrDeleteGroup
	ErrCreateGroup
	ErrCreateProxy
	ErrUpdateGroup
	ErrUpdateProxy
	ErrDeleteProxy
	ErrCreateProxyNotGroup
	ErrCreateProxyCheckGroup
	ErrCreateToken
	ErrBuildToken
	ErrDeleteToken
	ErrValidateToken
	ErrDeleteEmulator
	ErrQueryExistEmulatorUUID
	ErrCreateEmulator
	ErrUpdateEmulator
)

var codeMsg = map[RetCode]string{
	SUCCESS:                   "成功",
	FAILED:                    "失败",
	FORBIDDEN:                 "无权限",
	ErrorDatabaseRead:         "查询错误",
	ErrorDatabaseWrite:        "保存失败",
	ErrInvalidParams:          "参数错误",
	ErrInvalidJSONParams:      "参数不是合法的JSON",
	ErrorPrivilege:            "权限错误",
	ErrorResourceNotExist:     "资源不存在",
	ErrorCallOtherSrv:         "调用第三方服务异常",
	ErrGetList:                "列表数据查询失败",
	ErrGetDetail:              "详情数据查询失败",
	ErrDeleteGroup:            "删除组失败",
	ErrCreateGroup:            "创建分组失败",
	ErrUpdateGroup:            "更新分组失败",
	ErrUpdateProxy:            "更新代理失败",
	ErrCreateProxy:            "创建代理失败",
	ErrDeleteProxy:            "删除代理失败",
	ErrCreateProxyNotGroup:    "创建代理失败，不存在激活状态的分组，请先创建分组",
	ErrCreateProxyCheckGroup:  "创建代理失败，校验是否存在分组失败",
	ErrCreateToken:            "创建Token失败",
	ErrBuildToken:             "随机生成Token失败",
	ErrDeleteToken:            "删除Token失败",
	ErrValidateToken:          "校验Token失败",
	ErrDeleteEmulator:         "删除模拟器失败",
	ErrQueryExistEmulatorUUID: "查询已存在的UUID失败",
	ErrCreateEmulator:         "创建模拟器失败",
	ErrUpdateEmulator:         "更新模拟器失败",
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
