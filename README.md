# Overview 
Repository to store Lambda-Go functions for AWS

# Structure
High level directory should be the language used and the sub-directory
should be the team that is responsible for maintaining the function.

# Wercker
* Assumes that the sub-directory has a `main.go` at the root
* Assumes that the Lambda function name will coincide with the sub-directory name
* Assumes that the Lambda function and corresponding infrastructure was configured in Terraform

# Go-ification
As part of the Ryan Smith indoctrination of Golang

## TODO:
Add `doc.go` even though it will remain unrendered in our private GitHub repository