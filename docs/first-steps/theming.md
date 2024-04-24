# Theming

## CSS Variables

Godocument makes use of CSS variables to give users more control of their theme. Variables are either viewed as *utility* variables or *element-specific* variables.

<md-important>To adjust the themeing for your site, edit the variables found at the top of `/static/css/index.css`.</md-important>

## Utility Variables

Utility variables are not directly used in markup. Rather, they are intended to be used in *element-specific* variables. Here are the color utility variables for this site:

```css
--white: #fafafa;
--black: #1f1f1f;
--gray-800: #333333;
--gray-700: #555555;
--gray-600: #777777;
--gray-500: #999999;
--gray-400: #bbbbbb;
--gray-300: #dddddd;
--gray-200: #f0f0f0;
--gray-100: #f5f5f5;
--light-gray: #d0d0d0;
--gray: #555555;
--darkest-gray: #222222;
--purple: #ba8ef7;
--dark-purple: #712fec;
--green: #3bec74;
--dark-green: #057d2f;
--pink: #b370b1;
--yellow: #ffea6b;
--dark-yellow: #7f7108;
--orange: #ffa763;
--dark-orange: #c64719;
--blue: #2494da;
--dark-blue: #1b6dbf;
--red: #ff4d3f;
--dark-red: #c82216;
```

## Element-Specific Variables

Element-specific variables make use of *utility variables*. Here are the element-specific variables that control the colors in the codeblocks found on this site:

```css
/* light code blocks */
--code-bg-color: var(--gray-200);
--code-token-property: var(--dark-purple);
--code-string: var(--dark-green);
--code-token-selector: var(--dark-orange);
--code-function: var(--dark-yellow);
--code-keyword: var(--dark-purple);
--code-operator: var(--black);
--code-punctuation: var(--gray-700);
--code-important: var(--dark-orange);
--code-comment: var(--gray-700);

/* dark code blocks */
--dark-code-bg-color: var(--gray-800);
--dark-code-token-property: var(--purple);
--dark-code-string: var(--green);
--dark-code-token-selector: var(--orange);
--dark-code-function: var(--yellow);
--dark-code-keyword: var(--purple);
--dark-code-operator: var(--white);
--dark-code-punctuation: var(--gray-300);
--dark-code-comment: var(--gray-300)
```

## Variables in Tailwind

Godocument makes use of Tailwind's ability to use CSS variables within Tailwind classes. For example, here is the markup for the `<header>` at the top of this page:

```html
<header id="header" class="flex flex-row justify-between items-center border-b z-30 p-4 sticky top-0 w-full bg-[var(--header-bg-color)] dark:bg-[var(--dark-header-bg-color)] border-[var(--b-color)] dark:border-[var(--dark-b-color)]" style="grid-area: header;">
    <div class="header-logo-wrapper flex flex-row shrink-0 items-center w-[250px]">
        <img class="logo items-center flex-row" src="/static/img/logo.svg" alt="logo" id="logo">
    </div>
    <svg class="sun cursor-pointer hidden lg:block shrink-0 dark:lg:hidden" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
        <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5V3m0 18v-2M7.05 7.05 5.636 5.636m12.728 12.728L16.95 16.95M5 12H3m18 0h-2M7.05 16.95l-1.414 1.414M18.364 5.636 16.95 7.05M16 12a4 4 0 1 1-8 0 4 4 0 0 1 8 0Z"/>
    </svg>
    <svg class="moon cursor-pointer hidden shrink-0 dark:lg:block" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24">
        <path stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 21a9 9 0 0 1-.5-17.986V3c-.354.966-.5 1.911-.5 3a9 9 0 0 0 9 9c.239 0 .254.018.488 0A9.004 9.004 0 0 1 12 21Z"/>
    </svg>      
    <svg id="bars" class="cursor-pointer block lg:hidden shrink-0" aria-hidden="true" xmlns="http://www.w3.org/2000/svg" width="28" height="28" fill="none" viewBox="0 0 24 24">
        <path stroke="currentColor" stroke-linecap="round" stroke-width="2" d="M5 7h14M5 12h14M5 17h14"/>
    </svg>
</header>
```

Take note of the classes on the `<header>` element itself. You'll see classes such as `bg-[var(--header-bg-color)]` or `dark:border-[var(--dark-b-color)]`.

Albeit, the syntax is *ugly*. But it does come with its perks. 

You can adjust the colors of the elements on the page using variables instead of having to change the markup for each induvidual element. 

## Logo

To change the logo for your site, simply replace `/static/img/logo.svg` with your logo. There is only one caveat, logos with a large height may shift the navbar in unexpected ways. For this reason, it is reccomended to use a logo which is wide, not tall.

## favicon.ico

During development, the server will look for your favicon.ico at `/static/favicon.ico`. When you go to build your static assest, the favicon will be placed in `/out/favicon.ico`. To change your favicon.ico, simply replace the icon found at `/static/favicon.ico`.