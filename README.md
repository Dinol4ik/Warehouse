# WarehouseAPI

Данный проект представляет собой JSON RPC API, реализующее возможность взаимодействия со складами товаров.
#### !!! Для запуска проекта в корневой папке должен лежать файл config.yml!!!
Пример:

    logger:
    level: -1
    
    httpServer:
    Address: '0.0.0.0:8080'
    
    postgres:
    host: 'postgres'
    user: 'user'
    dbname: 'lamoda'
    password: 'postgres'
## Установка

    go mod download

## Старт приложения

     go run .\cmd\api\main.go    

## Применение миграций
     Для миграций был использован инструмент goose. Чтобы его использовать рекомендуется установить переменные среды:
           GOOSE_DBSTRING=host=0.0.0.0 user=postgres password=postgres dbname=postgres sslmode=disable (с вашими данными бд)
           GOOSE_DRIVER=postgres
           
     Далее можно применять основные команды:
    
           goose up (применение миграций)
           goose down (откат последней миграции)
           gooose create <название_миграции> sql (создание новой миграции)
           goose status (просмотр примененных миграциий)

    *Ленивый запуск - goose postgres "user=user dbname=lamoda password=postgres sslmode=disable" up or down or create ..
    ** Для тестирования сервиса нужно выполнить миграции с автозаполнением таблиц команда выше

## Запуск тестов

    go test .\internal\controller\. -v

## Запуск в Docker через Make

     Make up


## Запуск клиента для тестирования сервиса

    go run .\cmd\client\client.go

# JSON RPC

Далее идет описание методов, реализованных сервисом WarehouseService. Вызовы методов производились другим сервисом,
использующим стандартный клиент из пакета net/rpc. В args1 кладется json с данными для вызова метода, результат
вызова которого кладется в переменную reply, так же представляющую из себя json

## ReserveItemHandler

### client.Call("WarehouseService.ReserveItemHandler", args1, &reply)

    args1 = "{"items":[{"code":"MP002XW08D8E","amount":10}]}"

### reply

    "{"successes":["MP002XW08D8E"],"errors":null}"


## UnReserveItemHandler

### client.Call("WarehouseService.UnReserveItemHandler", args1, &reply)

    args1 = "{"items":[{"code":"MP002XW08D8E","amount":8}]}"

### reply

    "{"successes":["MP002XW08D8E"],"errors":null}"

## FetchWarehouseItemsHandler

### client.Call("WarehouseService.FetchWarehouseItemsHandler", args1, &reply)

    args1 = "{"id":"88802b69-2f8a-4f15-9f27-70b0c7446121"}"

### reply

    "{"items":[{"name":"Водолазка oodji","size":"S","code":"MP002XW09XTH","amount":25,"reserved":13},
               {"name":"Лоферы Pazolini","size":"XL","code":"MP002XW0YOGZ","amount":123,"reserved":100}"

## FetchItemsByCodesHandler

### client.Call("WarehouseService.FetchItemsByCodesHandler", args1, &reply)

    args1 = "{"codes":["MP002XM0AYK1"]}"

### reply
    "{"items":[{"name":"Кроссовки Big Triangle 2.0","size":"36","code":"MP002XM0AYK1","amount":20,"reserved":0,"warehouseId":"fb2ab07f-95c2-4013-8de6-efe7649ecd50","warehouseAvailable":true}]}"
    