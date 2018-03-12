package awshelpers

import (
	"time"
)

type ASGEventMessage struct {
	Version    string                `json:"version"`
	ID         string                `json:"id"`
	DetailType string                `json:"detail-type"`
	Source     string                `json:"source"`
	Account    string                `json:"account"`
	Time       time.Time             `json:"time"`
	Region     string                `json:"region"`
	Resources  []string              `json:"resources"`
	Detail     ASGEventMessageDetail `json:"detail"`
}

type ASGEventMessageDetail struct {
	StatusCode           string                 `json:"StatusCode"`
	AutoScalingGroupName string                 `json:"AutoScalingGroupName"`
	ActivityID           string                 `json:"AcitivityId"`
	Details              ASGEventMessageDetails `json:"Details"`
	RequestID            string                 `json:"RequestId"`
	EndTime              time.Time              `json:"EndTime"`
	EC2InstanceID        string                 `json:"EC2InstanceId"`
	StartTime            time.Time              `json:"StartTime"`
	Cause                string                 `json:"Cause"`
}

type ASGEventMessageDetails struct {
	AvailabilityZone string `json:"Availability Zone"`
	SubnetID         string `json:"Subnet ID"`
}

/* Sample Termination Body
Is there an AWS predefined struct for this?
{
  "version": "0",
  "id": "156d01c9-a6c3-4d7e-b883-5758266b95af",
  "detail-type": "EC2 Instance Terminate Successful",
  "source": "aws.autoscaling",
  "account": "123456789012",
  "time": "2015-11-11T21:36:57Z",
  "region": "us-east-1",
  "resources": [
    "arn:aws:autoscaling:us-east-1:123456789012:autoScalingGroup:eb56d16b-bbf0-401d-b893-d5978ed4a025:autoScalingGroupName/sampleTermASG",
    "arn:aws:ec2:us-east-1:123456789012:instance/i-b188560f"
  ],
  "detail": {
    "StatusCode": "InProgress",
    "AutoScalingGroupName": "sampleTermASG",
    "ActivityId": "56472e79-538a-4ba7-b3cc-768d889194b0",
    "Details": {
      "Availability Zone": "us-east-1b",
      "Subnet ID": "subnet-95bfcebe"
    },
    "RequestId": "56472e79-538a-4ba7-b3cc-768d889194b0",
    "EndTime": "2015-11-11T21:36:57.498Z",
    "EC2InstanceId": "i-b188560f",
    "StartTime": "2015-11-11T21:36:12.649Z",
    "Cause": "At 2015-11-11T21:36:03Z a user request update of AutoScalingGroup constraints to min: 0, max: 1, desired: 0 changing the desired capacity from 1 to 0.  At 2015-11-11T21:36:12Z an instance was taken out of service in response to a difference between desired and actual capacity, shrinking the capacity from 1 to 0.  At 2015-11-11T21:36:12Z instance i-b188560f was selected for termination."
  }
} */
