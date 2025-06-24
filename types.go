package MGitBridge

// Basic result types for iOS compatibility

// HelpResult represents the result of the help operation
type HelpResult struct {
	Success  bool   `json:"success"`
	HelpText string `json:"help_text"`
	Message  string `json:"message"`
}

// LogResult represents the result of logging tests
type LogResult struct {
	Success bool   `json:"success"`
	Result  string `json:"result"`
	Message string `json:"message"`
}

// MathResult represents the result of simple math operations
type MathResult struct {
	Success bool   `json:"success"`
	Result  int    `json:"result"`
	Message string `json:"message"`
}

// CloneResult represents the result of a clone operation
type CloneResult struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	RepoID    string `json:"repo_id"`
	RepoName  string `json:"repo_name"`
	LocalPath string `json:"local_path"`
}

// RepositoryInfo represents information about a repository
type RepositoryInfo struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Access string `json:"access"`
}

// CommitResult represents the result of a commit operation
type CommitResult struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	GitHash    string `json:"git_hash"`
	MGitHash   string `json:"mgit_hash"`
	CommitMsg  string `json:"commit_message"`
}

// PushResult represents the result of a push operation
type PushResult struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	CommitHash string `json:"commit_hash"`
}

// PullResult represents the result of a pull operation
type PullResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Changes int    `json:"changes"`
}

// MGitSignature represents a signature in an MGit commit (simplified for iOS)
type MGitSignature struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Pubkey string `json:"pubkey,omitempty"`
	When   string `json:"when"` // Using string instead of time.Time for iOS compatibility
}

// MCommitInfo represents simplified MGit commit information for iOS
type MCommitInfo struct {
	MGitHash     string        `json:"mgit_hash"`
	GitHash      string        `json:"git_hash"`
	Message      string        `json:"message"`
	Author       MGitSignature `json:"author"`
	Committer    MGitSignature `json:"committer"`
	ParentHashes []string      `json:"parent_hashes"`
	TreeHash     string        `json:"tree_hash"`
}