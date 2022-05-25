# Information for developing

## Local development environment

Change api domain in ```src/store/index.ts``` to the go server address (usually localhost:8080)

Launch dev server with ```yarn serve```

Launch api server with ```go run .``` in api folder

## Building multi architecture and pushing to Docker hub

Create new multiplatform builder

```docker buildx create --name mybuilder```

Activate it

```docker buildx use mybuilder```

Build and push the images

```docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t yellowtech/uptime-shark:latest --push .```