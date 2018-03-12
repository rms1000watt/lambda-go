package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/veritone/lambda-go/awshelpers"
)

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

func handleEvent(context context.Context, event events.CloudWatchEvent) (err error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	cwSvc := cloudwatch.New(sess)

	asgDetailMessage := awshelpers.ASGEventMessageDetail{}

	err = json.Unmarshal(event.Detail, &asgDetailMessage)
	if err != nil {
		log.Print(err)
		log.Printf("Could not parse the event details in the Cloudwatch Event")
		return
	}

	alarmsByInst, err := awshelpers.GetAlarmsByInstanceID(asgDetailMessage.EC2InstanceID, cwSvc)
	if err != nil {
		log.Print(err)
		log.Print("Could not retrieve alarms by instance ID")
		return
	}

	for _, alarms := range alarmsByInst {
		input := &cloudwatch.DeleteAlarmsInput{
			AlarmNames: []*string{
				aws.String(*alarms.AlarmName),
			},
		}
		log.Print("Deleted alarm: " + *alarms.AlarmName)
		cwSvc.DeleteAlarms(input)
	}

	return
}

func main() {
	lambda.Start(handleEvent)
}
