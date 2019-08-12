package generator

import (
	"fmt"
	"strings"
	"github.com/pkg/errors"
	"github.com/mkideal/log"
)

// 消息
/*
// 用户信息
message UserInfo {
    int32 Uid = 1;
    string SessionId = 2;
    string UserName = 3;
    string RoleName = 4;
    string Phone = 5;
    int32 IsCenter = 6;
}
 */
type Message struct {
	Name    string
	Comment string
	Fields  []MessageField
}

type MessageField struct {
	Name string
	Type string
	Comment string
}


// 枚举
/*
//返回的状态值
enum ReturnStatus {
    StatusFail = 0; //处理失败
    StatusSuccess = 1; //成功
    StatusArgumentInvalid = 2; //参数错误
    StatusAuthVerifyFail = 3; //访问校验失败
}
 */
type Enum struct {
	Name string
	Comment string
	Fields []EnumField
}

type EnumField struct {
	Name string
	Val int
	Comment string // 备注
}

// rpc服务
/*
service bg {
    rpc IOTCardRecharge (IOTCardRechargeModel) returns (CmsResp);
    rpc IOTCardRefresh (IOTCardModel) returns (CmsResp);
    rpc IOTCardImport (IOTCardImportModel) returns (CmsResp);
    rpc DingTalkSystemMsg (DingTalkSystemMsgModel) returns (CmsResp);
    rpc MachineList(MachineListReq) returns(MachineListResp);
}
 */
type Service struct {
	Name string
	Funcs []FuncField
	Comment string //
}

type FuncField struct {
	Name string
	Param  string
	Return string
	Comment string
}

// proto
type Proto3 struct {
	packageName string
	syntax string
	Services []*Service
	Messages []*Message
	Enums []*Enum
	constTypes []ConstType
}


func (p *Proto3)getProtoFileHead() string {
	return fmt.Sprintf("syntax = \"%s\";\npackage %s;\n\n", p.syntax, p.packageName)
}

func (p *Proto3) getServiceStr() (res string, err error) {
	if len(p.Services) >1 {
		return "", errors.New("Service数量超过1")
	}

	if len(p.Services) == 0 {return}

	return GeneServiceStr(p.Services[0])
}

func (p *Proto3) getMessageStr()(res string, err error) {
	for _, message := range p.Messages {
		tmp := ""
		tmp, err = GeneMessageStr(message)
		if err != nil {return }
		res += tmp
	}
	return
}

func (p *Proto3) getEnumStr()(res string, err error) {

	// 需先将constType ==> []Enum
	p.geneEnumFromConstType()
	for _, enum := range p.Enums {
		tmp := ""
		tmp, err = GeneEnumStr(enum)
		if err != nil {return }
		res += tmp
	}
	return
}


func (p *Proto3) geneEnumFromConstType() {
	p.Enums = make([]*Enum, 0)
	for _, item := range p.constTypes {
		var enum *Enum

		for _, enu := range p.Enums {
			if enu.Name == item.Type {
				enum =  enu
				break
			}
		}
		if enum == nil {
			enum = &Enum{Name: item.Type}
			p.Enums = append(p.Enums, enum)
		}
		enum.Fields = append(enum.Fields, EnumField {
			Name:item.Name,
			Val:item.Val,
			Comment:item.Comment,
		})
	}
}


func (p *Proto3) ToProto3String() (string, error) {
	head := p.getProtoFileHead()
	service, err := p.getServiceStr()
	if err != nil {
		return "", err
	}
	message, err := p.getMessageStr()
	if err != nil {
		return "", err
	}

	enum, err := p.getEnumStr()
	if err != nil {
		return "", err
	}
	return head + service + message + enum, nil
}






func GeneServiceStr(service *Service) (string, error) {

	if service == nil {
		log.Warn("service is nil")
		return "", nil
	}
	res := fmt.Sprintf("// %s\nservice %s {\n", service.Comment, service.Name)
	for _, field := range service.Funcs {
		res += fmt.Sprintf("    rpc %s (%s) returns (%s);\n", field.Name, field.Param, field.Return)
	}
	res += "}\n\n"
	return res, nil
}


func GeneEnumStr(enum *Enum) (string, error) {
	res := fmt.Sprintf("%s\nenum %s {\n", enum.Comment, enum.Name)
	for _, field := range enum.Fields {
		res += fmt.Sprintf("    %s = %d; %s\n", field.Name, field.Val, field.Comment)
	}
	res += "}\n\n"
	return res, nil
}

func GeneMessageStr(msg *Message) (string, error) {

	res := fmt.Sprintf("// %s\nmessage %s {\n", msg.Comment, msg.Name)

	for index, field := range msg.Fields {
		res +=  "    "
		fieldType := field.Type
		isArr := strings.Contains(fieldType, "[]")
		if isArr { // 是数组
			fieldType = strings.Split(fieldType,"[]")[1]
			res += "repeated "
		}

		// 类型
		if typ, ok := varTypeToProto3Map[lowerStr(fieldType)]; ok {
			res += typ
		} else {
			res += fieldType
		}
		// 字段名称
		res += fmt.Sprintf(" %s = %d; \n", ToUnderLine(field.Name), index+1)
	}
	res += "}\n\n"
	return res, nil
}