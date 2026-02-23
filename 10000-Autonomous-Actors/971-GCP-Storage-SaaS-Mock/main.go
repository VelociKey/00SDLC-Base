package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"connectrpc.com/connect"

	"Olympus2/90000-Enablement-Labs/P0000-pkg/000-whisper"
	storagev1 "OlympusGCP-Storage/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/storage/v1"
	"OlympusGCP-Storage/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/storage/v1/storagev1connect"
)

func main() {
	slog.Info("SaaS-Storage-Mock: Initializing...")
	w := whisper.New("saas_storage", "saas_storage.lpsv")
	defer w.Close()

	client := storagev1connect.NewStorageServiceClient(
		http.DefaultClient,
		"http://localhost:8091",
	)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Phase 1: Bucket Creation
	slog.Info("SaaS-Storage-Mock: Phase 1 - Provisioning Assets Bucket")
	start := time.Now()
	bucketName := "olympus-assurance-assets"
	_, err := client.CreateBucket(ctx, connect.NewRequest(&storagev1.CreateBucketRequest{
		Name: bucketName,
	}))
	if err != nil {
		w.Log("CREATE_BUCKET", "FAILURE", bucketName, err.Error(), time.Since(start))
		log.Fatalf("FAIL: Could not create bucket: %v", err)
	}
	w.Log("CREATE_BUCKET", "SUCCESS", bucketName, "Bucket provisioned", time.Since(start))

	// Phase 2: Object Upload
	slog.Info("SaaS-Storage-Mock: Phase 2 - Uploading Application Binary")
	start = time.Now()
	objectName := "app-v1.bin"
	payload := []byte("SIMULATED_BINARY_CONTENT_2026")
	_, err = client.UploadObject(ctx, connect.NewRequest(&storagev1.UploadObjectRequest{
		Bucket: bucketName,
		Name:   objectName,
		Data:   payload,
	}))
	if err != nil {
		w.Log("UPLOAD", "FAILURE", objectName, err.Error(), time.Since(start))
		log.Fatalf("FAIL: Upload failed: %v", err)
	}
	w.Log("UPLOAD", "SUCCESS", objectName, fmt.Sprintf("Uploaded %d bytes", len(payload)), time.Since(start))

	// Phase 3: Signed URL Generation
	slog.Info("SaaS-Storage-Mock: Phase 3 - Generating Download Access")
	start = time.Now()
	resp, err := client.GetDownloadURL(ctx, connect.NewRequest(&storagev1.GetDownloadURLRequest{
		Bucket: bucketName,
		Name:   objectName,
	}))
	if err != nil {
		w.Log("GET_URL", "FAILURE", objectName, err.Error(), time.Since(start))
		log.Fatalf("FAIL: URL generation failed: %v", err)
	}
	fmt.Printf("   SUCCESS: Generated GCS Access URL: %s\n", resp.Msg.Url)
	w.Log("GET_URL", "SUCCESS", objectName, resp.Msg.Url, time.Since(start))

	slog.Info("SaaS-Storage-Mock: All Storage layers PASSED")
}
