# Installation

## git clone

With git installed, clone the repository into your apps directory:

```bash
mkdir godocument-website
cd godocument-website
git clone https://github.com/phillip-england/godocument .
```

## Serving the Cloned Repo

With go 1.22.0 or greated installed, run:

```bash
go run main.go
```

After serving the application, visit localhost:8080/quickstart/installation and you should see this exact page being served.

## Resetting the Application

<span class='content-warning'>Resetting the application will DELETE all files in `./docs`, `./out`, and `./godocument.config.json`</span>