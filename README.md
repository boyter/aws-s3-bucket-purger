# aws-s3-bucket-purger

A program that will purge any AWS S3 Bucket of objects and versions quickly.

To run,

```
AwsRegion=ap-southeast-2 S3Bucket=my.bucket aws-s3-bucket-purger
```

The above using the default level of concurrency can purge a bucket of 15 million objects/versions in under an
hour, limited by CPU and bandwidth.

To remove specific keys set the `S3Prefix` value,

```
S3Prefix=A AwsRegion=ap-southeast-2 S3Bucket=my.bucket aws-s3-bucket-purger
```

Will remove objects/versions that start with the prefix A.


You can control the default level of concurrency by setting `LoadConcurrency` with a larger integer value.

```
LoadConcurrency=600 AwsRegion=ap-southeast-2 S3Bucket=my.bucket aws-s3-bucket-purger
```

