###Project "Tanks Game"

##### Description
Игра в танки. Пользователь подключается к серверу и ждет пока кто то еще подключиться к серверу, 
который тоже должен будет подключиться к серверу. То есть каждая пара подключенная к серверу будет играть между собой.

 Design [Board](https://miro.com/welcomeonboard/OWLgM1ynjkgtlGaWwsf4BtbBBZAQRoz2vBgMBjpnB5Wzm0RiYAtHOZ6vsSBApGCi)

### todo

1. проверка на онлайн. каждый проверяет одного. очередь для проверки.
 
    **Онлайн юзер**:
    1. проверяет одного юзера.
    2. если проверяющийся юзер вышел из сети
        1. посылает команду для его удаления
        2. делает Ack
        3. Получает другого юзера на проверку
            Другой юзер это:
            1. Юзер который проверялся  юзером который вышел из сети
            2. Он попадает обратно в очередь 
            3. И попадает к онлайн юзеру
