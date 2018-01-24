package libtsm

//#include "adapter.h"
import "C"

import (
	"syscall"
	"unsafe"
)

type VTE struct {
	vte *C.struct_tsm_vte

	hasLog     bool
	logTicket  logReceiverTicket
	logTicketC *C.uint64_t

	writeTicket  vteWriteTicket
	writeTicketC *C.uint64_t
}

type VTEPalette string

const (
	VTE_PALETTE_DEFAULT         VTEPalette = ""
	VTE_PALETTE_SOLARIZED                  = "solarized"
	VTE_PALETTE_SOLARIZED_BLACK            = "solarized-black"
	VTE_PALETTE_SOLARIZED_WHITE            = "solarized-white"
	VTE_PALETTE_SOFT_BLACK                 = "soft-black"
)

type KeyboardModifier uint

const (
	SHIFT   KeyboardModifier = C.TSM_SHIFT_MASK
	LOCK                     = C.TSM_LOCK_MASK
	CONTROL                  = C.TSM_CONTROL_MASK
	ALT                      = C.TSM_LOCK_MASK
	LOGO                     = C.TSM_LOGO_MASK
)

type KeyboardModifierSet KeyboardModifier

// `screen` must outlive the returned VTE
func NewVTE(screen *Screen, write VTEWriteCallback, log LogReceiver) (*VTE, error) {
	vte := VTE{hasLog: log != nil}

	if vte.hasLog {
		vte.logTicket = insertLogReceiver(log)
		vte.logTicketC = (*C.uint64_t)(C.malloc(C.sizeof_uint64_t))
		*vte.logTicketC = C.uint64_t(vte.logTicket)
	}

	vte.writeTicket = insertVTEWriteCallback(write)
	vte.writeTicketC = (*C.uint64_t)(C.malloc(C.sizeof_uint64_t))
	*vte.writeTicketC = C.uint64_t(vte.writeTicket)

	if code := C.go_call_tsm_vte_new(&vte.vte, screen.s, vte.logTicketC, vte.writeTicketC); code != 0 {
		return nil, syscall.Errno(-code)
	}

	return &vte, nil
}

func (vte *VTE) Close() {
	C.tsm_vte_unref(vte.vte)

	if vte.hasLog {
		deleteLogReceiver(vte.logTicket)
		C.free(unsafe.Pointer(vte.logTicketC))
	}

	deleteVTEWriteCallback(vte.writeTicket)
	C.free(unsafe.Pointer(vte.writeTicketC))
}

func (vte *VTE) SetPalette(palette VTEPalette) error {
	var paletteCStr *C.char
	if palette != "" {
		paletteCStr = C.CString(string(palette))
		defer C.free(unsafe.Pointer(paletteCStr))
	}
	if code := C.tsm_vte_set_palette(vte.vte, paletteCStr); code != 0 {
		return syscall.Errno(-code)
	}
	return nil
}

func (vte *VTE) GetDefaultAttr() ScreenAttr {
	a := &C.struct_tsm_screen_attr{}
	C.tsm_vte_get_def_attr(vte.vte, a)
	return newScreenAttr(a)
}

func (vte *VTE) Reset() {
	C.tsm_vte_reset(vte.vte)
}

func (vte *VTE) HardReset() {
	C.tsm_vte_hard_reset(vte.vte)
}

func (vte *VTE) Input(utf8 []byte) {
	cbytes := C.CBytes(utf8)
	defer C.free(cbytes)
	C.tsm_vte_input(vte.vte, (*C.char)(cbytes), C.size_t(len(utf8)))
}

func (vte *VTE) HandleKeyboard(keysym, ascii uint32, mods KeyboardModifierSet, unicode uint32) bool {
	return bool(C.tsm_vte_handle_keyboard(vte.vte, C.uint32_t(keysym), C.uint32_t(ascii), C.uint(mods), C.uint32_t(unicode)))
}
