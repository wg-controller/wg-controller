# build stage
FROM node:lts-alpine as build-stage

ARG IMAGE_TAG
ARG IMAGE_NAME

WORKDIR /app
COPY package*.json ./
RUN npm install --production
COPY . .

# Set env vars in .env file
RUN echo "VITE_APP_IMAGE_TAG=${IMAGE_TAG}" > .env
RUN echo "VITE_APP_IMAGE_NAME=${IMAGE_NAME}" >> .env

RUN npm run build

# production stage
FROM nginx:stable-alpine as production-stage
COPY --from=build-stage /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
COPY default.conf /etc/nginx/conf.d/default.conf
CMD ["nginx", "-g", "daemon off;"]