#!/bin/bash

cd lambda

aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws

sam build
sam deploy