# Introduction

## What is Godocument?
Godocument is a static site generator inspired by [Docusaurus](https://docusaurus.io/) and powered by [HTMX](https://htmx.org). Documenting your code should be *simple*.

<mkd-important text='Godocument requires Go version 1.22.0 or greater'></mkd-important>

## Hello, World


The simplest Godocument website can be created using the following steps:


1. Make a directory

```bash
mkdir <your-apps-name>
cd <your-apps-name>
```

2. Clone the repo

```bash
git clone https://github.com/phillip-england/godocument .
```

3. Build static assets

```bash
go run main.go --build
```

That's it! Your example is deployment-ready and can be found at `./out`. You can easily deploy on Github Pages, Amazon S3, or a CDN of your choice.

To test your static assests locally, run:

```bash
go run main.go --static
```