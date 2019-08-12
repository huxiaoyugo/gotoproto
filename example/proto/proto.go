package proto

// rpc_proto: 接口
type bg interface {
	IOTCardRecharge(IOTCardRechargeModel) CmsResp
	IOTCardRefresh(IOTCardModel) CmsResp
	IOTCardImport(IOTCardImportModel) CmsResp
	// 发送钉钉消息
	DingTalkSystemMsg(DingTalkSystemMsgModel) CmsResp
}

// rpc_proto: 物联网卡充值
type IOTCardRechargeModel struct {
	Iccid        string
	IsAutoEffect bool
}

// rpc_proto: 物联网卡
type IOTCardModel struct {
	Iccid string
}

// rpc_proto: 导入物联网卡
type IOTCardImportModel struct {
	Data string
}

// rpc_proto: 钉钉消息模型
type DingTalkSystemMsgModel struct {
	AdminName string
	Msg       string
}

type ReturnStatus int

// rpc_proto: 返回状态
const (
	StatusFail            ReturnStatus = 0 //处理失败
	StatusSuccess         ReturnStatus = 1 //成功
	StatusArgumentInvalid ReturnStatus = 2 //参数错误
	StatusAuthVerifyFail  ReturnStatus = 3 //访问校验失败
)

// rpc_proto: 通用返回模型
type CmsResp struct {
	Status ReturnStatus
	Desc   string
}
