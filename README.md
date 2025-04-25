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
code by "tag".  And since we are talking about code, we can use functions as
our "tags".  In addition, we provide some functionality for specifying how that
function is rendered.  And that is it, "do one thing well".
