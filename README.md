# Tagifai

Automatically generates tags by analysing your uploaded photo using the [Clarifia](https://clarifai.com/) image classification models. Just a simple experiment on working with 3rd-party REST APIs in Go.

## Requirements

* go1.10.3
* dep v0.5.0
* A free Clarifai account with an API_KEY

## Local Setup

	cd $GOPATH/src
	git clone <repo>
	cd <project>
	dep ensure
	go run main.go

## Production Deployment

This app is configured to run on Heroku.

* You need the Heroku CLI to interact with the service. `$ brew install heroku/brew/heroku`
* Login with your credentials: `$ heroku login`
* Navigate to your Go project repo
	* `heroku create <app-name>`: creates a remote repo in Heroku and generates a domain name pointing to the `master` branch.
	* Write your Go code as usual.
	* When ready, commit your changes and push to Heroku: `$ git push heroku <localbranch>:master`

### Environment Variables

This app requires the following envars. In Heroku, you manually enter them as *Config Vars* in the app's dashboard, but on your local machine, you save them in a `.env` file at the project root directory. 

> Treat these as secrets - **you must never commit them source control.**

	DEPLOYMENT = local|staging|production
	CLARIFAI_API_KEY = xxxxxxxxxxxxxxxxx
	IMG_ANALYSIS_ENDPOINT = https://url/to/clarifai/model

## Developer Notes

### Tools Used

* `dep`: _Go_ dependency managment tool. Just like _npm_ for Node.js, dependencies are specified in a `Gopkg.toml`, a versioned snapshot is saved in `Gopkg.lock`, and a `/vendor` directory holds the actual source code.

#### Go on Heroku

* Static Assets & Files included in your git repository are available at runtime and can be served from the local filesystem using Go’s http.FileServer or framework equivalent. 

* In Heroku, builds use the contents of the vendor/ directory when compiling as if they are in your $GOPATH

#### Using `dep`

	brew install dep 

	dep init

	dep ensure -add <dep>
	dep ensure -add <dep1> <depN>

	dep status
	dep check

	dep ensure -update
	dep ensure -update <dep>

