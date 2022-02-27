package handler

import (
	"encoding/json"
	dblayer "filestore-server/db"
	"filestore-server/meta"
	"filestore-server/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// html
		data, err := ioutil.ReadFile("./static/view/home.html")
		if err != nil {
			io.WriteString(w, "internal error")
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		// 接收文件流及存储到本地目录
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("failed to get data, err:%s\n", err.Error())
			return
		}
		defer file.Close()

		// 创建文件元信息
		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "/tmp/" + head.Filename,
			UploadAt: time.Now().Format("2006-12-01 11:00:00"),
		}

		// 在本地创建文件，用来存储客户发送过来的文件
		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("failed to create file, err:%s\n", err.Error())
			return
		}
		defer newFile.Close()
		// 内存文件，拷贝到新的buffer中
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("failed to save data into file, err:%s\n", err.Error())
			return
		}
		newFile.Seek(0, 0) // !!! 移动到文件的起始位置
		fileMeta.FileSha1 = util.FileSha1(newFile)
		// meta.UpdateFileMeta(fileMeta)
		_ = meta.UpdateFileMetaDB(fileMeta)

		// 更新用户文件表记录
		r.ParseForm()
		username := r.Form.Get("username")
		suc := dblayer.OnUserFileUploadFinished(username, fileMeta.FileSha1,
			fileMeta.FileName, fileMeta.FileSize)
		if suc {
			// 重定向
			http.Redirect(w, r, "/file/upload/suc", http.StatusOK)
		} else {
			w.Write([]byte("Upload Failed!"))
		}
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

// 上传完成
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished!")
}

// 查询文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileHash := r.Form.Get("filehash")
	// fMeta := meta.GetFileMeta(fileHash)
	fMeta, err := meta.GetFileMetaDB(fileHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// 查询批量的文件元信息
func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	limitCnt, _ := strconv.Atoi(r.Form.Get("limit"))
	username := r.Form.Get("username")
	userFiles, err := dblayer.QueryUserFileMetas(username, limitCnt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(userFiles)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// 尝试秒传接口
func TryFastUploadHandler(w http.ResponseWriter, r *http.Request) {
	// 解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filename := r.Form.Get("filename")
	filesize, _ := strconv.Atoi(r.Form.Get("filesize"))

	// 从文件表中查询相同hash的文件记录
	fileMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// 查不到记录则返回秒传失败
	if &fileMeta == nil {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "秒传失败，请访问普通上传接口",
		}
		w.Write(resp.JSONBytes())
		return
	}

	// 上传过则将文件信息写入用户文件表， 返回成功
	suc := dblayer.OnUserFileUploadFinished(
		username, filehash, filename, int64(filesize))
	if suc {
		resp := util.RespMsg{
			Code: 0,
			Msg:  "秒传成功",
		}
		w.Write(resp.JSONBytes())
		return
	}
	resp := util.RespMsg{
		Code: -2,
		Msg:  "秒传失败，请稍后重试",
	}
	w.Write(resp.JSONBytes())
	return
}

// 文件下载
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	// fm := meta.GetFileMeta(fsha1)
	fm, err := meta.GetFileMetaDB(fsha1)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	f, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 下载文件需修的header
	w.Header().Set("Content-type", "application/octect-stream")
	w.Header().Set("Content-Descrption", "attachment;filename=\""+fm.FileName+"\"")
	w.Write(data)
}

// 更新文件元信息
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	opType := r.Form.Get("op")
	filesha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	// 只支持修改文件名
	if opType != "rename" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// 只支持POST
	if r.Method != "POST" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	curFileMeta := meta.GetFileMeta(filesha1)
	curFileMeta.FileName = newFileName
	// meta.UpdateFileMeta(curFileMeta)
	meta.UpdateFileMetaDB(curFileMeta)
	// 以json形式返回元文件信息
	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// 删除文件
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileSha1 := r.Form.Get("filehash")
	fMeta := meta.GetFileMeta(fileSha1)
	os.Remove(fMeta.Location)
	meta.RemoveFileMeta(fileSha1)
	w.WriteHeader(http.StatusOK)
}
