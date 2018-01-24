#include "adapter.h"

int go_call_tsm_screen_new(struct tsm_screen** screen,
                           uint64_t* log_cb_ticket) {
  return tsm_screen_new(
      screen, log_cb_ticket != NULL ? (tsm_log_t)logReceiverDispatch : NULL,
      log_cb_ticket);
}

int go_call_tsm_vte_new(struct tsm_vte** vte, struct tsm_screen* screen,
                        uint64_t* log_cb_ticket, uint64_t* write_cb_ticket) {
  return tsm_vte_new(
      vte, screen, (tsm_vte_write_cb)vteWriteCbDispatch, write_cb_ticket,
      log_cb_ticket != NULL ? (tsm_log_t)logReceiverDispatch : NULL,
      log_cb_ticket);
}

char go_va_list_extract_char(void* arg) {
  va_list* ap = (va_list*)arg;
  return (char)va_arg(*ap, int);
}

unsigned go_va_list_extract_uint(void* arg) {
  va_list* ap = (va_list*)arg;
  return va_arg(*ap, unsigned);
}

char* go_va_list_extract_cstr(void* arg) {
  va_list* ap = (va_list*)arg;
  return va_arg(*ap, char*);
}

int go_va_list_extract_int(void* arg) {
  va_list* ap = (va_list*)arg;
  return va_arg(*ap, int);
}

void go_screen_attr_get_bitfields(const struct tsm_screen_attr* attr,
                                  bool* bold, bool* underline, bool* inverse,
                                  bool* protect, bool* blink) {
  *bold = attr->bold;
  *underline = attr->underline;
  *inverse = attr->inverse;
  *protect = attr->protect;
  *blink = attr->blink;
}

void go_screen_attr_set_bitfields(struct tsm_screen_attr* attr, bool bold,
                                  bool underline, bool inverse, bool protect,
                                  bool blink) {
  attr->bold = bold;
  attr->underline = underline;
  attr->inverse = inverse;
  attr->protect = protect;
  attr->blink = blink;
}

int go_call_screen_draw(struct tsm_screen* screen, uint64_t ticket) {
  return tsm_screen_draw(screen, (tsm_screen_draw_cb)screenDrawCbDispatch,
                         &ticket);
}
