FROM node:20-alpine AS builder
LABEL authors="daniilkolesnik"
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build  # this creates /app/build

FROM alpine:latest AS frontend-builder
WORKDIR /app
COPY --from=builder /app/build ./build
