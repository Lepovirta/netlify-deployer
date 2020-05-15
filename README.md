# netlify-deployer

Deploy website changes from a local directory to [Netlify](http://netlify.com/).

Netlify CLI is nice, but it's heavy for a lot of the automation tasks I use it for.
I created this tool to implement the parts that I need regularly in deployment automation.

## Installation

Use [Go](https://golang.org/) CLI tools to install:

```
go get github.com/lepovirta/netlify-deployer
```

## Usage

```
netlify-deployer
```

The tool uses these environment variables:

* `NETLIFY_AUTH_TOKEN`: [Authentication token for logging into Netlify](https://docs.netlify.com/cli/get-started/#obtain-a-token-in-the-netlify-ui)
* `NETLIFY_SITE_ID`: [The ID of the site](https://docs.netlify.com/cli/get-started/#link-with-an-environment-variable)
* `NETLIFY_DIRECTORY`: The directory of files to deploy
* `NETLIFY_DRAFT`: Set `true` to make a draft deployment. Set `false` to make a production deployment. Default: `true`.
* `NETLIFY_DEPLOYMESSAGE`: A short message to include in the deployment log
* `NETLIFY_LOGLEVEL`: Log level for application logs. Default: `warn`
* `NETLIFY_LOGFORMAT`: Format of the application logs. Either `text` for text logs or `json` for JSON logs. Default: `text`

The tool prints out the unique URL for the deployment, which you can use to view site.

## Building

Make sure you have [Go 1.13+](https://golang.org/dl/) installed.
After that, run these commands:

```
go mod download
go build
```

This should produce an executable binary called `netlify-deployer`.

## License

MIT License

See [LICENSE](LICENSE) for more details.
