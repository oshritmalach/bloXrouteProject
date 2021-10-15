# bloXrouteProject

1. docker pull rabbitmq
2. docker run -d --hostname my-rabbit --name some-rabbit rabbitmq:3 

---- run sever:
1. open traminal --> path/bloXrouteProject
2. go run consuming.go 
3. to see logs on run time: tail - f logFile.log

---- run client:
1. open traminal --> path/bloXrouteProject
2. go run main.go 
3. typing options:
   {"action": "AddItem","textMessage":"newItem"}
   {"action": "GetItem","id":2}
   {"action": "RemoveItem","id":1}
   {"action":"GetAllItems"}
   
*You can run several clients
