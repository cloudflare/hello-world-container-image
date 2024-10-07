# Cloudflare Container Example

This repo contains a simple example application which can be run in Cloudflare's
container runtime.

The example app serves a single HTML page with an incrementing visit counter and a
toggleable health status. The health status can be modified by sending a POST
request to `/health` with the `status` query parameter set to either `healthy`
or `unhealthy`. The health status can be queried by sending a GET request to
`/health`: if the application is healthy it will return an empty response with
a 200 status code; otherwise, it will send an empty response with a 503 status
code.
