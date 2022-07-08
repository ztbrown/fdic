/bin/echo -e "\e[44;97m Compiling ... \e[0m"
if GOOS=linux GOBIN=`pwd`/bin go install -v ./cmd/fdic; then
  /bin/echo -e "\e[42;97m SUCCESS \e[0m"
else
  /bin/echo -e "\e[101;97m FAILED \e[0m"
  exit 1
fi
