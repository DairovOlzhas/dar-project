###Project "Tanks Game"

##### Description
Игра в танки. Пользователь подключается к серверу и ждет пока кто то еще подключиться к серверу, 
который тоже должен будет подключиться к серверу. То есть каждая пара подключенная к серверу будет играть между собой.

##### Rule of game



##### Design 

###### Backend side:
1. Game logic
    1. Start menu
    2. Game process
        1. Loading map
        2. Creating tanks
        3. Listening for commands from frontend
    3. End of game


###### Frontend side:
1) Read pressing buttons from keyboard, then send to server
2) Request for image or get response from server without request
3) Draw getted image
4) token for identification users