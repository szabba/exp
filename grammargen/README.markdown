Generating text with context free grammars.

* [x] Code to load a grammar from JSON
* [x] JSON repr of the default grammar
* [x] file to load JSON from as an option
    * default = `""` -- memory
    * `"-"` -- from stdin
* [ ] make the errors nicer
* [ ] option to output JSON
* [ ] seed option
    * `"time"` means use current unix time
    * a 64-bit signed integer is the literal seed
* [ ] option to specify how many samples to generate