# go-mux-server

## Description

The use of the Application is it, that all HTML files which might be hiding in all subdirectories will be found and listed on to the index page. At first that was it, but after receiving some feedback that some developers might use diffrent directories for different Versions, I added an version Parameter for starting the Application.  

## TechStack

| Technologie | Description                               | Version |
| ----------- | ----------------------------------------- | ------- |
| Go          | Statically typed backend language         | 15.3.   |

I also used the github.com/gorilla/mux repository as server/multiplexer.

## Features

For now the applications "walks through" each subdirectory and saves the path of each HTML or PDF files into seperated list. Thereafter the filenames will be printed as lists on the index page (as links) and each can be opened just through clicking their names.  

## Installation

After downloading the app, run "go build ." to create an executable binary. Thereafter move the binary to where you want to use it. If your do not need the version distinction, I would recommend to remove those conditions before building an executive binary.

## Contact information

For further informations, feedback or similiar you can contact via :

Mail : yosef.hagos@googlemail.com
