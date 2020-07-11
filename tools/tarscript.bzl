"""
run a script with argv[1] as output, and tar the output
author: adamyi
"""

def tarscript(name, src):
    cmd = "rm -rf $(@D)/%s; mkdir $(@D)/%s; bash -c \"$(location %s) $(@D)/%s\"; tar -C $(@D)/%s -cf $(@D)/%s.tar .; rm -rf $(@D)/%s" % (name, name, src, name, name, name, name)
    native.genrule(
        name = name,
        outs = ["%s.tar" % name],
        srcs = [src],
        cmd = cmd,
    )
