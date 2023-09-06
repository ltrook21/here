package main

import (
	"fmt"
	"os"
	"os/exec"
	"encoding/json"
	"strings"
	"runtime"
)

type CallerIdentity struct {
	Arn string `json:"Arn"`
}


func getRole() {
	cmd := exec.Command("aws","sts","get-caller-identity")
	output, err := cmd.Output()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	//fmt.Println(string(output))
	
	// parse json output
	var callerIdentity CallerIdentity
	err = json.Unmarshal(output, &callerIdentity)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// split arn to extract role name:
	arnParts := strings.Split(callerIdentity.Arn, "/")
	if len(arnParts) != 3 {
		fmt.Println("Unexpected ARN format:", callerIdentity.Arn)
		return
	}

	roleName := arnParts[1]

	fmt.Println("Assumed Role:", roleName)
}

func getCtxNs() {
	cmd := exec.Command("kubectl", "config", "view", "--minify", "--output", "jsonpath=context: {..context.user} namespace: {..namespace}")
	
	// Check if the OS is macOS (Mac).
	if runtime.GOOS == "darwin" {
		// Add Homebrew's bin directory to the PATH for macOS.
		homebrewPath := "/opt/homebrew/bin" // Replace with your actual Homebrew bin directory.
		currentPath := os.Getenv("PATH")
		newPath := fmt.Sprintf("%s:%s", homebrewPath, currentPath)
		cmd.Env = append(os.Environ(), fmt.Sprintf("PATH=%s", newPath))
	}


	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(output))
}

func main() {
	getRole()
	getCtxNs()
}

