/*
 * Copyright (c) 2021, NVIDIA CORPORATION.  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package nvpci

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/NVIDIA/go-nvlib/pkg/nvpci/bytes"
)

// MockNvpci mock pci device.
type MockNvpci struct {
	*nvpci
}

var _ Interface = (*MockNvpci)(nil)

// NewMockNvpci create new mock PCI and remove old devices.
func NewMockNvpci() (mock *MockNvpci, rerr error) {
	rootDir, err := os.MkdirTemp(os.TempDir(), "")
	if err != nil {
		return nil, err
	}
	defer func() {
		if rerr != nil {
			os.RemoveAll(rootDir)
		}
	}()

	mock = &MockNvpci{
		New(WithPCIDevicesRoot(rootDir)).(*nvpci),
	}

	return mock, nil
}

// Cleanup remove the mocked PCI devices root folder.
func (m *MockNvpci) Cleanup() {
	os.RemoveAll(m.pciDevicesRoot)
}

func validatePCIAddress(addr string) error {
	r := regexp.MustCompile(`0{4}:[0-9a-f]{2}:[0-9a-f]{2}\.[0-9]`)
	if !r.Match([]byte(addr)) {
		return fmt.Errorf(`invalid PCI address should match 0{4}:[0-9a-f]{2}:[0-9a-f]{2}\.[0-9]: %s`, addr)
	}

	return nil
}

// AddMockA100 Create an A100 like GPU mock device.
func (m *MockNvpci) AddMockA100(address string, numaNode int, sriov *SriovInfo) error {
	err := validatePCIAddress(address)
	if err != nil {
		return err
	}

	deviceDir := filepath.Join(m.pciDevicesRoot, address)
	err = os.MkdirAll(deviceDir, 0755)
	if err != nil {
		return err
	}

	err = createNVIDIAgpuFiles(deviceDir)
	if err != nil {
		return err
	}

	iommuGroup := 20
	_, err = os.Create(filepath.Join(deviceDir, strconv.Itoa(iommuGroup)))
	if err != nil {
		return err
	}
	err = os.Symlink(filepath.Join(deviceDir, strconv.Itoa(iommuGroup)), filepath.Join(deviceDir, "iommu_group"))
	if err != nil {
		return err
	}

	numa, err := os.Create(filepath.Join(deviceDir, "numa_node"))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(numa, "%v", numaNode)
	if err != nil {
		return err
	}

	if sriov != nil && sriov.PhysicalFunction != nil {
		totalVFs, err := os.Create(filepath.Join(deviceDir, "sriov_totalvfs"))
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(totalVFs, "%d", sriov.PhysicalFunction.TotalVFs)
		if err != nil {
			return err
		}

		numVFs, err := os.Create(filepath.Join(deviceDir, "sriov_numvfs"))
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(numVFs, "%d", sriov.PhysicalFunction.NumVFs)
		if err != nil {
			return err
		}
		for i := 1; i <= int(sriov.PhysicalFunction.NumVFs); i++ {
			err = m.createVf(address, i, iommuGroup, numaNode)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func createNVIDIAgpuFiles(deviceDir string) error {
	vendor, err := os.Create(filepath.Join(deviceDir, "vendor"))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(vendor, "0x%x", PCINvidiaVendorID)
	if err != nil {
		return err
	}

	class, err := os.Create(filepath.Join(deviceDir, "class"))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(class, "0x%x", PCI3dControllerClass)
	if err != nil {
		return err
	}

	device, err := os.Create(filepath.Join(deviceDir, "device"))
	if err != nil {
		return err
	}
	_, err = device.WriteString("0x20bf")
	if err != nil {
		return err
	}

	_, err = os.Create(filepath.Join(deviceDir, "nvidia"))
	if err != nil {
		return err
	}
	err = os.Symlink(filepath.Join(deviceDir, "nvidia"), filepath.Join(deviceDir, "driver"))
	if err != nil {
		return err
	}

	config, err := os.Create(filepath.Join(deviceDir, "config"))
	if err != nil {
		return err
	}
	_data := make([]byte, PCICfgSpaceStandardSize)
	data := bytes.New(&_data)
	data.Write16(0, PCINvidiaVendorID)
	data.Write16(2, uint16(0x20bf))
	data.Write8(PCIStatusBytePosition, PCIStatusCapabilityList)
	_, err = config.Write(*data.Raw())
	if err != nil {
		return err
	}

	bar0 := []uint64{0x00000000c2000000, 0x00000000c2ffffff, 0x0000000000040200}
	resource, err := os.Create(filepath.Join(deviceDir, "resource"))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(resource, "0x%x 0x%x 0x%x", bar0[0], bar0[1], bar0[2])
	if err != nil {
		return err
	}

	pmcID := uint32(0x170000a1)
	resource0, err := os.Create(filepath.Join(deviceDir, "resource0"))
	if err != nil {
		return err
	}
	_data = make([]byte, bar0[1]-bar0[0]+1)
	data = bytes.New(&_data).LittleEndian()
	data.Write32(0, pmcID)
	_, err = resource0.Write(*data.Raw())
	if err != nil {
		return err
	}

	return nil
}

func (m *MockNvpci) createVf(pfAddress string, id, iommu_group, numaNode int) error {
	functionID := pfAddress[len(pfAddress)-1]
	// we are verifying the last character of pfAddress is integer.
	functionNumber, err := strconv.Atoi(string(functionID))
	if err != nil {
		return fmt.Errorf("can't conver physical function pci address function number %s to integer: %v", string(functionID), err)
	}

	vfFunctionNumber := functionNumber + id
	vfAddress := pfAddress[:len(pfAddress)-1] + strconv.Itoa(vfFunctionNumber)

	deviceDir := filepath.Join(m.pciDevicesRoot, vfAddress)
	err = os.MkdirAll(deviceDir, 0755)
	if err != nil {
		return err
	}

	err = createNVIDIAgpuFiles(deviceDir)
	if err != nil {
		return err
	}

	vfIommuGroup := strconv.Itoa(iommu_group + id)

	_, err = os.Create(filepath.Join(deviceDir, vfIommuGroup))
	if err != nil {
		return err
	}
	err = os.Symlink(filepath.Join(deviceDir, vfIommuGroup), filepath.Join(deviceDir, "iommu_group"))
	if err != nil {
		return err
	}

	numa, err := os.Create(filepath.Join(deviceDir, "numa_node"))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(numa, "%v", numaNode)
	if err != nil {
		return err
	}

	err = os.Symlink(filepath.Join(m.pciDevicesRoot, pfAddress), filepath.Join(deviceDir, "physfn"))
	if err != nil {
		return err
	}

	return nil
}
