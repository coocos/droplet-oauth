FROM node:18.0.0 AS frontend

WORKDIR /dist
COPY package.json package-lock.json tailwind.config.js ./
RUN npm ci
COPY static ./static
COPY templates ./templates
RUN npx tailwind -i ./static/input.css -o ./static/output.css --minify

FROM golang:1.18.2-alpine AS backend

WORKDIR /app
COPY go.mod go.sum main.go ./
COPY server ./server
RUN go build -o droplet

FROM alpine:3.16.0

USER guest
WORKDIR /droplet-oauth
COPY --from=frontend --chown=guest:guest /dist/static/ ./static/
COPY --from=backend --chown=guest:guest /app/droplet ./
COPY --chown=guest:guest ./templates ./templates

EXPOSE 8000
ENTRYPOINT ["./droplet"]
