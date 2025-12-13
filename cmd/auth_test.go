package main

import (
	"bytes"
	"encoding/json"
	"go_purple/internal/auth"
	"go_purple/internal/user"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDb() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic(err)
	}

	return db
}

func initData(db *gorm.DB) {
	db.Create(&user.User{
		Name:     "Ivan",
		Email:    "a1@test.ru",
		Password: "$2a$10$m1LlU5w8kxYNwPFsfvFVcuVAHcJeazvMvmOc2Ru1vrYYz5vpbmdbG",
	})
}

func removeData(db *gorm.DB) {
	db.Unscoped().Where("email = ?", "a1@test.ru").Delete(&user.User{})
}

func TestLoginSuccess(t *testing.T) {
	// Prepare
	db := initDb()
	initData(db)
	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "a1@test.ru",
		Password: "1",
	})

	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected 200 got %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	var resData auth.LoginResponse
	err = json.Unmarshal(body, &resData)

	if err != nil {
		t.Fatal(err)
	}

	if resData.Token == "" {
		t.Fatal("Token is empty")
	}

	removeData(db)

}

func TestLoginFail(t *testing.T) {

	// Prepare
	db := initDb()
	initData(db)
	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "a2@test.ru",
		Password: "2",
	})

	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("Expected 401 got %d", res.StatusCode)
	}

	removeData(db)
}
