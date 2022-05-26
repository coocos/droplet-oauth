# droplet-oauth

A simple example of using Go & DigitalOcean API via OAuth to list your droplets.

![Screenshot](/docs/screenshot.png)

## Usage

The application uses [DigitalOcean OAuth API](https://docs.digitalocean.com/reference/api/oauth-api/) to request an access token used to list droplets, so you need to register an OAuth application via the DigitalOcean control panel. Once you have created one, export its details and a random session key as environment variables (or alternatively use an `.env` file):

```conf
CLIENT_SECRET=client-secret-here
CLIENT_ID=client-id-here
SESSION_KEY=crypto-random-string-here
REDIRECT_URL=https://droplet.yourdomain.dev/redirect
```

Finally, start the stack:

```
docker compose up -d
```
