package libtsm

//#include "adapter.h"
import "C"

import (
	"syscall"
	"unsafe"
)

type Age uint32

type ScreenFlag uint

type ScreenFlagSet ScreenFlag

const (
	INSERT_MODE ScreenFlag = C.TSM_SCREEN_INSERT_MODE
	AUTO_WRAP   ScreenFlag = C.TSM_SCREEN_AUTO_WRAP
	REL_ORIGIN  ScreenFlag = C.TSM_SCREEN_REL_ORIGIN
	INVERSE     ScreenFlag = C.TSM_SCREEN_INVERSE
	HIDE_CURSOR ScreenFlag = C.TSM_SCREEN_HIDE_CURSOR
	FIXED_POS   ScreenFlag = C.TSM_SCREEN_FIXED_POS
	ALTERNATE   ScreenFlag = C.TSM_SCREEN_ALTERNATE
)

type Symbol uint32

type Screen struct {
	s *C.struct_tsm_screen

	hasLog     bool
	logTicketC *C.uint64_t
	logTicket  logReceiverTicket
}

// set either `Code` or `R`, `G`, `B`
type ScreenColor struct {
	Code int
	R    byte
	G    byte
	B    byte
}

type ScreenAttr struct {
	Foreground ScreenColor
	Background ScreenColor

	Bold      bool
	Underline bool
	Inverse   bool
	Protect   bool
	Blink     bool
}

func NewScreen(log LogReceiver) (*Screen, error) {
	s := Screen{hasLog: log != nil}

	if s.hasLog {
		s.logTicket = insertLogReceiver(log)
		s.logTicketC = (*C.uint64_t)(C.malloc(C.sizeof_uint64_t))
		*s.logTicketC = C.uint64_t(s.logTicket)
	}

	if code := C.go_call_tsm_screen_new(&s.s, s.logTicketC); code != 0 {
		return nil, syscall.Errno(-code)
	}

	return &s, nil
}

func (s *Screen) Close() {
	C.tsm_screen_unref(s.s)

	if s.hasLog {
		deleteLogReceiver(s.logTicket)
		C.free(unsafe.Pointer(s.logTicketC))
	}
}

func (s *Screen) Width() uint {
	return uint(C.tsm_screen_get_width(s.s))
}

func (s *Screen) Height() uint {
	return uint(C.tsm_screen_get_height(s.s))
}

func (s *Screen) Resize(width, height uint) error {
	if code := C.tsm_screen_resize(s.s, C.uint(width), C.uint(height)); code != 0 {
		return syscall.Errno(-code)
	} else {
		return nil
	}
}

func (s *Screen) SetMargins(top, bottom uint) error {
	if code := C.tsm_screen_set_margins(s.s, C.uint(top), C.uint(bottom)); code != 0 {
		return syscall.Errno(-code)
	} else {
		return nil
	}
}

func (s *Screen) SbUp(num uint) {
	C.tsm_screen_sb_up(s.s, C.uint(num))
}

func (s *Screen) SbDown(num uint) {
	C.tsm_screen_sb_down(s.s, C.uint(num))
}

func (s *Screen) SbPageUp(num uint) {
	C.tsm_screen_sb_page_up(s.s, C.uint(num))
}

func (s *Screen) SbPageDown(num uint) {
	C.tsm_screen_sb_page_down(s.s, C.uint(num))
}

func (s *Screen) SbReset() {
	C.tsm_screen_sb_reset(s.s)
}

func (s *Screen) SetDefaultAttr(attr ScreenAttr) {
	C.tsm_screen_set_def_attr(s.s, attr.toCStruct())
}

func (s *Screen) Reset() {
	C.tsm_screen_reset(s.s)
}

func (s *Screen) SetFlags(flags ScreenFlagSet) {
	C.tsm_screen_set_flags(s.s, C.uint(flags))
}

// *removes* these flags
func (s *Screen) ResetFlags(flags ScreenFlagSet) {
	C.tsm_screen_reset_flags(s.s, C.uint(flags))
}

func (s *Screen) GetFlags() ScreenFlagSet {
	return ScreenFlagSet(C.tsm_screen_get_flags(s.s))
}

func (s *Screen) GetCursorX() uint {
	return uint(C.tsm_screen_get_cursor_x(s.s))
}

func (s *Screen) GetCursorY() uint {
	return uint(C.tsm_screen_get_cursor_y(s.s))
}

func (s *Screen) SetTabstop() {
	C.tsm_screen_set_tabstop(s.s)
}

func (s *Screen) ResetTabstop() {
	C.tsm_screen_reset_tabstop(s.s)
}

func (s *Screen) ResetAllTabstops() {
	C.tsm_screen_reset_all_tabstops(s.s)
}

func (s *Screen) Write(sym Symbol, attr ScreenAttr) {
	C.tsm_screen_write(s.s, C.tsm_symbol_t(sym), attr.toCStruct())
}

func (s *Screen) Newline() {
	C.tsm_screen_newline(s.s)
}

func (s *Screen) ScrollUp(num uint) {
	C.tsm_screen_scroll_up(s.s, C.uint(num))
}

func (s *Screen) ScrollDown(num int) {
	C.tsm_screen_scroll_down(s.s, C.uint(num))
}

func (s *Screen) MoveTo(x, y uint) {
	C.tsm_screen_move_to(s.s, C.uint(x), C.uint(y))
}

func (s *Screen) MoveUp(num uint, scroll bool) {
	C.tsm_screen_move_up(s.s, C.uint(num), C.bool(scroll))
}

func (s *Screen) MoveLeft(num uint) {
	C.tsm_screen_move_left(s.s, C.uint(num))
}

func (s *Screen) MoveRight(num uint) {
	C.tsm_screen_move_right(s.s, C.uint(num))
}

func (s *Screen) MoveLineEnd() {
	C.tsm_screen_move_line_end(s.s)
}

func (s *Screen) MoveLineHome() {
	C.tsm_screen_move_line_home(s.s)
}

func (s *Screen) TabRight(num uint) {
	C.tsm_screen_tab_right(s.s, C.uint(num))
}

func (s *Screen) TabLeft(num uint) {
	C.tsm_screen_tab_left(s.s, C.uint(num))
}

func (s *Screen) InsertLines(num uint) {
	C.tsm_screen_insert_lines(s.s, C.uint(num))
}

func (s *Screen) DeleteLines(num uint) {
	C.tsm_screen_delete_lines(s.s, C.uint(num))
}

func (s *Screen) InsertChars(num uint) {
	C.tsm_screen_insert_chars(s.s, C.uint(num))
}

func (s *Screen) DeleteChars(num uint) {
	C.tsm_screen_delete_chars(s.s, C.uint(num))
}

func (s *Screen) EraseCursor() {
	C.tsm_screen_erase_cursor(s.s)
}

func (s *Screen) EraseChars(num uint) {
	C.tsm_screen_erase_chars(s.s, C.uint(num))
}

func (s *Screen) EraseCursorToEnd(protect bool) {
	C.tsm_screen_erase_cursor_to_end(s.s, C.bool(protect))
}

func (s *Screen) EraseHomeToCursor(protect bool) {
	C.tsm_screen_erase_home_to_cursor(s.s, C.bool(protect))
}

func (s *Screen) EraseScreenToCursor(protect bool) {
	C.tsm_screen_erase_screen_to_cursor(s.s, C.bool(protect))
}

func (s *Screen) EraseCursorToScreen(protect bool) {
	C.tsm_screen_erase_cursor_to_screen(s.s, C.bool(protect))
}

func (s *Screen) EraseScreen(protect bool) {
	C.tsm_screen_erase_screen(s.s, C.bool(protect))
}

func (s *Screen) SelectionReset() {
	C.tsm_screen_selection_reset(s.s)
}

func (s *Screen) SelectionStart(posx, posy uint) {
	C.tsm_screen_selection_start(s.s, C.uint(posx), C.uint(posy))
}

func (s *Screen) SelectionTarget(posx, posy uint) {
	C.tsm_screen_selection_target(s.s, C.uint(posx), C.uint(posy))
}

func (s *Screen) SelectionCopy() (string, error) {
	var cstr *C.char
	ret := C.tsm_screen_selection_copy(s.s, &cstr)
	if ret < 0 {
		return "", syscall.Errno(-ret)
	}
	defer C.free(unsafe.Pointer(cstr))
	return C.GoStringN(cstr, ret), nil
}

func (s *Screen) Draw(cb ScreenDrawCallback) Age {
	ticket := insertScreenDrawCb(cb)
	age := C.go_call_screen_draw(s.s, C.uint64_t(ticket))
	return Age(age)
}

func newScreenAttr(a *C.struct_tsm_screen_attr) ScreenAttr {
	out := ScreenAttr{
		Foreground: ScreenColor{
			Code: int(a.fccode),
			R:    byte(a.fr),
			G:    byte(a.fg),
			B:    byte(a.fb),
		},
		Background: ScreenColor{
			Code: int(a.bccode),
			R:    byte(a.br),
			G:    byte(a.bg),
			B:    byte(a.bb),
		},
	}
	C.go_screen_attr_get_bitfields(a, (*C.bool)(&out.Bold), (*C.bool)(&out.Underline), (*C.bool)(&out.Inverse), (*C.bool)(&out.Protect), (*C.bool)(&out.Blink))
	return out
}

func (a ScreenAttr) toCStruct() *C.struct_tsm_screen_attr {
	out := &C.struct_tsm_screen_attr{
		fccode: C.int8_t(a.Foreground.Code),
		bccode: C.int8_t(a.Background.Code),
		fr:     C.uint8_t(a.Foreground.R),
		fg:     C.uint8_t(a.Foreground.G),
		fb:     C.uint8_t(a.Foreground.B),
		br:     C.uint8_t(a.Background.R),
		bg:     C.uint8_t(a.Background.G),
		bb:     C.uint8_t(a.Background.B),
	}
	C.go_screen_attr_set_bitfields(out, C.bool(a.Bold), C.bool(a.Underline), C.bool(a.Inverse), C.bool(a.Protect), C.bool(a.Blink))

	return out
}

func (s ScreenFlagSet) HasFlag(flag ScreenFlag) bool {
	return s&ScreenFlagSet(flag) != 0
}

func (s *ScreenFlagSet) AddFlag(flag ScreenFlag) {
	*s |= ScreenFlagSet(flag)
}

func (s *ScreenFlagSet) ClearFlag(flag ScreenFlag) {
	*s &= ^ScreenFlagSet(flag)
}

func (s *ScreenFlagSet) ToggleFlag(flag ScreenFlag) {
	*s ^= ScreenFlagSet(flag)
}
