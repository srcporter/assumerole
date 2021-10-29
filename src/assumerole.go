package main

import (
    "log"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/aws/credentials/stscreds"
    "github.com/aws/aws-sdk-go/service/ecr"
)

/*

Example roles and policies

Assume we have two AWS account IDs:
12345: this account is the "main" where EKS / Central is running
67890: this account has the ECR "target" we want to access

In account 67890:
- Create an IAM role "remote-ecr-reader"* that allows read-only access to ECR
- In that IAM role, establish trust of the "assumerole-to-remote-ecr" role in 12345.

In account 12345:
- Create an IAM role "assumerole-to-remote-ecr" that allows AssumeRole to the arn for "remote-ecr-reader" in 67890

For additional ECR accounts:
1. target account must have a "remote-ecr-reader" role with trust established for 12345:role/assumerole-to-remote-ecr
2. main account 12345 must add a clause to allow AssumeRole to the arn created in #1
Repeat for any  new ECR target account

Lastly, any code that will perform the AssumeRole (like this hack, or ACS Central):
- Must be told to use the "assumerole-to-remote-ecr" arn
- Must be provided access to this role via an EC2 instance role, a serviceaccount-linked role,
    or secret/access key for a user account with this role


Example policy attached to role  "assumerole-to-remote-ecr" in 12345 that allows it to assumerole in two remote accounts:

{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "sts:AssumeRole",
            "Resource": "arn:aws:iam::67890:role/remote-ecr-reader"
        },
        {
            "Effect": "Allow",
            "Action": "sts:AssumeRole",
            "Resource": "arn:aws:iam::111213:role/remote-ecr-reader"
        }
    ]
}

Example policy attached to role "remote-ecr-reader" in each target account

{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": [
                "ecr:GetAuthorizationToken",
                "ecr:GetRegistryPolicy",
                "ecr:DescribeImageScanFindings",
                "ecr:GetLifecyclePolicyPreview",
                "ecr:GetDownloadUrlForLayer",
                "ecr:DescribeRegistry",
                "ecr:DescribeImageReplicationStatus",
                "ecr:GetAuthorizationToken",
                "ecr:ListTagsForResource",
                "ecr:ListImages",
                "ecr:BatchGetImage",
                "ecr:DescribeImages",
                "ecr:DescribeRepositories",
                "ecr:BatchCheckLayerAvailability",
                "ecr:GetRepositoryPolicy",
                "ecr:GetLifecyclePolicy"
            ],
            "Resource": "*"
        }
    ]
}

Example trust relationship attached to role "remote-ecr-reader" in each target account:
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "AWS": "arn:aws:iam::12345:role/assumerole-to-remote-ecr"
      },
      "Action": "sts:AssumeRole",
      "Condition": {}
    }
  ]
}

*names of roles aren't important but must be consistent

*/

func main() {

    awsConfig := &aws.Config{
        //Endpoint: aws.String(endpoint.URL),
        Region: aws.String("us-west-2"),
        }
    sess, err := session.NewSession(awsConfig)
    if (err != nil) {
        log.Fatal(err)
    }

    // this ARN is the role in the remote account where ECR is
    //   the role should have ECR read-only access and needs to have trust to local role
    // the IAM role in the local account that is runnign this code must
    //   be able to assumerole to the remote role
    remote_role_arn := "arn:aws:iam::044192981154:role/cporter-ecr-readonly"
    creds := stscreds.NewCredentials(sess, remote_role_arn)

    // create a session with the assumed role credentials
    service := ecr.New(sess, &aws.Config{Credentials: creds})

    // the remote repo we want to access
    repoName := "cporter-ecr"
    regId := "044192981154"
    imageSearch := &ecr.DescribeImagesInput{RegistryId: &regId, RepositoryName: &repoName}

    imagesDescription, describe_err := service.DescribeImages (imageSearch)
    if (describe_err != nil) {
        log.Fatal(describe_err)
    }
    log.Println("First Image in Repo:")
    log.Println(imagesDescription.ImageDetails[0].ImageDigest)

}
