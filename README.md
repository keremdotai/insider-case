# insider-case

**Author :** Kerem Avci ([@keremdotmu](https://github.com/keremdotmu) and [@mkeremavci](https://github.com/mkeremavci))


## Table of Contents

1. [Project Overview](#project-overview)
2. [Design Details](#design-details)
3. [Requirements](#requirements)
4. [Environment Variables](#environment-variables)
5. [Build the Application](#build-the-application)
6. [Run the Application](#run-the-application)
7. [Accessing Logs](#accessing-logs)
8. [API Documentation](#api-documentation)


## Project Overview

In this project, it is supposed to implement automatic message sending system. This system must be designed to send two messages every 2 seconds, which are not sent yet. Besides, the system should provide two endpoints which are responsible for starting/stopping automatic message sending and listing all sent messages, respectively.

Messages are generated and inserted into the database for every 0.2 second with random recipent phone number and message content. You can also add an item into PostgreSQL database with external insertion. The data model for an item is given below;

    recipent : string
        Phone number must be in the form of +905xxxxxxxxx, which is length of 14.

    content : string
        The message content with maximum 140 characters.

    sent : bool, default false
        The status indicating whether the message is sent.

    in_queue : bool, default false
        The status indicating whether the message has been in the queue.


A queue provides messages to the sender threads. Its capacity is 100. A thread is responsible for pushing message items that are not sent yet into the queue. If the system has to restart for any reason, we keep the in_queue attribute in the database so that the items which have been in the queue can be added back to the queue. When the system starts working, the queue is first filled with these items.

The project utilizes two sender threads that are responsible for sending messages to the given webhook for every 2 seconds. These threads start running when the system starts. Until a stop request is recieved, these threads continue sending messages.

> [!CAUTION]
> 
> [webhook.site](https://webhook.site/) offers a limited number of requests in its free version.
> After a while, the system will start receiving the error response **"HTTP 429 Too Many Requests"** repeatedly.
>


## Requirements

The softwares and tools required to build and run the project are given below:

- Docker


## Environment Variables

Before building and running the project, you must provide **.env** file inside the main directory of project. You can copy the content of **.env.template** file and replace values with actual ones.

    cp .env.template .env
    nano .env

Descriptions of the environment variables are given in **.env.template** file. Please consider them before assigning values.


## Build the Application

> [!IMPORTANT] 
>
> You must create folders that will be used as volumes for PostgreSQL and Redis databases. If you set `POSTGRES_VOLUME` and `REDIS_VOLUME` variables in **.env** file with proper paths, you can create folders with the following command;
>
>       source .env && mkdir -p ${POSTGRES_VOLUME} ${REDIS_VOLUME}
>
> Otherwise, you should create folders manually.
>

You can build the application with the following command;

    docker compose build


## Run the Application

You can run the application with the following command;

    docker compose up -d


## Accessing Logs

You can access the logs with the following command;

    docker compose logs -f


## API Documentation

While the server is running, you can access the API documentation from `http://localhost:${SERVER_PORT}/swagger/index.html`. Please replace `SERVER_PORT` with the corresponding environment variable.

You can also access the API documentation via [swagger.io](https://editor.swagger.io/), by copying the content of `src/docs/swagger.yaml` into the editor. Do not forget to replace `SERVER_PORT` in **line 9** with the corresponding environment variable.
