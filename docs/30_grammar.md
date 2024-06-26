# Backus–Naur Form Grammar for Valid Identifiers

This document will provide the grammar for the identifiers that can be used inside the configurations and as names
of modules and add-ons for our tool.

```text
<valid tag>     ::= "add-on-" <identifier> "-" <identifier> "-" <valid semver>
                  | "module-" <identifier> "-" <identifier> "-" <valid semver>

<addon-name>    ::= <identifier> "/" <identifier>

<module-name>   ::= <identifier> "/" <identifier> "/" <identifier>

<group-name>    ::= <identifier>

<cluster-name>  ::= <identifier>

<identifier>    ::= <alphanumerics> | <alphanumerics> "-" <identifier>

<alphanumerics> ::= <letter> | <letter> <characters>

<characters>    ::= <character> | <character> <characters>

<character>     ::= <digit> | <letter>

<digit>         ::= "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9"

<letter>        ::= "a" | "b" | "c" | "d" | "e" | "f" | "g" | "h" | "i" | "j"
                  | "k" | "l" | "m" | "n" | "o" | "p" | "q" | "r" | "s" | "t"
                  | "u" | "v" | "w" | "x" | "y" | "z"
```

`<valid semver>` is a valid semver as described by the [grammar definition] of [semver]

[grammar definition]: https://semver.org/spec/v2.0.0.html#backusnaur-form-grammar-for-valid-semver-versions "semver grammar"
[semver]: https://semver.org/spec/v2.0.0.html "semantic versioning v2.0.0 site"
