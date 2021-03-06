package actions

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/ipfs/go-datastore"
	"github.com/qri-io/dataset"
	"github.com/qri-io/dataset/dsfs"
	"github.com/qri-io/qri/repo"
)

// Select loads a dataset value specified by case.Sensitve.dot.separated.paths
func Select(r repo.Repo, ref repo.DatasetRef, path string) (interface{}, error) {
	ds, err := dsfs.LoadDataset(r.Store(), datastore.NewKey(ref.Path))
	if err != nil {
		return nil, err
	}

	if path == "" {
		return ds.Encode(), nil
	}

	v, err := pathValue(ds.Encode(), path)
	if err != nil {
		return nil, err
	}
	return v.Interface(), nil
}

// ApplyPath gets a dataset value by applying a case.Sensitve.dot.separated.path
func ApplyPath(ds *dataset.DatasetPod, path string) (interface{}, error) {
	var value reflect.Value
	value, err := pathValue(ds, path)
	if err != nil {
		return nil, err
	}
	return value.Interface(), nil
}

func pathValue(ds *dataset.DatasetPod, path string) (elem reflect.Value, err error) {
	elem = reflect.ValueOf(ds)

	for _, sel := range strings.Split(path, ".") {
		if elem.Kind() == reflect.Ptr {
			elem = elem.Elem()
		}

		switch elem.Kind() {
		case reflect.Struct:
			elem = elem.FieldByNameFunc(func(str string) bool {
				return strings.ToLower(str) == sel
			})
		case reflect.Slice:
			index, err := strconv.Atoi(sel)
			if err != nil {
				return elem, fmt.Errorf("invalid index value: %s", sel)
			}
			elem = elem.Index(index)
		case reflect.Map:
			for _, key := range elem.MapKeys() {
				// we only support strings as keys
				if strings.ToLower(key.String()) == sel {
					return elem.MapIndex(key), nil
				}
			}
			return elem, fmt.Errorf("invalid selection path: %s", path)
		}

		if elem.Kind() == reflect.Invalid {
			return elem, fmt.Errorf("invalid selection path: %s", path)
		}
	}

	return elem, nil
}
