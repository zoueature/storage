package s3

import (
	"context"
	"github.com/zoueature/config"
	"os"
	"testing"
)

func TestClient_Upload(t *testing.T) {
	cli, _ := NewS3Storage(config.StorageConfig{
		AccessKey:    "AKIAVWMHJF5WBBBZPUFV",
		AccessSecret: "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",
		Bucket:       "omi-shorts",
		Region:       "ap-southeast-1",
		Domain:       "https://omi-shorts.s3.ap-southeast-1.amazonaws.com",
		CDN:          "https://dr3xwe1z6y1px.cloudfront.net",
	})
	ctx := context.Background()
	f, _ := os.Open("1.png")
	s, err := cli.Upload(ctx, f)
	if err != nil {
		t.Fatal(err)
	}
	println(s)
	url, err := cli.SignAccessURlTryCDN(ctx, "1351e3308f3b5325ceecfa2e0c538370.jpg")
	if err != nil {
		t.Fatal(err)
	}
	println(url)
}
