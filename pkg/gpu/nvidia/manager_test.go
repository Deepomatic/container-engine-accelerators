// Copyright 2017 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package nvidia

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSharedNvidiaGPUManager(t *testing.T) {
	// Large prime number of duplicates devices, to check declared devices modulo
	const duplicationFactor = 97
	// Expects a valid GPUManager to be created.
	testGpuManager := NewSharedNvidiaGPUManager("/home/kubernetes/bin/nvidia", "/usr/local/nvidia", duplicationFactor)
	as := assert.New(t)
	as.NotNil(testGpuManager)

	testGpuManager.defaultDevices = []string{nvidiaCtlDevice, nvidiaUVMDevice, nvidiaUVMToolsDevice}
	// Tests discoverGPUs()
	if _, err := os.Stat(nvidiaCtlDevice); err == nil {
		err = testGpuManager.discoverGPUs()
		as.Nil(err)
		sharedGpus := reflect.ValueOf(testGpuManager).Elem().FieldByName("devices").Len()
		as.Equal(0, sharedGpus%duplicationFactor, fmt.Sprintf("The number of shared GPUs should be a multiple of the duplicationFactor: %v", duplicationFactor))

		for id := range testGpuManager.devices {
			as.FileExists(filepath.Join(devDirectory, testGpuManager.GetDeviceFilename(id)), "The real underlying device filepath should exist for each shared device.")
		}
	}
}
