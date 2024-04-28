<meta name="description" content="Learn to enhance your GoDocument markdown files by integrating Web Components and custom HTML elements with Goldmark's `html.WithUnsafe` renderer. This guide covers the initialization and usage of custom elements such as `md-important`, `md-warning`, and `md-correct` to effectively emphasize key content in your documentation.">



# Custom Components

## Creating a New Component

<span class='md-important'>All Javascript for Godocument is located in `/static/js/index.js`.</span>

To create a new custom component, first create a class to represent the component. 

Let's call this component, `simple-component`. The constructor should take in a parameter called `customComponent`, which will be an object built from the `CustomComponent` class.

The `CustomComponent` class is already included in Godocument.

```js
class SimpleComponent {
    constructor(customComponent) {

    }
}
```

In `onLoad()`, be sure to instantiate both `CustomComponent` and `SimpleComponent`:

```js
function onLoad() {
    // ...

    // defining custom components
    let customComponents = new CustomComponent()
    new SimpleComponent(customComponents)
}
```

In the constructor of `SimpleComponent` call `customComponent.registerComponent(className, htmlContent)` as follows:

```js
class SimpleComponent {
    constructor(customComponent) {
        customComponent.registerComponent("simple-component", `
            <div>
                <h1>Hello, World</h1>
                <p>{text}</p>
            </div>
        `)
    }
}
```

Take note of this line:

```js
<p>{text}</p>
```

Register component will replace `{text}` with the `innerHTML` of the component.

To utilize `simple-component`, place the following markup in any of your `.md` files:

```md
<span class='simple-component'>I am a simple component!</span>
```

## Included Components

Godocument comes with a few components already built. Here they are:

### MdImportant

<span class='md-important'>I am created by using this markup in a `.md` file: `<span class='md-important'>I am important!</span>`</span>

### MdWarning

<span class='md-warning'>I am created by using this markup in a `.md` file: `<span class='md-warning'>I am a warning!</span>`</span>

### MdCorrect

<span class='md-correct'>I am created by using this markup in a `.md` file: `<span class='md-correct'>That is correct!</span>`</span>

## Goldmark html.WithUnsafe

Since we are using [Goldmark](https://github.com/yuin/goldmark) to convert `.md` files into workable HTML, we have to use the `html.WithUnsafe()` renderer option when instantiating Goldmark in our project. This will allow us to place HTML elements directly in our `.md` files.

This is only considered *unsafe* if the content within our `.md` files is controlled by our users. In our case, since we will be the ones writing the markup directly, it is not considered unsafe.

<span class='md-warning'>Removing `html.WithUnsafe()` from Goldmark's rendering options will cause Goldmark to ignore any HTML markup within our `.md` files.</span>