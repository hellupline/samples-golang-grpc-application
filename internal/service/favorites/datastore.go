package favorites

import (
	"encoding/json"
	"time"

	"github.com/hellupline/samples-golang-grpc-server/storage"
)

const namespace = "favorites"

type Datastore struct {
	storage *storage.Storage
}

func NewDatastore(s *storage.Storage) *Datastore {
	return &Datastore{storage: s}
}

func (d *Datastore) FavoriteSearch() ([]*Favorite, error) {
	objs := make([]*Favorite, 0)
	if err := d.storage.Scan(namespace, "", func(v []byte) error {
		var obj *Favorite
		if err := json.Unmarshal(v, &obj); err != nil {
			return err
		}
		objs = append(objs, obj)
		return nil
	}); err != nil {
		return nil, err
	}
	return objs, nil
}

func (d *Datastore) FavoriteCreate(o Favorite) (*Favorite, error) {
	now := time.Now().UTC()
	obj := o
	obj.CreatedAt = now
	obj.UpdatedAt = now
	obj.Count = 0
	if err := d.storage.Put(namespace, obj.Name, func() ([]byte, error) {
		return json.Marshal(obj)
	}); err != nil {
		return nil, err
	}
	return &obj, nil
}

func (d *Datastore) FavoriteRead(key string) (*Favorite, error) {
	var obj Favorite
	if err := d.storage.Get(namespace, key, func(v []byte) error {
		return json.Unmarshal(v, &obj)
	}); err != nil {
		return nil, err
	}
	return &obj, nil
}

func (d *Datastore) FavoriteUpdate(key string, o Favorite) (*Favorite, error) {
	current, err := d.FavoriteRead(key)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	obj := o
	obj.Name = current.Name
	obj.CreatedAt = current.CreatedAt
	obj.UpdatedAt = now
	obj.Count = current.Count + 1
	if err := d.storage.Put(namespace, key, func() ([]byte, error) {
		return json.Marshal(obj)
	}); err != nil {
		return nil, err
	}
	return &obj, nil
}

func (d *Datastore) FavoriteDelete(key string) error {
	return d.storage.Delete(namespace, key)
}

type Favorite struct {
	Name string

	CreatedAt time.Time
	UpdatedAt time.Time

	Count uint
}
