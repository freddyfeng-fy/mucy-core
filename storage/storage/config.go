package storage

type DiskName string

const (
	Local DiskName = "local" // 本地
	S3    DiskName = "s3"    // S3
	Oss   DiskName = "oss"   // 阿里云
)
