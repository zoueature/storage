package s3

import (
	"context"
	"github.com/zoueature/config"
	"os"
	"testing"
)

func TestClient_Upload(t *testing.T) {
	cli := NewS3Storage(config.StorageConfig{
		AccessKey:    os.Getenv("ACCESS_KEY"),
		AccessSecret: os.Getenv("ACCESS_SECRET"),
		Bucket:       os.Getenv("BUCKET"),
		Region:       os.Getenv("REGION"),
		Domain:       os.Getenv("DOMAIN"),
	})
	ctx := context.Background()
	f, _ := os.Open("1.png")
	s, err := cli.Upload(ctx, f)
	if err != nil {
		t.Fatal(err)
	}
	println(s)
	url, err := cli.SignAccessURL(ctx, "0f5d04c818080a18946a39a541eaa0ca")
	if err != nil {
		t.Fatal(err)
	}
	println(url)
}
