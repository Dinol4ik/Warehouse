В данном разделе кратко описаны использованные непосредственно(не indirect) в проекте сторонние пакеты.
***

## github.com/jmoiron/sqlx

    sqlx - это библиотека, которая предоставляет набор расширений для стандартной библиотеки
    go database/sql. Версии sql.DB для sqlx, sql.TX, sql.Stmt и др. все они оставляют базовые интерфейсы нетронутыми,
    так что их интерфейсы являются надмножеством стандартных. Особенно в данном пакете выделяется возможность присваивать
    результаты запросов соответствующим структурам, не вынуждая итерироваться по пришедшим строкам бд.

### ★ Stars ★  13.8k
### 유 Used by 27.6k 유
****

## github.com/DATA-DOG/go-sqlmock

    sqlmock - это макетная библиотека, реализующая sql/driver. Который имеет одну-единственную
    цель - имитировать любое поведение драйвера sql в тестах, не требуя реального подключения к базе данных.

### ★ Stars  5.3k ★
### 유 Used by 14.3k 유
***

## github.com/gorilla/mux

    Пакет gorilla/mux реализует маршрутизатор запросов и диспетчер для сопоставления входящих
    запросов с их соответствующим обработчиком. 

### ★ Stars ★  18.2k
### 유 Used by 116k 유
***

## github.com/lib/pq

    Пакет является драйвером над database/sql для подключения к базе данных postgres. 

### ★ Stars ★  8.1k
### 유 Used by 117k 유
***

## github.com/spf13/viper

    Viper - это комплексное конфигурационное решение для приложений Go. Он предназначен для работы в 
    рамках приложения и может обрабатывать все типы конфигурационных требований и форматов. Он поддерживает:
        установку стандартных значений, чтение из большинства подходящих под конфиг файлов - JSON, YAML 
        и также чтение из переменных окружения. 

### ★ Stars ★  23.3k
### 유 Used by 102k 유
***

## github.com/stretchr/testify

    testify - набор пакетов, которые предоставляют множество инструментов для подтверждения того, что ваш код
    будет вести себя так, как вы задумали. Поддерживает создание моков и более удобные инструменты сверки ожидания с выводом

### ★ Stars ★  20k
### 유 Used by 424k 유
***


## go.uber.org/zap

    Невероятно быстрый, структурированный, многоуровневый инструмент логирования в Go. Пакет имеет множество возможностей
    конфигурации логера, облегченную (строго структурированную) версию логера, а так же 'засахаренную версию',
    которая все равно быстрее в 4 - 10 раз чем другие пакеты

### ★ Stars ★  19k
***