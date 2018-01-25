load("//bazel:libtsm.bzl", "LIBTSM_BUILD_FILE")

def jonnrb_go_libtsm_dependencies():
  _maybe(native.git_repository,
      name = "io_bazel_rules_go",
      commit = "a390e7f7eac912f6e67dc54acf67aa974d05f9c3",
      remote = "https://github.com/bazelbuild/rules_go.git",
  )

  _maybe(native.new_git_repository,
      name = "jonnrb_bazel_libtsm",
      build_file_content = LIBTSM_BUILD_FILE,
      commit = "e8b2001cbaf9ba7176f955e1ee4bd9e75e2a6703",
      remote = "https://github.com/Aetf/libtsm.git",
  )

def _maybe(repo_rule, name, **kwargs):
  if name not in native.existing_rules():
    repo_rule(name=name, **kwargs)
