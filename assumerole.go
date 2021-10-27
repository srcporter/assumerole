package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/aws/aws-sdk-go/service/s3/s3manager"
    "github.com/aws/aws-sdk-go/aws/endpoints"
    "fmt"
    "os"
    "time"
    "bytes"
)

func exitErrorf(msg string, args ...interface{}) {
    fmt.Fprintf(os.Stderr, msg+"\n", args...)
    os.Exit(1)
}


func main() {

    // Resolve endpoint for S3 in us-east-1
    resolver := endpoints.DefaultResolver()
    endpoint, err := resolver.EndpointFor(endpoints.S3ServiceID, endpoints.UsEast1RegionID)
    if err != nil {
        fmt.Println("failed to resolve endpoint", err)
        return
    }
    fmt.Println("Resolved URL:", endpoint.URL)

    //endpoint := fmt.Sprintf("s3.%s.amazonaws.com", "us-east-1")
    //endpoint := ""
    awsConfig := &aws.Config{
        //Endpoint: aws.String(endpoint.URL),
        Region: aws.String("us-east-1"),
        }
    sess, err := session.NewSession(awsConfig)

    fmt.Println("Endpoint: ", awsConfig.Endpoint) 

    // Create S3 service client
    svc := s3.New(sess)

    bresult, err := svc.ListBuckets(nil)
    if err != nil {
        exitErrorf("Unable to list buckets, %v", err)
    }

    fmt.Println("Buckets:")

    for _, b := range bresult.Buckets {
        fmt.Printf("* %s created on %s\n",
        aws.StringValue(b.Name), aws.TimeValue(b.CreationDate))
    }

    // upload something

    file, err := os.Open("test.txt")
    if err != nil {
        return
    }
    defer file.Close()

    fileInfo, _ := file.Stat()
    var size = fileInfo.Size()
    buffer := make([]byte, size)
    file.Read(buffer)

    uploader := s3manager.NewUploader(sess)
    upParams := &s3manager.UploadInput{
        Bucket: aws.String("chris-stackrox-backup"),
        Key:    aws.String("text.txt"),
        Body:   bytes.NewReader(buffer),
    }
    uresult, err := uploader.Upload(upParams)
    if err != nil {
        return
    }
    fmt.Printf("new object created at %s\n", uresult.Location)

    time.Sleep(3600 * time.Second)

}
