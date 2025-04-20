"""Rules for building Vite projects with Bazel"""

load("@aspect_rules_js//js:defs.bzl", "js_binary")
load("@bazel_skylib//rules:write_file.bzl", "write_file")

def vite_project(name, srcs, deps = [], **kwargs):
    """Build a Vite project
    
    Args:
        name: Name of the target
        srcs: Source files for the Vite project
        deps: NPM dependencies
        **kwargs: Additional arguments to pass to the underlying rules
    """
    
    # Create a helper script to change to the package directory
    write_file(
        name = name + "_chdir",
        out = name + "_chdir.js",
        content = ["process.chdir(require('path').join(process.env.BUILD_WORKSPACE_DIRECTORY, native.package_name()));"],
    )
    
    # Combine all inputs
    all_srcs = srcs + deps + [":" + name + "_chdir"]
    
    # Create a target to run vite build
    js_binary(
        name = name,
        entry_point = "@npm//vite:node_modules/vite/bin/vite.js",
        data = all_srcs,
        args = [
            "build",
            "--node_options=--require=./" + name + "_chdir.js",
        ],
        **kwargs
    )
    
    # Create a dev server target
    js_binary(
        name = name + "_dev",
        entry_point = "@npm//vite:node_modules/vite/bin/vite.js",
        data = all_srcs,
        args = [
            "--node_options=--require=./" + name + "_chdir.js",
        ],
        **kwargs
    )