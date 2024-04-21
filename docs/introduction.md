# Introduction

## What is Godocument?
Godocument is a static site generator inspired by [Docusaurus](https://docusaurus.io/) and powered by [HTMX](https://htmx.org). Documenting your code should be *simple* and not require knowledge of a javascript framework.

## Hello, World
The simplest Godocument website can be created using the following steps.

1. Make your project directory:

```bash
mkdir <your-projects-name>
cd <your-projects-name>
```

2. Clone the repo:

```bash
git clone https://github.com/phillip-england/godocument .
```

3. With Go 1.22.0 or later run:

```bash
go run main.go --reset
```

4. To view in the browser at localhost:8080 run:

```bash
go run main.go
```

5. To build static assets run:

```bash
go run main.go --build
```

That's it! Your example is deployment-ready and can be found at `./out`. You can easily deploy on Github Pages, Amazon S3, or any server that can serve static content.