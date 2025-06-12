# Задание 1

# Context diagram (AS IS)

```puml
@startuml
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Context.puml

title Smart Home - Context Diagram (AS IS)

Person(homeowner, "Пользователь умного дома", "Житель дома, управляет отоплением и мониторит температуру")
Person(installer, "Специалист-установщик", "Технический специалист, выполняет установку и подключение системы отопления к датчикам")

System(smartHomeMonolith, "Smart Home", "Система управления отоплением и мониторинга температуры")

System_Ext(sensors, "Датчики", "Физические устройства, установленные в помещениях дома: датчики температуры и т.д.")
System_Ext(heatingDevices, "Системы отопления", "Физические устройства отопления: котлы, радиаторы, термостаты")
System_Ext(webBrowser, "Веб-браузер", "Фронтенд-интерфейс пользователя")
System_Ext(installerWebBrowser, "Веб-браузер", "Фронтенд-интерфейс установщика")

' Взаимодействие пользователя
Rel(homeowner, webBrowser, "Использует веб-интерфейс")
Rel(webBrowser, smartHomeMonolith, "HTTP запросы", "REST API")

Rel(installer, installerWebBrowser, "Использует веб-интерфейс установщика для регистрации датчиков")
Rel(installerWebBrowser, smartHomeMonolith, "HTTP запросы", "REST API")

' Взаимодействие установщика
Rel(installer, sensors, "Устанавливает и настраивает датчики", "Физическое подключение")
Rel(installer, heatingDevices, "Подключает системы отопления", "Физическое подключение")

' Взаимодействие системы с устройствами
Rel(smartHomeMonolith, sensors, "Запрашивает данные о температуре", "API устройств")
Rel(smartHomeMonolith, heatingDevices, "Управляет отоплением", "API устройств")

@enduml
```

# Домены

## 1. Домен "Управление Устройствами" (Device Management)

Ответственность:

- Прямое управление системами отопления
- Синхронная отправка команд к устройствам
- Контроль состояния нагревательных устройств

## 2. Домен "Мониторинг Температуры" (Temperature Monitoring)

Ответственность:

- Синхронные запросы к датчикам температуры
- Получение и хранение показаний температуры
- Предоставление температурных данных пользователям

## 3. Домен "Управление Установкой" (Installation Management)

Ответственность:

- Регистрация новых датчиков и устройств в системе
- Управление процессом профессиональной установки
- Конфигурация подключенных устройств

## Context diagram (TO BE)

```puml
@startuml
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Context.puml

title Smart Home - Context Diagram (TO BE)

Person(homeowner, "Пользователь", "Владелец умного дома, самостоятельно управляет устройствами, настраивает сценарии и просматривает телеметрию")

System(smartHomePlatform, "Smart Home", "Система управления отоплением и мониторинга температуры")

System_Ext(webApp, "Веб-приложение", "Веб-интерфейс для управления умным домом")

System_Ext(sensors, "Датчики", "Физические устройства, установленные в помещениях дома: датчики температуры и т.д.")
System_Ext(heatingDevices, "Системы отопления", "Физические устройства отопления: котлы, радиаторы, термостаты")

System_Ext(partnerEcosystem, "Партнерская экосистема", "API и SDK партнеров для интеграции устройств по стандартным протоколам")

' Взаимодействие пользователя
Rel(homeowner, webApp, "Управляет умным домом")
Rel(webApp, smartHomePlatform, "REST API / WebSocket", "HTTPS/WSS")

' Взаимодействие с устройствами (самостоятельное подключение)
Rel(homeowner, sensors, "Самостоятельно устанавливает и настраивает")
Rel(homeowner, heatingDevices, "Самостоятельно устанавливает и настраивает")

' Интеграция платформы с устройствами
Rel(smartHomePlatform, sensors, "Управление и телеметрия", "API устройств")
Rel(smartHomePlatform, heatingDevices, "Управление и телеметрия", "API устройств")

' Партнерская экосистема
Rel(smartHomePlatform, partnerEcosystem, "Интеграция устройств", "REST API/SDK")
Rel(partnerEcosystem, heatingDevices, "Поддерживает устройства")

@enduml
```
