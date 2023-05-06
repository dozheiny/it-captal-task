# IT-Captal task
IT Captal task is task for focusing on Authentication and role base management. 

## Table of Contents
- [Getting Started](#getting-started) 
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Deployment](#deployment)
- [See API documents](#see-api-documents)
- [Generate Documents](#generate-documents)
- [Built With](#built-with)
- [Contributing](#contributing)
- [License](#license)

## Getting Started
These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

## Prerequisites
To run this project, you will need to have the following installed on your system:

Golang\
MongoDB\
Redis\
Minio\
Docker

## Installation
To get started, first clone the repository:

`git clone https://github.com/dozheiny/it-captal-task.git`


`cd backend-template`

Then, run the following command to install the required dependencies:

`go mod download`

After that, you can set the required environment variables by creating a .env file in the root of the project and adding the necessary variables. You can use the .env.example file as a template.

# Usage
To run the project locally, you can use the following command:


`go run ./... .`
This will start the server and make it accessible at http://localhost:8080.

# Deployment
To deploy the project using Docker, you can use the following commands:

`docker compose up -d --build`

This will build the Docker image and start the container, making the server accessible at http://localhost:8080.

# See API Documents
For see API Documents after making server accessible; you can find documents at http://localhost:8080/swagger/ .

# Generate Documents
To generate a document with swagger, use the following command:

`swag init -g 'cmd/api/main.go' --output docs/swagger --parseDependency`

For more usage, follow the official document.

# Built With
- Golang - Programming language used
- MongoDB - Database used
- Redis - Caching service used
- Docker - Deployment tool used
- Swagger - Document tool used

# Contributing
Please read CONTRIBUTING.md for details on our code of conduct, and the process for submitting pull requests.

# License
This project is licensed under the MIT License; see the LICENSE file for details.