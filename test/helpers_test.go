package main

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/ssh"
)

// checkBinaries checks if docker & docker-compose binaries exist in the instance
func checkBinaries(t *testing.T, h ssh.Host) (string, error) {
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

// checkAppFile ensures application layer files exists (docker-compose, nginx, etc)
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

// checkDiskSize verifies the disk size is more than an input variable (sizeAtLeast). It
// is a variable because in some cloud-providers (like Scaleway) it varies a lot
func checkDiskSize(t *testing.T, h ssh.Host, sizeAtLeast int, mountPoint string) (string, error) {
	diskSizeCmd := fmt.Sprintf("df | grep %s | awk '{print $2}'", mountPoint)
	diskSizeR, err := ssh.CheckSshCommandE(t, h, diskSizeCmd)

	if err != nil {
		return "", err
	}

	diskSizeInt, err := strconv.Atoi(strings.TrimSuffix(diskSizeR, "\n"))
	if err != nil {
		return "", err
	}

	if diskSizeInt <= sizeAtLeast {
		return "", fmt.Errorf("Expected disk size is not at least %v. Got '%v'", sizeAtLeast, diskSizeInt)
	}

	t.Log("Validator seems to have the right disk size")

	return "", nil
}
