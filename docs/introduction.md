# Introduction

## Quick and Effortless

Godocument makes building your documentation website as easy as possible. All you have to do is clone the repo, write your markdown files, and setup your config. We'll handle the rest. This tool was built with developers in mind. The main driving force is *simplicity.* Sure, we might not have all the bells and whistles you'd find in other tools, but if getting it done quick with the least amount of pain as possible sounds good to you, I think you'll feel right at home.

## How Does it Work?

Godocument takes your godocument.json.config file and uses it to generate a series of endpoints. Here is an example of what the config might look like:


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