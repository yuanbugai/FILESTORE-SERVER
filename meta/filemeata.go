package meta

import (
	mysql "awesomeProject4/db"
	"reflect"
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
func (a FileMeta) IsEmpty() bool {
	return reflect.DeepEqual(a, FileMeta{})
}

// GetFileMeta GetFileMeta:通过sha1值获取文件的元信息对象
func GetFileMeta(FileSha1 string) FileMeta {
	return fileMetas[FileSha1]
}
func GetFileMetadb(filesha1 string) (FileMeta, error) {
	meta, err := mysql.GetFileMeta(filesha1)
	if err != nil {
		return FileMeta{}, err
	}
	fmeta := FileMeta{
		FileSha1: meta.Filehash,
		FileName: meta.Filename,
		FileSize: meta.Filesize,
		Location: meta.Fileaddr,
	}
	return fmeta, nil

}

// RemoveFileMeta : 删除元信息
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
