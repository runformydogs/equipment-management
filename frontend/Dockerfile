FROM nginx:alpine

COPY . /usr/share/nginx/html
RUN apk add --no-cache sed

RUN sed -i "s|const API_BASE_URL = 'http://localhost:8080'|const API_BASE_URL = window.location.origin + '/api'|" /usr/share/nginx/html/js/api.js

COPY nginx.conf /etc/nginx/conf.d/default.conf