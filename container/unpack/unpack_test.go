package unpack

import (
	"context"
	"testing"
)

func TestSmallGetImageManifestsURL(t *testing.T) {
	got, err := GetImageManifestsURL("busybox")
	expected := "https://registry-1.docker.io/v2/library/busybox/manifests/latest"
	if err != nil {
		t.Errorf("got err=%v", err)
	}
	if got != expected {
		t.Errorf("got=%s, want=%s", got, expected)
	}
}

func TestMediumUnpack(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	image := "gcr.io/google-containers/busybox:1.27"
	target := "./rootfs"
	// defer os.RemoveAll(target)
	if err := UnpackImage(ctx, image, target); err != nil {
		t.Errorf("err=%v", err)
	}
}

func TestMediumUnpackLatest(t *testing.T) {
	t.Skip("[known issue] this cannot download latest tag images")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	image := "gcr.io/google-containers/busybox"
	target := "./rootfs"
	// defer os.RemoveAll(target)
	if err := UnpackImage(ctx, image, target); err != nil {
		t.Errorf("err=%v", err)
	}
}

func TestMediumUnpackDockerhub(t *testing.T) {
	t.Skip("[known issue] no support for authentication causes failure to download images from dockerhub")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	image := "busybox"
	target := "./rootfs"
	// defer os.RemoveAll(target)
	if err := UnpackImage(ctx, image, target); err != nil {
		t.Errorf("err=%v", err)
	}
}
