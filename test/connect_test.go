package test

import (
	"context"
	"log"
	"testing"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func TestConnect(t *testing.T) {
	endpoint := "172.20.165.191:9000"
	accessKeyID := "hraItU9c4j3kmovXwDtI"
	secretAccessKey := "lz5okY4duYCuKAAwmAJ8eQetMF7ATUc99mfNbscM"

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})

	if err != nil {
		t.Fatal(err)
	}

	// Create a context for the operation
	ctx := context.Background()

	// List buckets to test the connection
	buckets, err := minioClient.ListBuckets(ctx)
	if err != nil {
		t.Fatal("Failed to list buckets:", err)
	}

	for _, bucket := range buckets {
		log.Printf("Bucket: %s\n", bucket.Name)
	}

	log.Println("Connection to MinIO is successful")
}
