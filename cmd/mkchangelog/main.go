package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	version := os.Args[1]
	
	// Get git log since last tag
	cmd := exec.Command("git", "log", "--oneline", "--no-merges", fmt.Sprintf("v%s..HEAD", getPrevVersion(version)))
	output, _ := cmd.Output()
	
	changes := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(changes) == 1 && changes[0] == "" {
		changes = []string{"  * Bug fixes and improvements"}
	} else {
		for i, change := range changes {
			changes[i] = "  * " + strings.SplitN(change, " ", 2)[1]
		}
	}
	
	fmt.Printf("tailscale (%s-1) unstable; urgency=medium\n\n", version)
	for _, change := range changes {
		fmt.Println(change)
	}
	fmt.Printf("\n -- Tailscale Inc <info@tailscale.com>  %s\n", time.Now().Format("Mon, 02 Jan 2006 15:04:05 -0700"))
}

func getPrevVersion(current string) string {
	cmd := exec.Command("git", "tag", "-l", "--sort=-version:refname")
	output, _ := cmd.Output()
	tags := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, tag := range tags {
		if strings.TrimPrefix(tag, "v") != current {
			return strings.TrimPrefix(tag, "v")
		}
	}
	return "0.0.0"
}