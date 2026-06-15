package services

import (
	"forum/internal/models"
	"testing"
)

func TestCreateUserAndGetByEmail(t *testing.T) {
	db := setupTestDB(t)

	user := models.User{
		Username: "salas",
		Email:    "salas@gmail.com",
		Password: "hashed-password",
	}
	if err := CreateUser(db, user); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}

	found, err := GetUserByEmail(db, "salas@gmail.com")
	if err != nil {
		t.Fatalf("GetUserByEmail: %v", err)
	}
	if found == nil {
		t.Fatal("expected user, got nil")
	}
	if found.Username != "salas" {
		t.Errorf("username = %q, want salas", found.Username)
	}
	if found.Email != "salas@gmail.com" {
		t.Errorf("email = %q, want salas@gmail.com", found.Email)
	}
}

func TestGetUserByEmailNotFound(t *testing.T) {
	db := setupTestDB(t)

	found, err := GetUserByEmail(db, "missing@gmail.com")
	if err != nil {
		t.Fatalf("GetUserByEmail: %v", err)
	}
	if found != nil {
		t.Fatal("expected nil user for missing email")
	}
}

func TestGetUserByUsername(t *testing.T) {
	db := setupTestDB(t)

	user := models.User{
		Username: "salas",
		Email:    "salas@gmail.com",
		Password: "hashed-password",
	}
	if err := CreateUser(db, user); err != nil {
		t.Fatalf("CreateUser: %v", err)
	}

	found, err := GetUserByUsername(db, "salas")
	if err != nil {
		t.Fatalf("GetUserByUsername: %v", err)
	}
	if found == nil {
		t.Fatal("expected user, got nil")
	}
	if found.Email != "salas@gmail.com" {
		t.Errorf("email = %q, want salas@gmail.com", found.Email)
	}
}

func TestCreateUserDuplicateEmail(t *testing.T) {
	db := setupTestDB(t)

	first := models.User{Username: "user1", Email: "same@gmail.com", Password: "hash1"}
	second := models.User{Username: "user2", Email: "same@gmail.com", Password: "hash2"}

	if err := CreateUser(db, first); err != nil {
		t.Fatalf("CreateUser first: %v", err)
	}
	if err := CreateUser(db, second); err == nil {
		t.Fatal("expected error for duplicate email")
	}
}

func TestCreateUserDuplicateUsername(t *testing.T) {
	db := setupTestDB(t)

	first := models.User{Username: "sameuser", Email: "one@gmail.com", Password: "hash1"}
	second := models.User{Username: "sameuser", Email: "two@gmail.com", Password: "hash2"}

	if err := CreateUser(db, first); err != nil {
		t.Fatalf("CreateUser first: %v", err)
	}
	if err := CreateUser(db, second); err == nil {
		t.Fatal("expected error for duplicate username")
	}
}
