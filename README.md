# parse-log

A simple program to extract all the data on a log file, store on a database and make some calculations with this data.

## Build and Up the mongoDB container

    docker-compose up --build -d
        
## Build and Run the application
  Need to be on source code folder, with in that case is ``api/``

    go run main.go <file-path>

## Linux user/shell

    make all input=<file-path>
    
---
**NOTE**

Type ``make`` or ``make help`` to see all the commands.

---