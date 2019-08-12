package generator

import (
	"go/ast"
	"go/token"
	"github.com/pkg/errors"
	"github.com/mkideal/log"
	"strconv"
)


// 常量的值只允许为int类型的
type ConstType struct {
	Name string
	Type string
	Val  int
	Comment string
}



func geneConstFromAstGenDecl(genDecl *ast.GenDecl) (list []ConstType, err error) {

	if genDecl.Tok != token.CONST {
		return list, errors.New("token is not const")
	}

	var modelChnName = getCommentTagName(genDecl)
	// 没有注释标记的话就直接跳过
	if modelChnName == "" {return}
	/*
		为了简化，做一下限制
		必须要保证type val 都得明确声明，且不允许使用iota
	 */

	 for index, spec := range genDecl.Specs {
	 	// 注释
	 	// 先判断是否有type
	 	var constType = ""
		switch valSpec := spec.(type) {
		case *ast.ValueSpec:

			constComment := getComment(valSpec.Doc) + getComment(valSpec.Comment)

			if valSpec.Type == nil { // 判断是否有type
				log.Warn("pos:%v const genDecl.Specs[%d] dont have type",  genDecl.TokPos, index,)
				continue
			}
			// 获取type的名称
			if ident, ok := valSpec.Type.(*ast.Ident); ok {
				constType = ident.Name
			} else {
				log.Warn("pos:%v const genDecl.Specs[%d] type is not *ast.Ident, but %s", genDecl.TokPos, index, ident)
				continue
			}

			// 开始获取常量名和val
			count := len(valSpec.Names)
			for i:=0; i< count; i++ {

				constName := valSpec.Names[i].Name
				valBasicLit := valSpec.Values[i]
				if ident, ok := valBasicLit.(*ast.BasicLit); ok {
					if ident.Kind != token.INT {
						log.Warn("pos:%v const genDecl.Specs[%d] values[%d] is not int, but token(%v)", genDecl.TokPos, index,i, ident.Kind)
						break
					}
					// 到此处，才满足了所有的条件
					constVal := 0
					constVal,err = strconv.Atoi(ident.Value)
					if err != nil {
						log.Error("pos:%v const genDecl.Specs[%d] values[%d].value to int error:%v", genDecl.TokPos, index,i, err)
						return
					}
					// 创建constType 添加到返回的list中
					list = append(list, ConstType{
						Name: constName,
						Type: constType,
						Val: constVal,
						Comment: constComment,
					})
				} else {
					log.Warn("pos:%v const genDecl.Specs[%d] values[%d] is not *ast.Ident, but %v", genDecl.TokPos, index,i, ident)
					break
				}
			}

		default:
			log.Warn("pos:%v const genDecl.Specs[%d] is not *ast.ValueSpec, but %v", genDecl.TokPos, index, valSpec)
			continue
		}
	 }
	return
}





