# json-shortlink

Tiny personal short link resolver written in Go and providing Docker integration.

Features live reloading of link configuration without restarting service.

## Usage

1. Create .env based on .env.dist
2. Configure _config/links.json
3. Run via `docker compose up -d`

> Alternatively the application can be built and will run without docker with regular go commands.
> The binary only requires the _config/links.json and .env as the HTML and CSS are included within the compiled output itself.

## Example configuration

Redirects from `your.site/example` to `https://example.com`
and `your.site/tier/kadse` to `https://kadse.io`.

```json
{
  "/example": "https://example.com",
  "/tier/kadse": "https://kadse.io"
}
```
