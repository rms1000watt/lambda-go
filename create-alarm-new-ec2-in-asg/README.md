# Overview
This Lambda function will create an alarm associated with an EC2 instance if
a Cloudwatch Event marks the instance as `running` and an alarm with the 
instance ID does not already exist