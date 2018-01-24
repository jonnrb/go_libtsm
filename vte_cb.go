package libtsm

//#include "adapter.h"
import "C"

import (
	"sync"
	"unsafe"
)

type VTEWriteCallback func([]byte)

type vteWriteTicket uint64

var (
	vteWriteMu         sync.RWMutex
	vteWriteCbMap      = make(map[vteWriteTicket]VTEWriteCallback)
	vteWriteNextTicket vteWriteTicket
)

//export vteWriteCbDispatch
func vteWriteCbDispatch(_ *C.struct_tsm_vte, u8 *C.char, len C.size_t, ticketPtr unsafe.Pointer) {
	ticket := vteWriteTicket(*(*uint64)(ticketPtr))
	cb := getVTEWriteCallback(ticket)
	cb(C.GoBytes(unsafe.Pointer(u8), C.int(len)))
}

func insertVTEWriteCallback(cb VTEWriteCallback) vteWriteTicket {
	vteWriteMu.Lock()
	defer vteWriteMu.Unlock()

	var ticket vteWriteTicket
	for ok := true; ok; _, ok = vteWriteCbMap[ticket] {
		ticket = vteWriteNextTicket
		vteWriteNextTicket += 1
	}

	vteWriteCbMap[ticket] = cb
	return ticket
}

func getVTEWriteCallback(ticket vteWriteTicket) VTEWriteCallback {
	vteWriteMu.Lock()
	defer vteWriteMu.Unlock()
	return vteWriteCbMap[ticket]
}

func deleteVTEWriteCallback(ticket vteWriteTicket) {
	vteWriteMu.Lock()
	defer vteWriteMu.Unlock()
	delete(vteWriteCbMap, ticket)
}
