// Copyright (c)  WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
//
// WSO2 Inc. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package utils

import (
	"fmt"
	"github.com/go-logr/logr"
	registryclient "github.com/heroku/docker-registry-client/registry"
	"strings"
)

type RegAuth struct {
	RegistryUrl string
	Username    string
	Password    string
}

func IsImageExists(auth RegAuth, image string, tag string, logger logr.Logger) (bool, error) {
	hub, err := registryclient.New(auth.RegistryUrl, auth.Username, auth.Password)
	if err != nil {
		logger.Error(err, "Error connecting to the docker registry", "registry-url", auth.RegistryUrl)
		return false, err
	}

	// remove registry name if exists in the image name
	imageWithoutReg := image
	splits := strings.Split(image, "/")
	if len(splits) == 3 {
		imageWithoutReg = fmt.Sprintf("%s/%s", splits[1], splits[2])
	}

	tags, err := hub.Tags(imageWithoutReg)
	if err != nil {
		logger.Error(err, "Error getting tags from the image in the docker registry", "registry-url", auth.RegistryUrl, "image", image)
		return false, err
	}
	for _, foundTag := range tags {
		if foundTag == tag {
			logger.Info("Found the image tag from the registry", "image", imageWithoutReg, "tag", foundTag)
			return true, nil
		}
	}
	return false, nil
}
