package MGitBridge

import (
	"time"

	"github.com/go-git/go-git/v5/plumbing/object"
)

// convertToGitSignature converts our iOS-compatible Signature to go-git's object.Signature
func convertToGitSignature(sig *Signature) *object.Signature {
	// Parse the ISO 8601 time string
	when, err := time.Parse(time.RFC3339, sig.When)
	if err != nil {
		// Fallback to current time if parsing fails
		when = time.Now()
	}
	
	return &object.Signature{
		Name:  sig.Name,
		Email: sig.Email,
		When:  when,
	}
}

// convertToMGitSignature converts go-git's object.Signature to our iOS-compatible MGitSignature
func convertToMGitSignature(sig object.Signature, pubkey string) *MGitSignature {
	return &MGitSignature{
		Name:   sig.Name,
		Email:  sig.Email,
		Pubkey: pubkey,
		When:   sig.When.Format(time.RFC3339), // Convert time.Time to ISO 8601 string for iOS
	}
}