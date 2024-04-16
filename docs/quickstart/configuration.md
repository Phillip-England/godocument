# Configuration

## Project Structure

Godocument has the following file structure.

```bash
├── /docs
├── /html
├── /internal
├── /static
├── main.go
├── godocument.config.json
```

Let's focus on `/docs` and `godocument.config.json` for this quickstart.

## godocument.config.json

The `godocument.config.json` file for this application (up to this point) looks like this:

```json
{
    "docs": {
        "Introduction": "/introduction.md",
        "Quickstart": {
            "Installation": "/quickstart/installation.md",
            "Configuration": "/quickstart/configuration.md",
            ...
        }
        ...
    }
}
```

Each line under the `docs` object should point to an existing markdown file found in `./docs`. The application with throw an error if you attempt to point to a markdown file which does not exist.

<span class='content-important'>The order in which your organize your `godocument.config.json` file will determine the order of your elements in your apps navigation. It will also determine the order of your pages.</span>



## /docs

`./docs` for this app (up to this point in your reading) should look like this:

```bash
./docs
├── introduction.md
└── quickstart
    ├── configuration.md
    └── installation.md
```

<span class='content-important'>You should set up the structure of `/docs` to mirror the structure of `godocument.config.json`. Not doing so can lead to unforseen bugs.</span>