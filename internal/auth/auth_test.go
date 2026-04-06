package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)


func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {

	tests := []struct{
		name	    string
		tokenString	string
		tokenSecret string
		expiresIn	string
		id			uuid.UUID
		wantErr		bool
	}{
		{
			name:		 "Valid JWT",
			tokenString: "tokenString",
			tokenSecret: "tokenSecret",
			expiresIn:	 "1m",
			wantErr:	 false,
		},
		// {
		// 	name:		 "Invalid JWT",
		// 	tokenString: "tokenString",
		// 	tokenSecret: "tokenSecret",
		// 	expiresIn:	 "0.001s",
		// 	wantErr:	 true,
		// },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			userID := uuid.New()
			expiresIn, _ := time.ParseDuration(tt.expiresIn)
			tokenString, err := MakeJWT(userID, tt.tokenSecret, expiresIn)
			if (err != nil) != tt.wantErr {
				t.Errorf("MakeJWT() error = %v, wantErr %v", err, tt.wantErr)
			}

			validatedUser, err := ValidateJWT(tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
			} 
			if tt.wantErr != true && validatedUser != userID {
				t.Errorf("ValidateJWT() expects user = %v, got %v", userID, validatedUser)
			}
		})
	}
}
