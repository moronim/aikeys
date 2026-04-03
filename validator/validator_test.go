package validator

import "testing"

func TestValidateOpenAIKey(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		wantValid bool
	}{
		{"valid sk- key", "sk-abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLM", true},
		{"valid sk-proj- key", "sk-proj-abcdefghijklmnopqrstuvwxyz012345678", true},
		{"too short", "sk-abc", false},
		{"wrong prefix", "pk-abcdefghijklmnopqrstuvwxyz0123456789ABCDEF", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Validate("OPENAI_API_KEY", tt.value)
			if result.Valid != tt.wantValid {
				t.Errorf("Validate(OPENAI_API_KEY, %q).Valid = %v, want %v", tt.value, result.Valid, tt.wantValid)
			}
			if !tt.wantValid && result.Error == "" {
				t.Error("invalid value should produce an error message")
			}
		})
	}
}

func TestValidateAnthropicKey(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		wantValid bool
	}{
		{"valid", "sk-ant-abcdefghijklmnopqrstuvwxyz012345", true},
		{"wrong prefix", "sk-abcdefghijklmnopqrstuvwxyz0123456789", false},
		{"too short", "sk-ant-abc", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Validate("ANTHROPIC_API_KEY", tt.value)
			if result.Valid != tt.wantValid {
				t.Errorf("Validate(ANTHROPIC_API_KEY, %q).Valid = %v, want %v", tt.value, result.Valid, tt.wantValid)
			}
		})
	}
}

func TestValidateHuggingFaceToken(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		wantValid bool
	}{
		{"valid", "hf_abcDEF123", true},
		{"wrong prefix", "token_abc123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Validate("HF_TOKEN", tt.value)
			if result.Valid != tt.wantValid {
				t.Errorf("Validate(HF_TOKEN, %q).Valid = %v, want %v", tt.value, result.Valid, tt.wantValid)
			}
		})
	}
}

func TestValidateReplicateToken(t *testing.T) {
	result := Validate("REPLICATE_API_TOKEN", "r8_abcdef123456")
	if !result.Valid {
		t.Error("valid Replicate token should pass")
	}

	result = Validate("REPLICATE_API_TOKEN", "wrong_prefix")
	if result.Valid {
		t.Error("invalid Replicate token should fail")
	}
	if result.Error == "" {
		t.Error("invalid Replicate token should have an error message")
	}
}

func TestValidateWandbKey(t *testing.T) {
	result := Validate("WANDB_API_KEY", "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2")
	if !result.Valid {
		t.Error("valid 40-char hex W&B key should pass")
	}

	result = Validate("WANDB_API_KEY", "not-hex-at-all")
	if result.Valid {
		t.Error("invalid W&B key should fail")
	}
}

func TestValidateLangchainKey(t *testing.T) {
	result := Validate("LANGCHAIN_API_KEY", "lsv2_abc_123_def")
	if !result.Valid {
		t.Error("valid LangChain key should pass")
	}

	result = Validate("LANGCHAIN_API_KEY", "wrong_prefix")
	if result.Valid {
		t.Error("invalid LangChain key should fail")
	}
}

func TestValidateLangchainTracing(t *testing.T) {
	result := Validate("LANGCHAIN_TRACING_V2", "true")
	if !result.Valid {
		t.Error("'true' should be valid for LANGCHAIN_TRACING_V2")
	}

	result = Validate("LANGCHAIN_TRACING_V2", "yes")
	if result.Valid {
		t.Error("'yes' should be invalid for LANGCHAIN_TRACING_V2")
	}
}

func TestValidateUnknownKey(t *testing.T) {
	result := Validate("MY_CUSTOM_SECRET", "anything-goes")
	if !result.Valid {
		t.Error("unknown keys should always pass validation")
	}
	if result.Error != "" {
		t.Error("unknown keys should have no error")
	}
}

func TestValidateKeyWithNoPattern(t *testing.T) {
	// HF_HOME has no pattern defined
	result := Validate("HF_HOME", "/some/path")
	if !result.Valid {
		t.Error("key without pattern should always pass")
	}
}

func TestValidateOrgID(t *testing.T) {
	result := Validate("OPENAI_ORG_ID", "org-abc123xyz")
	if !result.Valid {
		t.Error("valid org ID should pass")
	}

	result = Validate("OPENAI_ORG_ID", "notorg-123")
	if result.Valid {
		t.Error("invalid org ID should fail")
	}
}

func TestValidateProjectID(t *testing.T) {
	result := Validate("OPENAI_PROJECT_ID", "proj_abc123")
	if !result.Valid {
		t.Error("valid project ID should pass")
	}

	result = Validate("OPENAI_PROJECT_ID", "project-123")
	if result.Valid {
		t.Error("invalid project ID should fail")
	}
}

func TestValidateErrorIncludesKeyName(t *testing.T) {
	result := Validate("OPENAI_API_KEY", "bad-value")
	if !result.Valid && result.Error == "" {
		t.Fatal("expected error message for invalid key")
	}
	if result.Valid {
		t.Fatal("expected validation to fail")
	}
	if len(result.Error) == 0 {
		t.Error("error message should not be empty")
	}
	// Error should mention the key name for clear CLI output
	if !contains(result.Error, "OPENAI_API_KEY") {
		t.Errorf("error message %q should contain key name", result.Error)
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
