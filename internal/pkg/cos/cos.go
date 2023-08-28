/*******************************************************************************
 * Copyright 2017 Dell Inc.
 * Copyright (c) 2019 Intel Corporation
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/
package cos

import (
	"context"
	"errors"
	"fmt"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Cos struct {
	client *cos.Client
}

func NewCos(uri, ak, sk string) *Cos {
	u, _ := url.Parse(uri)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  ak,
			SecretKey: sk,
		},
	})
	return &Cos{
		client: client,
	}
}

func (c *Cos) Get(name string) ([]byte, error) {
	if c.client == nil {
		return []byte{}, errors.New("cos client is null")
	}
	resp, err := c.client.Object.Get(context.Background(), name, nil)
	if err != nil {
		return []byte{}, fmt.Errorf("get cos object err:%s", err.Error())
	}
	bs, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return bs, nil
}

func (c *Cos) DownloadFiled(name, filepath string) error {
	if c.client == nil {
		return errors.New("cos client is null")
	}
	_, err := c.client.Object.Download(context.Background(), name, filepath, nil)
	return err
}
