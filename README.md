# go-ecommerce-with-auth

## Project Description
This is a simple e-commerce project with authentication and authorization. 
The project is built with Go and MongoDB. 

## Developing (Run the project at watch mode)
 ```bash 
    make -f Makefile.linux dev front 
 ```

 ```bash 
    make -f Makefile.linux dev back 
 ```

 For more commands, see Makefile.linux or Makefile.win


## Starting From Scratch (Step by step Installaion)


### 1. Clone the repository

```bash
git clone
```

### 2. Install dependencies

```bash 
go get
```

### 3. Run dependencies
    
```bash
run.sh
```

### 4. Run the application

```bash
go run.sh
```

### 5. Open the browser and go to

```bash
http://localhost:4000
```


## Step by step init new project
- go mod init github.com/username/projectname
- folder named 'cmd' is the entry point of the application ('src' in java and also others)
- make sure the extentions are working properly, see (golang-extensions-guide)[https://medium.com/backend-habit/setting-golang-plugin-on-vscode-for-autocomplete-and-auto-import-30bf5c58138a]


### Installations: 
- chi v5: 
```bash
go get github.com/go-chi/chi/v5
```

## common problems and how to solve them:
-  [go.mod file not found in current directory or any parent directory;](https://stackoverflow.com/questions/66894200/error-message-go-go-mod-file-not-found-in-current-directory-or-any-parent-dire)

verify go mod is enabled:
```bash
go env -w GO111MODULE=auto

```
- [how-to-solve-stderr-go-mod-tidy-go-mod-file-indicates-go-...](https://stackoverflow.com/questions/71881727/how-to-solve-stderr-go-mod-tidy-go-mod-file-indicates-go-1-17-but-maximum-su?rq=2)



## Docs I Found useful for this course:
maybe implement - https://stackoverflow.com/questions/21151765/cannot-unmarshal-string-into-go-value-of-type-int64
