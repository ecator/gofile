package server

import (
	"os"
	"path/filepath"
	"time"
)

type fileInfo struct {
	FileName        string `json:"fileName"`
	UploadTimestamp int64  `json:"uploadTimestamp"`
	ExpireTimestamp int64  `json:"expireTimestamp"`
	Size            int64  `json:"size"`
	Token           string `json:"token"`
}

var fileMap map[string]fileInfo
var tokenMap map[string][]string

func init() {
	fileMap = make(map[string]fileInfo)
	tokenMap = make(map[string][]string)
}

// 根据用户token返回关联的文件token列表
func getFileTokens(token string) []string {
	fileTokens, ok := tokenMap[token]
	if !ok {
		return []string{}
	}
	return fileTokens
}

func addFileInfo(userToken string, fileToken string, file fileInfo) {
	fileTokens := getFileTokens(userToken)
	tokenMap[userToken] = append(fileTokens, fileToken)
	fileMap[fileToken] = file
}

func delFileInfo(userToken string, fileToken string) {
	fileTokens := getFileTokens(userToken)
	target := -1
	for i := 0; i < len(fileTokens); i++ {
		if fileTokens[i] == fileToken {
			target = i
			break
		}
	}
	if target >= 0 {
		fileTokens = append(fileTokens[0:target], fileTokens[target+1:]...)
		tokenMap[userToken] = fileTokens
	}
	delete(fileMap, fileToken)
}

// 物理删除文件
func delFile(fileToken string) bool {
	err := os.Remove(filepath.Join(serverInfo.DataDir, fileToken))
	if err == nil {
		return true
	} else {
		return false
	}
}

// 遍历回收过期的文件
func fileCollect() {
	for userToken, fileTokens := range tokenMap {
		delFileTokens := []string{}
		for _, fileToken := range fileTokens {
			_fileInfo := fileMap[fileToken]
			if _fileInfo.ExpireTimestamp <= time.Now().Unix()*1000 && delFile(fileToken) {
				delFileTokens = append(delFileTokens, fileToken)
			}
		}
		for _, fileToken := range delFileTokens {
			delFileInfo(userToken, fileToken)
		}
	}
}
