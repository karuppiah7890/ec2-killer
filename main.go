package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/karuppiah7890/ec2-killer/pkg/config"
	"github.com/karuppiah7890/ec2-killer/pkg/slack"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconf "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

// TODO: Write tests for all of this

var version string

func main() {
	log.Printf("version: %v", version)

	c, err := config.NewConfigFromEnvVars()
	if err != nil {
		log.Fatalf("\nerror occurred while getting configuration from environment variables: %v", err)
	}

	awsconfig, err := awsconf.LoadDefaultConfig(context.TODO())
	if err != nil {
		logAndExit("error occurred while loading aws configuration", err, c)
	}

	ec2Client := ec2.NewFromConfig(awsconfig)

	// Find all instances that has been running for more than 30 minutes
	instanceIds, err := GetIdsOfInstancesToTerminate(ec2Client, c.GetVpcId(), c.GetSecurityGroupId(), c.GetEc2InstanceMaxLife())
	if err != nil {
		logAndExit("Couldn't retrieve running instances", err, c)
	}

	log.Printf("\nInstance IDs of running instances are: %v", instanceIds)

	if len(instanceIds) > 0 {
		message := fmt.Sprintf("Instance IDs of running instances are: %v", instanceIds)
		slackErr := slack.SendMessage(c.GetSlackToken(), c.GetSlackChannel(), message)
		if slackErr != nil {
			log.Printf("\nerror occurred while sending message to slack: %v", slackErr)
		}

		// Terminate all the instances that has been running for more than 30 minutes
		instanceIds, err = TerminateInstances(ec2Client, instanceIds)
		if err != nil {
			logAndExit("Couldn't terminate instances", err, c)
		}

		// Send Slack alerts for all the instances that got killed to know which instances got killed
		message = fmt.Sprintf("Instance IDs of terminated instances are: %v", instanceIds)
		slackErr = slack.SendMessage(c.GetSlackToken(), c.GetSlackChannel(), message)
		if slackErr != nil {
			log.Printf("\nerror occurred while sending message to slack: %v", slackErr)
		}
	}
}

func GetIdsOfInstancesToTerminate(client *ec2.Client, vpcId string, securityGroupId string, maxLife time.Duration) ([]string, error) {
	instanceIds := []string{}
	var nextToken *string = nil

	for {
		input := ec2.DescribeInstancesInput{
			NextToken: nextToken,
			Filters: []types.Filter{
				{
					Name: aws.String("instance-state-name"),
					Values: []string{
						"running",
						"pending",
						"stopping",
						"stopped",
					},
				},
				{
					Name: aws.String("vpc-id"),
					Values: []string{
						vpcId,
					},
				},
				{
					Name: aws.String("instance.group-id"),
					Values: []string{
						securityGroupId,
					},
				},
			},
		}
		output, err := client.DescribeInstances(context.TODO(), &input)

		if err != nil {
			return nil, err
		}

		for _, reservation := range output.Reservations {
			for _, instance := range reservation.Instances {
				instanceRunningDuration := time.Since(*instance.LaunchTime)

				if instanceRunningDuration > maxLife {
					instanceIds = append(instanceIds, *instance.InstanceId)
				}
			}
		}

		nextToken = output.NextToken
		if nextToken == nil {
			break
		}
	}

	return instanceIds, nil
}

func TerminateInstances(client *ec2.Client, instanceIds []string) ([]string, error) {
	terminatedInstanceIds := []string{}
	input := ec2.TerminateInstancesInput{
		InstanceIds: instanceIds,
	}

	output, err := client.TerminateInstances(context.TODO(), &input)
	if err != nil {
		return nil, err
	}

	for _, instanceStateChange := range output.TerminatingInstances {
		log.Printf("\nInstance Id: %s. Previous State: %s. Current State: %s", *instanceStateChange.InstanceId, instanceStateChange.PreviousState.Name, instanceStateChange.CurrentState.Name)
		terminatedInstanceIds = append(terminatedInstanceIds, *instanceStateChange.InstanceId)
	}

	return terminatedInstanceIds, nil
}

func logAndExit(errorContext string, err error, c *config.Config) {
	message := fmt.Sprintf("Critical :rotating_light:! An error occured in ec2-killer bot in %s environment :rotating_light: : ```\n%s: %v\n```", c.GetEnvironmentName(), errorContext, err)
	slackErr := slack.SendMessage(c.GetSlackToken(), c.GetSlackChannel(), message)
	if slackErr != nil {
		log.Printf("\nerror occurred while sending message to slack: %v", slackErr)
	}
	log.Fatalf("\n%s: %v", errorContext, err)
}
