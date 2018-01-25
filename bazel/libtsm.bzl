LIBTSM_BUILD_FILE = """
package(default_visibility = ["//visibility:private"])

cc_library(
    name = "libtsm",
    hdrs = ["src/tsm/libtsm.h"],
    strip_include_prefix = "src/tsm",
    srcs = glob([
        "config.h",
        "src/tsm/libtsm-int.h",
        "src/tsm/*.c",
    ]),
    deps = [
        ":external",
        ":shared",
    ],
    visibility = ["//visibility:public"],
)

cc_library(
    name = "shared",
    hdrs = glob(["src/shared/*.h"]),
    srcs = [
        "src/shared/shl-htable.c",
        "src/shared/shl-ring.c",
    ],
    strip_include_prefix = "src/shared",
)

cc_library(
    name = "external",
    hdrs = [
        "external/wcwidth.h",
        "external/xkbcommon-keysyms.h",
    ],
    srcs = ["external/wcwidth.c"],
)

CONFIG_FILE = r\"\"\"
/* config.h.  Generated from config.h.in by configure.  */
/* config.h.in.  Generated from configure.ac by autoheader.  */

/* Enable debug mode */
#define BUILD_ENABLE_DEBUG 1

/* Have xkbcommon library */
/* #undef BUILD_HAVE_XKBCOMMON */

/* Define to 1 if you have the <dlfcn.h> header file. */
#define HAVE_DLFCN_H 1

/* Define to 1 if you have the <inttypes.h> header file. */
#define HAVE_INTTYPES_H 1

/* Define to 1 if you have the <memory.h> header file. */
#define HAVE_MEMORY_H 1

/* Define to 1 if you have the <stdint.h> header file. */
#define HAVE_STDINT_H 1

/* Define to 1 if you have the <stdlib.h> header file. */
#define HAVE_STDLIB_H 1

/* Define to 1 if you have the <strings.h> header file. */
#define HAVE_STRINGS_H 1

/* Define to 1 if you have the <string.h> header file. */
#define HAVE_STRING_H 1

/* Define to 1 if you have the <sys/stat.h> header file. */
#define HAVE_SYS_STAT_H 1

/* Define to 1 if you have the <sys/types.h> header file. */
#define HAVE_SYS_TYPES_H 1

/* Define to 1 if you have the <unistd.h> header file. */
#define HAVE_UNISTD_H 1

/* Define to the sub-directory where libtool stores uninstalled libraries. */
#define LT_OBJDIR ".libs/"

/* No Debug */
/* #undef NDEBUG */

/* Name of package */
#define PACKAGE "libtsm"

/* Define to the address where bug reports for this package should be sent. */
#define PACKAGE_BUGREPORT "http://bugs.freedesktop.org/enter_bug.cgi?product=kmscon"

/* Define to the full name of this package. */
#define PACKAGE_NAME "libtsm"

/* Define to the full name and version of this package. */
#define PACKAGE_STRING "libtsm 3"

/* Define to the one symbol short name of this package. */
#define PACKAGE_TARNAME "libtsm"

/* Define to the home page for this package. */
#define PACKAGE_URL "http://www.freedesktop.org/wiki/Software/libtsm"

/* Define to the version of this package. */
#define PACKAGE_VERSION "3"

/* Define to 1 if you have the ANSI C header files. */
#define STDC_HEADERS 1

/* Enable extensions on AIX 3, Interix.  */
#ifndef _ALL_SOURCE
# define _ALL_SOURCE 1
#endif
/* Enable GNU extensions on systems that have them.  */
#ifndef _GNU_SOURCE
# define _GNU_SOURCE 1
#endif
/* Enable threading extensions on Solaris.  */
#ifndef _POSIX_PTHREAD_SEMANTICS
# define _POSIX_PTHREAD_SEMANTICS 1
#endif
/* Enable extensions on HP NonStop.  */
#ifndef _TANDEM_SOURCE
# define _TANDEM_SOURCE 1
#endif
/* Enable general extensions on Solaris.  */
#ifndef __EXTENSIONS__
# define __EXTENSIONS__ 1
#endif


/* Version number of package */
#define VERSION "3"

/* Enable large inode numbers on Mac OS X 10.5.  */
#ifndef _DARWIN_USE_64_BIT_INODE
# define _DARWIN_USE_64_BIT_INODE 1
#endif

/* Number of bits in a file offset, on hosts where this is settable. */
/* #undef _FILE_OFFSET_BITS */

/* Define for large files, on AIX-style hosts. */
/* #undef _LARGE_FILES */

/* Define to 1 if on MINIX. */
/* #undef _MINIX */

/* Define to 2 if the system does not provide POSIX.1 features except with
   this defined. */
/* #undef _POSIX_1_SOURCE */

/* Define to 1 if you need to in order for 'stat' and other things to work. */
/* #undef _POSIX_SOURCE */
\"\"\"

genrule(
    name = "gen_config",
    outs = ["config.h"],
    cmd = "cat > $(location :config.h) << EOF" + CONFIG_FILE + "EOF",
)
"""
