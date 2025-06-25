package MGitBridge

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/go-git/go-git/v5"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

// NSLog provides iOS-style logging that's visible in Xcode Console
func NSLog(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	log.Printf("MGitModule: %s", message)
}

// Help returns MGit help information for iOS
// This function provides the same help text as the main MGit project
func Help() string {
	// Test comprehensive logging during help execution
	log.Printf("MGitBridge: Help() function called")
	fmt.Printf("MGitBridge: Help() - testing fmt.Printf\n")
	
	// Test different log levels and methods
	log.Printf("MGitBridge: Generating help text")
	fmt.Println("MGitBridge: About to generate help text")
	
	// Test accessing environment
	pwd, _ := os.Getwd()
	log.Printf("MGitBridge: Current working directory: %s", pwd)
	
	// Test system information
	log.Printf("MGitBridge: Runtime info - NumCPU: %d, NumGoroutine: %d", 
		runtime.NumCPU(), runtime.NumGoroutine())
	
	// Test various data types in logging
	log.Printf("MGitBridge: Integer: %d", 42)
	log.Printf("MGitBridge: Float: %.2f", 3.14159)
	log.Printf("MGitBridge: Boolean: %t", true)
	log.Printf("MGitBridge: String: %s", "test-string")
	
	// Get the help text (same as MGit's printUsage)
	helpText := getMGitHelpText()
	
	log.Printf("MGitBridge: Help text generated successfully")
	log.Printf("MGitBridge: Help text length: %d characters", len(helpText))
	
	// Test final logging
	fmt.Printf("MGitBridge: Help() returning result\n")
	log.Printf("MGitBridge: Help() function complete")
	
	return helpText
}

// getMGitHelpText returns the same help text as MGit's printUsage() function
func getMGitHelpText() string {
	var buf bytes.Buffer
	
	buf.WriteString("mgit - A go-git wrapper\n")
	buf.WriteString("Usage: mgit <command> [args]\n")
	buf.WriteString("Commands:\n")
	buf.WriteString("  init                        Initialize a new repository\n")
	buf.WriteString("  clone [-jwt <token>] <url>  Clone a repository\n")
	buf.WriteString("  add <files...>              Add files to staging\n")
	buf.WriteString("  commit -m <msg>             Commit staged changes\n")
	buf.WriteString("  push                        Push commits to remote\n")
	buf.WriteString("  pull                        Pull changes from remote\n")
	buf.WriteString("  status                      Show repository status\n")
	buf.WriteString("  branch                      List branches\n")
	buf.WriteString("  branch <n>               Create a new branch\n")
	buf.WriteString("  checkout <ref>              Checkout a branch or commit\n")
	buf.WriteString("  log                         Show commit history\n")
	buf.WriteString("  show [commit]               Show commit details and changes\n")
	buf.WriteString("  config                      Get and set configuration values\n")
	buf.WriteString("  verify                      Verify MGit commit chain integrity\n")
	
	return buf.String()
}

// TestLogging performs comprehensive logging tests for iOS debugging
func TestLogging() string {
	// Test every conceivable logging method
	log.Printf("=== MGitBridge: Starting comprehensive logging test ===")
	
	// Standard library logging
	fmt.Print("MGitBridge: fmt.Print test\n")
	fmt.Println("MGitBridge: fmt.Println test")
	fmt.Printf("MGitBridge: fmt.Printf test with arg: %s\n", "test-value")
	
	// Log package with different methods
	log.Print("MGitBridge: log.Print test")
	log.Println("MGitBridge: log.Println test")
	log.Printf("MGitBridge: log.Printf test with arg: %s", "test-value")
	
	// Test writing to different outputs
	fmt.Fprint(os.Stdout, "MGitBridge: fmt.Fprint to stdout\n")
	fmt.Fprint(os.Stderr, "MGitBridge: fmt.Fprint to stderr\n")
	fmt.Fprintf(os.Stdout, "MGitBridge: fmt.Fprintf to stdout with arg: %s\n", "test-value")
	fmt.Fprintf(os.Stderr, "MGitBridge: fmt.Fprintf to stderr with arg: %s\n", "test-value")
	
	// Test different log levels (simulated)
	log.Printf("MGitBridge: [DEBUG] This is a debug message")
	log.Printf("MGitBridge: [INFO] This is an info message")
	log.Printf("MGitBridge: [WARN] This is a warning message")
	log.Printf("MGitBridge: [ERROR] This is an error message")
	
	// Test runtime information
	log.Printf("MGitBridge: Runtime - GOOS: %s, GOARCH: %s", runtime.GOOS, runtime.GOARCH)
	log.Printf("MGitBridge: Runtime - Version: %s", runtime.Version())
	log.Printf("MGitBridge: Runtime - NumCPU: %d", runtime.NumCPU())
	log.Printf("MGitBridge: Runtime - NumGoroutine: %d", runtime.NumGoroutine())
	
	// Test environment access
	if pwd, err := os.Getwd(); err == nil {
		log.Printf("MGitBridge: Current directory: %s", pwd)
	}
	
	result := "Comprehensive logging test completed. Check console/logs for output."
	log.Printf("MGitBridge: Test result: %s", result)
	log.Printf("=== MGitBridge: Logging test complete ===")
	
	return result
}

// SimpleAdd performs a simple addition (for basic functionality testing)
func SimpleAdd(a, b int) int {
	log.Printf("MGitBridge: SimpleAdd called with a=%d, b=%d", a, b)
	result := a + b
	log.Printf("MGitBridge: SimpleAdd result: %d", result)
	return result
}

// Add stages files for commit using go-git (exactly mimics CLI mgit add)
func Add(repoPath string, filePaths string) *AddResult {
	log.Printf("MGitBridge: Add() called for repo: %s", repoPath)
	log.Printf("MGitBridge: Files to add: %s", filePaths)
	
	result := &AddResult{
		Success: false,
		Message: "",
		Error:   "",
	}
	
	// Open the repository using go-git
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Printf("MGitBridge: Error opening repository: %s", err)
		result.Error = fmt.Sprintf("Error opening repository: %s", err)
		return result
	}
	
	// Get the worktree
	w, err := repo.Worktree()
	if err != nil {
		log.Printf("MGitBridge: Error getting worktree: %s", err)
		result.Error = fmt.Sprintf("Error getting worktree: %s", err)
		return result
	}
	
	log.Printf("MGitBridge: Adding file to staging: %s", filePaths)
	
	// Add the file to staging (exactly like CLI does with w.Add(file))
	_, err = w.Add(filePaths)
	if err != nil {
		log.Printf("MGitBridge: Error adding file: %s", err)
		result.Error = fmt.Sprintf("Error adding file %s: %s", filePaths, err)
		return result
	}
	
	log.Printf("MGitBridge: File successfully staged: %s", filePaths)
	
	result.Success = true
	result.Message = "Changes staged for commit"
	
	return result
}

// Push pushes changes to the remote repository
func Push(repoPath, token string) *PushResult {
	result := &PushResult{Success: false, Message: "", CommitHash: ""}

	// Open the repository
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		result.Message = fmt.Sprintf("Failed to open repository: %v", err)
		return result
	}

	// Get the current HEAD commit hash before pushing
	ref, err := repo.Head()
	if err != nil {
		result.Message = fmt.Sprintf("Failed to get HEAD reference: %v", err)
		return result
	}
	result.CommitHash = ref.Hash().String()

	// Get the remote origin
	remote, err := repo.Remote("origin")
	if err != nil {
		result.Message = fmt.Sprintf("No remote 'origin' found: %v", err)
		return result
	}

	// Get the remote URL
	config := remote.Config()
	if len(config.URLs) == 0 {
		result.Message = "No remote URL configured"
		return result
	}

	// Set up authentication using githttp.BasicAuth (consistent with clone)
	auth := &githttp.BasicAuth{
		Username: "", // Empty username works with MGit server
		Password: token, // JWT token from QR scan
	}

	// Push to the remote (main branch)
	err = repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       auth,
	})

	if err != nil {
		// Handle common push errors
		if err == git.NoErrAlreadyUpToDate {
			result.Success = true
			result.Message = "Already up to date"
			return result
		}
		result.Message = fmt.Sprintf("Push failed: %v", err)
		return result
	}

	result.Success = true
	result.Message = "Push completed successfully"
	return result
}