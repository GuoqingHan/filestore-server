package meta

import (
	mydb "filestore-server/db"
)

// 文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// 更新或插入文件元信息：内存方式
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

// 更新或插入文件元信息：数据库方式
func UpdateFileMetaDB(fmeta FileMeta) bool {
	return mydb.OnFileUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}

// 获取文件元信息
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

// 获取文件元信息：数据库方式
func GetFileMetaDB(fileSha1 string) (FileMeta, error) {
	tf, err := mydb.GetFileMeta(fileSha1)
	if err != nil {
		return FileMeta{}, err
	}
	fm := FileMeta{
		FileSha1: tf.FileHash,
		FileName: tf.FileName.String,
		FileSize: tf.FileSize.Int64,
		Location: tf.FileAddr.String,
	}
	return fm, err
}

// 删除文件元信息
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
