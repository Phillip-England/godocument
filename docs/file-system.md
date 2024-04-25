# File system

## Project Structure

The structure of the application is as follows:

```bash
├── docs/
├── html/
├── internal/
├── out/
├── static/
```

Let's go through each directory so you can have a solid grasp of what each directory is responsible for.

## /docs

`/docs` will contain all of your markdown files. During development, Godocument will use this directory along with your `godocument.config.json` file to generate a series of routes. These routes will be served at localhost:8080.

<md-warning>Inside of `/docs` you will see a file named `introduction.md`. This file is **required** and Godocument will crash without it.</md-warning>

## /html

`/html` contains all of the html templates and components needed for your site. Godocument uses [Golang's standard html templating](https://pkg.go.dev/html/template). This directory is pretty straight forward and only comes with one caveat: `/html/components/sitenav.html` is auto-generated every time you run `go run main.go`.

<md-important>To make changes to `/html/components/sitenav.html`, you can edit the funcs `GenerateDynamicNavbar()` and `workOnNavbar()` found at `/internal/filewriter/filewriter.go`. These funcs are both responsible for generating and writing the html to `/html/components/sitenav.html`.</md-important>


## /internal

`/internal` contains all the source files for Godocument. Godocument is written entirely in Go. Here is an overview of the core modules and their purpose.

### /internal/config

`/internal/config` is responsible for taking the `godocument.config.json` file and turning it into a workable data structure. 

### /internal/contentrouter

`/internal/contentrouter` makes use of the config generated from `/internal/config` and uses the data to generate a series of routes. These routes are used during development to test the application. When building static assets, these routes are read by `/internal/filewriter` to generate the `/out` directory.

### /internal/filewriter

`/internal/filewriter` is responsible for working with files. Any functionality that has to do with writing, deleting, or minifying files is handled in this module.

### /internal/handler

`/internal/handler` contains any http handlers that are not auto-generated in `/internal/contentrouter` such as serving the favicon or static assests from `/static` or `/out`.

### /internal/middleware

`/internal/middleware` currently handles logging requests out to the terminal, but is built in a way which makes it very easy to extend. Since Godocument is intended to generate static assest, middleware will not be used in production. However, if you have a use case which requires more extensive middleware, Godocument can be shipped as a server side rendered application by skipping the build step and instead running `go build -o app`.

### /internal/stypes

`/internal/stypes` contains all of the custom types used in Godocument. It's sole purpose is to be a place where all types can live to prevent cross-imports.

## /out

`/out` is where all of your static assests will be generated when running `go run main.go --build`. This directory can then be deployed on a server such as Github Pages or Amazon S3.

## /static

`/static` contains all of the static assests needed to test your application during development. When building, this directory is read by  `/internal/filewriter`. `/internal/filewriter` will then bundle and minify all the files and copy them over to `/out/static`. 