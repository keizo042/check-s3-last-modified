package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var (
	VERSION string
)

const (
	DEFAULT_REGION   = "ap-northeast-1"
	DEFAYLT_INTERVAL = 60 * 60 * 24
)

type Config struct {
	Region   *string
	Bucket   *string
	Interval *int
	ID       *string
	Secret   *string
	Token    *string
	Version  *bool
	Prefix   *string
}

func main() {
	var c Config
	c.Bucket = flag.String("bucket", "", "AWS S3 bucket name without s3://")
	c.ID = flag.String("id", "", "AWS Account Key ID")
	c.Secret = flag.String("secret", "", "AWS Secret Access Token")
	c.Token = flag.String("token", "", "AWS session token(optional)")
	c.Region = flag.String("region", DEFAULT_REGION, "AWS region")
	c.Version = flag.Bool("version", false, "version")
	c.Prefix = flag.String("prefix", "", "AWS s3 directory prefix, not path")

	c.Interval = flag.Int("interval", DEFAYLT_INTERVAL, "interval seconds until S3 last modified")
	flag.Usage = func() {
		fmt.Println("Description: check last modified object in AWS S3 folder  until  -interval seconds\n")
		fmt.Printf("Version: %s\n", VERSION)
		flag.PrintDefaults()
	}
	flag.Parse()

	if *c.Version {
		flag.Usage()
		return
	}
	if *c.Bucket == "" {
		fmt.Println("option: -bucket required")
		flag.Usage()
		os.Exit(1)
		return
	}
	if *c.ID == "" {
		fmt.Println("option: -id required")
		flag.Usage()
		os.Exit(1)
		return
	}
	if *c.Secret == "" {
		fmt.Println("option: -secret required")
		flag.Usage()
		os.Exit(1)
		return
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(*c.Region),
		Credentials: credentials.NewStaticCredentials(*c.ID, *c.Secret, *c.Token),
	})
	if err != nil {
		fmt.Printf("WARNING: error: %v\n", err)
		os.Exit(1)
	}
	svc := s3.New(sess)
	input := &s3.ListObjectsInput{
		Bucket: aws.String(*c.Bucket),
		Prefix: aws.String(*c.Prefix),
	}

	output, err := svc.ListObjects(input)
	if err != nil {
		fmt.Printf("WARNING: error: %v\n", err)
		os.Exit(1)
		return
	}
	interval := time.Duration(*c.Interval) * time.Second
	period := time.Now().Add(-interval)
	for _, obj := range output.Contents {
		if obj.LastModified.After(period) {
			fmt.Printf("OK: last modified at %s in s3://%s\n",
				obj.LastModified.String(),
				path.Join(*c.Bucket, *c.Prefix))
			return
		}
	}
	fmt.Printf("WARNING: not found modified object until %s\n in s3://%s\n",
		period.String(),
		path.Join(*c.Bucket, *c.Prefix))

	os.Exit(1)
	return

}
