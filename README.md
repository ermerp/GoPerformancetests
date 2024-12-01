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

# Docker

## Bank: benchmark

### build Container
```
docker build -t bank-go .
```

### start Container
```
docker-compose up
```

### reset Container
```
docker-compose down -v
```

# Driver

pgx:
```
go get github.com/jackc/pgx/v4
go get github.com/jackc/pgx/v4/pgxpool
```


