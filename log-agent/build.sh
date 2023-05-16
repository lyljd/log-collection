filename="log-agent"
if [ -e "$filename" ]; then
  rm "$filename"
fi
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "$filename"
echo "build $filename successfully"
