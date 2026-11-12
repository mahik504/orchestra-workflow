package resources

import "testing"

func TestRejectPrivateURL(t *testing.T) {
	if err := RejectPrivateURL("http://127.0.0.1/"); err == nil {
		t.Fatal("127.0.0.1 was accepted")
	}
	if err := RejectPrivateURL("http://169.254.169.254/"); err == nil {
		t.Fatal("metadata address was accepted")
	}
	if err := RejectPrivateURL("https://example.invalid/orchestra"); err != nil {
		t.Fatalf("a DNS failure is not a private address: %v", err)
	}
}
