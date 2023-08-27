# GoTH Stack

<div style="display:flex">
<img src="static/img/go.png" style="width:31%;margin-right:3.5%" /><img src="static/img/htmx.jpg" style="width:31%;margin-right:3.5%" /><img src="static/img/tailwindcss.svg" style="width:31%" />
</div>

## Getting started
1. Clone the repo
1. Run `go mod tidy`
1. Fine and replace `github.com/tomdoestech/goth` with your own module name
1. Download `Air` - https://github.com/cosmtrek/air
1. Run `air` to start the dev server

## Contributing
Contributions are welcome. Please open an issue to discuss your idea before opening a PR unless it's a fix for a bug or a typo.

## Principles
1. Simplicity
1. Speed 
1. Interactivity
1. Security

## Included features
The included features are minimal because this is just a starting point. If you think there is a feature missing that all projects need, please open an issue.
1. Register an account with email and password
1. Login with email and password
1. Logout

## Templates
The templates are written in [Go Templates](https://pkg.go.dev/text/template). The templates are located in the `templates` directory. The `templates/base` template is the base template that all other templates extend. The `templates/partial` directory contains partial templates that are included in other templates.

## Styles
The tailwindcss executable is for linux x64. If your system requires a different executable, please following this guide: https://tailwindcss.com/blog/standalone-cli

Basically all you need to do is download the executable from [here](https://github.com/tailwindlabs/tailwindcss/releases/tag/v3.3.3), mode it so it's runnable and then optionally rename it to `tailwindcss`.

Generate a new style.css file and watch for changes:
```bash
./tailwindcss -i static/css/input.css -o static/css/style.css --watch
```
Generate a new minified style.css file for production:
```bash
./tailwindcss -i static/css/input.css -o static/css/style.css --minify
```

## Security
The header partial `templates/partial/header` includes a meta tag with the CSP headers. The included shad256 hash is because HTMX will create an inline style tag, so this hash is for that tag. 

Please also read the [HTMX security guide](https://htmx.org/docs/security/).

## Metrics
By default there is a metrics server that starts at `http://localhost:9100/metrics` and serves prometheus metrics.