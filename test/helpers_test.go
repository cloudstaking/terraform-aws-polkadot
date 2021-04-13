package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/ssh"
)

// checkBinaries checks if docker & docker-compose binaries exist in the instance
func checkBinaries(t *testing.T, h ssh.Host, appLayer string) (string, error) {
	commandsDocker := []string{
		"docker",
		"docker-compose",
	}

	commandsHost := []string{
		"polkadot",
		"caddy",
		"nginx",
	}

	switch appLayer {
	case "host":
		for _, cmd := range commandsHost {
			_, err := ssh.CheckSshCommandE(t, h, fmt.Sprintf("command -v %s", cmd))

			if err != nil {
				return "", fmt.Errorf("It looks like %s is not installed '%w'", cmd, err)
			}

			t.Logf("Validator has %s installed", cmd)
		}
	case "docker":
		for _, cmd := range commandsDocker {
			_, err := ssh.CheckSshCommandE(t, h, fmt.Sprintf("command -v %s", cmd))

			if err != nil {
				return "", fmt.Errorf("It looks like %s is not installed '%w'", cmd, err)
			}

			t.Logf("Validator has %s installed", cmd)
		}
	}

	return "", nil
}

// checkNodeExporter checks if node_exporter is working
func checkNodeExporter(t *testing.T, h string, username string, password string) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:9100/metrics", h), nil)
	if err != nil {
		return "", fmt.Errorf("It looks node_exporter endpoint is not working '%w'", err)
	}

	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(username+":"+password)))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("It looks	validator's node_exporter endpoint is not working (or not ready): '%w'", err)
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("validator's node_exporter endpoint returned different than 200: %v", resp.StatusCode)
	}

	t.Log("Validator is serving node_exporter's metrics")

	return "", nil
}

// checkAppFile ensures application layer files exists (docker-compose, nginx, etc)
func checkAppFiles(t *testing.T, h ssh.Host, appLayer string) (string, error) {
	switch appLayer {
	case "host":
		polkadotCustomCliArgs := fmt.Sprintf("grep CLOUDSTAKING-TEST /etc/default/polkadot")
		_, err := ssh.CheckSshCommandE(t, h, polkadotCustomCliArgs)

		if err != nil {
			return "", fmt.Errorf("It looks like %s doesn't have doesn't have the right content '%w'", polkadotCustomCliArgs, err)
		}

		t.Logf("Validator has %s installed", polkadotCustomCliArgs)
	case "docker":
		appFilesExistCmd := fmt.Sprintf("ls /home/polkadot/docker-compose.yml /home/polkadot/nginx.conf")
		_, err := ssh.CheckSshCommandE(t, h, appFilesExistCmd)

		if err != nil {
			return "", fmt.Errorf("Files /home/polkadot/docker-compose.yml and /home/polkadot/nginx.conf doesn't exist: '%w'", err)
		}

		t.Log("Validator has /home/polkadot/{docker-compose.yml,nginx.conf} files")
	}

	return "", nil
}

// checkPolkadotSnapshot check snapshot size is "big enough"
func checkPolkadotSnapshot(t *testing.T, h ssh.Host) (string, error) {
	polkashotSizeCmd := fmt.Sprintf("sudo du /home/polkadot/.local | tail -n1 | awk '{print $1}'")
	polkashotFolderSize, err := ssh.CheckSshCommandE(t, h, polkashotSizeCmd)

	if err != nil {
		return "", fmt.Errorf("Error checking size of /home/polkadot/.local (polkashot snapshot): '%w'", err)
	}

	polkashotFolderSizeInt, err := strconv.Atoi(strings.TrimSuffix(polkashotFolderSize, "\n"))
	if err != nil {
		return "", err
	}

	// >5GB means snapshot is extracking good
	if polkashotFolderSizeInt <= 5000000 {
		return "", fmt.Errorf("Snapshot folder-size (/home/polkadot/.local) < 5GB. Problem downloading snapshot? Size: '%v'", polkashotFolderSizeInt)
	}

	t.Log("Snapshot folder (/home/polkadot/.local) > 5GB. Snapshot downloaded")

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
