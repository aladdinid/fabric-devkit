package cmd

import "testing"

func TestPullAndDelete(t *testing.T) {
	sources := []string{"hyperledger/fabric-ca:x86_64-1.1.0", "hyperledger/fabric-tools:x86_64-1.1.0"}
	pullAndRetagImages(sources)
	err := deleteImages(sources)
	if err != nil {
		t.Fatalf("Expected: no err Got: %v", err)
	}
}
