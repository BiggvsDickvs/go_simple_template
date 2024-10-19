# Important stuff
this project uses sqlite database that will be automatically created in root directory. if you want something else, f.e. postgresql, you will have to use different driver and add db credentials to .env file.
### FIRST OF ALL
create a .env file in root directory
then add following:
```
JWT_KEY=<your jwt key>
```
### INSTALL ALL DEPENDENCIES
simply run 
```
go mod tidy
```
### RUN PROJECT 
simply run
```
go run cmd/app/main.go
```
enjoy