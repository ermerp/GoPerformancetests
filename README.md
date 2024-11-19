# Readme

## Build

Navigate to the projects root folder (/GoPerformancetests)

### Windows (Git Bash)
```
./build.sh
```
Executables:
```
.\bank\bankGo.exe
.\mergesort\mergesortGo.exe
```

### Unix
```
chmod +x build.sh
./build.sh
```
Executables:
```
./bank/bankGo
./mergesort/mergesortGo
```

## Run

### Bank
#### var:
```
$env:BANK_IMPLEMENTATION = "sql"

or

$env:BANK_IMPLEMENTATION = "postgrest"
```
#### run:
```
go run .
```

## Driver

pgx:
```
go get github.com/jackc/pgx/v4
go get github.com/jackc/pgx/v4/pgxpool
```


