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

type MGitObjectType string

const (
	MGitCommitObject MGitObjectType = "commit"
	MGitTreeObject   MGitObjectType = "tree"
	MGitBlobObject   MGitObjectType = "blob"
)

// Signature represents the author or committer information including nostr pubkey
type Signature struct {
	Name   string
	Email  string
	Pubkey string
	When   string
}

// MGitSignature represents a signature in an MGit commit (simplified for iOS)
type MGitSignature struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Pubkey string `json:"pubkey,omitempty"`
	When   string `json:"when"` // Using string instead of time.Time for iOS compatibility
	Signature string    `json:"signature,omitempty"` // Nostr signature
}

// MCommitStruct represents an mcommit object
type MCommitStruct struct {
	Type         MGitObjectType       `json:"type"`
	MGitHash     string               `json:"mgit_hash"`
	GitHash      string               `json:"git_hash"`
	TreeHash     string               `json:"tree_hash"`
	ParentHashes []string             `json:"parent_hashes"` // MGit hashes of parents
	Author       *MGitSignature       `json:"author"`
	Committer    *MGitSignature       `json:"committer"`
	Message      string               `json:"message"`
	Metadata     map[string]string    `json:"metadata,omitempty"`
	NostrSig     string               `json:"nostr_sig,omitempty"` // Nostr signature
}

// NostrEvent represents a signed Nostr event for the commit
type NostrEvent struct {
	ID        string                 `json:"id"`
	Pubkey    string                 `json:"pubkey"`
	CreatedAt int64                  `json:"created_at"`
	Kind      int                    `json:"kind"`
	Tags      [][]string             `json:"tags"`
	Content   string                 `json:"content"`
	Sig       string                 `json:"sig"`
}

// MGitStorage handles the storage and retrieval of MGit objects
type MGitStorage struct {
	RootDir string // Usually ".mgit"
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

// AddResult represents the result of an add operation
type AddResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error"`
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