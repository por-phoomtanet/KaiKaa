package token

import "testing"

func TestGenerateAndParse(t *testing.T) {
	secret := "test-secret"
	tok, err := Generate(secret, "user1", "shop1")
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	claims, err := Parse(secret, tok)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if claims.UserID != "user1" {
		t.Errorf("UserID = %q, want %q", claims.UserID, "user1")
	}
	if claims.ShopID != "shop1" {
		t.Errorf("ShopID = %q, want %q", claims.ShopID, "shop1")
	}
}

func TestParseWrongSecret(t *testing.T) {
	tok, _ := Generate("secret-a", "u", "s")
	if _, err := Parse("secret-b", tok); err == nil {
		t.Error("expected error parsing with wrong secret, got nil")
	}
}

func TestParseTampered(t *testing.T) {
	tok, _ := Generate("secret", "u", "s")
	if _, err := Parse("secret", tok+"tampered"); err == nil {
		t.Error("expected error parsing tampered token, got nil")
	}
}

func TestParseGarbage(t *testing.T) {
	if _, err := Parse("secret", "not-a-jwt"); err == nil {
		t.Error("expected error parsing garbage, got nil")
	}
}
