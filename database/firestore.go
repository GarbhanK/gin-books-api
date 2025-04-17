package database

import (
	"context"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/garbhank/gin-books-api/models"
	"google.golang.org/api/iterator"
)

type DB interface {
	Conn()
	Insert()
	Get()
	Delete()
}

type Firestore struct {
	Client    firestore.Client
	projectId string
}

func NewFirestore() *Firestore {
	projectId := os.Getenv("GCP_PROJECT_ID")
	if projectId == "" {
		log.Fatalf("Warning: GCP_PROJECT_ID env variable is not set. Firestore client will not be initialized")
	}

	return &Firestore{
		projectId: projectId,
	}
}

func (f Firestore) Conn(ctx context.Context) error {
	// sets gcp project id
	if f.projectId == "" {
		log.Fatalf("Warning: GCP_PROJECT_ID env variable is not set. Firestore client will not be initialized")
	}

	fs_client, err := firestore.NewClient(ctx, f.projectId)
	if err != nil {
		log.Printf("failed to create client: %v", err)
		return err
	}

	// add client to db struct
	f.Client = *fs_client
	return nil
}

func (f Firestore) Get(ctx context.Context, col string, key string, val string) ([]models.Book, error) {
	// create Books slice
	var bookDocs []models.Book

	// iterate over books collection in firestore
	iter := f.Client.Collection(col).Where(key, "==", val).Documents(ctx)
	defer iter.Stop() // clean up resources

	// loop until all documents matching title are added to books array
	for {
		var booksBuffer models.Book

		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Failed to iterate:\n%v", err)
			return nil, err
		}

		log.Println(doc.Data())
		if err := doc.DataTo(&booksBuffer); err != nil {
			log.Printf("can't cast docsnap to Book:\n%v", err)
			return nil, err
		}

		bookDocs = append(bookDocs, booksBuffer)
	}

	return bookDocs, nil
}

func (f Firestore) Insert(ctx context.Context, col string, data models.InsertBookInput) (models.Book, error) {
	// create a DocumentReference
	_, _, err := f.Client.Collection(col).Add(ctx, data)
	if err != nil {
		log.Printf("Failed adding document:\n%v", err)
		return models.Book{}, err
	}

	return models.Book(data), nil
}

func (f Firestore) Drop(ctx context.Context, col string, key string, val string) error {
	bulkwriter := f.Client.BulkWriter(ctx)

	for {
		iter := f.Client.Collection(col).Where(key, "==", val).Documents(ctx)
		numDeleted := 0

		for {
			var bookBuffer models.Book
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to iterate:\n%v", err)
			}

			log.Println(doc.Data())

			if err := doc.DataTo(&bookBuffer); err != nil {
				log.Fatalf("Can't cast docsnap to Book:\n%v", err)
			}

			// lowercase titles for matching book titles
			valueLower := strings.ToLower(val)
			// TODO: looks like this won't work
			parsedFirebaseTitle := strings.ToLower(bookBuffer.Title)
			if parsedFirebaseTitle == valueLower {
				bulkwriter.Delete(doc.Ref)
				log.Printf("Deleted record: %s", val)
				numDeleted++
			}
		}

		// if there are no docs to delete, process over
		if numDeleted == 0 {
			bulkwriter.End()
			return nil
		}
		bulkwriter.Flush()
	}
}
