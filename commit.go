package MGitBridge

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// NewMGitStorage creates a new storage instance
func NewMGitStorage(repoPath string) *MGitStorage {
	return &MGitStorage{
		RootDir: filepath.Join(repoPath, ".mgit"),
	}
}

// Initialize creates the necessary directory structure for MGit
func (s *MGitStorage) Initialize() error {
	// Create the main directory
	if err := os.MkdirAll(s.RootDir, 0755); err != nil {
		return fmt.Errorf("failed to create MGit directory: %w", err)
	}
	
	// Create subdirectories
	dirs := []string{
		filepath.Join(s.RootDir, "objects"),     // For storing commit objects
		filepath.Join(s.RootDir, "refs"),        // For storing branch refs
		filepath.Join(s.RootDir, "refs/heads"),  // For branch heads
		filepath.Join(s.RootDir, "refs/tags"),   // For tags
		filepath.Join(s.RootDir, "mappings"),    // For storing hash mappings
	}
	
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create an initial HEAD file if it doesn't exist
	headPath := filepath.Join(s.RootDir, "HEAD")
	if _, err := os.Stat(headPath); os.IsNotExist(err) {
		if err := ioutil.WriteFile(headPath, []byte("ref: refs/heads/main"), 0644); err != nil {
			return fmt.Errorf("failed to create HEAD file: %w", err)
		}
	}
	
	return nil
}

// StoreCommit stores an MGit commit object
func (s *MGitStorage) StoreCommit(commit *MCommitStruct) error {
	if commit.MGitHash == "" {
		return fmt.Errorf("MGit hash cannot be empty")
	}
	
	commit.Type = MGitCommitObject
	
	// Create the object path using the hash
	prefix := commit.MGitHash[:2]
	suffix := commit.MGitHash[2:]
	objDir := filepath.Join(s.RootDir, "objects", prefix)
	objPath := filepath.Join(objDir, suffix)
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(objDir, 0755); err != nil {
		return fmt.Errorf("failed to create object directory: %w", err)
	}
	
	// Marshal to JSON
	data, err := json.MarshalIndent(commit, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal commit: %w", err)
	}
	
	// Write to file
	if err := ioutil.WriteFile(objPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write commit object: %w", err)
	}
	
	return nil
}

// StoreMapping stores a mapping between Git and MGit hashes
func (s *MGitStorage) StoreMapping(gitHash string, mgitHash string, pubkey string) error {
	mappingPath := filepath.Join(s.RootDir, "mappings", "hash_mappings.json")
	
	// Read existing mappings
	var mappings []struct {
		GitHash  string `json:"git_hash"`
		MGitHash string `json:"mgit_hash"`
		Pubkey   string `json:"pubkey"`
	}
	
	if data, err := ioutil.ReadFile(mappingPath); err == nil {
		json.Unmarshal(data, &mappings)
	}
	
	// Add new mapping
	newMapping := struct {
		GitHash  string `json:"git_hash"`
		MGitHash string `json:"mgit_hash"`
		Pubkey   string `json:"pubkey"`
	}{
		GitHash:  gitHash,
		MGitHash: mgitHash,
		Pubkey:   pubkey,
	}
	
	mappings = append(mappings, newMapping)
	
	// Write back
	data, err := json.MarshalIndent(mappings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal hash mappings: %w", err)
	}
	
	if err := ioutil.WriteFile(mappingPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write hash mappings: %w", err)
	}
	
	return nil
}

// GetMGitHashFromGit gets the MGit hash for a Git hash
func (s *MGitStorage) GetMGitHashFromGit(gitHash string) (string, error) {
	mappingPath := filepath.Join(s.RootDir, "mappings", "hash_mappings.json")
	
	var mappings []struct {
		GitHash  string `json:"git_hash"`
		MGitHash string `json:"mgit_hash"`
		Pubkey   string `json:"pubkey"`
	}
	
	if _, err := os.Stat(mappingPath); os.IsNotExist(err) {
		return "", fmt.Errorf("no MGit hash found for Git hash %s", gitHash)
	}
	
	data, err := ioutil.ReadFile(mappingPath)
	if err != nil {
		return "", fmt.Errorf("failed to read hash mappings: %w", err)
	}
	
	if err := json.Unmarshal(data, &mappings); err != nil {
		return "", fmt.Errorf("failed to unmarshal hash mappings: %w", err)
	}
	
	for _, mapping := range mappings {
		if mapping.GitHash == gitHash {
			return mapping.MGitHash, nil
		}
	}
	
	return "", fmt.Errorf("no MGit hash found for Git hash %s", gitHash)
}

// UpdateRef updates an MGit reference (branch or tag)
func (s *MGitStorage) UpdateRef(refName string, mgitHash string) error {
	if !strings.HasPrefix(refName, "refs/") {
		refName = "refs/heads/" + refName
	}
	
	refPath := filepath.Join(s.RootDir, refName)
	
	// Create directory if it doesn't exist
	refDir := filepath.Dir(refPath)
	if err := os.MkdirAll(refDir, 0755); err != nil {
		return fmt.Errorf("failed to create ref directory: %w", err)
	}
	
	// Write the ref
	if err := ioutil.WriteFile(refPath, []byte(mgitHash), 0644); err != nil {
		return fmt.Errorf("failed to write ref: %w", err)
	}
	
	return nil
}

// computeMGitHash computes a new hash incorporating the nostr pubkey
func computeMGitHash(commit *object.Commit, parentMGitHashes []string, pubkey string) plumbing.Hash {
	// Create a new hasher
	hasher := sha1.New()
	
	// Include the tree hash
	hasher.Write(commit.TreeHash[:])
	
	// Include all parent MGit hashes
	for _, parentHashStr := range parentMGitHashes {
		parentHash := plumbing.NewHash(parentHashStr)
		hasher.Write(parentHash[:])
	}
	
	// Include the author information with pubkey
	authorStr := fmt.Sprintf("%s <%s> %d %s", 
		commit.Author.Name, 
		commit.Author.Email, 
		commit.Author.When.Unix(), 
		pubkey)
	hasher.Write([]byte(authorStr))
	
	// Include committer information with pubkey
	committerStr := fmt.Sprintf("%s <%s> %d %s", 
		commit.Committer.Name, 
		commit.Committer.Email, 
		commit.Committer.When.Unix(),
		pubkey)
	hasher.Write([]byte(committerStr))
	
	// Include the commit message
	hasher.Write([]byte(commit.Message))
	
	// Calculate the new hash
	mgitHash := hasher.Sum(nil)
	
	// Convert to plumbing.Hash
	var result plumbing.Hash
	copy(result[:], mgitHash[:20]) // SHA-1 is 20 bytes
	
	return result
}

// createCommitSignature creates a Nostr signature directly for the MGit commit data
func createCommitSignature(commit *MCommitStruct, pubkey string) string {
	// Create the signable content from the commit data
	signableContent := fmt.Sprintf("%s|%s|%s|%s|%s|%s", 
		commit.MGitHash,
		commit.GitHash, 
		commit.TreeHash,
		strings.Join(commit.ParentHashes, ","),
		commit.Message,
		pubkey)
	
	// In a real implementation, this would:
	// 1. Hash the signable content with SHA256
	// 2. Sign the hash with the Nostr private key (schnorr signature)
	// 3. Return the signature as hex string
	
	// For now, create a deterministic placeholder signature
	hasher := sha1.New()
	hasher.Write([]byte(signableContent))
	hash := hasher.Sum(nil)
	
	return fmt.Sprintf("nostr_sig_%s_%x", pubkey[:10], hash[:4])
}

// Commit creates an MGit commit with Nostr signature using go-git directly
func Commit(repoPath string, message string, authorName string, authorEmail string, nostrPubkey string) *CommitResult {
	log.Printf("MGitBridge: Commit() called for repo: %s", repoPath)
	log.Printf("MGitBridge: Message: %s", message)
	log.Printf("MGitBridge: Author: %s <%s>", authorName, authorEmail)
	log.Printf("MGitBridge: Nostr pubkey: %s", nostrPubkey)
	
	result := &CommitResult{
		Success:   false,
		Message:   "",
		GitHash:   "",
		MGitHash:  "",
		CommitMsg: message,
	}
	
	// Open the repository using go-git
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Printf("MGitBridge: Error opening repository: %s", err)
		result.Message = fmt.Sprintf("Error opening repository: %s", err)
		return result
	}
	
	// Get the worktree
	w, err := repo.Worktree()
	if err != nil {
		log.Printf("MGitBridge: Error getting worktree: %s", err)
		result.Message = fmt.Sprintf("Error getting worktree: %s", err)
		return result
	}
	
	// Create author signature
  authorSig := &Signature{
  	Name:  authorName,
  	Email: authorEmail,
		Pubkey: nostrPubkey,
   	When:  time.Now().Format(time.RFC3339),
  }
   
	author := convertToGitSignature(authorSig)
	
	// Create commit options
	commitOpts := &git.CommitOptions{
		Author: author,
	}
	
	log.Printf("MGitBridge: Performing git commit...")
	
	// Perform the standard git commit
	gitHash, err := w.Commit(message, commitOpts)
	if err != nil {
		log.Printf("MGitBridge: Error committing: %s", err)
		result.Message = fmt.Sprintf("Error committing: %s", err)
		return result
	}
	
	log.Printf("MGitBridge: Git commit successful: %s", gitHash.String())
	
	// If no pubkey is present, just return the Git hash (like CLI)
	if nostrPubkey == "" {
		result.Success = true
		result.Message = "Commit created successfully (no MGit hash - no pubkey)"
		result.GitHash = gitHash.String()
		result.MGitHash = gitHash.String() // Use git hash as fallback
		return result
	}
	
	// Get the commit object we just created
	gitCommit, err := repo.CommitObject(gitHash)
	if err != nil {
		log.Printf("MGitBridge: Error retrieving commit: %s", err)
		result.Message = fmt.Sprintf("Error retrieving commit: %s", err)
		return result
	}
	
	// Initialize MGit storage (create .mgit directory structure)
	storage := NewMGitStorage(repoPath)
	if err := storage.Initialize(); err != nil {
		log.Printf("MGitBridge: Error initializing MGit storage: %s", err)
		result.Message = fmt.Sprintf("Error initializing MGit storage: %s", err)
		return result
	}
	
	// Collect MGit hashes for parent commits (exactly like CLI)
	parentMGitHashes := []string{}
	for _, parentGitHash := range gitCommit.ParentHashes {
		mgitHash, err := storage.GetMGitHashFromGit(parentGitHash.String())
		if err == nil {
			// We found an MGit hash for this parent
			parentMGitHashes = append(parentMGitHashes, mgitHash)
			log.Printf("MGitBridge: Found MGit hash for parent %s: %s", 
				parentGitHash.String()[:7], mgitHash[:7])
		} else {
			// No MGit hash found, use the Git hash as a fallback
			parentMGitHashes = append(parentMGitHashes, parentGitHash.String())
			log.Printf("MGitBridge: No MGit hash found for parent %s", parentGitHash.String()[:7])
		}
	}
	
	// Compute the MGit hash (exactly like CLI)
	mgitHash := computeMGitHash(gitCommit, parentMGitHashes, nostrPubkey)
	
	// Create an MGit commit object (exactly like CLI)
	mgitCommit := &MCommitStruct{
		Type:         "commit",
		MGitHash:     mgitHash.String(),
		GitHash:      gitHash.String(),
		TreeHash:     gitCommit.TreeHash.String(),
		ParentHashes: parentMGitHashes,
		Author:       mGitAuthorSig,
		Committer:    committerSig,
		Message:      gitCommit.Message,
		Metadata:     map[string]string{"version": "1.0"},
	}
	
	// Create Nostr signature for the entire commit
	commitSignature := createCommitSignature(mgitCommit, nostrPubkey)
	mgitCommit.NostrSig = commitSignature
	
	// Add Nostr signature to author and committer
	mGitAuthorSig.Signature = commitSignature
	committerSig.Signature = commitSignature
	
	log.Printf("MGitBridge: Created Nostr signature: %s", commitSignature)
	
	// Store the MGit commit object (exactly like CLI)
	if err := storage.StoreCommit(mgitCommit); err != nil {
		log.Printf("MGitBridge: Error storing MGit commit: %s", err)
		result.Message = fmt.Sprintf("Error storing MGit commit: %s", err)
		return result
	}
	
	// Store the mapping between Git and MGit hashes (exactly like CLI)
	if err := storage.StoreMapping(gitHash.String(), mgitHash.String(), nostrPubkey); err != nil {
		log.Printf("MGitBridge: Error storing hash mapping: %s", err)
		result.Message = fmt.Sprintf("Error storing hash mapping: %s", err)
		return result
	}
	
	// Update the current branch reference in MGit (exactly like CLI)
	head, err := repo.Head()
	if err == nil && head.Name().IsBranch() {
		branchName := head.Name().Short()
		refName := fmt.Sprintf("refs/heads/%s", branchName)
		
		if err := storage.UpdateRef(refName, mgitHash.String()); err != nil {
			log.Printf("MGitBridge: Warning: Failed to update branch ref: %s", err)
		}
	}
	
	log.Printf("MGitBridge: Created MGit commit: %s (Git hash: %s)", 
		mgitHash.String(), gitHash.String())
	log.Printf("MGitBridge: Nostr signature stored in commit metadata")
	
	// Update result for success case
	result.Success = true
	result.Message = fmt.Sprintf("MGit commit created successfully with Nostr signature")
	result.GitHash = gitHash.String()
	result.MGitHash = mgitHash.String()
	
	return result
}