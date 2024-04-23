# Introduction

## What is Godocument?
Godocument is a static site generator inspired by [Docusaurus](https://docusaurus.io/) and powered by [Htmx](https://htmx.org). Documenting your code should be *simple*.

<md-important>Godocument requires Go version 1.22.0 or greater</md-important>

## Hello, World


A simple Godocument website can be created using the following steps:


- Make a directory

```bash
mkdir <your-apps-name>
cd <your-apps-name>
```

- Clone the repo

```bash
git clone https://github.com/phillip-england/godocument .
```

- Add some new entries to `godocument.config.json`:

```json
{
    "docs": {
        "Introduction": "/introduction.md",
        "First Page": "/first-page.md",
        "First Section": {
            "Second Page": "/first-section/second-page.md"
        }
    }
}
```

- Inside of `/docs` create `first-page.md`

```bash
touch /docs/first-page.md
```

- Add the following lines to `/docs/first-page.md`

```md
# First Page

## Hello, World

This is the first page I've created using Godocument!
```

- Inside of `/docs` create a directory called `/first-section`

```bash
mkdir /docs/first-section
```

- Inside of `/docs/first-section` create a file called `second-page.md`

```bash
touch /docs/first-section/second-page.md
```

- Add the following lines to `/docs/first-section/second-page.md`

```md
# Second Page

## Hello, World

This is the second page I've created using Godocument!
```

- From the root of your application, run the following to view the results on `localhost:8080`:

```bash
go run main.go
```

- To build your static assets, run:

```bash
go run main.go --build
```

- To test your static assests locally, run:

```bash
go run main.go --static
```

That's it! Your example is deployment-ready and can be found at `./out`. You can easily deploy on Github Pages, Amazon S3, or a CDN of your choice.