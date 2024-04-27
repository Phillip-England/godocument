
<meta name="description" content="Learn how to set up and use Godocument, a static site generator inspired by Docusaurus and powered by Htmx, to easily document your code. Start building with Godocument using simple steps to create, configure, and deploy your documentation site.">

# Introduction

## What is Godocument?
Godocument is a static site generator inspired by [Docusaurus](https://docusaurus.io/) and powered by [Htmx](https://htmx.org). Documenting your code should be *simple*.

<md-important>Godocument requires Go version 1.22.0 or greater</md-important>

## Hello, World


A Godocument website can be created using the following steps:


- Make a directory

```bash
mkdir <your-apps-name>
cd <your-apps-name>
```

- Clone the repo

```bash
git clone https://github.com/phillip-england/godocument .
```

- Create a .env file:

```bash
touch .env
```

- Add the following environment variables:

```bash
PORT=8080 # for development
STATIC_PORT=8000 # for testing static assets
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

- Inside of `/docs`, create `first-page.md`

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

- Inside of `/docs/first-section`, create a file called `second-page.md`

```bash
touch /docs/first-section/second-page.md
```

- Add the following lines to `/docs/first-section/second-page.md`

```md
# Second Page

## Hello, World

This is the second page I've created using Godocument!
```

- From your application's root directory, run the following command to view the results on `localhost:8080`:

```bash
go run main.go
```

- To build your static assets, run:

```bash
go run main.go --build
```

That's it! Your example is deployment-ready and can be found at `/out`. You can easily deploy on Github Pages, Amazon S3, or a CDN of your choice.