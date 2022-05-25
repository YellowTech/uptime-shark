# build stage for vue
FROM node:lts as frontend-build-stage
WORKDIR /app
COPY package.json yarn.lock ./
ENV NODE_ENV=production
RUN yarn
COPY . .
RUN yarn build
RUN rm -f ./dist/js/*.map

# build stage for go
FROM golang:1.17-alpine as backend-build-stage
RUN apk add gcc musl-dev

WORKDIR /app
COPY ./api ./

RUN go mod download
RUN cd ./ent && go run -mod=mod entgo.io/ent/cmd/ent generate ./schema && cd /app && chmod +x ./cleanOmit.sh && ./cleanOmit.sh

RUN go build -o ./apiapp

# production stage
FROM alpine as final
WORKDIR /app
COPY --from=frontend-build-stage /app/dist /app/dist
COPY --from=backend-build-stage /app/apiapp /app/apiapp

EXPOSE 8080

CMD [ "./apiapp" ]