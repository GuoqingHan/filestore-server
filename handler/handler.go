package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// html
		data, err := ioutil.ReadFile("./static/view/index.html")
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
		// 在本地创建文件，用来存储客户发送过来的文件
		newFile, err := os.Create("/tmp/" + head.Filename)
		if err != nil {
			fmt.Printf("failed to create file, err:%s\n", err.Error())
			return
		}
		defer newFile.Close()
		// 内存文件，拷贝到新的buffer中
		_, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("failed to save data into file, err:%s\n", err.Error())
			return
		}
		// 重定向
		http.Redirect(w, r, "/file/upload/suc", http.StatusOK)
	}
}

// 上传完成
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished!")
}
