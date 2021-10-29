package main

import (
    "log"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/aws/credentials/stscreds"
    "github.com/aws/aws-sdk-go/service/ecr"
)

func main() {

    awsConfig := &aws.Config{
        //Endpoint: aws.String(endpoint.URL),
        Region: aws.String("us-west-2"),
        }
    sess, err := session.NewSession(awsConfig)
    if (err != nil) {
        log.Fatal(err)
    }
    rolearn := "arn:aws:iam::044192981154:role/cporter-ecr-readonly"
    creds := stscreds.NewCredentials(sess, rolearn)

    // connect to ecr
    service := ecr.New(sess, &aws.Config{Credentials: creds})

    repoName := "cporter-ecr"
    regId := "044192981154"
    imageSearch := &ecr.DescribeImagesInput{RegistryId: &regId, RepositoryName: &repoName}
    //imageSearch = ecr.SetRepositoryName("cporter-ecr")

    //var imagesDescription *DescribeImagesOutput
    imagesDescription, describe_err := service.DescribeImages (imageSearch)
    if (describe_err != nil) {
        log.Fatal(describe_err)
    }
    log.Println("First Image in Repo:")
    log.Println(imagesDescription.ImageDetails[0].ImageDigest)

}
