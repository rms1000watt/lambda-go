package awshelpers

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

// Iterate through responses and tokens and store to memory
func GetAllAlarms(svc *cloudwatch.CloudWatch) (alarms []*cloudwatch.MetricAlarm, err error) {
	const maxRecordsLimit = 100 // hardcoded limit by AWS 03/12/2018
	nextToken := ""
	nextInput := new(cloudwatch.DescribeAlarmsInput)

	for nextToken != "" || len(alarms) == 0 {
		alarmsReturn, _ := svc.DescribeAlarms(nextInput) // allowed to pass nil to DescribeAlarms

		if alarmsReturn != nil {
			if len(alarmsReturn.MetricAlarms) > 0 {
				alarms = append(alarms, alarmsReturn.MetricAlarms...)
			}
			if alarmsReturn.NextToken != nil {
				nextToken = *alarmsReturn.NextToken
			} else {
				nextToken = ""
			}
		}

		nextInput = &cloudwatch.DescribeAlarmsInput{
			NextToken:  aws.String(nextToken),
			MaxRecords: aws.Int64(maxRecordsLimit),
		}
	}

	return
}

func GetAlarmsByInstanceID(instanceID string, svc *cloudwatch.CloudWatch) (alarms []*cloudwatch.MetricAlarm, err error) {
	allAlarms, _ := GetAllAlarms(svc)

	for _, alarm := range allAlarms {
		if len(alarm.Dimensions) > 0 {
			for _, dimension := range alarm.Dimensions {
				if *dimension.Name == "InstanceId" && *dimension.Value == instanceID {
					alarms = append(alarms, alarm)
				}
			}
		}
	}

	return
}
