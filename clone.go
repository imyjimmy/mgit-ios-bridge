package MGitBridge

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
)

// cloneRepository implements the MGit clone functionality using go-git
func cloneRepository(url, destination, token string) error {
	// Create the destination directory if it doesn't exist
	if err := os.MkdirAll(destination, 0755); err != nil {
		return fmt.Errorf("error creating destination directory: %w", err)
	}

	// Fetch repository metadata first
	NSLog("Fetching repository metadata...")
	repoInfo, err := fetchRepositoryInfo(url, token)
	if err != nil {
		return fmt.Errorf("error fetching repository metadata: %w", err)
	}

	NSLog("Repository: %s, Access level: %s", repoInfo.Name, repoInfo.Access)

	// Clone the Git data using go-git instead of system git
	NSLog("Cloning Git repository with go-git...")
	if err := gitCloneWithGoGit(url, destination, token); err != nil {
		return fmt.Errorf("error cloning Git repository: %w", err)
	}

	// Fetch and set up MGit metadata
	NSLog("Setting up MGit metadata...")
	if err := fetchMGitMetadata(url, destination, token); err != nil {
		NSLog("Warning: Failed to fetch MGit metadata: %s", err.Error())
	}

	// Set up MGit configuration
	if err := setupMGitConfig(destination, repoInfo); err != nil {
		return fmt.Errorf("error setting up MGit config: %w", err)
	}

	NSLog("Clone completed successfully")
	return nil
}

// gitCloneWithGoGit performs the Git clone using go-git library (iOS compatible)
func gitCloneWithGoGit(url, destination, token string) error {
	repoID := extractRepoID(url)
	serverBaseURL := extractServerBaseURL(url)
	
	// Construct the Git URL for the repository
	gitURL := fmt.Sprintf("%s/api/mgit/repos/%s", serverBaseURL, repoID)
	
	NSLog("Cloning from: %s", gitURL)
	NSLog("Destination: %s", destination)
	
	// Set up authentication using Bearer token
	auth := &githttp.BasicAuth{
		Username: "token", // Username can be anything when using token auth
		Password: token,
	}
	
	// Clone options
	cloneOptions := &git.CloneOptions{
		URL:               gitURL,
		Auth:              auth,
		RemoteName:        "origin",
		ReferenceName:     "", // Clone default branch
		SingleBranch:      false, // Clone all branches
		NoCheckout:        false, // Do checkout working directory
		Depth:             0, // Full clone, not shallow
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	}
	
	// Perform the clone
	repo, err := git.PlainClone(destination, false, cloneOptions)
	if err != nil {
		NSLog("Clone failed with error: %s", err.Error())
		return fmt.Errorf("error cloning repository with go-git: %w", err)
	}
	
	// Verify the clone was successful
	workTree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("error getting worktree: %w", err)
	}
	
	// Get HEAD to verify we have commits
	head, err := repo.Head()
	if err != nil {
		NSLog("Warning: Could not get HEAD reference: %s", err.Error())
	} else {
		NSLog("Successfully cloned, HEAD at: %s", head.Hash().String()[:7])
	}
	
	// Log worktree status
	status, err := workTree.Status()
	if err != nil {
		NSLog("Warning: Could not get worktree status: %s", err.Error())
	} else {
		NSLog("Worktree status: %d files", len(status))
	}
	
	NSLog("Git clone completed successfully")
	return nil
}

// fetchRepositoryInfo fetches information about the repository
func fetchRepositoryInfo(url, token string) (*RepositoryInfo, error) {
	repoID := extractRepoID(url)
	serverBaseURL := extractServerBaseURL(url)
	
	infoURL := fmt.Sprintf("%s/api/mgit/repos/%s/info", serverBaseURL, repoID)
	
	req, err := http.NewRequest("GET", infoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error response from server (status %d): %s", resp.StatusCode, string(bodyBytes))
	}
	
	var repoInfo RepositoryInfo
	if err := json.NewDecoder(resp.Body).Decode(&repoInfo); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}
	
	return &repoInfo, nil
}

// extractRepoID extracts the repository ID from a URL
func extractRepoID(url string) string {
	url = strings.TrimSuffix(strings.TrimSuffix(url, "/"), ".git")
	parts := strings.Split(url, "/")
	return parts[len(parts)-1]
}

// extractServerBaseURL extracts the server base URL from a repository URL
func extractServerBaseURL(url string) string {
	repoID := extractRepoID(url)
	baseURL := strings.TrimSuffix(url, "/"+repoID)
	baseURL = strings.TrimSuffix(baseURL, repoID)
	return baseURL
}

// fetchMGitMetadata fetches the MGit metadata and sets it up in the repository
func fetchMGitMetadata(url, destination, token string) error {
	repoID := extractRepoID(url)
	serverBaseURL := extractServerBaseURL(url)
	
	metadataURL := fmt.Sprintf("%s/api/mgit/repos/%s/metadata", serverBaseURL, repoID)
	
	req, err := http.NewRequest("GET", metadataURL, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error response from server (status %d): %s", resp.StatusCode, string(bodyBytes))
	}
	
	var mappings []interface{}
	if err := json.NewDecoder(resp.Body).Decode(&mappings); err != nil {
		return fmt.Errorf("error parsing metadata response: %w", err)
	}
	
	// Create the .mgit directory structure
	mgitDir := filepath.Join(destination, ".mgit")
	mappingsDir := filepath.Join(mgitDir, "mappings")
	if err := os.MkdirAll(mappingsDir, 0755); err != nil {
		return fmt.Errorf("error creating .mgit/mappings directory: %w", err)
	}
	
	// Write the hash_mappings.json file
	mappingsPath := filepath.Join(mappingsDir, "hash_mappings.json")
	mappingsJSON, err := json.MarshalIndent(mappings, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializing mappings: %w", err)
	}
	
	if err := os.WriteFile(mappingsPath, mappingsJSON, 0644); err != nil {
		return fmt.Errorf("error writing hash_mappings.json file: %w", err)
	}
	
	// Also write to nostr_mappings.json for compatibility
	nostrMappingsPath := filepath.Join(mgitDir, "nostr_mappings.json")
	if err := os.WriteFile(nostrMappingsPath, mappingsJSON, 0644); err != nil {
		return fmt.Errorf("error writing nostr_mappings.json file: %w", err)
	}
	
	NSLog("Successfully fetched and stored MGit metadata (%d mappings)", len(mappings))
	return nil
}

// setupMGitConfig sets up the MGit configuration for the cloned repository
func setupMGitConfig(destination string, repoInfo *RepositoryInfo) error {
	// Create basic MGit config structure
	mgitDir := filepath.Join(destination, ".mgit")
	if err := os.MkdirAll(mgitDir, 0755); err != nil {
		return fmt.Errorf("error creating .mgit directory: %w", err)
	}
	
	configPath := filepath.Join(mgitDir, "config")
	
	// Create a simple config file
	configContent := fmt.Sprintf(`[repository]
	id = %s
	name = %s
	access = %s

[mgit]
	version = 1.0
	initialized = true
`, repoInfo.ID, repoInfo.Name, repoInfo.Access)
	
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("error writing MGit config: %w", err)
	}
	
	NSLog("MGit config created successfully")
	return nil
}
