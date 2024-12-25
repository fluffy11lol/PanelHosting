package models

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestUserSerialization(t *testing.T) {
	tests := []struct {
		name string
		user User
	}{
		{
			name: "complete user",
			user: User{
				ID:       "123",
				Username: "testuser",
				Email:    "test@example.com",
				Password: "hashedpassword",
				Active:   true,
				Token:    "jwt-token",
			},
		},
		{
			name: "minimal user",
			user: User{
				Username: "minimal",
				Password: "pass",
			},
		},
		{
			name: "inactive user",
			user: User{
				ID:       "456",
				Username: "inactive",
				Email:    "inactive@example.com",
				Password: "pass",
				Active:   false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling
			data, err := json.Marshal(tt.user)
			if err != nil {
				t.Errorf("Failed to marshal user: %v", err)
			}

			// Test JSON unmarshaling
			var unmarshaled User
			err = json.Unmarshal(data, &unmarshaled)
			if err != nil {
				t.Errorf("Failed to unmarshal user: %v", err)
			}

			// Compare fields
			if unmarshaled.ID != tt.user.ID {
				t.Errorf("ID mismatch: got %v, want %v", unmarshaled.ID, tt.user.ID)
			}
			if unmarshaled.Username != tt.user.Username {
				t.Errorf("Username mismatch: got %v, want %v", unmarshaled.Username, tt.user.Username)
			}
			if unmarshaled.Email != tt.user.Email {
				t.Errorf("Email mismatch: got %v, want %v", unmarshaled.Email, tt.user.Email)
			}
			if unmarshaled.Password != tt.user.Password {
				t.Errorf("Password mismatch: got %v, want %v", unmarshaled.Password, tt.user.Password)
			}
			if unmarshaled.Active != tt.user.Active {
				t.Errorf("Active mismatch: got %v, want %v", unmarshaled.Active, tt.user.Active)
			}
			if unmarshaled.Token != tt.user.Token {
				t.Errorf("Token mismatch: got %v, want %v", unmarshaled.Token, tt.user.Token)
			}
		})
	}
}

func TestUserFieldTags(t *testing.T) {
	user := User{}
	userType := reflect.TypeOf(user)

	expectedTags := map[string]struct {
		jsonTag string
		dbTag   string
	}{
		"ID":       {"id", "id"},
		"Username": {"username", "username"},
		"Email":    {"email", "email"},
		"Password": {"password", "password"},
		"Active":   {"active", "active"},
		"Token":    {"token", ""},
	}

	for fieldName, expected := range expectedTags {
		field, found := userType.FieldByName(fieldName)
		if !found {
			t.Errorf("Field %s not found in User struct", fieldName)
			continue
		}

		jsonTag := field.Tag.Get("json")
		if jsonTag != expected.jsonTag {
			t.Errorf("Field %s json tag = %s, want %s", fieldName, jsonTag, expected.jsonTag)
		}

		dbTag := field.Tag.Get("db")
		if dbTag != expected.dbTag {
			t.Errorf("Field %s db tag = %s, want %s", fieldName, dbTag, expected.dbTag)
		}
	}
}
