# Custom Components

## Web Components

Godocument uses [Web Components](https://developer.mozilla.org/en-US/docs/Web/API/Web_components) for creating custom elements to use within your `.md` files.

Web Components are defined in `/static/js/index.js` and can be used directly in our `.md` files to keep things clean and quick to use.

## Goldmark html.WithUnsafe

Since we are using [Goldmark](https://github.com/yuin/goldmark) to convert `.md` files into workable HTML, we have to use the `html.WithUnsafe()` renderer option when instantiating Goldmark in our project. This will allow us to place HTML elements directly in our `.md` files.

This is only considered *unsafe* if the content within our `.md` files are controlled by our users. In our case, since we will be the ones writing the markup directly, it is not considered unsafe.

<md-important>Removing `html.WithUnsafe()` from Goldmark's rendering options will cause Goldmark to ignore any HTML markup within our `.md` files.</md-important>

## Component Initialization

Custom components require some initialization to work as expected. This is due to the way Goldmark parses `.md` files. These initialization steps are to ensure the component renders properly as well as enabling users to utilize backticks to convey inline code examples.

Within a components constructor, do the following to initialize the component properly and get the text from within the component:

```js
class MdImportant extends HTMLElement {
    constructor() {
        super()
        this.hook()
    }
    hook() {
        let text = initMdComponent(this)   
        this.innerHTML = `
            <div class='bg-[var(--md-bg-color)] dark:bg-[var(--dark-md-bg-color)] p-4 rounded-md border-l-4 border-[var(--md-important-border-color)] dark:border-[var(--dark-md-important-border-color)] flex flex-col gap-2'>
                <span class='flex flex-row item-center gap-2 dark:text-[var(--dark-md-important-text-color)] text-[var(--md-important-text-color)]'>
                    <span class='flex items-center dark:text-[var(--dark-md-important-text-color)] text-[var(--md-important-text-color)]'>
                        <svg class="h-6 w-6" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
                            <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 13V8m0 8h.01M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"/>
                        </svg>
                    </span>
                    <p class='font-bold'>Important</p>                   
                </span>
                <p class='custom-inline-code'>${text}</p>
            </div>
        `
    }    
}
```

We are specifically focusing on this line:

```js
...
let text = initMdComponent(this)
...
```

### Removing the Parent

When Goldmark find HTML within a `.md` file, it will wrap the HTML with a `<p>` tag. To ensure things render the way you expect, it is important to remove the `<p>` tag when instantiating our components.

A special function named `initMdComponent()` serves to enable custom components behave as expected.

```js
function initMdComponent(node) {
    node.parentElement.replaceWith(node)
    let text = node.innerHTML
    text = replaceBackticksWithCodeTags(text)
    return text
}
```

### Enabling Inline Code

Goldmark will not wrap backticks with `<code>` tags within our HTML markup found within `.md` files. To prevent users from manually having to type out `<code>some code example</code>` each time, `initMdComponent()` handles this as well.

`initMdComponent` references another function, `replaceBackticksWithCodeTags()`. Here is that functions declaration:

```js
function replaceBackticksWithCodeTags(text) {
    for (let i = 0; i < text.length; i++) {
        if (text[i] == '`') {
            text = text.slice(0, i) + '<code>' + text.slice(i + 1)
            i++
            while (i < text.length && text[i] != '`') {
                i++
            }
            text = text.slice(0, i) + '</code>' + text.slice(i + 1)
        }
    }
    return text
}
```


## MdImportant

<md-important>I am created by using this markup in a `.md` file: `<md-important>I am important!</md-important>`</md-important>

## MdWarning

<md-warning>I am created by using this markup in a `.md` file: `<md-warning>I am a warning!</md-warning>`</md-warning>

## MdCorrect

<md-correct>I am created by using this markup in a `.md` file: `<md-correct>That is correct!</md-correct>`</md-correct>