package awshelpers

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

func IsInstanceInASG(instanceID string, asgName string, svc *autoscaling.AutoScaling) (response bool) {
	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{
			aws.String(asgName),
		},
	}

	result, err := svc.DescribeAutoScalingGroups(input)
	if err != nil {
		return false
	}

	for _, Instance := range result.AutoScalingGroups[0].Instances {
		if *Instance.InstanceId == instanceID {
			return true
		}
	}

	return
}

func GetASGMinSize(asgName string, svc *autoscaling.AutoScaling) (instanceCount int64, err error) {
	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{
			aws.String(asgName),
		},
	}

	result, err := svc.DescribeAutoScalingGroups(input)

	if result != nil && len(result.AutoScalingGroups) > 0 && result.AutoScalingGroups[0].MinSize != nil {
		return *result.AutoScalingGroups[0].MinSize, err
	}

	return
}

func GetASGHealthyInstancesCount(asgName string, svc *autoscaling.AutoScaling) (instanceCount int64, err error) {
	input := &autoscaling.DescribeAutoScalingGroupsInput{
		AutoScalingGroupNames: []*string{
			aws.String(asgName),
		},
	}

	result, err := svc.DescribeAutoScalingGroups(input)

	if result != nil && len(result.AutoScalingGroups) > 0 && result.AutoScalingGroups[0].Instances != nil {
		for _, Instance := range result.AutoScalingGroups[0].Instances {
			if Instance.HealthStatus != nil {
				if *Instance.HealthStatus == "Healthy" {
					instanceCount++
				}
			}
		}
	}

	return
}
