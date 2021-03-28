# netlify-deployer

Deploy website changes from a local directory to [Netlify](http://netlify.com/).

Netlify CLI is nice, but it's heavy for a lot of the automation tasks I use it for.
I created this tool to implement the parts that I need regularly in deployment automation.

## Installation

Use [Go](https://golang.org/) CLI tools to install:

```
go get gitlab.com/lepovirta/netlify-deployer
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

## Docker images

Docker images for this tool are hosted in the [Gitlab Container Registry](https://gitlab.com/lepovirta/netlify-deployer/container_registry).
This is the name of the image that is available

```
registry.gitlab.com/lepovirta/netlify-deployer
```

The image has the following tags available:

* `latest`: Latest version from the master branch with just the netlify-deployer binary included (i.e. `Dockerfile.minimal`)
* `ci`: Latest version from the master branch with the netlify-deployer binary and additional tools useful with CI platforms such as Gitlab CI (i.e. `Dockerfile.ci`)

## Integration with Gitlab CI

Here's how you can deploy your site from Gitlab CI to Netlify using the `netlify-deployer` Docker image.
This integration does the following:

* Publish your site when running the pipeline on master branch
* Publish a draft site when running the pipeline on any other branch
* Post a link to your merge requests when the draft site is available (requires a Gitlab access token)

First, set up your Netlify access credentials in your repository's CI/CD settings:

1. Go to your repository page in Gitlab
2. Go to `Settings` > `CI / CD` > `Variables`
3. Add a new variable `NETLIFY_SITE_ID` and use your Netlify site's ID as the value
4. Add a new variable `NETLIFY_AUTH_TOKEN` and use your [Netlify access token](https://docs.netlify.com/cli/get-started/#obtain-a-token-in-the-netlify-ui) as the value
5. (Optional) If you want to post draft links to your merge requests, you need to also add a new variable `GITLAB_ACCESS_TOKEN` and use a Gitlab access token as the value

Next, your `.gitlab-ci.yml` should be made to look something like this:

```yaml
stages:
- # Other stages go here
- build-site
- publish-site

variables:
  # You can use whatever directory here you want.
  NETLIFY_DIRECTORY: public

build-site:
  stage: build-site
  script:
  - # Use whatever commands you need here to generate your site.
  - # Place the site to $NETLIFY_DIRECTORY
  artifacts:
    paths:
    - $NETLIFY_DIRECTORY

# You can run this same job in both master and MR branches.
# It will automatically publish the site as a draft when run on non-master branches.
# If you want to publish your site on another branch other than master,
# set the variable NETLIFY_MAIN_BRANCH to point to some other branch.
publish-site:
  stage: publish-site
  image: registry.gitlab.com/lepovirta/netlify-deployer:ci
  script:
  - gitlab-deploy-site
```

For a more detailed guide, see this blog post: [Website previews for custom Netlify deployments using GitLab CI](https://lepovirta.org/posts/2021-03-28-website-previews-for-custom-netlify-deployments-using-gitlab-ci.html)

## Building

Make sure you have [Go 1.13+](https://golang.org/dl/) installed.
After that, run these commands:

```
go mod download
go build
```

This should produce an executable binary called `netlify-deployer`.

After you've built the app, you can build your own Docker image using the Dockerfiles in the repository:

```
docker build -f Dockerfile.ci -t netlify-deployer-ci .
docker build -f Dockerfile.minimal -t netlify-deployer .
```

## License

MIT License

See [LICENSE](LICENSE) for more details.
