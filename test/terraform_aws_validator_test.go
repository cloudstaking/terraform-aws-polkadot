package main

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/ssh"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

// TestValidatorWithAddicionalVolume deploys a validator with an addional volume
func TestValidatorWithAddicionalVolume(t *testing.T) {
	t.Parallel()

	instanceName := fmt.Sprintf("terratest-%s", random.UniqueId())

	// Pick a random AWS region to test in. This helps ensure your code works in all regions.
	awsRegion := aws.GetRandomStableRegion(t, nil, nil)

	// Generate SSH keypairs
	keyPair := ssh.GenerateRSAKeyPair(t, 2048)

	terraformOptions := &terraform.Options{
		TerraformDir: "../examples/addicional-volume",
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
	description := fmt.Sprintf("SSHing to validator %s to check /srv size", publicInstanceIP)

	// Testing partition was created with the right size
	retry.DoWithRetry(t, description, maxRetries, timeBetweenRetries, func() (string, error) {
		return checkVolumeSize(t, publicHost)
	})

	description = fmt.Sprintf("SSHing to validator %s to check if docker & docker-compose are installed", publicInstanceIP)
	retry.DoWithRetry(t, description, maxRetries, timeBetweenRetries, func() (string, error) {
		return checkDockerBinaries(t, publicHost)
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
		return checkDockerBinaries(t, publicHost)
	})

	description = fmt.Sprintf("SSHing into validator (%s) to check if snapshot folder exist and >5GB", publicInstanceIP)
	retry.DoWithRetry(t, description, maxRetries, timeBetweenRetries, func() (string, error) {
		return checkPolkadotSnapshot(t, publicHost)
	})

	description = fmt.Sprintf("SSHing into validator (%s) to check if application files exist", publicInstanceIP)
	retry.DoWithRetry(t, description, maxRetries, timeBetweenRetries, func() (string, error) {
		return checkAppFiles(t, publicHost)
	})
}

/////////////
// Helpers //
/////////////

// checkDockerBinaries Function to check if docker and docker-compose binaries exist in the host
func checkDockerBinaries(t *testing.T, h ssh.Host) (string, error) {
	dockerExistCmd := fmt.Sprintf("command -v docker")
	_, err := ssh.CheckSshCommandE(t, h, dockerExistCmd)

	if err != nil {
		return "", fmt.Errorf("It looks like docker is not installed '%w'", err)
	}

	dockerComposeExistCmd := fmt.Sprintf("command -v docker-compose")
	_, err = ssh.CheckSshCommandE(t, h, dockerComposeExistCmd)

	if err != nil {
		return "", fmt.Errorf("Docker is not yet installed '%w'", err)
	}

	t.Log("Validator got docker & docker-compose installed")

	return "", nil
}

// checkAppFile ensure application layer files exists
func checkAppFiles(t *testing.T, h ssh.Host) (string, error) {
	appFilesExistCmd := fmt.Sprintf("ls /srv/docker-compose.yml /srv/nginx.conf")
	_, err := ssh.CheckSshCommandE(t, h, appFilesExistCmd)

	if err != nil {
		return "", fmt.Errorf("Files /srv/docker-compose.yml and /srv/nginx.conf doesn't exist: '%w'", err)
	}

	t.Log("Validator has /srv/{docker-compose.yml,nginx.conf} files")

	return "", nil
}

// checkPolkadotSnapshot check snapshot size is "big enough"
func checkPolkadotSnapshot(t *testing.T, h ssh.Host) (string, error) {
	polkashotSizeCmd := fmt.Sprintf("sudo du /srv/kusama/ | tail -n1 | awk '{print $1}'")
	polkashotFolderSize, err := ssh.CheckSshCommandE(t, h, polkashotSizeCmd)

	if err != nil {
		return "", fmt.Errorf("Error checking size of /srv/kusama/ directory: '%w'", err)
	}

	polkashotFolderSizeInt, err := strconv.Atoi(strings.TrimSuffix(polkashotFolderSize, "\n"))
	if err != nil {
		return "", err
	}

	// >5GB means snapshot is extracking good
	if polkashotFolderSizeInt <= 5000000 {
		return "", fmt.Errorf("Snapshot folder-size (/srv/kusama/) < 5GB. Problem downloading snapshot? Size: '%v'", polkashotFolderSizeInt)
	}

	t.Log("Snapshot folder (/srv/kusama) > 5GB. Snapshot downloaded")

	return "", nil
}

// checkVolumeSize verify the size of the volume attached
func checkVolumeSize(t *testing.T, h ssh.Host) (string, error) {
	diskSizeCmd := fmt.Sprintf("df | grep /srv | awk '{print $2}'")
	diskSizeR, err := ssh.CheckSshCommandE(t, h, diskSizeCmd)

	if err != nil {
		return "", err
	}

	diskSizeInt, err := strconv.Atoi(strings.TrimSuffix(diskSizeR, "\n"))
	if err != nil {
		return "", err
	}

	if diskSizeInt <= 190000000 && diskSizeInt >= 210000000 {
		return "", fmt.Errorf("Expected disk size to be within 190000000 and 210000000 but got '%v'", diskSizeInt)
	}

	t.Log("Validator seems to have the right disk size in /srv")

	return "", nil
}
