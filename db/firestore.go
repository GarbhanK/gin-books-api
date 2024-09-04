package db

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/firestore"
)

func CreateFirestoreClient(ctx context.Context) *firestore.Client {
	// sets gcp project id
	projectID := os.Getenv("GCP_PROJECT_ID")

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	return client
}
