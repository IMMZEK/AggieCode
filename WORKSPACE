# WORKSPACE

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# === Node.js Rules ===
# Fetches the Node.js rules repository
http_archive(
    name = "rules_nodejs",
    sha256 = "b9a2d85a187199d46e2bd73726977a1b06901a982783f7c6f135c664a8618e50", # Use a known SHA256 for security
    urls = ["https://github.com/bazelbuild/rules_nodejs/releases/download/6.0.0/rules_nodejs-6.0.0.tar.gz"], # Check for the latest version
)

# Load the Node.js rules setup function
load("@rules_nodejs//nodejs:repositories.bzl", "nodejs_register_toolchains")

# Register Node.js toolchains (e.g., for Node.js version 18)
# You might adjust the version based on your project needs
nodejs_register_toolchains(
    name = "nodejs",
    node_version = "18.18.0", # Specify the desired Node.js version
)

# === Go Rules ===
# Fetches the Go rules repository
http_archive(
    name = "io_bazel_rules_go",
    sha256 = "b5c016429d828b0659170c77e1da03a6010f1b1f1728e49f8c960939dfafe678", # Use a known SHA256
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.46.0/rules_go-v0.46.0.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.46.0/rules_go-v0.46.0.tar.gz", # Check for the latest version
    ],
)

# Fetches Gazelle, which helps manage Go dependencies and BUILD files
http_archive(
    name = "bazel_gazelle",
    sha256 = "10c84a6d9d621a10280b3220ff7728b16350e09e267065e3715b1c88bf1686bc", # Use a known SHA256
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.35.0/bazel-gazelle-v0.35.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.35.0/bazel-gazelle-v0.35.0.tar.gz", # Check for the latest version
    ],
)

# Load the Go rules setup functions
load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

# Define Go dependencies
go_rules_dependencies()

# Register Go toolchains
go_register_toolchains()

# Define Gazelle dependencies
gazelle_dependencies()

# === Setup Node.js Dependencies ===
# This might need adjustment based on how you manage node_modules (e.g., yarn_install, npm_install)
# For now, let's assume npm install is run separately or managed via BUILD files.
# You might use rules_nodejs's npm_install or yarn_install here later.
# Example (needs adjustment):
# load("@rules_nodejs//nodejs:repositories.bzl", "npm_install")
# npm_install(
#     name = "npm", # Corresponds to the :npm target in BUILD files
#     package_json = "//:package.json", # Assuming a root package.json, adjust if needed
#     package_lock_json = "//:package-lock.json",
# )

