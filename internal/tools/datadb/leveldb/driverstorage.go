package leveldb

import (
	"github.com/winc-link/hummingbird/internal/pkg/errort"
	"github.com/winc-link/hummingbird/internal/pkg/logger"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	//"gitlab.com/tedge/edgex/internal/pkg/errort"
	//"gitlab.com/tedge/edgex/internal/pkg/logger"
)

type DriverStorageClient struct {
	client        *leveldb.DB
	loggingClient logger.LoggingClient
}

func NewDriverStorageClient(filePath, subPath string, lc logger.LoggingClient) (*DriverStorageClient, error) {
	_, fileErr := os.Stat(filePath)
	if fileErr != nil || !os.IsExist(fileErr) {
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	var (
		err    error
		client *leveldb.DB
	)

	if client, err = leveldb.OpenFile(filePath+subPath, nil); err != nil {
		lc.Errorf("leveldb openFile error: %s", err)
		return nil, err
	}
	lc.Infof("driver storage create path: %s, driver service id: %s", filePath+subPath, subPath)
	return &DriverStorageClient{
		client:        client,
		loggingClient: lc,
	}, nil
}

func (dsc *DriverStorageClient) CloseSession() {
	dsc.client.Close()
}

func (dsc *DriverStorageClient) All() (map[string][]byte, error) {
	kvs := make(map[string][]byte)
	iter := dsc.client.NewIterator(nil, &opt.ReadOptions{
		DontFillCache: true,
	})
	for iter.Next() {
		value, _ := dsc.client.Get(iter.Key(), nil)
		kvs[string(iter.Key())] = value
	}
	iter.Release()
	return kvs, iter.Error()
}

func (dsc *DriverStorageClient) Get(keys []string) (map[string][]byte, error) {
	kvs := make(map[string][]byte, len(keys))
	for _, key := range keys {
		value, err := dsc.client.Get([]byte(key), nil)
		if err != nil {
			dsc.loggingClient.Error("get value with key(%s) error: %s", keys, err)
			kvs[key] = []byte("")
			continue
		}
		kvs[key] = value
	}
	return kvs, nil
}

func (dsc *DriverStorageClient) Put(kvs map[string][]byte) error {
	batch := new(leveldb.Batch)
	defer batch.Reset()

	for k, v := range kvs {
		batch.Put([]byte(k), v)
	}
	if err := dsc.client.Write(batch, &opt.WriteOptions{
		NoWriteMerge: true,
		Sync:         true,
	}); err != nil {
		return errort.NewCommonEdgeX(errort.KindDatabaseError, "batch transaction write", err)
	}
	return nil
}

func (dsc *DriverStorageClient) Delete(keys []string) error {
	batch := new(leveldb.Batch)
	defer batch.Reset()
	for _, key := range keys {
		batch.Delete([]byte(key))
	}
	if err := dsc.client.Write(batch, nil); err != nil {
		return errort.NewCommonEdgeX(errort.KindDatabaseError, "batch transaction delete failed", err)
	}
	return nil
}
