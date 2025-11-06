package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <version>\n", os.Args[0])
		os.Exit(1)
	}
	version := os.Args[1]

	// Get git log since last tag
	cmd := exec.Command("git", "log", "--oneline", "--no-merges", fmt.Sprintf("v%s..HEAD", getPrevVersion(version)))
	output, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running git log: %v\n", err)
		os.Exit(1)
	}

	changes := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(changes) == 1 && changes[0] == "" {
		changes = []string{"  * Bug fixes and improvements"}
	} else {
		for i, change := range changes {
			parts := strings.SplitN(change, " ", 2)
			if len(parts) < 2 {
				changes[i] = "  * " + change
			} else {
				changes[i] = "  * " + parts[1]
			}
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
	output, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting git tags: %v\n", err)
		return "0.0.0"
	}

	tags := strings.Split(strings.TrimSpace(string(output)), "\n")
	foundCurrent := false
	for _, tag := range tags {
		tagVersion := strings.TrimPrefix(tag, "v")
		if tagVersion == current {
			foundCurrent = true
			continue
		}
		if foundCurrent {
			return tagVersion
		}
	}
	for _, tag := range tags {
		if strings.TrimPrefix(tag, "v") != current {
			return strings.TrimPrefix(tag, "v")
		}
	}
	return "0.0.0"
}
