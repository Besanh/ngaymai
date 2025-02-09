# NGAYMAI PROBLEM ASSESSMENT
All rights reserved © AnhLe
Please contact anhle3532@gmail.com if you want to use AnhLe's source code

## 1. Architecture and entities
[Architecture](https://app.diagrams.net/#G1GWZHzJ0_Ez5r0NW_qkkRE_Vi-52d-SN4#%7B%22pageId%22%3A%22LIlMXQeiWCyiD9asdB7_%22%7D)

## 2. Install library swagger

`go get -u github.com/swaggo/gin-swagger`

`go get -u github.com/swaggo/files`

`go install github.com/swaggo/swag/cmd/swag@latest`

Init folder docs

`swag init`

## 3. Create file .env similar to .env_example, then putting appropriate configuration
## 4. Build service use docker
I have already created a Dockerfile, you can take advantage of this file to build and deploy in docker

# 5. Link docs api
{{domain}}/swagger/index.html