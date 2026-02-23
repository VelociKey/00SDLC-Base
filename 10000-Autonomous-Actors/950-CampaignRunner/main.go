package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: campaign-runner <campaign-name>")
		fmt.Println("Example: campaign-runner Storage-Correlation")
		os.Exit(1)
	}

	campaignName := os.Args[1]
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	root, _ := os.Getwd()
	// Navigate up to root if in submodule
	for {
		if _, err := os.Stat(filepath.Join(root, "OlympusAssurance")); err == nil {
			break
		}
		parent := filepath.Dir(root)
		if parent == root {
			break
		}
		root = parent
	}

	campaignDir := ""
	filepath.WalkDir(filepath.Join(root, "OlympusAssurance", "C0200-Execution-Campaigns"), func(path string, d os.DirEntry, err error) error {
		if d.IsDir() && strings.Contains(strings.ToLower(d.Name()), strings.ToLower(campaignName)) {
			campaignDir = path
			return filepath.SkipDir
		}
		return nil
	})

	if campaignDir == "" {
		logger.Error("Campaign not found", "name", campaignName)
		os.Exit(1)
	}

	logger.Info("🛡️  Starting Campaign", "name", campaignName, "dir", campaignDir)

	// In a real implementation, we would parse a campaign_manifest.jebnf here.
	// For now, we will use a heuristic to find the Manager bin and SaaS Mock.

	cluster := strings.Split(campaignName, "-")[0]
	managerBin := filepath.Join(root, "OlympusGCP-"+cluster, "10000-Autonomous-Actors", "900-"+cluster+"Manager", cluster+"Manager.exe")
	if _, err := os.Stat(managerBin); err != nil {
		// Fallback for different naming patterns
		managerBin = filepath.Join(root, "OlympusGCP-"+cluster, "10000-Autonomous-Actors", "900-"+cluster+"Manager", "main.go")
	}

	logger.Info("🚀 Launching Manager", "bin", managerBin)
	var managerCmd *exec.Cmd
	if strings.HasSuffix(managerBin, ".go") {
		managerCmd = exec.Command("go", "run", "main.go")
		managerCmd.Dir = filepath.Dir(managerBin)
	} else {
		managerCmd = exec.Command(managerBin)
		managerCmd.Dir = filepath.Dir(managerBin)
	}

	logPath := filepath.Join(root, "Olympus2/C0500-Agent-Intelligence-Outputs/LPSV", strings.ToLower(campaignName)+".log")
	logFile, _ := os.Create(logPath)
	defer logFile.Close()
	managerCmd.Stdout = logFile
	managerCmd.Stderr = logFile

	if err := managerCmd.Start(); err != nil {
		logger.Error("Failed to start manager", "error", err)
		os.Exit(1)
	}

	time.Sleep(5 * time.Second)

	// Run SaaS Mock / Validator
	saasBin := ""
	filepath.WalkDir(filepath.Join(root, "OlympusAssurance", "10000-Autonomous-Actors"), func(path string, d os.DirEntry, err error) error {
		if !d.IsDir() && (strings.Contains(d.Name(), cluster) && strings.Contains(d.Name(), "Mock")) {
			saasBin = path
			return filepath.SkipDir
		}
		return nil
	})

	if saasBin != "" {
		logger.Info("🧪 Executing SaaS Mock", "bin", saasBin)
		var saasCmd *exec.Cmd
		if strings.HasSuffix(saasBin, ".go") {
			saasCmd = exec.Command("go", "run", "main.go")
			saasCmd.Dir = filepath.Dir(saasBin)
		} else {
			saasCmd = exec.Command(saasBin)
			saasCmd.Dir = filepath.Dir(saasBin)
		}
		saasCmd.Stdout = os.Stdout
		saasCmd.Stderr = os.Stderr
		if err := saasCmd.Run(); err != nil {
			logger.Error("SaaS Mock failed", "error", err)
		}
	} else {
		logger.Warn("No SaaS Mock found for cluster", "cluster", cluster)
	}

	managerCmd.Process.Kill()
	logger.Info("Campaign Finished", "log", logPath)
}
