// Below is the output from $(vdl help -style=godoc ...)

/*
The vdl tool manages veyron VDL source code.  It's similar to the go tool used
for managing Go source code.

Usage:
   vdl [flags] <command>

The vdl commands are:
   generate    Compile packages and dependencies, and generate code
   compile     Compile packages and dependencies, but don't generate code
   listinfo    List package and dependency info in transitive order
   help        Display help for commands

The vdl flags are:
   -max_errors=-1: Stop processing after this many errors, or -1 for unlimited.
   -v=false: Turn on verbose logging.

vdl Generate

Generate compiles packages and their transitive dependencies, and generates code
in the specified languages.

Usage:
   vdl generate [flags] <packages>

<packages> are a list of packages to process, specified as arguments for each
command.  The format is similar to the go tool.  In its simplest form each
package is an import path; e.g. "veyron/lib/vdl".  A package that is an absolute
path or that contains a "." is interpreted as a file system path and denotes the
package in that directory.  A package that ends with "..." does a wildcard match
against all directories with that prefix.  The special import path "all" expands
to all package directories found in all the GOPATH trees.

For more information use "go help packages" to see the standard go package
documentation.

The generate flags are:
   -go_fmt=true: Format generated Go code
   -go_out_dir=: Go output directory.  There are three modes:
      ""         : Generate output in-place in the source tree
      "dir"      : Generate output rooted at dir
      "src->dst" : Generate output rooted at x, with translation from src to dst
   Assume your source tree is organized as follows:
   GOPATH=/home/me/code/v
      /home/me/code/v/src/veyron2/vdl/test_base/base1.vdl
      /home/me/code/v/src/veyron2/vdl/test_base/base2.vdl
   Here's example output under the different modes:
   --go_out_dir=""
      /home/me/code/v/src/veyron2/vdl/test_base/base1.vdl.go
      /home/me/code/v/src/veyron2/vdl/test_base/base2.vdl.go
   --go_out_dir="/tmp/foo"
      /tmp/foo/veyron2/vdl/test_base/base1.vdl.go
      /tmp/foo/veyron2/vdl/test_base/base2.vdl.go
   --go_out_dir="v/src->bar/src"
      /home/me/code/bar/src/veyron2/vdl/test_base/base1.vdl.go
      /home/me/code/bar/src/veyron2/vdl/test_base/base2.vdl.go
   When the src->dst form is used, src must match the suffix of the path just
   before the package path, and dst is the replacement for src.
   -java_out_dir=v/src->java/src/main/java: Same semantics as --go_out_dir but applies to java code generation.
   -java_pkg_prefix=com: Package prefix that will be added to the VDL package prefixes when generating Java files.
   -lang=[go]: Comma-separated list of languages to generate, currently supporting "go", "java"
   -status=true: Show package names while we compile

vdl Compile

Compile compiles packages and their transitive dependencies, but does not
generate code.  This is useful to sanity-check that your VDL files are valid.

Usage:
   vdl compile [flags] <packages>

<packages> are a list of packages to process, specified as arguments for each
command.  The format is similar to the go tool.  In its simplest form each
package is an import path; e.g. "veyron/lib/vdl".  A package that is an absolute
path or that contains a "." is interpreted as a file system path and denotes the
package in that directory.  A package that ends with "..." does a wildcard match
against all directories with that prefix.  The special import path "all" expands
to all package directories found in all the GOPATH trees.

For more information use "go help packages" to see the standard go package
documentation.

The compile flags are:
   -status=true: Show package names while we compile

vdl Listinfo

Listinfo returns information about packages and their transitive dependencies,
in transitive order.  This is the same order the generate and compile commands
use for processing.  If "vdl listinfo A" is run and A depends on B, which
depends on C, the returned order will be C, B, A.  If multiple packages are
specified the ordering is over all combined dependencies.

Reminder: cyclic dependencies between packages are not allowed.  Cyclic
dependencies between VDL files within the same package are also not allowed.
This is more strict than regular Go; it makes it easier to generate code for
other languages like C++.

Usage:
   vdl listinfo <packages>

<packages> are a list of packages to process, specified as arguments for each
command.  The format is similar to the go tool.  In its simplest form each
package is an import path; e.g. "veyron/lib/vdl".  A package that is an absolute
path or that contains a "." is interpreted as a file system path and denotes the
package in that directory.  A package that ends with "..." does a wildcard match
against all directories with that prefix.  The special import path "all" expands
to all package directories found in all the GOPATH trees.

For more information use "go help packages" to see the standard go package
documentation.

vdl Help

Help displays usage descriptions for this command, or usage descriptions for
sub-commands.

Usage:
   vdl help [flags] <command>

<command> is an optional sequence of commands to display detailed per-command
usage.  The special-case "help ..." recursively displays help for this command
and all sub-commands.

The help flags are:
   -style=text: The formatting style for help output, either "text" or "godoc".
*/
package main

import "veyron2/vdl/vdl/cmds"

func main() {
	cmds.Root().Main()
}
