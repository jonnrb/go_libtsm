#include <libtsm.h>

extern void logReceiverDispatch(void*, char*, int, char*, char*, unsigned,
                                char*, void*);

extern int screenDrawCbDispatch(struct tsm_screen*, uint32_t, uint32_t*, size_t,
                                unsigned int, unsigned int, unsigned int,
                                struct tsm_screen_attr*, tsm_age_t, void*);

extern void vteWriteCbDispatch(struct tsm_vte* vte, char* u8, size_t len,
                               void* data);

int go_call_tsm_screen_new(struct tsm_screen** screen, uint64_t* log_cb_ticket);

int go_call_tsm_vte_new(struct tsm_vte** vte, struct tsm_screen* screen,
                        uint64_t* log_cb_ticket, uint64_t* write_cb_ticket);

char go_va_list_extract_char(void* arg);
unsigned go_va_list_extract_uint(void* arg);
char* go_va_list_extract_cstr(void* arg);
int go_va_list_extract_int(void* arg);

void go_screen_attr_get_bitfields(const struct tsm_screen_attr* attr,
                                  bool* bold, bool* underline, bool* inverse,
                                  bool* protect, bool* blink);

void go_screen_attr_set_bitfields(struct tsm_screen_attr* attr, bool bold,
                                  bool underline, bool inverse, bool protect,
                                  bool blink);

int go_call_screen_draw(struct tsm_screen* screen, uint64_t ticket);
