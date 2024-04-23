# Core Technologies

Godocument depends on a few technologies. To use Godocument, you do not have to know these techs, *but* to extend the application and customize it to your liking, these are the things you will need to be familiar with.

## Tailwind

Godocument primarly uses [Tailwind](https://tailwindcss.com/) for styling. Some vanilla css is used as well, but the majority of styling is done using Tailwind.

If you intend to modify any classes in your application, you will need to download the binary from [here](https://tailwindcss.com/blog/standalone-cli) and place it on your $PATH.

Assuming you named your binary "tailwindcss", run the following command from the root of your application to use Tailwind:

```bash
tailwindcss -i "./static/css/input.css" -o "./static/css/output.css" --watch
```

## Htmx

Godocument uses [Htmx](https://htmx.org) to provide a better user experience. The primary feature used is [hx-boost](https://htmx.org/attributes/hx-boost).

In `/html/templates/base.html` you will notice the hx-boost attribute on the `<body>` tag:

```html
<body hx-boost="true" ...></body>
```

This attribute prevents the page from doing a full reload upon navigation, but it does come with a caveats: 

<md-important>hx-boost fundamentally changes the way we need to think about using Javascript for interactivity.</md-important>

To learn more about how this impacts the use of Javascript, [read more]("").

When clicking any navigational links in the application, hx-boost will take over. Instead of issuing a full page refresh, hx-boost will send an AJAX request to the href of the clicked anchor tag. Then, when it gets the responding HTML back, it will replace the content within the `<body>` tag with the new html.

This has the benefit of providing a SPA-like user experience without all the overhead of a traditional framework.

## Goldmark

Godocument uses [Goldmark](https://github.com/yuin/goldmark) to convert `.md` files into HTML. Goldmark is extendable, but Godocument uses only a few basic features.

In `/internal/contentrouter` we use Goldmark in the following way:

```go
// ...
md := goldmark.New(
    goldmark.WithParserOptions(
        parser.WithAutoHeadingID(),
    ),
    goldmark.WithRendererOptions(
        html.WithUnsafe(),
    ),
)
// ...
```

The `parser.WithAutoHeadingID()` makes it so the HTML generated from our headers in our `.md` files will automatically have an id attribute.

The `html.WithUnsafe()` enables us to write html directly into our markdown files. This is considered unsafe in situations where your content is dynamically controlled by users, but in our case, we will be writing all the markup ourselves.

I opted in the use `html.WithUnsafe()` because it enables us to create custom markdown components using Web Component.

## Web Components

Godocument uses [Web Components](https://developer.mozilla.org/en-US/docs/Web/API/Web_components) to create custom elements for markdown content. For example:

<md-important>I was created using a custom web component named, `<md-important>`</md-important>