package fileWriter

import (
	"os"
	"io"
	"github.com/mkideal/log"
)

// 写入方式
type SeekType  int
const (
	Head    SeekType = iota // 头部位置
	Current		// 当前位置
	Append		// 尾部位置
)

type FileWriter struct {
	file *os.File
}

func NewFileWritor(filename string, perm os.FileMode)(*FileWriter, error) {
	file, err :=  os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return nil, err
	}
	fileWriter := &FileWriter{
		file:file,
	}
	return fileWriter, nil
}

func (writer *FileWriter) Close() error {
	if writer.file == nil {
		return nil
	}
	if err := writer.file.Close(); err != nil {
		return err
	}
	return nil
}

// 追加内容
func (writer *FileWriter) AppendFile(content string) error {
	return writer.writeToFile(content, Append)
}

// 继续接入内容
func (writer *FileWriter) WriteFile(content string) error {
	return writer.writeToFile(content, Current)
}

// 从头部写入内容
func (writer *FileWriter) WriteFileFromHead(content string) error {
	return writer.writeToFile(content, Head)
}


func (writer *FileWriter) writeToFile(content string, flag SeekType) error {
	data :=  []byte(content)
	var errOut error
	// 将指针移动到文件末尾
	ret, errOut := writer.file.Seek(0, int(flag))
	log.Trace("seek current position %d", ret)
	if errOut != nil {
		log.Error("%v",errOut)
		return errOut
	}
	// 写入数据
	n, errOut := writer.file.Write(data)
	if errOut == nil && n < len(data) {
		errOut = io.ErrShortWrite
	}
	return errOut

}