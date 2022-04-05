package server

import (
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ecator/gofile/util"
	"github.com/julienschmidt/httprouter"
)

// 相对于webdir直接返回静态页面
func resStaticPage(w http.ResponseWriter, pagePath string) {
	buf, err := ioutil.ReadFile(filepath.Join(serverInfo.WebDir, pagePath))
	if err != nil {
		http.Error(w, err.Error(), 500)
	} else {
		w.Write(buf)
	}
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	resStaticPage(w, "404.html")
}

func handleIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fileCollect()
	getToken(w, r)
	requestFile := r.RequestURI[1:]
	if requestFile == "" {
		requestFile = "index.html"
	}
	resStaticPage(w, requestFile)
}

func handleFileDown(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fileCollect()
	fileToken := ps.ByName("token")
	file, err := os.Open(filepath.Join(serverInfo.DataDir, fileToken))
	if err != nil {
		handleNotFound(w, r)
		return
	}
	defer file.Close()
	_fileInfo := fileMap[fileToken]
	fileHeader := make([]byte, 512)
	file.Read(fileHeader)
	fileStat, _ := file.Stat()
	fileName := url.QueryEscape(_fileInfo.FileName)
	if fileName == "" {
		fileName = strconv.FormatInt(time.Now().Unix(), 10)
	}
	fileName = _fileInfo.FileName
	//w.Header().Set("Content-Type", http.DetectContentType(fileHeader))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))
	w.Header().Add("Content-Disposition", "attachment; filename=\""+fileName+"\"")
	file.Seek(0, 0)
	io.Copy(w, file)

}

func handleFileUp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	token := getToken(w, r)
	fileToken := util.Md5FromStr(strconv.Itoa(rand.Intn(1000)+2000) + token + r.UserAgent())
	_fileInfo := fileInfo{}
	jData := new(jsonData)
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		jData.Code = 500
		jData.Result = err.Error()
		resJson(w, jData)
		return
	}
	defer file.Close()
	saveFile, err := os.Create(filepath.Join(serverInfo.DataDir, fileToken))
	if err != nil {
		jData.Code = 500
		jData.Result = err.Error()
		resJson(w, jData)
		return
	}
	defer saveFile.Close()
	if _, err := io.Copy(saveFile, file); err != nil {
		jData.Code = 500
		jData.Result = err.Error()
		resJson(w, jData)
		return
	}

	// 添加映射关系
	_fileInfo.FileName = fileHeader.Filename
	_fileInfo.Size = fileHeader.Size
	_fileInfo.UploadTimestamp = time.Now().Unix() * 1000
	_fileInfo.ExpireTimestamp = _fileInfo.UploadTimestamp + (int64(serverInfo.ExpireSec * 1000))
	_fileInfo.Token = fileToken
	addFileInfo(token, fileToken, _fileInfo)

	// 返回状态
	jData.Code = 0
	jData.Result = _fileInfo
	resJson(w, jData)
}

// 根据用户token返回关联的文件列表
func handleFileInfos(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	token := getToken(w, r)
	jData := new(jsonData)
	jData.Code = 0
	fileInfos := []fileInfo{}
	fileTokens := getFileTokens(token)
	for _, v := range fileTokens {
		_fileInfo, ok := fileMap[v]
		if ok {
			fileInfos = append(fileInfos, _fileInfo)
		}
	}
	jData.Result = fileInfos

	resJson(w, jData)
}

// 根据文件token删除文件
func handleFileDel(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	token := getToken(w, r)
	fileToken := ps.ByName("token")
	var fileInfo interface{}
	jData := new(jsonData)
	jData.Code = 0
	if delFile(fileToken) {
		fileInfo = fileMap[fileToken]
		delFileInfo(token, fileToken)
	} else {
		jData.Code = 404
		fileInfo = "not found"
	}
	jData.Result = fileInfo
	resJson(w, jData)
}

// 根据文件token设置过期时间戳
func handleFileExpireTimestamp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fileToken := ps.ByName("token")
	r.ParseForm()
	expireTimestampStr := r.Form.Get("expireTimestamp")
	jData := new(jsonData)
	jData.Code = 0
	jData.Result = "ok"
	expireTimestamp := 0
	var err error
	if expireTimestampStr != "" {
		expireTimestamp, err = strconv.Atoi(expireTimestampStr)
		if err != nil {
			jData.Code = 500
			jData.Result = err.Error()
		}
	}
	if err == nil {
		if setFileExpireTimestamp(fileToken, int64(expireTimestamp)) == false {
			jData.Code = 404
			jData.Result = "not found"
		}
	}

	resJson(w, jData)
}
