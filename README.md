# drew-fetch-receipts-processor

Apologies if this go code isn't completely idomatic, I decided to try using go to do this, which was fun :), but it was my first time using it.

## How to run

cd into `src` and run `go run .`

### Notes on running

- Defaults to running on `localhost:8080`
  - This can be changed by adding a `.env` file in the root of the project and setting `LISTEN_ADDR` to whatever address you want

## How to test

cd into `src` and run `go test .`

## Important Notes

- Core logic is in [src/app.go](src/app.go)

- I used [openapi-generator](https://openapi-generator.tech/) to generate an api interface,
  which I implement in [src/app.go](src/app.go)

- Everything in [server/openapi](server/openapi/) is generated
  with the exception of some manually added field validation code that the generator failed to generate.
  - Specifically [model_receipt.go](server/openapi/model_receipt.go) and [model_item.go](server/openapi/model_item.go): The sections that have been edited are delineated clearly by comments

### Other notes

- Is using openapi generator a bit overkill for a service with 2 endpoints? Yeah probably. I wanted to tinker with the go version to see what it spits out; I've used the java and kotlin generators before and was curious how to use it in go :)
