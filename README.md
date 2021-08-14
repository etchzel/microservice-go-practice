# Steps

1. Create database: sekolah, collection name: siswa
2. Install MongoDB for Go: go get go.mongodb.org/mongo-driver/mongo
3. Install HTTP Router GOrilla Mux: go get -u github.com/gorilla/mux
4. Create and start the server: go run server.go

# Setup For Installing Modules

1. Make sure GOPATH is set (check with go env)
2. Make sure GO111MODULE is on (check with go env)
3. If not on, on terminal, set by using:

        set GO111MODULE=on

4. Install the package using go get or go install
5. Create go.mod file using:

        go mod init app-name

6. Run the following command:

        go mod tidy

7. Done
