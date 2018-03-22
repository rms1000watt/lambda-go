package main

import (
	"context"
	"log"

	"github.com/veritone/lambda-go/awshelpers"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func handleSNSMessage(context context.Context, event events.SNSEvent) (err error) {
	if len(event.Records) == 0 {
		log.Print("SNS Record does not exist")
		return
	}

	records := event.Records
	lastRecord := records[len(records)-1].SNS.Message

	log.Print(lastRecord)

	lastRecordMessage := awshelpers.ConvertStringToAlarmMessage(lastRecord)

	instanceID, err := awshelpers.GetInstanceID(*lastRecordMessage)
	if err != nil {
		log.Print(err)
		log.Print("Could not parse the Instance ID from the Cloudwatch Alarm")
		return err
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	ec2Svc := ec2.New(sess)
	asgSvc := autoscaling.New(sess)

	asgName, err := awshelpers.GetInstanceASG(instanceID, ec2Svc)
	if err != nil {
		log.Print(err)
		log.Printf("Could not find Auto Scaling Group associated with Instance ID: %s", instanceID)
		return err
	}

	asgInstanceMinSize, err := awshelpers.GetASGMinSize(asgName, asgSvc)
	if err != nil {
		log.Print(err)
		log.Printf("Could not parse the minimum size for the Auto Scaling Group: %s", asgName)
		return err
	}

	asgInstanceCount, err := awshelpers.GetASGHealthyInstancesCount(asgName, asgSvc)
	if err != nil {
		log.Print(err)
		log.Printf("Could not find healthy instances as part of the Auto Scaling Group: %s", asgName)
		return err
	}

	if awshelpers.IsInstanceInASG(instanceID, asgName, asgSvc) {

		log.Printf("Instance (%s) is in the AutoScaling Group (%s)", instanceID, asgName)

		if asgInstanceCount > 1 && asgInstanceMinSize > 1 {

			log.Printf("Healthy instances (%d) and AutoScaling Group (%s) Minimum Instance Size (%d)", asgInstanceCount, asgName, asgInstanceMinSize)
			log.Printf("Attempting to terminate the instance (%s) in the AutoScaling Group (%s)...", instanceID, asgName)

			// Terminate Instance in ASG
			termInput := &autoscaling.TerminateInstanceInAutoScalingGroupInput{
				InstanceId:                     aws.String(instanceID),
				ShouldDecrementDesiredCapacity: aws.Bool(false),
			}
			termOutput, err := asgSvc.TerminateInstanceInAutoScalingGroup(termInput)
			if err != nil {
				log.Print(err)
				log.Print("Failed to terminate the instance (" + instanceID + ") in the AutoScaling Group: " + asgName)
				return err
			}

			if termOutput != nil {
				log.Print("Activity ID for terminating the instance (" + instanceID + "): " + *termOutput.Activity.ActivityId)
			}

			// Wait for Terminate
			waitInput := &ec2.DescribeInstancesInput{
				InstanceIds: []*string{
					aws.String(instanceID),
				},
			}

			ec2Svc.WaitUntilInstanceTerminated(waitInput)
			if err != nil {
				log.Print(err)
				log.Print("The instance (" + instanceID + ") did not terminate gracefully")
				return err
			}
		}
	}

	return
}

func main() {
	lambda.Start(handleSNSMessage)
}
