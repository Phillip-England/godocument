# Installation

## The Repo

A blank Godocument template can be found at [https://github.com/phillip-england/godocument](https://github.com/phillip-england/godocument).

<mkd-important text='All the commands in this guide will assume you are using a Unix-based terminal.'></mkd-important>

To get started, create a directory and clone the repo within it:

```bash
mkdir <your-app-name>
cd <your-app-name>
git clone https://github.com/phillip-england/godocument .
```

## The File System

The structure of the application is as follows:

```bash
├── docs/
├── godocument.config.json
├── html/
├── internal/
├── out/
├── static/
├── main.go
```

Let's go through each directory so you can have a solid grasp of what each directory is responsible for.

### docs/

`docs/` will contain all of your markdown files. During development, Godocument will use this directory along with your `godocument.config.json` file to generate a series of routes. These routes will be served at localhost:8080.

Inside of `docs/` you will see a file named `introduction.md`. This file is **required** and Godocument will crash without it.

<mkd-warning text='`/docs` and `/docs/introduction.md` are required for Godocument. Renaming or moving either of these will result in an error.'></mkd-warning>

