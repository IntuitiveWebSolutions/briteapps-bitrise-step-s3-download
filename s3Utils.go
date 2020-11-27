package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"path/filepath"
	"strings"
)

func get_aws_build_bucket() string {
	aws_build_bucket := "briteapps-builds-output"
	return aws_build_bucket
}

func Download_directory_into(download_from string,	download_to string, sess *session.Session) {
	svc := s3.New(sess)

	params := &s3.ListObjectsInput{
		Bucket: aws.String(get_aws_build_bucket()),
		Prefix: aws.String(download_from),
	}

	resp, err := svc.ListObjects(params)
	if err != nil {
		exitErrorf("Unable to list buckets, %v", err)
	}

	for _, item := range resp.Contents {
		//fmt.Println("Name:         ", *item.Key)
		//fmt.Println("Last modified:", *item.LastModified)
		//fmt.Println("Size:         ", *item.Size)
		//fmt.Println("Storage class:", *item.StorageClass)
		//fmt.Println("")
		new_str := strings.Replace(*item.Key, download_from, "", -1)
		new_path := download_to + new_str
		println("downloading to: " +new_path)
		ensureDir(new_path)
		file, _ := os.Create(new_path)
		downloader := s3manager.NewDownloader(sess)
		numBytes, _ := downloader.Download(file,
			&s3.GetObjectInput{
				Bucket: aws.String(get_aws_build_bucket()),
				Key:    item.Key,
			})

		println(numBytes)

	}
}
func ensureDir(fileName string) {
	dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, os.ModePerm)
		if merr != nil {
			panic(merr)
		}
	}
}