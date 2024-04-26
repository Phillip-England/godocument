# Building

## go run main.go --build

This command will build all of your static assets and place them in `/out`. In order to use this command, ensure `SERVER_URL` environment variable is set in your `.env` file.

Godocument will take all the relative paths you use during development and modify them into absolute paths for production.

<md-important>If you do not set a value for `SERVER_URL` in `.env`, Godocument will use `localhost:8080` when setting absolute paths.</md-important>

## Bundling Stylesheets

Godocument will bundle your stylesheets to reduce the number of network calls needed on the initial page load. During development, two network calls are made to introduce styling into your pages. Let's look at the `<head>` links in `/html/templates/base.html`:

```html
<link rel="stylesheet" type="text/css" href="/static/css/index.css">
<link rel="stylesheet" type="text/css" href="/static/css/output.css">
```

After bundling, these will be converted into:

```html
<link rel="stylesheet" type="text/css" href="<SERVER_URL>/static/css/index.css">
```

The vanilla CSS will be stacked on top of the Tailwind CSS, giving Tailwind priority.

## Minification

Godocument uses [Minify](https://github.com/tdewolff/minify) to compress and minify static assets, which helps to reduce the file sizes of HTML, CSS, and Javascript files. These optimizations improves loading times and bandwidth usage in production.