workspace(name = "jonnrb_go_libtsm")

load("//bazel:deps.bzl", "jonnrb_go_libtsm_dependencies")

jonnrb_go_libtsm_dependencies()

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_rules_dependencies",
    "go_register_toolchains",
)

go_rules_dependencies()

go_register_toolchains()
