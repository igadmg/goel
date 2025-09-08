# El - Embedded Language Documentation

El is a lightweight templating engine for Go, inspired by PHP's embedded code style. It allows you to embed Go code directly within your templates, making it easy to generate dynamic content.

## Syntax

El templates use `<?` and `?>` as delimiters to embed Go code within the template. Everything outside these delimiters is treated as plain text.
`<?` tag can have optional language specification parameter. Usually that should be used in first tag, currently only `<?go` is supported (they are not used anyway, but in future can be used to generate code on other languages).

### Basic Syntax

```
Hello, <?= name ?>!
```

In this example, name is a Go variable that will be evaluated and inserted into the template.

### Embedding Go Code

You can embed any valid Go code within the `<?` and `?>` delimiters:


```
<? 
for i := 0; i < 10; i++ { 
?>
    <p>Item <?= i ?></p>
<? 
} 
?>
```


This will generate a list of 10 items, each wrapped in a `<p>` tag.

Go code should be wrapped to a function accepting `wr io.Writer` parameter. And have proper go package an import definitions.

```
<?go
package main

import (
	"fmt"
	"io"
)

func generate(wr io.Writer) {
    for i := 0; i < 10; i++ { 
?>
    <p>Item <?= i ?></p>
<? 
    }
} 
?>
```

Later that function `generate` can be called from Go code. 

### Outputting Variables

To output a variable, use the `<?= ... ?>` syntax:

```
<p>Welcome, <?= user.Name ?>!</p>
```

### Control Structures

You can use Go's control structures directly within the templates:

```
<?
if user.IsAdmin {
?>
    <p>Welcome, Admin!</p>
<?
} else {
?>
    <p>Welcome, User!</p>
<?
}
?>
```

### Including Other Templates

You can include other by calling template function from anywhere in your template or go code and passing `wr` parameter there.

```
<?go
package main

import (
	"fmt"
	"io"
)

func generateItems(wr io.Writer) {
    for i := 0; i < 10; i++ { 
?>
    <p>Item <?= i ?></p>
<? 
    }
} 
?>
```

```
<?go
package main

import (
	"fmt"
	"io"
)

func generateTemplate(wr io.Writer) {
    generateItems(wr)
} 
?>
```

## Usage

To use the El tool from Go code put

`//go:generate go run github.com/igadmg/goel`

line anywhere in your code. It will scan project for `*.go.el` files and generate `*.go` files from them. It will process any `*.el` file and generate new file omitting `.el` extension. So `template.go.el` will be converted to `tempate.go`, and `template.cpp.el` will be converted to `template.cpp` (but C++ language model is not supported yet) 

### Example

Given a template file `greeting.go.el`:

```
<?go
package template

import (
    "fmt"
    "io"
)

func HelloTemplate(wr io.Writer) {
?>
<h1>Hello, <?= name ?>!</h1>
<?
}
?>
```

Running the el tool will generate a Go file `greeting.go` that have a function `HelloTemplate` that can be used to render the template.


### Tags

<? ?> - defines simple embedding of code into text template
<?= var ?> - define template code to output var into template
<?^= var ?> - define template code to output var into template with first letter converted to Title case
