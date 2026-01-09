# Start the templates + server with live reloading present at
# :7331
templ generate --watch --proxy="http://127.0.0.1:8081" --cmd="go run ."
