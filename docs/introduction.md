# Introduction

## Static Site Generator

Godocument is a static site generator built with *simplicity* in mind. My goal is to make it as quick and easy as possible to get a documentation website up and running. 

Not every developer is a web developer. But every developer who writes open-source software will need to document their code. Godocument makes that process painless. Why spend all your time coding up your documentation when you could spend it improving the project your documenting?

## Why?

This project was originally built for myself. While learning [Templ](https://templ.guide), I came across [Docusaurus](https://docusaurus.io/docs), another static site generator. **Docusaurus inspired this project.** I thought to myself, "Docusaurus looks great, but I prefer to avoid React and Node if I can." 

Why make something complex if it doesn't have to be? Static site generation shouldn't require you to be locked into the React world. What if there was a way to get most of the benefits of Docusaurus, but without all the framework madness?

And Godocument was born.

### Testing

Godocument's workflow is straightforward: clone the repo, setup your config, write your markdown files, build the static assets, and deploy.

Everything starts with your godocument.config.json file. Here is what the config file looks like for this site. Yes, I use my own software. 😉

#### Again

Godocument's workflow is straightforward: clone the repo, setup your config, write your markdown files, build the static assets, and deploy.

Everything starts with your godocument.config.json file. Here is what the config file looks like for this site. Yes, I use my own software. 😉

##### Again

Godocument's workflow is straightforward: clone the repo, setup your config, write your markdown files, build the static assets, and deploy.

Everything starts with your godocument.config.json file. Here is what the config file looks like for this site. Yes, I use my own software. 😉

## How Does it Work?

Godocument's workflow is straightforward: clone the repo, setup your config, write your markdown files, build the static assets, and deploy.

Everything starts with your godocument.config.json file. Here is what the config file looks like for this site. Yes, I use my own software. 😉

```json
{
    "docs": {
        "Introduction": "/introduction.md",
        "Getting Started": {
            "Installation": "/getting-started/installation.md",
            "Usage": "/getting-started/usage.md",
            "More Information": {
                "Configuration": "/getting-started/more-information/configuration.md"
            }
        },
        "Learn More": {
            "Advanced Usage": "/learn-more/advanced-usage.md",
            "API Reference": "/learn-more/api-reference.md",
            "Contributing": "/learn-more/contributing.md",
            "FAQ": "/learn-more/faq.md"
        }
    }
}
```

```css
html.dark pre[class*="language-"] {
	padding: 1em;
	margin: .5em 0;
	overflow: auto;
}
```

```go
func writeNavbarHTML(html string) {
	f, err := os.Create(config.GeneratedNavPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.WriteString("<!-- This file is auto-generated. Do not modify. -->\n")
	if err != nil {
		panic(err)
	}
	_, err = f.WriteString(html)
	if err != nil {
		panic(err)
	}
}
```

Godocument takes this config file and uses it to build a series of routes. You can use these routes during development to test your application. Once you get things looking how you want, you can read the static html from the routes and generate an ./out directory containing your static assets.

From there, you can deploy your application on a CDN of your choice. If that sounds like a plan, keep reading.