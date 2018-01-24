load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "adapter.c",
        "adapter.h",
        "log.go",
        "render.go",
        "screen.go",
        "vte.go",
        "vte_cb.go",
    ],
    cdeps = ["@jonnrb_bazel_libtsm//:libtsm"],
    cgo = True,
    importpath = "github.com/jonnrb/go_libtsm",
    visibility = ["//visibility:public"],
)
