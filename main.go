package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"os"
	"strconv"
	"sync"
)

type VersionKey struct {
	Name      string
	VersionId string
}

// Deletes data from bucket including versions
// Example usage would be as follows
// AwsRegion=ap-southeast-2 S3Bucket=my.bucket S3Prefix=191 go run main.go
// Which will load delete all records with the bucket with the prefix 191
func main() {
	purge()
}

func purge() {
	svc, err := session.NewSession(&aws.Config{
		Region: aws.String(getEnvString("AwsRegion", ""))},
	)

	if err != nil {
		log.Println(err)
		return
	}

	var wg sync.WaitGroup
	s3client := s3.New(svc)
	s3ListQueue := make(chan string, 100000)
	s3VersionListQueue := make(chan VersionKey, 100000)

	// Get the keys from S3
	wg.Add(1)
	go func() {
		err = s3client.ListObjectsPages(&s3.ListObjectsInput{
			Bucket: aws.String(getEnvString("S3Bucket", "")),
			Prefix: aws.String(getEnvString("S3Prefix", "")),
		}, func(page *s3.ListObjectsOutput, lastPage bool) bool {
			for _, value := range page.Contents {
				s3ListQueue <- *value.Key
			}
			return true
		})

		if err != nil {
			log.Println(err)
		}

		close(s3ListQueue)
		wg.Done()
	}()

	// Get all the versions from S3
	wg.Add(1)
	go func() {
		err = s3client.ListObjectVersionsPages(&s3.ListObjectVersionsInput{
			Bucket: aws.String(getEnvString("S3Bucket", "")),
			Prefix: aws.String(getEnvString("S3Prefix", "")),
		}, func(page *s3.ListObjectVersionsOutput, lastPage bool) bool {
			for _, value := range page.Versions {
				s3VersionListQueue <- VersionKey{
					Name:      *value.Key,
					VersionId: *value.VersionId,
				}
			}

			for _, value := range page.DeleteMarkers {
				s3VersionListQueue <- VersionKey{
					Name:      *value.Key,
					VersionId: *value.VersionId,
				}
			}

			return true
		})

		if err != nil {
			log.Println(err)
		}

		close(s3VersionListQueue)
		wg.Done()
	}()

	for i := 0; i < getEnvInt("LoadConcurrency", 300); i++ {
		wg.Add(1)
		go func(input chan string) {
			for key := range input {
				log.Println(fmt.Sprintf("purge::key:%s", key))

				s3client.DeleteObject(&s3.DeleteObjectInput{
					Bucket: aws.String(getEnvString("S3Bucket", "")),
					Key:    aws.String(key),
				})
			}
			wg.Done()
		}(s3ListQueue)
	}

	for i := 0; i < getEnvInt("LoadConcurrency", 300); i++ {
		wg.Add(1)
		go func(input chan VersionKey) {
			for key := range input {
				log.Println(fmt.Sprintf("purge::key:%s", key))

				s3client.DeleteObject(&s3.DeleteObjectInput{
					Bucket:    aws.String(getEnvString("S3Bucket", "")),
					Key:       aws.String(key.Name),
					VersionId: aws.String(key.VersionId),
				})
			}
			wg.Done()
		}(s3VersionListQueue)
	}

	log.Println(fmt.Sprintf("starting::bucket:%s::prefix:%s", getEnvString("S3Bucket", ""), getEnvString("S3Prefix", "")))
	wg.Wait()
	log.Println(fmt.Sprintf("finished::bucket:%s::prefix:%s", getEnvString("S3Bucket", ""), getEnvString("S3Prefix", "")))
}

func getEnvString(variable string, def string) string {
	val := os.Getenv(variable)

	if val != "" {
		return val
	}

	return def
}

func getEnvInt(variable string, def int) int {
	tmp := os.Getenv(variable)

	if tmp != "" {
		val, err := strconv.Atoi(tmp)

		if err == nil {
			return val
		}
	}

	return def
}
