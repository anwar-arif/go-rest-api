package utils

import (
	"crypto/md5"
	"reflect"
	"testing"
)

func TestGenerateSalt(t *testing.T) {
	salt := GenerateSalt()
	if salt == nil {
		t.Errorf("unexpected salt: %v", salt)
	}
	if len(*salt) != SaltSize {
		t.Errorf("unexpected salt size")
	}
}

func TestGetHasher(t *testing.T) {
	hasher := GetHasher()
	expectedHasher := md5.New()

	if reflect.DeepEqual(hasher, expectedHasher) == false {
		t.Errorf("unexpected hasher method")
	}
}

func TestHashPassword(t *testing.T) {
	salt := "my_complex_salt"
	password := "my_complex_password"
	hashedPassword := "af88f052a6dd551480ea5749e93eb1af"

	if hashedPassword != HashPassword(password, salt) {
		t.Errorf("hashed password with salt is not equal")
	}
}

func TestIsSamePassword(t *testing.T) {
	passwords := []struct {
		Salt            string
		CurrentPassword string
		HashedPassword  string
		Expected        bool
	}{
		{"my_complex_salt", "my_complex_password", "af88f052a6dd551480ea5749e93eb1af", true},
		{"my_complex_salt", "my_complex_password", "af88f052a6dd551480ea5749e93eb1az", false},
		{"my_complex_salt", "my_complex_password_incorrect", "af88f052a6dd551480ea5749e93eb1af", false},
	}

	for _, p := range passwords {
		isSame := IsSamePassword(p.HashedPassword, p.CurrentPassword, p.Salt)
		if isSame != p.Expected {
			t.Errorf("passwords do not match, expected %v, got %v", p.Expected, isSame)
		}
	}
}
