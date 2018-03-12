# Overview
Lambda function to detach and terminate an EC2 instance if the disk is more than 
95% full and the EC2 is associated with an active AutoScaling Group and that 
AutoScaling Group's minimum size is larger than one and there are currently more 
than one active healthy instances