# aws-s3-bucket-purger

A program that will purge any AWS S3 Bucket of objects and versions quickly. Considerably faster then using the AWS CLI tools or doing it though the Console.

The problem.

You have a S3 bucket with millions of files, and those files have versions. Deleting on the CLI takes hours to days, and doesn't do incremental deleting. The console hangs when you click the "Empty bucket" option. This program solves that issue for you by spawning a configurable number but by default hundreds of threads to clean the bucket for you.

Features

 - Fast! Can delete the contents of a S3 bucket with 15 million objects in under an hour given enough CPU and network
 - Resumable. It starts deleting as it finds object keys, not afterwards so if you network crashes you can resume
 - Stateless, so you can run it on different machines and it will pick up from where it was before
 - Supports S3 Key Prefix filtering, so you can fan out if you have billions of objects.
 - MIT Licensed

To run, either install go and use `go run main.go` or download the supplied binary. You configure it by passing in environment variables which makes it easy to run on AWS infrastructure, such as a task on ECS for example.

```
AwsRegion=ap-southeast-2 S3Bucket=my.bucket aws-s3-bucket-purger
```

The above using the default level of concurrency can purge a bucket of 15 million objects/versions in under an
hour, limited by CPU and bandwidth.

To remove specific keys filtered by prefix set the `S3Prefix` value,

```
S3Prefix=A AwsRegion=ap-southeast-2 S3Bucket=my.bucket aws-s3-bucket-purger
```

Will remove objects/versions that start with the prefix A.

Need even more speed? You can control the default level of concurrency by setting `LoadConcurrency` with a larger integer value. Keep in mind this value is set to be lower then the current S3 API limits so you may run into issues if you set this too aggressively.

```
LoadConcurrency=600 AwsRegion=ap-southeast-2 S3Bucket=my.bucket aws-s3-bucket-purger
```

