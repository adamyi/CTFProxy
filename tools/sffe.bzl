"""
group sffe files
author: adamyi
"""

def sffe_files(files, name = "sffe"):
    native.filegroup(
        name = name,
        srcs = files,
        visibility = ["//infra/sffe:__pkg__"],
    )
