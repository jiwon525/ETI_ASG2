# syntax=docker/dockerfile:1

FROM nginx

# Set destination for COPY
WORKDIR /app

COPY ./Index.html /usr/share/nginx/html
COPY ./ViewStudentList.html /usr/share/nginx/html
COPY ./ViewClassList.html /usr/share/nginx/html
COPY ./UpdateClass.html /usr/share/nginx/html
COPY ./SearchClass.html /usr/share/nginx/html
COPY ./DeleteClass.html /usr/share/nginx/html
COPY ./createclasssheet.css /usr/share/nginx/html
COPY ./class.css /usr/share/nginx/html

# This is for documentation purposes only.
# To actually open the port, runtime parameters
# must be supplied to the docker command.
EXPOSE 80