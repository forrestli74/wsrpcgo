curl --insecure \
     --request "POST" \
     --location "https://localhost:8080/twirp/tmp.RoomService/AddWriter" \
     --header "Content-Type:application/json" \
     --data '{"proposed_ids": ["a"]}' \
     --verbose
