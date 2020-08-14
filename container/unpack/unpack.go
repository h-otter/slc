package unpack

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/docker/distribution/manifest/schema2"
	"github.com/docker/docker/pkg/archive"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

const ImageManifestsURLFormat = "https://%s/v2/%s/%s/manifests/%s"
const ImageBlobURLFormat = "https://%s/v2/%s/%s/blobs/%s"

func GetImageManifestsURL(image string) (string, error) {
	registry := "registry-1.docker.io"
	library := "library"
	imageName := ""
	tag := "latest"

	t := strings.Split(image, ":")
	switch len(t) {
	case 1:
	case 2:
		tag = t[1]

	default:
		return "", errors.New("wrong format of image")
	}

	t = strings.Split(t[0], "/")
	switch len(t) {
	case 1:
		imageName = t[0]

	case 2:
		library = t[0]
		imageName = t[1]

	case 3:
		registry = t[0]
		library = t[1]
		imageName = t[2]

	default:
		return "", errors.New("wrong format of image")
	}

	return fmt.Sprintf(ImageManifestsURLFormat, registry, library, imageName, tag), nil
}

func GetImageBlobURL(image, layerDigest string) (string, error) {
	registry := "registry-1.docker.io"
	library := "library"
	imageName := ""

	t := strings.Split(image, ":")
	switch len(t) {
	case 1:
	case 2:

	default:
		return "", errors.New("wrong format of image")
	}

	t = strings.Split(t[0], "/")
	switch len(t) {
	case 1:
		imageName = t[0]

	case 2:
		library = t[0]
		imageName = t[1]

	case 3:
		registry = t[0]
		library = t[1]
		imageName = t[2]

	default:
		return "", errors.New("wrong format of image")
	}

	return fmt.Sprintf(ImageBlobURLFormat, registry, library, imageName, layerDigest), nil
}

func GetManifest(image string) (*schema2.Manifest, error) {
	manifestURL, err := GetImageManifestsURL(image)
	if err != nil {
		errors.Wrapf(err, "GetImageManifestsURL(%s)", image)
	}

	req, err := http.NewRequest("GET", manifestURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get manifest with error message, maybe it is caused by no support for authorization: %s", buf)
	}

	manifest := &schema2.Manifest{}
	if err := json.Unmarshal(buf, manifest); err != nil {
		return nil, errors.Wrap(err, "")
	}

	return manifest, nil
}

func GetBlob(image string, layerDigest string) (io.ReadCloser, error) {
	blobURL, err := GetImageBlobURL(image, layerDigest)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	resp, err := http.Get(blobURL)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	return resp.Body, nil
}

func UnpackImage(ctx context.Context, image string, target string) error {
	manifest, err := GetManifest(image)
	if err != nil {
		return errors.Wrap(err, "GetManifest()")
	}

	tmp, err := ioutil.TempDir("", "slc-initialization")
	if err != nil {
		return errors.Wrap(err, `ioutil.TempDir("", "slc-initialization")`)
	}
	defer os.RemoveAll(tmp)

	eg := &errgroup.Group{}
	blobs := sync.Map{}
	for _, layer := range manifest.Layers {
		layerDigest := layer.Digest.String()
		eg.Go(func() error {
			log.Printf("[INFO] downloading %s", layerDigest)

			b, err := GetBlob(image, layerDigest)
			if err != nil {
				return errors.Wrap(err, "GetBlob()")
			}
			blobs.Store(layerDigest, b)

			log.Printf("[INFO] downloaded %s", layerDigest)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return errors.Wrap(err, "failed downloads")
	}

	os.MkdirAll(target, 0755)
	for _, layer := range manifest.Layers {
		layerDigest := layer.Digest.String()
		log.Printf("[INFO] extracting %s", layerDigest)

		b, _ := blobs.Load(layerDigest)
		blob := b.(io.ReadCloser)
		defer blob.Close()

		if err := archive.Untar(blob, target, &archive.TarOptions{}); err != nil {
			return errors.Wrap(err, "archive.Untar()")
		}
	}

	return nil
}
