package libtsm

//#include <libtsm.h>
import "C"

import (
	"sync"
	"unsafe"
)

type ScreenDrawCallback func(id uint32, data []byte, width, posx, posy uint, attr ScreenAttr, age Age) int

type screenDrawCbTicket uint64

var (
	screenDrawMu         sync.Mutex
	screenDrawMap        = make(map[screenDrawCbTicket]ScreenDrawCallback)
	screenDrawNextTicket screenDrawCbTicket
)

//export screenDrawCbDispatch
func screenDrawCbDispatch(con *C.struct_tsm_screen, id C.uint32_t, ch *C.uint32_t, len C.size_t, width, posx, posy C.uint, attr *C.struct_tsm_screen_attr, age C.tsm_age_t, data unsafe.Pointer) C.int {
	ticket := screenDrawCbTicket(*(*uint64)(data))
	cb := getScreenDrawCb(ticket)
	cb(uint32(id), C.GoBytes(unsafe.Pointer(ch), C.int(len)), uint(width), uint(posx), uint(posy), newScreenAttr(attr), Age(age))
	return 0
}

func insertScreenDrawCb(cb ScreenDrawCallback) screenDrawCbTicket {
	screenDrawMu.Lock()
	defer screenDrawMu.Unlock()

	var ticket screenDrawCbTicket
	for ok := true; ok; _, ok = screenDrawMap[ticket] {
		ticket = screenDrawNextTicket
		screenDrawNextTicket += 1
	}

	screenDrawMap[ticket] = cb
	return ticket
}

func getScreenDrawCb(ticket screenDrawCbTicket) ScreenDrawCallback {
	screenDrawMu.Lock()
	defer screenDrawMu.Unlock()

	cb := screenDrawMap[ticket]
	delete(screenDrawMap, ticket)
	return cb
}
