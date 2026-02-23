package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"connectrpc.com/connect"
	firebasev1 "OlympusGCP-Firebase/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/firebase/v1"
	"OlympusGCP-Firebase/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/firebase/v1/firebasev1connect"
	"Olympus2/90000-Enablement-Labs/P0000-pkg/000-whisper"
)

func main() {
	logger := whisper.New("SaaS-Firebase-Mock", "gcp_firebase.lpsv")
	client := firebasev1connect.NewFirebaseServiceClient(
		http.DefaultClient,
		"http://localhost:8099",
	)

	slog.Info("SaaS-Firebase-Mock: Initializing...")

	// Phase 1: Auth
	slog.Info("SaaS-Firebase-Mock: Phase 1 - Auth User Creation")
	userReq := &firebasev1.CreateUserRequest{
		Email:       "tester@olympus.dev",
		Password:    "super-secret-123",
		DisplayName: "Olympus Tester",
	}
	start1 := time.Now()
	res1, err := client.CreateUser(context.Background(), connect.NewRequest(userReq))
	if err != nil {
		slog.Error("FAIL: Could not create user", "error", err)
		return
	}
	logger.Log("Auth_CreateUser", "VERIFIED", userReq.Email, res1.Msg.Uid, time.Since(start1))

	// Phase 2: Firestore
	slog.Info("SaaS-Firebase-Mock: Phase 2 - Firestore Doc Set")
	docReq := &firebasev1.SetDocumentRequest{
		Collection: "profiles",
		DocId:      res1.Msg.Uid,
		DataJson:   `{"bio": "Software Engineer", "interests": ["Go", "Cloud"]}`,
	}
	start2 := time.Now()
	_, err = client.SetDocument(context.Background(), connect.NewRequest(docReq))
	if err != nil {
		slog.Error("FAIL: Could not set document", "error", err)
		return
	}
	logger.Log("Firestore_SetDoc", "VERIFIED", "profiles/"+res1.Msg.Uid, "SUCCESS", time.Since(start2))

	slog.Info("SaaS-Firebase-Mock: All Firebase Assurance checks passed.")
	time.Sleep(1 * time.Second)
}
