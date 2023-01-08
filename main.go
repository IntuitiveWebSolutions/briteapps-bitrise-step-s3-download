package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func main() {
	aws_access_key_id := os.Getenv("aws_access_key_id")
	aws_secret_access_key := os.Getenv("aws_secret_access_key")
	aws_build__internal_id := os.Getenv("aws_build_internal_id")
	ba_ios_icon_hex_color := os.Getenv("ba_ios_icon_hex_color")
	working_dir := os.Getenv("working_dir")
	workspace_dir := working_dir

	aws_build_bucket := "briteapps-builds-output"
	//aws_build__internal_id := "intuitive_web_solutions/2020-11-04_18-16-13_ee806e7a-cd50-4b8f-90fa-619440b775e8"
	dist_loc := aws_build__internal_id + "/out/dist/"
	println(aws_build_bucket)
	println(dist_loc)
	fmt.Println("This is the value specified for the input 'aws_access_key_id':", aws_access_key_id)
	fmt.Println("This is the value specified for the input 'aws_access_key_id':", aws_secret_access_key)

	//
	// --- Step Outputs: Export Environment Variables for other Steps:
	// You can export Environment Variables for other Steps with
	//  envman, which is automatically installed by `bitrise setup`.
	// A very simple example:

	cmdLog, err := exec.Command("bitrise", "envman", "add", "--key", "EXAMPLE_STEP_OUTPUT", "--value", "the value you want to share").CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to expose output with envman, error: %#v | output: %s", err, cmdLog)
		os.Exit(1)
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"), Credentials: credentials.NewStaticCredentials(aws_access_key_id, aws_secret_access_key, "")},
	)

	// Create S3 service client
	//list_buckets(sess, err)

	//download_from := aws_build__internal_id + "/out/cordova"
	cordova_workspace := workspace_dir + "/cordova"

    fmt.Println("Start download IN folder")
	Download_directory_into(aws_build__internal_id + "/in", cordova_workspace, sess)

	fmt.Println("Start download OUT-CORDOVA folder")
	Download_directory_into(aws_build__internal_id + "/out/cordova", cordova_workspace, sess)

    fmt.Println("Start download OUT-DIST folder")
	Download_directory_into(aws_build__internal_id+"/out/dist", cordova_workspace+"/www", sess)

	icon_file := cordova_workspace + "/icon.png"
	splash_file := cordova_workspace + "/splash.png"

    fmt.Println("Start ConvertPossibleJpegToPNG splash_file")
	ConvertPossibleJpegToPNG(splash_file)

	// "#5094D0"
	fmt.Println("Start OverlayImageWithColor icon_file", ba_ios_icon_hex_color)
	OverlayImageWithColor(icon_file, ba_ios_icon_hex_color)

	// You can find more usage examples on envman's GitHub page
	//  at: https://github.com/bitrise-io/envman

	//
	// --- Exit codes:
	// The exit code of your Step is very important. If you return
	//  with a 0 exit code `bitrise` will register your Step as "successful".
	// Any non zero exit code will be registered as "failed" by `bitrise`.
	os.Exit(0)
}

func list_buckets(sess *session.Session, err error) {
	svc := s3.New(sess)
	result, err := svc.ListBuckets(nil)
	if err != nil {
		exitErrorf("Unable to list buckets, %v", err)
	}

	fmt.Println("Buckets:")

	for _, b := range result.Buckets {
		fmt.Printf("* %s created on %s\n",
			aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
	}
}


func get_aws_build_bucket() string {
	aws_build_bucket := "briteapps-builds-output"
	return aws_build_bucket
}

func Download_directory_into(download_from string,	download_to string, sess *session.Session) {
	fmt.Printf("Download_directory_into:  %s -> %s \n", download_from, download_to)

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
		println("------------")
		println("copying from: " + *item.Key)
		new_str := strings.Replace(*item.Key, download_from, "", -1)
		new_path := download_to + new_str
		println("...to: " +new_path)
		ensureDir(new_path)
		file, _ := os.Create(new_path)
		downloader := s3manager.NewDownloader(sess)
		numBytes, err := downloader.Download(file,
			&s3.GetObjectInput{
				Bucket: aws.String(get_aws_build_bucket()),
				Key:    item.Key,
			})
		if err != nil {
		    exitErrorf("Unable to download item %q, %v", item.Key, err)
		}
		println("Bytes copied: " + numBytes)

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