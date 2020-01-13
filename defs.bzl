def _impl(ctx):
    content = ""

    for s in ctx.attr.defines:
        content += "%s=1\n" % s

    for k, v in ctx.attr.variables.items():
        content += "%s=%s\n" % (k, v)

    variable_file = ctx.actions.declare_file(ctx.label.name + "-variables")
    ctx.actions.write(variable_file, content = content)

    args = ctx.actions.args()
    args.add("-output", ctx.outputs.out.path)
    if not ctx.attr.cmake:
        args.add("-nocmake")
    args.add(variable_file)
    args.add(ctx.file.template)
    ctx.actions.run(
        executable = ctx.executable._configure,
        inputs = [variable_file, ctx.file.template],
        outputs = [ctx.outputs.out],
        arguments = [args],
    )

expand_configure = rule(
    implementation = _impl,
    attrs = {
        "template": attr.label(allow_single_file = True, mandatory = True),
        "defines": attr.string_list(default = []),
        "out": attr.output(mandatory = True),
        "variables": attr.string_dict(default = {}),
        "cmake": attr.bool(default = True),
        "_configure": attr.label(executable = True, default = "//cmd/configure_template", cfg = "host"),
    },
)