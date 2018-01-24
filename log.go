package libtsm

//#include "adapter.h"
import "C"

import (
	"sync"
	"unsafe"
)

type LogReceiver interface {
	Log(file string, line int, fn, subs string, sev uint, format string, a ...interface{})
}

type logReceiverTicket uint64

var (
	logMu          sync.RWMutex
	logReceiverMap = make(map[logReceiverTicket]LogReceiver)
	logNextTicket  logReceiverTicket
)

//export logReceiverDispatch
func logReceiverDispatch(ticketPtr unsafe.Pointer, fileCstr *C.char, line C.int, fnCstr, subsCstr *C.char, sev C.unsigned, fmtCstr *C.char, vaList unsafe.Pointer) {
	ticket := logReceiverTicket(*(*uint64)(ticketPtr))
	receiver := getLogReceiver(ticket)

	file, fn, subs, rawFmt := C.GoString(fileCstr), C.GoString(fnCstr), C.GoString(subsCstr), C.GoString(fmtCstr)

	fmt, a := parseVaList(rawFmt, vaList)

	receiver.Log(file, int(line), fn, subs, uint(sev), fmt, a...)
}

func insertLogReceiver(r LogReceiver) logReceiverTicket {
	logMu.Lock()
	defer logMu.Unlock()

	var ticket logReceiverTicket
	for ok := true; ok; _, ok = logReceiverMap[ticket] {
		ticket = logNextTicket
		logNextTicket += 1
	}

	logReceiverMap[ticket] = r
	return ticket
}

func getLogReceiver(ticket logReceiverTicket) LogReceiver {
	logMu.RLock()
	defer logMu.RUnlock()
	return logReceiverMap[ticket]
}

func deleteLogReceiver(ticket logReceiverTicket) {
	logMu.Lock()
	defer logMu.Unlock()
	delete(logReceiverMap, ticket)
}

func parseVaList(rawFmt string, vaList unsafe.Pointer) (string, []interface{}) {
	type State int
	const (
		_ State = iota
		SCAN
		READ
		ESCAPE
	)

	var cur State = SCAN
	ret := []interface{}{}
	fmt := []byte(rawFmt)
	for i, c := range fmt {
		switch cur {
		case SCAN:
			switch c {
			case '\\':
				cur = ESCAPE
			case '%':
				cur = READ
			default:
				cur = SCAN
			}
		case READ:
			ret = append(ret, extractVaType(fmt[i:i+1], vaList))
		case ESCAPE:
			cur = SCAN
		}
	}

	return string(fmt), ret
}

func extractVaType(c []byte, vaList unsafe.Pointer) interface{} {
	switch c[0] {
	case 'c':
		r := rune(C.go_va_list_extract_char(vaList))
		return &r
	case 'x':
		fallthrough
	case 'u':
		r := C.go_va_list_extract_uint(vaList)
		return &r
	case 's':
		r := C.GoString(C.go_va_list_extract_cstr(vaList))
		return r
	case 'i':
		// not a golang fmt specifier
		c[0] = 'd'
		fallthrough
	case 'd':
		r := C.go_va_list_extract_int(vaList)
		return &r
	default:
		panic("invalid format specifier")
	}
	return nil
}
