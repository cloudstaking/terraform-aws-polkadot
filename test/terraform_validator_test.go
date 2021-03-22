package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

// TestValidator deploys a validator and check the basic things like: disk space, tools
// (docker/docker-compose/etc)
func TestValidator(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("terratest-%s", random.UniqueId())

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	awsRegion := aws.GetRandomStableRegion(t, nil, nil)

	// Generate SSH keypairs
	keyPair := ssh.GenerateRSAKeyPair(t, 2048)

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/simple-validator",
		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": awsRegion,
		},
		Vars: map[string]interface{}{
			"ssh_key":       keyPair.PublicKey,
			"instance_name": instanceName,
		},
	}

	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the value of an output variable
	publicInstanceIP := terraform.Output(t, terraformOptions, "public_ip")

	publicHost := ssh.Host{
		Hostname:    publicInstanceIP,
		SshKeyPair:  keyPair,
		SshUserName: "ubuntu",
	}

	maxRetries := 30
	timeBetweenRetries := 30 * time.Second
	description := fmt.Sprintf("SSHing to validator %s to check disk size", publicInstanceIP)

	retry.DoWithRetry(t, description, maxRetries, timeBetweenRetries, func() (string, error) {
		return checkDiskSize(t, publicHost, 190000000)
	})

	description = fmt.Sprintf("SSHing to validator %s to check if docker & docker-compose are installed", publicInstanceIP)
	retry.DoWithRetry(t, description, maxRetries, timeBetweenRetries, func() (string, error) {
		return checkBinaries(t, publicHost)
	})

	description = fmt.Sprintf("SSHing in validator (%s) to check if application files exist", publicInstanceIP)
	retry.DoWithRetry(t, description, maxRetries, timeBetweenRetries, func() (string, error) {
		return checkAppFiles(t, publicHost)
	})
}

// TestValidatorWithPolkashots deploys a validator and enable polkashots
func TestValidatorWithPolkashots(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("terratest-%s", random.UniqueId())

	awsRegion := aws.GetRandomStableRegion(t, nil, nil)
	keyPair := ssh.GenerateRSAKeyPair(t, 2048)

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/polkashots",
		EnvVars: map[string]string{
			"AWS_DEFAULT_REGION": awsRegion,
		},
		Vars: map[string]interface{}{
			"ssh_key":       keyPair.PublicKey,
			"instance_name": instanceName,
		},
	}

	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	publicInstanceIP := terraform.Output(t, terraformOptions, "public_ip")

	publicHost := ssh.Host{
		Hostname:    publicInstanceIP,
		SshKeyPair:  keyPair,
		SshUserName: "ubuntu",
	}

	maxRetries := 30
	timeBetweenRetries := 30 * time.Second

	description := fmt.Sprintf("SSHing to validator %s to check if docker & docker-compose are installed", publicInstanceIP)
	retry.DoWithRetry(t, description, maxRetries, timeBetweenRetries, func() (string, error) {
		return checkBinaries(t, publicHost)
	})

	description = fmt.Sprintf("SSHing to validator (%s) to check if snapshot folder exist and >5GB", publicInstanceIP)
	retry.DoWithRetry(t, description, maxRetries, timeBetweenRetries, func() (string, error) {
		return checkPolkadotSnapshot(t, publicHost)
	})

	description = fmt.Sprintf("SSHing to validator (%s) to check if application files exist", publicInstanceIP)
	retry.DoWithRetry(t, description, maxRetries, timeBetweenRetries, func() (string, error) {
		return checkAppFiles(t, publicHost)
	})
}
