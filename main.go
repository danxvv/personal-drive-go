package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type File struct {
	Name     string
	Size     int64
	Content  []byte
	FileType string
}

func NewFile(name string, size int64, content []byte, fileType string) *File {
	return &File{
		Name:     name,
		Size:     size,
		Content:  content,
		FileType: fileType,
	}
}

func main() {
	fileName := "video.mov"
	f, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
	}
	file := NewFile(fileName, int64(len(f)), f, "video/quicktime")
	key := os.Getenv("S3_KEY")
	space := os.Getenv("S3_SPACE")
	fmt.Println(key)
	fmt.Println(space)

	s3config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(key, space, ""),
		Endpoint:    aws.String("fra1.digitaloceanspaces.com"),
		Region:      aws.String("fra1"),
	}
	sess, err := session.NewSession(s3config)
	if err != nil {
		fmt.Println(err)
	}
	svc := s3.New(sess)
	params := &s3.PutObjectInput{
		Bucket:        aws.String("persona-danxvv-drive"),
		Key:           aws.String(file.Name),
		Body:          strings.NewReader(string(file.Content)),
		ContentType:   aws.String(file.FileType),
		ContentLength: aws.Int64(file.Size),
	}
	resp, err := svc.PutObject(params)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
}
