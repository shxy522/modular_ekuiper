package ossuploader

import (
	"bytes"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"strings"
)

func newAliyunOss(endpoint, accessKeyId, accessKeySecret, bucketName string) *AliyunOss {
	aliOss := &AliyunOss{
		endpoint:        endpoint,
		accessKeyId:     accessKeyId,
		accessKeySecret: accessKeySecret,
		bucketName:      bucketName,
	}
	return aliOss
}

type AliyunOss struct {
	endpoint        string
	accessKeyId     string
	accessKeySecret string
	bucketName      string

	// config map[string]interface{}
}

func handleError(err error) {
	fmt.Println("Error:", err)
}
func login(oss2 *AliyunOss) (*oss.Client, error) {
	return oss.New(oss2.endpoint, oss2.accessKeyId, oss2.accessKeySecret)
}

// 创建BUcket
func (ossClient *AliyunOss) createBucket() {
	// 创建OSSClient实例。
	// client, err := oss.New(oss.endpoint, oss.accessKeyId, oss.accessKeySecret)
	client, err := login(ossClient)
	if err != nil {
		handleError(err)
	}
	// 创建存储空间。
	err = client.CreateBucket(ossClient.bucketName)
	if err != nil {
		handleError(err)
	}
}

// 上传文件
func (ossClient *AliyunOss) uploadFile(objectName string, localFileName string) {
	client, err := login(ossClient)
	if err != nil {
		handleError(err)
	}
	// 获取存储空间。
	bucket, err := client.Bucket(ossClient.bucketName)
	if err != nil {
		handleError(err)
	}
	// 上传文件。
	err = bucket.PutObjectFromFile(objectName, localFileName)
	if err != nil {
		handleError(err)
	}
}

// 上传字符串
func (ossClient *AliyunOss) UploadString(objectName string, uplaodStr string) error {
	client, err := login(ossClient)
	if err != nil {
		handleError(err)
		return err
	}
	// 获取存储空间。
	bucket, err := client.Bucket(ossClient.bucketName)
	if err != nil {
		handleError(err)
		return err
	}
	// 指定Object存储类型为低频访问。
	storageType := oss.ObjectStorageClass(oss.StorageIA)

	// 指定Object访问权限为私有。
	objectAcl := oss.ObjectACL(oss.ACLPrivate)
	// 上传字符串。
	err = bucket.PutObject(objectName, strings.NewReader(uplaodStr), storageType, objectAcl)
	if err != nil {
		handleError(err)
		return err
	}
	return nil
}
func (ossClient *AliyunOss) uploadByte(objectName string, item []byte) {
	client, err := login(ossClient)
	if err != nil {
		handleError(err)
	}
	// 获取存储空间。
	bucket, err := client.Bucket(ossClient.bucketName)
	if err != nil {
		handleError(err)
	}
	// 指定Object存储类型为低频访问。
	// 上传byte
	err = bucket.PutObject(objectName, bytes.NewReader([]byte(item)))
	if err != nil {
		handleError(err)
	}
}

// objectName为bucket上的文件名
// downloadFileName为下载后的文件名
func (ossClient *AliyunOss) downloadFile(objectName string, downloadedFileName string) {
	client, err := login(ossClient)
	if err != nil {
		handleError(err)
	}
	// 获取存储空间。
	bucket, err := client.Bucket(ossClient.bucketName)
	if err != nil {
		handleError(err)
	}
	// 下载文件。
	err = bucket.GetObjectToFile(objectName, downloadedFileName)
	if err != nil {
		handleError(err)
	}
}

// objectName为bucket上的文件名
func (ossClient *AliyunOss) deleteFile(objectName string) {
	client, err := login(ossClient)
	if err != nil {
		handleError(err)
	}
	// 获取存储空间。
	bucket, err := client.Bucket(ossClient.bucketName)
	if err != nil {
		handleError(err)
	}
	// 删除文件。
	err = bucket.DeleteObject(objectName)
	if err != nil {
		handleError(err)
	}
}
func (ossClient *AliyunOss) listFile() {
	client, err := login(ossClient)
	if err != nil {
		handleError(err)
	}
	// 获取存储空间。
	bucket, err := client.Bucket(ossClient.bucketName)
	if err != nil {
		handleError(err)
	}
	// 列举文件。
	marker := ""
	for {
		lsRes, err := bucket.ListObjects(oss.Marker(marker))
		if err != nil {
			handleError(err)
		}
		// 打印列举文件，默认情况下一次返回100条记录。
		for _, object := range lsRes.Objects {
			fmt.Println("Bucket: ", object.Key)
		}
		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}
}
