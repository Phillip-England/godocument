# Scripting

## Single Page Application

Godocument uses [Htmx](https://htmx.org) to provide a single page application user experience. Htmx enables us to do this without the complexity of a Javascript framework.

This functionality is enabled by a single attribute on our `<body>` tag, [hx-boost](https://htmx.org/attributes/hx-boost).

```html
<body hx-boost='true' ...>
```

Basically, when you click on a navigational link within your website, Htmx will take over and the following will happen:

1. An AJAX request will be sent to the `href` of the clicked `<a>`.
2. When the response is received, the `<body>` from the response will be isolated.
3. Our page's current `<body>` will be replaced with the newly recieved `<body>`.

All of this will be done without a full-page refresh. Since we are generating static pages, the wait time between when a request is clicked and when the response is receieved will be minimal, giving the illusion of a single page experience to the user.

## Htmx Implications

Using Htmx's `hx-boost` attribute has implications on how we need to think about using Javascript in our application.

<md-important>Since the `<body>` is the only thing changed when using `hx-boost`, the Javascript located in the `<head>` of our document will only be loaded once. However, Javascript located in the `<body>` will be ran on each request.</md-important>

This can create issues when declaring function, delcaring variables, and mounting event listeners.

