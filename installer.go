// runs the installer for the given project task.
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"kau/services"
)

func runInstallations(tasks []services.ProjectTask) error {
	checkedRuntimes := make(map[string]bool)

	for _, task := range tasks {
		rt := task.Runtime
		if rt.CheckCmd != "" && !checkedRuntimes[rt.CheckCmd] {
			checkedRuntimes[rt.CheckCmd] = true
			if err := ensureRuntime(rt); err != nil {
				return fmt.Errorf("runtime error for %s: %w", task.Ecosystem, err)
			}
		}

		fmt.Printf("Installing dependencies for %s in %s...\n", task.Ecosystem, task.RootDir)
		for _, cmd := range task.InstallSteps {
			if err := runCommand(cmd, task.RootDir); err != nil {
				return fmt.Errorf("error installing dependencies for %s: %w", task.Ecosystem, err)
			}
		}
	}
	return nil
}

func ensureRuntime(rt services.RuntimeSpec) error {
	if _, err := exec.LookPath(rt.CheckCmd); err == nil {
		check := exec.Command(rt.CheckCmd, rt.CheckArgs...)
		if check.Run() == nil {
			return nil
		}
	}

	fmt.Printf("%s (%s) is not installed.\n", rt.Name, rt.CheckCmd)
	if rt.WinWingetID == "" {
		return fmt.Errorf("missing runtime %s", rt.Name)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Install %s via winget? (y/n): ", rt.Name)
	input, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	ans := strings.TrimSpace(strings.ToLower(input))
	if ans != "y" && ans != "yes" {
		return fmt.Errorf("declined installation of %s", rt.Name)
	}

	return runCommand([]string{"winget", "install", "-e", "--id", rt.WinWingetID}, "")
}

func runCommand(cmdArgs []string, dir string) error {
	if len(cmdArgs) == 0 {
		return nil
	}

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
