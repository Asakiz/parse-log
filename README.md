# parse-log

If you are a linux user could use the Makefile instead the commands. Type ``make`` or ``make help`` on terminal to see the commands available.

## Build and Up the mongoDB container

    docker-compose up --build -d
        
## Build and Run the application
  Need to be on source code folder, with in that case is ``api/``

    go run main.go <input-file>
