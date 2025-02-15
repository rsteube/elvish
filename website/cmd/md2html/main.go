// Command md2html converts Markdown to HTML.
//
// It is tailored for converting the Markdown sources of the Elvish website
// (https://elv.sh) to HTML.
//
// It first applies a pre-processing step that expands the following macros:
//
//   - @module inserts elvdocs for a given module
//
//   - @ttyshot inserts a ttyshot
//
//   - @dl expands to a binary download link
//
// The processed Markdown source is then converted to HTML using a codec based
// on [md.HTMLCodec], with the following additional features:
//
//   - Autogenerated ID for each heading
//
//   - Self link for each heading
//
//   - Table of content (optional, turn on with <!-- toc -->)
//
//   - Implicit links to elvdoc targets when link destination is empty and link
//     text is code span - for example, [`put`]() has destination
//     builtin.html#put (or just #put within doc for the builtin module itself)
//
//   - Section numbers for headings (optional, turn on with <-- number-sections
//     -->)
//
//   - Syntax highlighting of code blocks with language elvish or
//     elvish-transcript
//
// The comment block for optional features should appear before the main text,
// and can contain multiple features (like <!-- toc number-sections -->).
//
// A note on the implicit elvdoc target feature: ideally, we would like to use
// Markdown's shortcut link feature and let a simple [`put`] have an implicit
// target of builtin.html#put. Doing this has two prerequisites:
//
//   - The [src.elv.sh/pkg/md] package must be modified to support shortcut
//     links.
//
//   - The [src.elv.sh/pkg/elvdoc] package must be modified to add declarations
//     of these shortcut targets. This is because shortcut links in Markdown
//     must be declared, otherwise [`put`] is just literal [<code>put</code>].
//
//     This is feasible for targets in the same file, but much more tricky for
//     targets in a different module. For example, if the elvdoc of a:foo
//     references b:bar, we need to insert a declaration for b:bar to the elvdoc
//     of a:foo.
//
// Hence we've settled on using an empty target for now. It is a bit ugly but
// hopefully not too ugly.
package main

import (
	"os"
	"strings"

	"src.elv.sh/pkg/md"
)

func main() {
	var expanded strings.Builder
	f := filterer{}
	f.filter(os.Stdin, &expanded)

	codec := &htmlCodec{}
	codec.preprocessInline = func(ops []md.InlineOp) {
		addImplicitElvdocTargets(f.module, ops)
	}
	codec.ConvertCodeBlock = convertCodeBlock
	md.Render(expanded.String(), md.SmartPunctsCodec{Inner: codec})
	os.Stdout.WriteString(codec.String())
}
