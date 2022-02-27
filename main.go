package main

import (
	"filestore-server/handler"
	"fmt"
	"net/http"
)

func main() {
	//处理静态资源映射
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	//url映射
	//文件接口相关
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/download", handler.DownloadHandler)
	http.HandleFunc("/file/update", handler.FileMetaUpdateHandler)
	http.HandleFunc("/file/delete", handler.FileDeleteHandler)
	http.HandleFunc("/file/query", handler.HTTPInterceptor(handler.FileQueryHandler))
	http.HandleFunc("/file/fastupload", handler.TryFastUploadHandler)

	//分块上传接口
	http.HandleFunc("/file/mpupload/init", handler.HTTPInterceptor(handler.InitialMultipartUploadHandler))
	http.HandleFunc("/file/mpupload/uppart", handler.HTTPInterceptor(handler.UploadPartHandler))
	http.HandleFunc("/file/mpupload/complete", handler.HTTPInterceptor(handler.CompleteUploadHandler))

	http.HandleFunc("/user/signup", handler.SignupHandler)
	http.HandleFunc("/user/signin", handler.SignInHandler)
	http.HandleFunc("/user/info", handler.HTTPInterceptor(handler.UserInfoHandler))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Faied to start server, err: ", err)
	}
}
