package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type BlobStorage struct {
	endpoint       *url.URL
	shared_key     *azblob.SharedKeyCredential
	container_name string
}

func NewBlobStorageConn(container_url, account_name, key, container string) (*BlobStorage, error) {
	endp, err := url.Parse(container_url)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse url: %q error: %w", container_url, err)
	}
	skc, err := azblob.NewSharedKeyCredential(account_name, key)
	if err != nil {
		return nil, fmt.Errorf("couldn't create new shared key cred: %w", err)
	}
	return &BlobStorage{
		endpoint:       endp,
		shared_key:     skc,
		container_name: container,
	}, nil
}

// downloads filename
func (b *BlobStorage) Download(filename string, content io.Writer) (err error) {
	var (
		client *azblob.Client
	)
	if client, err = azblob.NewClientWithSharedKeyCredential(b.endpoint.String(), b.shared_key, nil); err != nil {
		err = fmt.Errorf("get client error: %w", err)
		return
	}
	rsp, err := client.DownloadStream(
		context.Background(),
		b.container_name, filename,
		nil,
	)
	if err != nil {
		err = fmt.Errorf("get response error: %w", err)
		return
	}
	defer rsp.Body.Close()
	_, err = io.Copy(content, rsp.Body)
	return nil
}

// uploads local file to azure by name
func (b *BlobStorage) UploadFile(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}
	return b.Upload(name, f)
}

// uploads output file to azure storage
func (b *BlobStorage) Upload(name string, content io.Reader) (err error) {
	var (
		client *azblob.Client
	)
	defer func() {
		if err != nil {
			err = fmt.Errorf("upload %q error: %w", b.endpoint.JoinPath(name), err)
		}
	}()
	if client, err = azblob.NewClientWithSharedKeyCredential(b.endpoint.String(), b.shared_key, nil); err != nil {
		err = fmt.Errorf("get client error: %w", err)
		return
	}
	_, err = client.UploadStream(
		context.Background(),
		b.container_name,
		name,
		content,
		nil,
		//&azblob.UploadStreamOptions{
		//	HTTPHeaders: &blob.HTTPHeaders{
		//		BlobContentType: ptr("application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"),
		//	},
		//},
	)
	return
}

//func ptr(x string) *string {
//	return &x
//}
