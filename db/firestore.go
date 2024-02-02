package db

import (
	"context"
	// "flag"
	"log"

	// "google.golang.org/api/iterator"

	"cloud.google.com/go/firestore"
)

var firestoreDB *firestore.Client

func CreateFirestoreClient(ctx context.Context) {
	// sets gcp project id
	projectID := "YOUR_PROJECT_ID"

	// overwrite with project flags
	// flag.StringVar(&projectID, "project", projectID, "The GCP project ID.")
	// flag.Parse()

	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	// close client when finished
	// defer client.Close()

	firestoreDB = client
}

