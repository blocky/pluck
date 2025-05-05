# pluck

Why is it so hard to put working code into docs? A few common options are:

1. Embed the code in the text.  But trying to actually run that code tends to
   be problematic.  The result is that often the code does not even
   run--let alone run correctly.

2. [Literate programming](https://en.wikipedia.org/wiki/Literate_programming)
   is a great option.  (In fact, I tried for a bit to write my dissertation
   using this approach.)  Often, the challenge with this option is that it needs
   purpose built system (such as [knitr](https://yihui.org/knitr/) or
   [jupyter](https://jupyter.org/)) and so it can be challenging if you have
   constraints that do not fit in those ecosystems.

3. Extract code by line number.  This seems to be a pretty common option. It
   is lightweight and allows the code to be tested. Unfortunately, it can
   result in some pretty odd results if the code changes without updating the
   line numbers.

Pluck seeks to provide many of the benefits of extracting by line, but extract
code by "tag".  And since we are talking about code, we can use functions and
type definitions as our "tags".

## Install

To install the latest version, run:

```bash
go install github.com/blocky/pluck/cmd/pluck@latest
```

And give it is try! Let's create a go file:

```bash
cat  <<EOF > go-file.go
package myPackge

type AType struct {
    FieldOfAType int
}

func (f *AType) AMethodOfAType() error {
    return nil
}

func AFunction() {}
EOF
```

And let's extract the type from the file:

```bash
pluck --input go-file.go --pick type:AType
```

While will produce the following code snippet:

```
type AType struct {
    FieldOfAType int
}
```

We can even grab multiple items for example the two functions:

```bash
pluck --input go-file.go --pick function:AType.AMethodOfAType --pick function:AFunction
```

Which produces the code snippet:

```
func (f *AType) AMethodOfAType() error {
    return nil
}

func AFunction() {}
```

And that is about it... Enjoy!
