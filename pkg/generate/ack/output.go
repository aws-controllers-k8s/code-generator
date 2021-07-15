// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package ack

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/ghodss/yaml"

	"github.com/aws-controllers-k8s/code-generator/pkg/version"
)

const (
	outputFileName = "ack-generate-metadata.yaml"
)

// UpdateReason is the reason a package got modified.
type UpdateReason string

const (
	// UpdateReasonAPIGeneration should be used when an API package
	// is modified by the APIs generator (ack-generate apis).
	UpdateReasonAPIGeneration UpdateReason = "API generation"

	// UpdateReasonConversionFunctionsGeneration Should be used when
	// an API package is modified by conversion functions generator.
	// TODO(hilalymh) ack-generate conversion-functions
	UpdateReasonConversionFunctionsGeneration UpdateReason = "Conversion functions generation"
)

// GenerationMetadata represents the parameters used to generate/update the
// API version directory.
//
// This type is public because soon it will be used by conversion generators
// to load APIs generation metadata.
// TODO(hilalymh) Add functions to load/edit metadata files.
type GenerationMetadata struct {
	// The APIs version e.g v1alpha2
	APIVersion string `json:"api_version"`
	// The checksum of all the combined files generated within the APIs directory
	APIDirectoryChecksum string `json:"api_directory_checksum"`
	// Last modification reason
	LastModification lastModificationInfo `json:"last_modification"`
	// AWS SDK Go version used generate the APIs
	AWSSDKGoVersion string `json:"aws_sdk_go_version"`
	// Informatiom about the ack-generate binary used to generate the APIs
	ACKGenerateInfo ackGenerateInfo `json:"ack_generate_info"`
	// Information about the generator config file used to generate the APIs
	GeneratorConfigInfo generatorConfigInfo `json:"generator_config_info"`
}

// ack-generate binary information
type ackGenerateInfo struct {
	Version   string `json:"version"`
	GoVersion string `json:"go_version"`
	BuildDate string `json:"build_date"`
	BuildHash string `json:"build_hash"`
}

// generator.yaml information
type generatorConfigInfo struct {
	OriginalFileName string `json:"original_file_name"`
	FileChecksum     string `json:"file_checksum"`
}

// last modification information
type lastModificationInfo struct {
	// UTC Timestamp
	Timestamp string `json:"timestamp"`
	// Modification reason
	Reason UpdateReason `json:"reason"`
}

// CreateGenerationMetadata gathers information about the generated code and save
// a yaml version in the API version directory
func CreateGenerationMetadata(
	apiVersion string,
	apisPath string,
	modificationReason UpdateReason,
	awsSDKGo string,
	generatorFileName string,
) error {
	filesDirectory := filepath.Join(apisPath, apiVersion)
	hash, err := hashDirectoryContent(filesDirectory)
	if err != nil {
		return err
	}

	generatorFileHash, err := hashFile(generatorFileName)
	if err != nil {
		return err
	}

	generationMetadata := &GenerationMetadata{
		APIVersion:           apiVersion,
		APIDirectoryChecksum: hash,
		LastModification: lastModificationInfo{
			Timestamp: time.Now().UTC().String(),
			Reason:    modificationReason,
		},
		AWSSDKGoVersion: awsSDKGo,
		ACKGenerateInfo: ackGenerateInfo{
			Version:   version.Version,
			BuildDate: version.BuildDate,
			BuildHash: version.BuildHash,
			GoVersion: version.GoVersion,
		},
		GeneratorConfigInfo: generatorConfigInfo{
			OriginalFileName: filepath.Base(generatorFileName),
			FileChecksum:     generatorFileHash,
		},
	}

	data, err := yaml.Marshal(generationMetadata)
	if err != nil {
		return err
	}

	outputFileName := filepath.Join(filesDirectory, outputFileName)
	err = ioutil.WriteFile(
		outputFileName,
		data,
		os.ModePerm,
	)
	if err != nil {
		return err
	}
	return nil
}

// LoadGenerationMetadata read the generation metadata for a given api version and
// apis path.
func LoadGenerationMetadata(apisPath, apiVersion string) (*GenerationMetadata, error) {
	filePath := filepath.Join(apisPath, apiVersion, outputFileName)
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var generationMetadata GenerationMetadata
	err = yaml.Unmarshal(b, &generationMetadata)
	if err != nil {
		return nil, err
	}
	return &generationMetadata, nil
}

// hashDirectoryContent returns the sha1 checksum of a given directory. It will walk
// the file tree of a directory and combine and the file contents before hashing it.
func hashDirectoryContent(directory string) (string, error) {
	h := sha1.New()
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// ignore yaml files (output.yaml and generator.yaml)
			fileExtension := filepath.Ext(info.Name())
			if fileExtension == ".yaml" {
				return nil
			}

			fileReader, err := os.Open(path)
			if err != nil {
				return err
			}
			_, err = io.Copy(h, fileReader)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	hash := hex.EncodeToString(h.Sum(nil))
	return hash, nil
}

// hashFile returns the sha1 hash of a given file
func hashFile(filename string) (string, error) {
	h := sha1.New()
	fileReader, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(h, fileReader)
	if err != nil {
		return "", err
	}
	hash := hex.EncodeToString(h.Sum(nil))
	return hash, nil
}
