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

- Entry point for program is [main.go](src/main.go)
- Application is setup in [app.go](src/app.go)
- Core logic is in [src/receipts_service.go](src/receipts_service.go)
- I used [openapi-generator](https://openapi-generator.tech/) to generate an api interface,
  which I implement in [src/receipts_service.go](src/receipts_service.go)

- Everything in [server/openapi](server/openapi/) is generated
  with the exception of some manually added field validation code that the generator failed to generate.
  - Specifically [model_receipt.go](server/openapi/model_receipt.go) and [model_item.go](server/openapi/model_item.go): The sections that have been edited are delineated clearly by comments

### Other notes

- Is using openapi generator a bit overkill for a service with 2 endpoints? Yeah probably. I wanted to tinker with the go version to see what it spits out; I've used the java and kotlin generators before and was curious how to use it in go :)
