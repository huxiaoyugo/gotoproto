package generator

import (
	"errors"
	"fmt"
	"github.com/huxiaoyugo/gotoproto/fileWriter"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"github.com/lunny/log"
)



type protoParser struct {
	proto3 *Proto3
}

func NewProtoParser(packageName string) (*protoParser) {
	parse := new(protoParser)
	parse.proto3 = &Proto3{
		packageName:packageName,
		syntax:"proto3",
	}
	return parse
}


func(p *protoParser)Parse(fileName string, src interface{}) error{
	err := p.parseToProto3(fileName, src)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (p*protoParser) ToFile(toFileSrc string) (error){

	writer, err := fileWriter.NewFileWritor(toFileSrc, 0666)
	if err != nil {
		return err
	}
	defer writer.Close()

	str, err := p.proto3.ToProto3String()
	if err != nil {
		log.Error(err)
		return err
	}

	err = writer.WriteFile(str)
	if err != nil {
		log.Error("%v", err)
	}
	return err
}



func(p *protoParser)parseToProto3(fileName string, src interface{}) (outErr error) {

	if p.proto3 == nil {
		p.proto3 = new(Proto3)
	}
	proto3 := p.proto3
	if fileName == "" && src == "" {
		err := errors.New("fileName and src at least one not nil")
		log.Error("%v", err)
		return
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, src, parser.ParseComments)
	if err != nil {
		log.Error("%v", err)
		return
	}

outer:
	for _, item := range f.Decls {
		switch geneTyp := item.(type) {
		case *ast.GenDecl:
			switch geneTyp.Tok {
			case token.CONST:
				constList, err := geneConstFromAstGenDecl(geneTyp)
				if err != nil {
					outErr = err
					break outer
				}
				proto3.constTypes = append(proto3.constTypes, constList...)
			case token.TYPE:
				messageList, serviceList, err := getTypeFromGeneDecl(geneTyp)
				if err != nil {
					outErr = err
					break outer
				}
				proto3.Services = append(proto3.Services, serviceList...)
				proto3.Messages = append(proto3.Messages, messageList...)
			}
		default:
			continue
		}
	}

	if outErr != nil {
		log.Error("%v", outErr)
		return
	}
	return
}

func getTypeFromGeneDecl(val *ast.GenDecl) (messageArr []*Message, serviceArr []*Service, err error) {

	if val.Tok != token.TYPE {
		return
	}
	var modelChnName = getCommentTagName(val)
	// 没有注释标记的话就直接跳过
	if modelChnName == "" {
		return
	}

	// 遍历
	for _, spec := range val.Specs {
		if typSpec, ok := spec.(*ast.TypeSpec); ok {

			switch typ := typSpec.Type.(type) {
			// 结构体
			case *ast.StructType:
				// 模型名称
				message := Message{
					Name:    typSpec.Name.Name,
					Comment: modelChnName,
					Fields:  make([]MessageField, 0),
				}
				for _, field := range typ.Fields.List {
					if len(field.Names) == 0 {
						continue
					}
					messageField := MessageField{
						Name:    field.Names[0].Name,
						Comment: getComment(field.Doc) + getComment(field.Comment),
					}

					// 获取字段的类型
					switch ty := field.Type.(type) {
					case *ast.Ident:
						messageField.Type = ty.Name
					case *ast.ArrayType: // 字段类型为数组
						if eltType, ookk := ty.Elt.(*ast.Ident); ookk {
							messageField.Type = "[]" + eltType.Name
						} else {
							err = errors.New("数组类型的字段元素不能为指针")
							log.Error(err)
							return
						}
					default:
						err = errors.New(fmt.Sprintf("%s cant not be handle", ty))
						log.Error(err)
						return
					}
					message.Fields = append(message.Fields, messageField)
				}
				messageArr = append(messageArr, &message)
				// 接口
			case *ast.InterfaceType:
				// 模型名称
				service := Service{
					Name:    typSpec.Name.Name,
					Comment: modelChnName,
					Funcs:   make([]FuncField, 0),
				}

				for _, method := range typ.Methods.List {

					//method.Comment.
					//method.Comment
					methodComment := getComment(method.Doc) + getComment(method.Comment)
					methodName := method.Names[0].Name
					var methodParam, methodResult string
					if funcType, ok := method.Type.(*ast.FuncType); ok {
						if len(funcType.Params.List) != 1 {
							err = errors.New("接口方法有且只能有一个参数")
							log.Error(err)
							return
						}
						if len(funcType.Results.List) != 1 {
							err = errors.New("接口方法有且只能有一个返回值")
							log.Error(err)
							return
						}
						if ident, ok := funcType.Params.List[0].Type.(*ast.Ident); ok {
							methodParam = ident.Name
						}
						if ident, ok := funcType.Results.List[0].Type.(*ast.Ident); ok {
							methodResult = ident.Name
						}
						service.Funcs = append(service.Funcs, FuncField{
							Name:    methodName,
							Param:   methodParam,
							Return:  methodResult,
							Comment: methodComment,
						})
					}
				}
				serviceArr = append(serviceArr, &service)
			}
		}
	}
	return
}

func getCommentTagName(val *ast.GenDecl) (commentTagName string) {
	// 检查是否有rpc_proto的注释
	if val.Doc == nil {
		return
	}

	for _, comment := range val.Doc.List {
		if strings.Contains(comment.Text, RPC_PROTO_TAG) {
			ar := strings.Split(comment.Text, RPC_PROTO_TAG)
			if len(ar) > 1 {
				commentTagName = strings.Split(comment.Text, RPC_PROTO_TAG)[1]
				break
			}
		}
	}
	commentTagName = strings.Trim(commentTagName, " ")
	commentTagName = strings.Trim(commentTagName, ":")
	commentTagName = strings.Trim(commentTagName, "：")
	return
}

func getComment(group *ast.CommentGroup) string {
	if group == nil {
		return ""
	}
	res := ""
	for _, item := range group.List {
		res += item.Text + " "
	}
	return res
}


func ParseToAst(fileName string, src interface{}) {
	if fileName == "" && src == "" {
		err := errors.New("fileName and src at least one not nil")
		log.Error(err)
		return
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, src, parser.ParseComments)
	if err != nil {
		log.Error(err)
		return
	}
	ast.Print(fset, f)
}
