package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)


const (
	AWS_S3_REGION = ""
	AWS_S3_BUCKET = ""
)

var sess = connectAWS()

func connectAWS() *session.Session {
	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(AWS_S3_REGION)
		}
	)
	if err != nil {
		panic(err)
	}
	return sess
}

func upload() {
	file, header, err := r.FormFile("file")
	if err != nil {
		// Do your error handling here
		return
	}
	defer file.Close()

	filename := header.Filename
	uploader := s3manager.NewUploader(sess)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(AWS_S3_BUCKET), // Bucket to be used
		Key:    aws.String(filename),      // Name of the file to be saved
		Body:   file,                      // File
	})
	if err != nil {
		// Do your error handling here
		return
	}
}

func download() {
	f, err := os.Create(filename)
	if err != nil {
		// Do your error handling here
		return
	}

	downloader := s3manager.NewDownloader(sess)
	_, err = downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(AWS_S3_BUCKET),
		Key:    aws.String(filename),
	})
	if err != nil {
		// Do your error handling here
		return
	}
}

func list() {
	svc := s3.New(sess)
	input := &s3.ListObjectsInput{
		Bucket: aws.String(AWS_S3_BUCKET),
	}

	result, err := svc.ListObjects(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				fmt.Println(s3.ErrCodeNoSuchBucket, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			fmt.Println(err.Error())
		}
		return
	}
	w.Header().Set("Content-Type", "text/html")
	for _, item := range result.Contents {
		fmt.Fprintf(w, "<li>File %s</li>", *item.Key)
	}
}