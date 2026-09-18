package addresses

import (
	"context"
	"errors"
	"testing"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/willie68/go-micro/pkg/pmodel"
)

type fakeStorage struct {
	items map[string]pmodel.Address
	err   error
}

func (f *fakeStorage) Addresses() ([]pmodel.Address, error) {
	if f.err != nil {
		return nil, f.err
	}
	out := make([]pmodel.Address, 0, len(f.items))
	for _, v := range f.items {
		out = append(out, v)
	}
	return out, nil
}

func (f *fakeStorage) Has(id string) bool {
	_, ok := f.items[id]
	return ok
}

func (f *fakeStorage) Read(id string) (*pmodel.Address, error) {
	if f.err != nil {
		return nil, f.err
	}
	adr, ok := f.items[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return &adr, nil
}

func (f *fakeStorage) Create(adr pmodel.Address) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	if adr.ID == "" {
		adr.ID = "new-id"
	}
	f.items[adr.ID] = adr
	return adr.ID, nil
}

func (f *fakeStorage) Update(adr pmodel.Address) error {
	if f.err != nil {
		return f.err
	}
	f.items[adr.ID] = adr
	return nil
}

func (f *fakeStorage) Delete(id string) error {
	if f.err != nil {
		return f.err
	}
	delete(f.items, id)
	return nil
}

func TestAddressesCRUD(t *testing.T) {
	ast := assert.New(t)
	stg := &fakeStorage{items: map[string]pmodel.Address{
		"1": {ID: "1", Name: "Doe"},
	}}
	svc := NewAddresses(stg)
	ctx := context.Background()

	list, err := svc.Addresses(ctx)
	ast.NoError(err)
	ast.Len(list, 1)

	ast.True(svc.Has(ctx, "1"))
	ast.False(svc.Has(ctx, "missing"))

	got, err := svc.Read(ctx, "1")
	ast.NoError(err)
	ast.Equal("Doe", got.Name)

	id, err := svc.Create(ctx, pmodel.Address{Name: "Smith"})
	ast.NoError(err)
	ast.Equal("new-id", id)

	err = svc.Update(ctx, pmodel.Address{ID: "1", Name: "Updated"})
	ast.NoError(err)
	got, err = svc.Read(ctx, "1")
	ast.NoError(err)
	ast.Equal("Updated", got.Name)

	err = svc.Delete(ctx, "1")
	ast.NoError(err)
	ast.False(svc.Has(ctx, "1"))
}

func TestProvide(t *testing.T) {
	ast := assert.New(t)
	inj := do.New()
	do.ProvideValue(inj, &fakeStorage{items: map[string]pmodel.Address{}})
	ast.NoError(Provide(inj))
	svc := do.MustInvoke[*Addresses](inj)
	ast.NotNil(svc)
}
