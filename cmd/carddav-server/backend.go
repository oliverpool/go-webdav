package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/emersion/go-vcard"
	"github.com/emersion/go-webdav"
	"github.com/emersion/go-webdav/carddav"
)

var _ carddav.Backend = &bwbBackend{}

type bwbBackend struct {
	currentUserPrincipal   string // must begin and end with a slash
	addressBookHomeSetPath string // a slash will be appended if needed
	addresses              []carddav.AddressObject
}

func (b *bwbBackend) CurrentUserPrincipal(ctx context.Context) (string, error) {
	return b.currentUserPrincipal, nil
}

func (b *bwbBackend) AddressBookHomeSetPath(ctx context.Context) (string, error) {
	upPath, err := b.CurrentUserPrincipal(ctx)
	if err != nil {
		return "", err
	}
	return path.Join(upPath, b.addressBookHomeSetPath) + "/", nil
}

func (b *bwbBackend) hasPermission(ctx context.Context, path string) error {
	ab, err := b.AddressBookHomeSetPath(ctx)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(path, ab) {
		return webdav.NewHTTPError(http.StatusForbidden, errors.New("carddav: outside of home set path"))
	}
	return nil
}

func (b *bwbBackend) ListAddressBooks(ctx context.Context) ([]carddav.AddressBook, error) {
	log.Println("ListAddressBooks")
	defer log.Println("> ListAddressBooks")

	upPath, err := b.AddressBookHomeSetPath(ctx)
	if err != nil {
		return nil, err
	}

	return []carddav.AddressBook{
		{
			Path:                 upPath + "default/",
			Name:                 "My contacts",
			Description:          "Default address book",
			MaxResourceSize:      1024,
			SupportedAddressData: nil,
		},
	}, nil
}

func (b *bwbBackend) GetAddressBook(ctx context.Context, path string) (*carddav.AddressBook, error) {
	log.Println("GetAddressBook")
	defer log.Println("> GetAddressBook")

	if err := b.hasPermission(ctx, path); err != nil {
		return nil, err
	}

	log.Println(path)
	abs, err := b.ListAddressBooks(ctx)
	if err != nil {
		return nil, err
	}
	for _, ab := range abs {
		if ab.Path == path {
			return &ab, nil
		}
	}
	return nil, webdav.NewHTTPError(404, fmt.Errorf("Not found"))
}

func (b *bwbBackend) CreateAddressBook(ctx context.Context, ab *carddav.AddressBook) error {
	return webdav.NewHTTPError(http.StatusNotImplemented, errors.New("carddav: create address book not supported"))
}

func (*bwbBackend) DeleteAddressBook(ctx context.Context, path string) error {
	return webdav.NewHTTPError(http.StatusNotImplemented, errors.New("carddav: delete address book not supported"))
}

func (b *bwbBackend) GetAddressObject(ctx context.Context, path string, req *carddav.AddressDataRequest) (*carddav.AddressObject, error) {
	log.Println("GetAddressObject")

	if err := b.hasPermission(ctx, path); err != nil {
		return nil, err
	}
	for _, ao := range b.addresses {
		if ao.Path == path {
			return &ao, nil
		}
	}

	log.Println("### Not found", path)
	return nil, webdav.NewHTTPError(404, fmt.Errorf("Not found"))
}

func (b *bwbBackend) ListAddressObjects(ctx context.Context, path string, req *carddav.AddressDataRequest) ([]carddav.AddressObject, error) {
	if err := b.hasPermission(ctx, path); err != nil {
		return nil, err
	}
	aos := make([]carddav.AddressObject, 0, len(b.addresses))
	for _, ao := range b.addresses {
		if strings.HasPrefix(ao.Path, path) {
			aos = append(aos, ao)
		}
	}
	return aos, nil
}

func (*bwbBackend) QueryAddressObjects(ctx context.Context, path string, query *carddav.AddressBookQuery) ([]carddav.AddressObject, error) {
	log.Println("QueryAddressObjects")
	defer log.Println("> QueryAddressObjects")
	panic("TODO: implement")
}

func (*bwbBackend) PutAddressObject(ctx context.Context, path string, card vcard.Card, opts *carddav.PutAddressObjectOptions) (*carddav.AddressObject, error) {
	// 403 works on iOS (update gets reverted shortly after the update)
	return nil, webdav.NewHTTPError(http.StatusForbidden, errors.New("carddav: create address book not supported"))
}

func (*bwbBackend) DeleteAddressObject(ctx context.Context, path string) error {
	// 403 works on iOS (deletion gets reverted shortly after the update)
	return webdav.NewHTTPError(http.StatusForbidden, errors.New("carddav: create address book not supported"))
}
