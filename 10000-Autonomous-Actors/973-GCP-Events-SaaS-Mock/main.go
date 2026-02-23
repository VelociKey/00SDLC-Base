package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"connectrpc.com/connect"
	eventsv1 "OlympusGCP-Events/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/events/v1"
	"OlympusGCP-Events/40000-Communication-Contracts/430-Protocol-Definitions/000-gen/events/v1/eventsv1connect"
	"Olympus2/90000-Enablement-Labs/P0000-pkg/000-whisper"
)

func main() {
	logger := whisper.New("SaaS-Events-Mock", "gcp_events.lpsv")
	client := eventsv1connect.NewEventsServiceClient(
		http.DefaultClient,
		"http://localhost:8094",
	)

	log.Printf("INFO SaaS-Events-Mock: Initializing...")

	// Phase 1: Pub/Sub
	log.Printf("INFO SaaS-Events-Mock: Phase 1 - Pub/Sub Publish")
	pubReq := &eventsv1.PublishRequest{
		Topic: "user-signups",
		Data:  `{"user_id": "user_456", "email": "test@olympus.dev"}`,
	}
	start1 := time.Now()
	res1, err := client.Publish(context.Background(), connect.NewRequest(pubReq))
	if err != nil {
		log.Fatalf("FAIL: Could not publish message: %v", err)
	}
	logger.Log("PubSub_Publish", "VERIFIED", "user-signups", res1.Msg.MessageId, time.Since(start1))

	// Phase 2: Cloud Tasks
	log.Printf("INFO SaaS-Events-Mock: Phase 2 - Cloud Tasks")
	taskReq := &eventsv1.CreateTaskRequest{
		Queue:        "email-delivery",
		Payload:      `{"to": "test@olympus.dev", "template": "welcome"}`,
		DelaySeconds: 10,
	}
	start2 := time.Now()
	res2, err := client.CreateTask(context.Background(), connect.NewRequest(taskReq))
	if err != nil {
		log.Fatalf("FAIL: Could not create task: %v", err)
	}
	logger.Log("CloudTasks_Create", "VERIFIED", "email-delivery", res2.Msg.TaskName, time.Since(start2))

	// Phase 3: Cloud Scheduler
	log.Printf("INFO SaaS-Events-Mock: Phase 3 - Cloud Scheduler")
	jobReq := &eventsv1.CreateJobRequest{
		Name:      "daily-cleanup",
		Schedule:  "0 0 * * *",
		TargetUrl: "http://localhost:8080/cleanup",
	}
	start3 := time.Now()
	res3, err := client.CreateJob(context.Background(), connect.NewRequest(jobReq))
	if err != nil {
		log.Fatalf("FAIL: Could not create job: %v", err)
	}
	logger.Log("Scheduler_Create", "VERIFIED", "daily-cleanup", res3.Msg.JobId, time.Since(start3))

	log.Printf("INFO SaaS-Events-Mock: All Events Assurance checks passed.")
	time.Sleep(1 * time.Second)
}
