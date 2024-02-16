package db

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
)

func CreateFirestoreClient(ctx context.Context) *firestore.Client {
	// sets gcp project id
	projectID := "learn-gcp-2112"

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	return client
}
