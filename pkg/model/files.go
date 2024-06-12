package model

type FileWriteInfo struct {
	FileNo     int32 `json:"fileNumber"`
	BytesWrite int64 `json:"bytesWrite"`
}
