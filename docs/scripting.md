<meta name="description" content="Explore the integration of Htmx in Godocument to achieve a Single Page Application experience without the complexity of a full JavaScript framework. Learn about the implications of using hx-boost, managing JavaScript execution, and handling event listeners to ensure smooth navigation and interaction within your documentation site.">


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
3. Our page's current `<body>` will be replaced with the newly received `<body>`.

All of this will be done without a full-page refresh. Since we are generating static pages, the wait time between when a request is clicked and when the response is receieved will be minimal, giving the illusion of a single page experience to the user.

## Htmx Implications

Using Htmx's `hx-boost` attribute has implications on how we need to think about using Javascript in our application.

<span class='md-important'>Since the `<body>` is the only thing changed when using `hx-boost`, the Javascript located in the `<head>` of our document will only be loaded once. However, Javascript located in the `<body>` will be ran on each request.</span>

This can create issues when declaring functions, declaring variables, and mounting event listeners.

## loaded attribute

Godocument makes use of an attribute, `loaded`, on the `<html>` tag to avoid reinstantiating variables multiple times.

```html
<html lang="en" loaded="false" ..>
```

After the page is loaded on the initial visit, this attribute is set to `true`. This will prevent our variables and functions from being instantiated more than once.

<span class='md-warning'>Failing to set `loaded="true"` on `<html>` will result in unexpected behavior</span>

## onLoad function

Godocument makes use of the function `onLoad()` to run the appropriate Javascript on all page loads.

```js
function onLoad() {

    // elements
    const body = qs(document, 'body')
    const sitenav = qs(document, '#sitenav')
    const sitenavItems = qsa(sitenav, '.item')
    const sitenavDropdowns = qsa(sitenav, '.dropdown')
    const pagenav = qs(document, '#pagenav')
    const pagenavLinks = qsa(pagenav, 'a')
    const article = qs(document, '#article')
    const articleTitles = qsa(article, 'h2, h3, h4, h5, h6')
    const header = qs(document, '#header')
    const headerBars = qs(header, '#bars')
    const overlay = qs(document, '#overlay')
    const sunIcons = qsa(document, '.sun')
    const moonIcons = qsa(document, '.moon')
    const htmlDocument = qs(document, 'html')

    // hooking events and running initializations
    window.scrollTo(0, 0, { behavior: 'auto' })
    new SiteNav(sitenav, sitenavItems, sitenavDropdowns, header, overlay)
    new PageNav(pagenav, pagenavLinks, articleTitles)
    new Header(headerBars, overlay, sitenav)
    new Theme(sunIcons, moonIcons, htmlDocument)

    // web components
    doOnce(() => {
        customElements.define('md-important', MdImportant)
        customElements.define('md-warning', MdWarning)
        customElements.define('md-correct', MdCorrect)
    })

    // init
    Prism.highlightAll();

    // reveal body
    zez.applyState(body, 'loaded')

}
```

`onLoad()` is mounted to the `window` and `<body>` using the `DOMContentLoaded` and `htmx:afterOnLoad` events. Click [here](https://htmx.org/events/) to read more about Htmx events.

```js
eReset(window, 'DOMContentLoaded', onLoad) // initial page load
eReset(document.getElementsByTagName('body')[0], "htmx:afterOnLoad", onLoad) // after htmx swaps
```

<span class='md-important'>`DOMContentLoaded` will handle the initial page load, while `hmtx:afterOnLoad` will handle all other navigations.</span>

## Managing Events

Since the page is never refreshed using Htmx, we need to make sure we are unmounting events and remounting them on every navigation. `eReset()` is a handy function that does just that.

```js
function eReset(node, eventType, callback) {
    node.removeEventListener(eventType, callback)
    node.addEventListener(eventType, callback)
}
```

Instead of calling `element.addEventListener()`, it is better to use `eReset()` to ensure events are properly managed between page navigations.

<span class='md-warning'>Failing to unhook events upon navigation will result in the same events being hooked multiple times to the target element, which can have unexpected consequences and lead to poor memory management.</span>


