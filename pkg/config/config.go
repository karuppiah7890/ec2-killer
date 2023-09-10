package config

import (
	"fmt"
	"os"
	"time"
)

// All configuration is through environment variables

const SECURITY_GROUP_ID_ENV_VAR = "SECURITY_GROUP_ID"
const VPC_ID_ENV_VAR = "VPC_ID"
const EC2_INSTANCE_MAX_LIFE_ENV_VAR = "EC2_INSTANCE_MAX_LIFE"
const DEFAULT_EC2_INSTANCE_MAX_LIFE = "30m"
const AWS_REGION_ENV_VAR = "AWS_REGION"
const AWS_ACCESS_KEY_ID_ENV_VAR = "AWS_ACCESS_KEY_ID"
const AWS_SECRET_ACCESS_KEY_ENV_VAR = "AWS_SECRET_ACCESS_KEY"
const ENVIRONMENT_NAME_ENV_VAR = "ENVIRONMENT_NAME"
const DEFAULT_ENVIRONMENT_NAME = "Production"
const SLACK_TOKEN_ENV_VAR = "SLACK_TOKEN"
const SLACK_CHANNEL_ENV_VAR = "SLACK_CHANNEL"

type Config struct {
	awsRegion          string
	awsAccessKeyId     string
	awsSecretAccessKey string
	environmentName    string
	slackToken         string
	slackChannel       string
	ec2InstanceMaxLife time.Duration
	vpcId              string
	securityGroupId    string
}

func NewConfigFromEnvVars() (*Config, error) {
	awsRegion, err := getAwsRegion()
	if err != nil {
		return nil, fmt.Errorf("error occurred while getting AWS Region: %v", err)
	}

	awsAccessKeyId, err := getAwsAccessKeyId()
	if err != nil {
		return nil, fmt.Errorf("error occurred while getting AWS Access Key ID: %v", err)
	}

	awsSecretAccessKey, err := getAwsSecretAccessKey()
	if err != nil {
		return nil, fmt.Errorf("error occurred while getting AWS Secret Access Key: %v", err)
	}

	environmentName := getEnvironmentName()

	slackToken, err := getSlackToken()
	if err != nil {
		return nil, fmt.Errorf("error occurred while getting slack token: %v", err)
	}

	slackChannel, err := getSlackChannel()
	if err != nil {
		return nil, fmt.Errorf("error occurred while getting slack channel: %v", err)
	}

	vpcId, err := getVpcId()
	if err != nil {
		return nil, fmt.Errorf("error occurred while getting vpc id: %v", err)
	}

	securityGroupId, err := getSecurityGroupId()
	if err != nil {
		return nil, fmt.Errorf("error occurred while getting security group id: %v", err)
	}

	ec2InstanceMaxLife, err := getEc2InstanceMaxLife()
	if err != nil {
		return nil, fmt.Errorf("error occurred while getting ec2 instance max life: %v", err)
	}

	return &Config{
		awsRegion:          awsRegion,
		awsAccessKeyId:     awsAccessKeyId,
		awsSecretAccessKey: awsSecretAccessKey,
		environmentName:    environmentName,
		slackToken:         slackToken,
		slackChannel:       slackChannel,
		ec2InstanceMaxLife: ec2InstanceMaxLife,
		vpcId:              vpcId,
		securityGroupId:    securityGroupId,
	}, nil
}

// Get AWS VPC ID
func getVpcId() (string, error) {
	vpcId, ok := os.LookupEnv(VPC_ID_ENV_VAR)
	if !ok {
		return "", fmt.Errorf("%s environment variable value is a required value. Please define it", VPC_ID_ENV_VAR)
	}

	return vpcId, nil
}

func getSecurityGroupId() (string, error) {
	securityGroupId, ok := os.LookupEnv(SECURITY_GROUP_ID_ENV_VAR)
	if !ok {
		return "", fmt.Errorf("%s environment variable value is a required value. Please define it", SECURITY_GROUP_ID_ENV_VAR)
	}

	return securityGroupId, nil
}

func getEc2InstanceMaxLife() (time.Duration, error) {
	maxLifeStr, ok := os.LookupEnv(EC2_INSTANCE_MAX_LIFE_ENV_VAR)
	if !ok {
		maxLifeStr = DEFAULT_EC2_INSTANCE_MAX_LIFE
	}

	maxLife, err := time.ParseDuration(maxLifeStr)
	if err != nil {
		return 0, fmt.Errorf("error occurred while parsing ec2 instance max life value %s: %v", maxLifeStr, err)
	}

	return maxLife, nil
}

// Get AWS Region
func getAwsRegion() (string, error) {
	awsRegion, ok := os.LookupEnv(AWS_REGION_ENV_VAR)
	if !ok {
		return "", fmt.Errorf("%s environment variable value is a required value. Please define it", AWS_REGION_ENV_VAR)
	}

	return awsRegion, nil
}

// Get AWS Access Key ID
func getAwsAccessKeyId() (string, error) {
	awsAccessKeyId, ok := os.LookupEnv(AWS_ACCESS_KEY_ID_ENV_VAR)
	if !ok {
		return "", fmt.Errorf("%s environment variable value is a required value. Please define it", AWS_ACCESS_KEY_ID_ENV_VAR)
	}

	return awsAccessKeyId, nil
}

// Get AWS Secret Access Key
func getAwsSecretAccessKey() (string, error) {
	awsSecretAccessKey, ok := os.LookupEnv(AWS_SECRET_ACCESS_KEY_ENV_VAR)
	if !ok {
		return "", fmt.Errorf("%s environment variable value is a required value. Please define it", AWS_SECRET_ACCESS_KEY_ENV_VAR)
	}

	return awsSecretAccessKey, nil
}

// Get optional environment name for the environment where
// the services are running. Default is "Production". This name will
// be used in the alert messages
func getEnvironmentName() string {
	environmentName, ok := os.LookupEnv(ENVIRONMENT_NAME_ENV_VAR)
	if !ok {
		return DEFAULT_ENVIRONMENT_NAME
	}

	return environmentName
}

func getSlackToken() (string, error) {
	slackToken, ok := os.LookupEnv(SLACK_TOKEN_ENV_VAR)
	if !ok {
		return "", fmt.Errorf("%s environment variable is not defined and is required. Please define it", SLACK_TOKEN_ENV_VAR)
	}
	return slackToken, nil
}

func getSlackChannel() (string, error) {
	slackChannel, ok := os.LookupEnv(SLACK_CHANNEL_ENV_VAR)
	if !ok {
		return "", fmt.Errorf("%s environment variable is not defined and is required. Please define it", SLACK_CHANNEL_ENV_VAR)
	}
	return slackChannel, nil
}

func (c *Config) GetEnvironmentName() string {
	return c.environmentName
}

func (c *Config) GetSlackToken() string {
	return c.slackToken
}

func (c *Config) GetSlackChannel() string {
	return c.slackChannel
}

func (c *Config) GetEc2InstanceMaxLife() time.Duration {
	return c.ec2InstanceMaxLife
}

func (c *Config) GetVpcId() string {
	return c.vpcId
}

func (c *Config) GetSecurityGroupId() string {
	return c.securityGroupId
}
