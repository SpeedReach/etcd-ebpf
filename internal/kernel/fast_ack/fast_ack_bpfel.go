// Code generated by bpf2go; DO NOT EDIT.
//go:build 386 || amd64 || arm || arm64 || loong64 || mips64le || mipsle || ppc64le || riscv64

package fast_ack

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"

	"github.com/cilium/ebpf"
)

type fast_ackLogEntry struct {
	Term    uint32
	Command [64]int8
}

// loadFast_ack returns the embedded CollectionSpec for fast_ack.
func loadFast_ack() (*ebpf.CollectionSpec, error) {
	reader := bytes.NewReader(_Fast_ackBytes)
	spec, err := ebpf.LoadCollectionSpecFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("can't load fast_ack: %w", err)
	}

	return spec, err
}

// loadFast_ackObjects loads fast_ack and converts it into a struct.
//
// The following types are suitable as obj argument:
//
//	*fast_ackObjects
//	*fast_ackPrograms
//	*fast_ackMaps
//
// See ebpf.CollectionSpec.LoadAndAssign documentation for details.
func loadFast_ackObjects(obj interface{}, opts *ebpf.CollectionOptions) error {
	spec, err := loadFast_ack()
	if err != nil {
		return err
	}

	return spec.LoadAndAssign(obj, opts)
}

// fast_ackSpecs contains maps and programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type fast_ackSpecs struct {
	fast_ackProgramSpecs
	fast_ackMapSpecs
}

// fast_ackSpecs contains programs before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type fast_ackProgramSpecs struct {
	FastReturn *ebpf.ProgramSpec `ebpf:"fast_return"`
}

// fast_ackMapSpecs contains maps before they are loaded into the kernel.
//
// It can be passed ebpf.CollectionSpec.Assign.
type fast_ackMapSpecs struct {
	Logs       *ebpf.MapSpec `ebpf:"logs"`
	NewEntries *ebpf.MapSpec `ebpf:"new_entries"`
	Term       *ebpf.MapSpec `ebpf:"term"`
}

// fast_ackObjects contains all objects after they have been loaded into the kernel.
//
// It can be passed to loadFast_ackObjects or ebpf.CollectionSpec.LoadAndAssign.
type fast_ackObjects struct {
	fast_ackPrograms
	fast_ackMaps
}

func (o *fast_ackObjects) Close() error {
	return _Fast_ackClose(
		&o.fast_ackPrograms,
		&o.fast_ackMaps,
	)
}

// fast_ackMaps contains all maps after they have been loaded into the kernel.
//
// It can be passed to loadFast_ackObjects or ebpf.CollectionSpec.LoadAndAssign.
type fast_ackMaps struct {
	Logs       *ebpf.Map `ebpf:"logs"`
	NewEntries *ebpf.Map `ebpf:"new_entries"`
	Term       *ebpf.Map `ebpf:"term"`
}

func (m *fast_ackMaps) Close() error {
	return _Fast_ackClose(
		m.Logs,
		m.NewEntries,
		m.Term,
	)
}

// fast_ackPrograms contains all programs after they have been loaded into the kernel.
//
// It can be passed to loadFast_ackObjects or ebpf.CollectionSpec.LoadAndAssign.
type fast_ackPrograms struct {
	FastReturn *ebpf.Program `ebpf:"fast_return"`
}

func (p *fast_ackPrograms) Close() error {
	return _Fast_ackClose(
		p.FastReturn,
	)
}

func _Fast_ackClose(closers ...io.Closer) error {
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Do not access this directly.
//
//go:embed fast_ack_bpfel.o
var _Fast_ackBytes []byte