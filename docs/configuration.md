# Configuration

## godocument.config.json

`godocument.config.json` is the configuration file for your application. It contains the necessary information to generate your website's routes. 

<md-important>The order of items in `godocument.config.json` will determine the order of your pages in your website.</md-important>

Here is the base configuration needed to generate a site using Godocument:

```json
{
    "docs": {
        "Introduction": "/introduction.md"
    }
}
```

<md-warning>The `/docs` directory and the `/docs/introduction.md` file are required for Godocument. Also, the json object `"docs"` must be named `"docs"` and the first entry beneath `"docs"` must be `"Introduction": "/introduction.md"`. Failing to meet these requirements will result in a panic.</md-warning>

## Pages

The entries in `godocument.config.json` can either be pages or sections. Let's start with pages.

To denote a page, simply create a key-value pair with the key being the name of the page and the value being the file path to the `.md` file for the page. You can name pages whatever you would like.

<md-important>All file paths in `godocument.config.json` are relative to `/docs`. This means you do not have to the include `/docs` in your file paths as Godocument assumes all your markdown files are in `/docs`.</md-important>

Here is how you would add a new page to the base configuration:

```json
{
    "docs": {
        "Introduction": "/introduction.md",
        "New Page": "/new-page.md"
    }
}
```

After adding the page to `godocument.config.json` you will need to create the associated file. From the root of your application, run:

```bash
touch /docs/new-page.md
```

Then, add the following lines to `/docs/new-page.md`:

```md
# New Page

I created a new page using Godocument!
```

From the root of your application, run `go run main.go` and view the results at `localhost:8080`.

## Sections

Sections are named areas of your website which contain a series of pages. Sections can also contain sub-sections. In `godocument.config.json`, a section can be denoted by creating an object. For example:

```json
{
    "docs": {
        "Introduction": "/introduction.md",
        "New Section": {
            
        }
    }
}
```

After creating a section, you can add pages within it:

```json
{
    "docs": {
        "Introduction": "/introduction.md",
        "New Section": {
            "About": "/new-section/about.md" 
        }
    }
}
```

You can also add sub-sections:

```json
{
    "docs": {
        "Introduction": "/introduction.md",
        "New Section": {
            "About": "/new-section/about.md",
            "More Info": {
                "Origins": "/new-section/more-info/origins.md"
            }
        }
    }
}
```

Create the corresponding files and directories:

```bash
mkdir /docs/new-section
touch /docs/new-section/about.md
mkdir /docs/new-section/more-info
touch /docs/new-section/more-info/origins.md
```

Add the following content to `/docs/new-section/about.md`

```md
# About

I created a page within a section using Godocument!
```

Then, add the following lines to `/docs/new-section/more-info/origin.md`:

```md
# Origins

I created a page within a sub-section using Godocument!
```

To test the results, run `go run main.go` from the root of your application and visit `localhost:8080`

## /docs structure

Godocument does not require you to structure your `/docs` directory in any particular way, **BUT** it is highly recommended to have your `/docs` directory mirror the structure of your `godocument.config.json` file.

For example, here is a `godocument.config.json` file which does not follow the proper conventions.

<md-warning>The example below does not follow the recommended conventions for `godocument.config.json`.</md-warning>

```json
{
    "docs":{
        "Introduction": "/introduction.md",
        "More Info": {
            "About": "/about.md"
        }
    }
}
```

It does not follow the conventions because `/about.md` should have a file path which mirrors the structure of `godocument.config.json`.

<md-correct>To correct the above `godocument.config.json` make the changes below.</md-correct>

```json
{
    "docs": {
        ...
        "More Info": {
            "About": "/more-info/about.md"
        }
    }
}
```

Such a change will ensure that the /docs directory mirrors the structure of godocument.config.json, as recommended.
