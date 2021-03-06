package meta

import (
	mysql "awesomeProject4/db"
)

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

// UpdateFileMeta UpdateFileMeta:新增/更新文件元信息
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}
func UpdateFileMetadb(fmeta FileMeta) bool {
	return mysql.OnFileUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}

// GetFileMeta GetFileMeta:通过sha1值获取文件的元信息对象
func GetFileMeta(FileSha1 string) FileMeta {
	return fileMetas[FileSha1]
}
