<meta name="description" content="Learn how to build static assets with Godocument. This section explains how to use the 'go run main.go --build' command to create optimized production-ready assets, bundle stylesheets, and leverage minification with Minify for enhanced performance and efficiency.">


# Building

## go run main.go --build

This command will build all of your static assets and place them in `/out`.

Godocument will take all the relative paths you use during development and modify them into absolute paths for production.

<md-warning>Using relative paths during development is a requirement. Not doing so will result in unexpected behaviour when building.</md-warning>

## Providing an Absolute Path

When you are ready to build your site for production, you will need to provide an absolute path.

```bash
go run main.go --build <absolute-path>
```

<md-important>If an absolute path is not provided, Godocument will serve your assets on whatever port 8080. This is useful for testing your application prior to deployment. If you do provide an absolute path, Godocument will not serve the assets locally.</md-important>

For example, let's say I wanted to build for `godocument.dev`, I would run:

```bash
go run main.go --build godocument.dev
```

<md-warning>Absolute paths should not include a "/" at the end, this will result in a panic.</md-warning>

## Bundling Stylesheets

Godocument will bundle your stylesheets to reduce the number of network calls needed on the initial page load. During development, two network calls are made to introduce styling into your pages. Let's look at the `<head>` links in `/html/templates/base.html`:

```html
<link rel="stylesheet" type="text/css" href="/static/css/index.css">
<link rel="stylesheet" type="text/css" href="/static/css/output.css">
```

After bundling, these will be converted into:

```html
<link rel="stylesheet" type="text/css" href="/<absolute-path>/static/css/index.css">
```

The vanilla CSS will be stacked on top of the Tailwind CSS, giving Tailwind priority.

## Minification

Godocument uses [Minify](https://github.com/tdewolff/minify) to compress and minify static assets, which helps to reduce the file sizes of HTML, CSS, and Javascript files. These optimizations improves loading times and bandwidth usage in production.