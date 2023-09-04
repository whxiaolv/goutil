package goutil

import (
	"os"
	"testing"
)

func Test_file(t *testing.T) {
	t.Run("testFile", testFile)
	//t.Run("testUrlEncode", testUrlEncode)
}

var (
	file File
)

const (
	Dir         = "data"
	File1       = "1_file"
	File2       = "2_file"
	File3       = "3_file"
	FileContent = "whxiaolv"
)

func testFile(t *testing.T) {
	file.WriteFile(File1, FileContent)

	//复制文件
	if _, err := file.CopyFile(File1, File2); err != nil {
		t.Errorf("testFile fail : %s", err.Error())
	}

	// 文件内容进行md5 加密
	fileContent, err := os.Open(File1)
	if file.MD5(FileContent) != file.GetFileMd5(fileContent) || err != nil {
		t.Errorf("testFile fail : %s", err.Error())
	}

	// 文件内容sha1 加密
	if file.Sha1(FileContent) != file.GetFileSha1Sum(fileContent) || err != nil {
		t.Errorf("testFile fail : %s", err.Error())
	}

	// 选择加密形式 sha1md5/加密
	file.GetFileSum(fileContent, "sha1")

	// 通过文件路径获取内容 sha/1md5 加密
	fileSumMd5, err := file.GetFileSumByName(File1, "")
	if err != nil || fileSumMd5 != file.GetFileMd5(fileContent) {
		t.Errorf("testFile fail : %s", err.Error())
	}

	// 写入文件内容,字符串 string 【清空文件数据】
	if !file.WriteFile(File3, "hello") {
		t.Errorf("testFile fail ")
	}

	// 写入文件二进制内容，字节 byte 【清空文件数据】
	if !file.WriteBinFile(File3, []byte("hello")) {
		t.Errorf("testFile fail ")
	}

	// 读取二进制文件
	binContent, err := file.ReadBinFile(File3)
	if err != nil || len(binContent) != 5 {
		t.Errorf("testFile fail : %s", err.Error())
	}

	// 判断文件存在
	if !file.FileExists(File3) {
		t.Errorf("testFile fail ")
	}

	// 通过偏移量写入文件
	if err := file.WriteFileByOffSet(File3, 5, []byte("hello")); err != nil {
		t.Errorf("testFile fail : %s", err.Error())
	}

	// 通过偏移量读取文件
	offsetContent, err := file.ReadFileByOffSet(File3, 5, 5)
	if err != nil || len(offsetContent) != 5 {
		t.Errorf("testFile fail : %s", err.Error())
	}

	// 文件是否存在
	if !file.IsExist(File3) {
		t.Errorf("testFile fail ")
	}

	// 删除指定目录下的空目录
	go file.RemoveEmptyDir(Dir)

}
