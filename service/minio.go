package service

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/viper"
	"log"
	"mime/multipart"
	"strings"
)

var MinioClient *minio.Client

func init() {
	viper.SetConfigName("configure")
	viper.SetConfigType("json")
	viper.AddConfigPath("$GOPATH/src/github.com/service-computing-2020/bbs_backend/config/")
	viper.AddConfigPath("config/")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	var err error
	endPoint := viper.GetString(`minio.endPoint`)
	accessKeyID := viper.GetString(`minio.accessKeyID`)
	secretAccessKey := viper.GetString(`minio.secretAccessKey`)
	secure := viper.GetBool(`minio.secure`)
	MinioClient, err = minio.New(endPoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: secure,
	})
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("%#v\n", MinioClient)
	}
}

// 传入api路径和文件扩展名，如 getUploadName('/api/users/1/avatar', '.png')
func getUploadName(path string, ext string) string {
	prefix := strings.ReplaceAll(path, "/", "-")
	return fmt.Sprintf("%s%s",prefix, ext)
}

// 同上
func getDownloadName(path string, ext string) string {
	return strings.ReplaceAll(path, "/", "-") + ext
}



func FileUpload(file multipart.File,header *multipart.FileHeader, bucketName string, path string, ext string)(info minio.UploadInfo, err error) {
	ctx := context.Background()

	return MinioClient.PutObject(ctx, bucketName, getUploadName(path, ext), file, header.Size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
}

func FileDownload(filename string, bucketName string, ext string) (*minio.Object, error) {
	ctx := context.Background()
	log.Println(getDownloadName(filename, ext))
	return MinioClient.GetObject(ctx, bucketName, getDownloadName(filename, ext), minio.GetObjectOptions{})
}


