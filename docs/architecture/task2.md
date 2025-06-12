# Задание 2

# Containers diagram

```puml
@startuml
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Container.puml

title Smart Home - Containers Diagram (TO BE)

Person(user, "Пользователь", "Владелец умного дома")
System_Ext(partnerApis, "Partner APIs", "API партнеров для интеграции устройств")

System_Boundary(platform, "Smart Home Platform") {
    Container(webApp, "Web Application", "SPA", "Веб-интерфейс для управления умным домом")

    Container(apiGateway, "API Gateway", "Nginx", "Единая точка входа, маршрутизация, аутентификация, rate limiting")

    Container(authService, "Auth Service", "Go/Node.js", "Управление пользователями, подписки SaaS, аутентификация")
    Container(deviceService, "Device Service", "Go/Node.js", "Управление всеми типами устройств, команды управления")
    Container(scenarioService, "Scenario Service", "Go/Node.js", "Пользовательские сценарии автоматизации, правила")
    Container(telemetryService, "Telemetry Service", "Go/Node.js", "Сбор, обработка и хранение телеметрии устройств")
    Container(integrationService, "Integration Service", "Go/Node.js", "Интеграция с партнерскими API и устройствами")
    Container(notificationService, "Notification Service", "Go/Node.js", "Push-уведомления, email, SMS алерты")

    ContainerDb(userDb, "User Database", "PostgreSQL", "Пользователи, подписки, биллинг")
    ContainerDb(deviceDb, "Device Database", "PostgreSQL", "Метаданные устройств, состояния, конфигурации")
    ContainerDb(telemetryDb, "Telemetry Database", "PostgreSQL", "База данных для хранения телеметрии")
    ContainerDb(scenarioDb, "Scenario Database", "PostgreSQL", "Сценарии автоматизации, правила")

    Container(messageQueue, "Message Queue", "Apache Kafka", "Event streaming между микросервисами")
    Container(realTimeQueue, "Real-time Queue", "Redis Pub/Sub", "Real-time уведомления и WebSocket коммуникация")
}

' Пользовательские интерфейсы
Rel(user, webApp, "Использует", "HTTPS")
Rel(webApp, apiGateway, "API calls", "HTTPS/WSS")

' API Gateway маршрутизация
Rel(apiGateway, authService, "Маршрутизация запросов", "HTTP/gRPC")
Rel(apiGateway, deviceService, "Маршрутизация запросов", "HTTP/gRPC")
Rel(apiGateway, scenarioService, "Маршрутизация запросов", "HTTP/gRPC")
Rel(apiGateway, telemetryService, "Маршрутизация запросов", "HTTP/gRPC")
Rel(apiGateway, notificationService, "Маршрутизация запросов", "HTTP/gRPC")

' Микросервисы и базы данных
Rel(authService, userDb, "Запись/чтение", "SQL")
Rel(deviceService, deviceDb, "Запись/чтение", "SQL")
Rel(telemetryService, telemetryDb, "Запись", "SQL")
Rel(scenarioService, scenarioDb, "Запись/чтение", "SQL")

' Event streaming
Rel(deviceService, messageQueue, "Device events", "Kafka Protocol")
Rel(scenarioService, messageQueue, "Scenario events", "Kafka Protocol")
Rel(telemetryService, messageQueue, "Telemetry events", "Kafka Protocol")
Rel(notificationService, messageQueue, "Consumes events", "Kafka Protocol")

' Real-time коммуникация
Rel(deviceService, realTimeQueue, "Real-time updates", "Redis Pub/Sub")
Rel(notificationService, realTimeQueue, "Push notifications", "Redis Pub/Sub")
Rel(apiGateway, realTimeQueue, "WebSocket updates", "Redis Pub/Sub")

' Внешние интеграции
Rel(integrationService, partnerApis, "Integration", "HTTP/MQTT")
Rel(deviceService, integrationService, "Partner devices", "gRPC")

@enduml
```

# Микросервисы

## 1. Auth Service

**Домен:** Управление пользователями и подписками

**Ответственность:** Аутентификация, авторизация, SaaS подписки

```puml
@startuml
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Component.puml

title Auth Service - Components Diagram (TO BE)

Person(user, "Пользователь", "Владелец умного дома")
Container(webApp, "Web Application", "SPA", "Веб-интерфейс")
Container(apiGateway, "API Gateway", "Nginx", "Единая точка входа")

Container_Boundary(authService, "Auth Service") {
    Component(authController, "Authentication Controller", "Go/Node.js", "REST API для аутентификации, login/logout")
    Component(userController, "User Management Controller", "Go/Node.js", "REST API для управления профилем пользователя")
    Component(subscriptionController, "Subscription Controller", "Go/Node.js", "REST API для управления подписками SaaS")

    Component(jwtService, "JWT Token Service", "Go/Node.js", "Создание, валидация и управление JWT токенами")
    Component(passwordService, "Password Service", "Go/Node.js", "Хеширование и проверка паролей")
    Component(oauthProvider, "OAuth Provider", "Go/Node.js", "Интеграция с Google, Facebook, GitHub OAuth")
    Component(emailService, "Email Service", "Go/Node.js", "Отправка email для восстановления пароля, верификации")
    Component(billingService, "Billing Service", "Go/Node.js", "Обработка платежей и биллинг подписок")

    Component(userRepository, "User Repository", "Go/Node.js", "Доступ к данным пользователей в БД")
    Component(subscriptionRepository, "Subscription Repository", "Go/Node.js", "Доступ к данным подписок в БД")
    Component(sessionRepository, "Session Repository", "Go/Node.js", "Управление пользовательскими сессиями")
}

ContainerDb(userDb, "User Database", "PostgreSQL", "Пользователи, подписки, сессии")
Container(messageQueue, "Message Queue", "Apache Kafka", "Event streaming")
System_Ext(emailProvider, "Email Provider", "SendGrid/AWS SES")
System_Ext(paymentProvider, "Payment Provider", "Stripe/PayPal")
System_Ext(oauthProviders, "OAuth Providers", "Google, Facebook, GitHub")

' Внешние связи
Rel(user, webApp, "Использует", "HTTPS")
Rel(webApp, apiGateway, "API calls", "HTTPS")
Rel(apiGateway, authController, "Authentication requests", "HTTP/gRPC")
Rel(apiGateway, userController, "User management", "HTTP/gRPC")
Rel(apiGateway, subscriptionController, "Subscription management", "HTTP/gRPC")

' Внутренние связи контроллеров
Rel(authController, jwtService, "Создает/валидирует токены", "")
Rel(authController, passwordService, "Проверяет пароли", "")
Rel(authController, oauthProvider, "OAuth аутентификация", "")
Rel(authController, userRepository, "Получает данные пользователей", "")
Rel(authController, sessionRepository, "Управляет сессиями", "")

Rel(userController, userRepository, "CRUD операции", "")
Rel(userController, emailService, "Отправляет уведомления", "")
Rel(userController, passwordService, "Хеширует новые пароли", "")

Rel(subscriptionController, subscriptionRepository, "CRUD операции", "")
Rel(subscriptionController, billingService, "Обрабатывает платежи", "")
Rel(subscriptionController, userRepository, "Связывает с пользователями", "")

' Связи сервисов с внешними системами
Rel(emailService, emailProvider, "Отправляет email", "SMTP/API")
Rel(billingService, paymentProvider, "Обрабатывает платежи", "HTTPS API")
Rel(oauthProvider, oauthProviders, "OAuth flow", "HTTPS")

' Связи с базой данных
Rel(userRepository, userDb, "Users table", "SQL")
Rel(subscriptionRepository, userDb, "Subscriptions table", "SQL")
Rel(sessionRepository, userDb, "Sessions table", "SQL")

' Event streaming
Rel(authController, messageQueue, "События аутентификации", "Kafka")
Rel(userController, messageQueue, "События пользователя", "Kafka")
Rel(subscriptionController, messageQueue, "События подписки", "Kafka")

@enduml
```

### Аутентификация пользователя

```puml

@startuml
!theme plain
title Последовательность аутентификации пользователя (TO BE)

actor Пользователь as User
autonumber
participant "Web App" as WebApp
participant "API Gateway" as Gateway
participant "Auth Controller" as AuthCtrl
participant "Password Service" as PwdService
participant "JWT Service" as JWTService
participant "User Repository" as UserRepo
participant "Session Repository" as SessionRepo
participant "User Database" as UserDB
participant "Message Queue" as MQ

== Аутентификация по логину и паролю ==

User -> WebApp: Ввод логина/пароля
WebApp -> Gateway: POST /auth/login\n{email, password}
Gateway -> AuthCtrl: HTTP Request\nAuthenticate user

AuthCtrl -> UserRepo: findByEmail(email)
UserRepo -> UserDB: SELECT * FROM users\nWHERE email = ?
UserDB --> UserRepo: User data or null
UserRepo --> AuthCtrl: User object or null

alt Пользователь не найден
    AuthCtrl --> Gateway: 401 Unauthorized\n{error: "Invalid credentials"}
    Gateway --> WebApp: HTTP 401
    WebApp --> User: Ошибка входа
else Пользователь найден
    AuthCtrl -> PwdService: validatePassword(inputPassword, hashedPassword)
    PwdService --> AuthCtrl: isValid (boolean)

    alt Неверный пароль
        AuthCtrl --> Gateway: 401 Unauthorized\n{error: "Invalid credentials"}
        Gateway --> WebApp: HTTP 401
        WebApp --> User: Ошибка входа
    else Пароль верный
        AuthCtrl -> JWTService: generateAccessToken(userId, userRole)
        JWTService --> AuthCtrl: accessToken (JWT)

        AuthCtrl -> JWTService: generateRefreshToken(userId)
        JWTService --> AuthCtrl: refreshToken

        AuthCtrl -> SessionRepo: createSession(userId, refreshToken, expiresAt)
        SessionRepo -> UserDB: INSERT INTO sessions\n(user_id, refresh_token, expires_at)
        UserDB --> SessionRepo: Session created
        SessionRepo --> AuthCtrl: Session object

        AuthCtrl -> MQ: publishEvent("user.authenticated", {userId, timestamp})
        MQ --> AuthCtrl: Event published

        AuthCtrl --> Gateway: 200 OK\n{accessToken, refreshToken, user}
        Gateway --> WebApp: HTTP 200 + Tokens
        WebApp -> WebApp: Сохранить токены в localStorage
        WebApp --> User: Успешный вход
    end
end

== Валидация токена для защищенных запросов ==

User -> WebApp: Действие требующее авторизации
WebApp -> Gateway: GET /devices\nAuthorization: Bearer <accessToken>
Gateway -> Gateway: Извлечь токен из заголовка
Gateway -> AuthCtrl: validateToken(accessToken)
AuthCtrl -> JWTService: verifyToken(accessToken)
JWTService --> AuthCtrl: {valid: true, userId, role} or {valid: false}

alt Токен невалидный или истек
    AuthCtrl --> Gateway: 401 Unauthorized\n{error: "Invalid token"}
    Gateway --> WebApp: HTTP 401
    WebApp -> WebApp: Удалить токены
    WebApp --> User: Перенаправление на страницу входа
else Токен валидный
    Gateway -> Gateway: Добавить userId в context
    Gateway -> "Device Service": GET /devices\nX-User-Id: userId
    note right: Продолжение запроса к целевому сервису
end

== Обновление токена ==

WebApp -> Gateway: POST /auth/refresh\n{refreshToken}
Gateway -> AuthCtrl: Refresh access token

AuthCtrl -> SessionRepo: findByRefreshToken(refreshToken)
SessionRepo -> UserDB: SELECT * FROM sessions\nWHERE refresh_token = ?
UserDB --> SessionRepo: Session data or null
SessionRepo --> AuthCtrl: Session object or null

alt Сессия не найдена или истекла
    AuthCtrl --> Gateway: 401 Unauthorized\n{error: "Invalid refresh token"}
    Gateway --> WebApp: HTTP 401
    WebApp -> WebApp: Удалить токены
    WebApp --> User: Перенаправление на страницу входа
else Сессия валидная
    AuthCtrl -> JWTService: generateAccessToken(userId, userRole)
    JWTService --> AuthCtrl: newAccessToken

    AuthCtrl -> SessionRepo: updateLastUsed(sessionId)
    SessionRepo -> UserDB: UPDATE sessions\nSET last_used = NOW()
    UserDB --> SessionRepo: Updated
    SessionRepo --> AuthCtrl: Success

    AuthCtrl --> Gateway: 200 OK\n{accessToken}
    Gateway --> WebApp: HTTP 200 + New token
    WebApp -> WebApp: Обновить токен в localStorage
end

@enduml

```

## 2. Device Service

**Домен:** Управление устройствами

**Ответственность:** CRUD устройств, команды управления, состояния устройств

```puml
@startuml
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Component.puml

title Device Service - Components Diagram (TO BE)

Person(user, "Пользователь", "Владелец умного дома")
Container(webApp, "Web Application", "SPA", "Веб-интерфейс")
Container(apiGateway, "API Gateway", "Nginx", "Единая точка входа")

Container_Boundary(deviceService, "Device Service") {
    Component(deviceController, "Device Controller", "Go/Node.js", "REST API для управления устройствами, CRUD операции")
    Component(sensorController, "Sensor Controller", "Go/Node.js", "REST API для управления датчиками температуры")
    Component(heatingController, "Heating Controller", "Go/Node.js", "REST API для управления системами отопления")
    Component(commandController, "Command Controller", "Go/Node.js", "REST API для отправки команд устройствам")

    Component(deviceManager, "Device Manager", "Go/Node.js", "Основная логика управления устройствами")
    Component(sensorService, "Sensor Service", "Go/Node.js", "Обработка данных датчиков температуры")
    Component(heatingService, "Heating Service", "Go/Node.js", "Управление системами отопления, термостатами")
    Component(commandService, "Command Service", "Go/Node.js", "Очередь и выполнение команд устройствам")
    Component(telemetryCollector, "Telemetry Collector", "Go/Node.js", "Сбор телеметрии от устройств")

    Component(deviceRepository, "Device Repository", "Go/Node.js", "Доступ к метаданным устройств в БД")
    Component(sensorRepository, "Sensor Repository", "Go/Node.js", "Доступ к данным датчиков в БД")
    Component(heatingRepository, "Heating Repository", "Go/Node.js", "Доступ к данным систем отопления")
    Component(commandRepository, "Command Repository", "Go/Node.js", "История команд и статусы выполнения")
}

ContainerDb(deviceDb, "Device Database", "PostgreSQL", "Метаданные устройств, состояния, конфигурации")
Container(messageQueue, "Message Queue", "Apache Kafka", "Event streaming")
Container(realTimeQueue, "Real-time Queue", "Redis Pub/Sub", "Real-time уведомления")
System_Ext(sensors, "Датчики", "Физические датчики температуры")
System_Ext(heatingDevices, "Системы отопления", "Котлы, радиаторы, термостаты")

' Внешние связи
Rel(user, webApp, "Использует", "HTTPS")
Rel(webApp, apiGateway, "API calls", "HTTPS")
Rel(apiGateway, deviceController, "Управление устройствами", "HTTP/gRPC")
Rel(apiGateway, sensorController, "Управление датчиками", "HTTP/gRPC")
Rel(apiGateway, heatingController, "Управление отоплением", "HTTP/gRPC")
Rel(apiGateway, commandController, "Управление командами", "HTTP/gRPC")

' Внутренние связи контроллеров
Rel(deviceController, deviceManager, "Основные операции", "")
Rel(deviceController, deviceRepository, "CRUD операции", "")

Rel(sensorController, sensorService, "Обработка датчиков", "")
Rel(sensorController, sensorRepository, "CRUD операции", "")
Rel(sensorController, telemetryCollector, "Сбор данных", "")

Rel(heatingController, heatingService, "Управление отоплением", "")
Rel(heatingController, heatingRepository, "CRUD операции", "")

Rel(commandController, commandService, "Выполнение команд", "")
Rel(commandController, commandRepository, "История команд", "")

' Связи сервисов
Rel(deviceManager, sensorService, "Управление датчиками", "")
Rel(deviceManager, heatingService, "Управление отоплением", "")
Rel(deviceManager, telemetryCollector, "Сбор телеметрии", "")

Rel(commandService, deviceManager, "Выполнение команд", "")

' Связи с базой данных
Rel(deviceRepository, deviceDb, "Devices table", "SQL")
Rel(sensorRepository, deviceDb, "Sensors table", "SQL")
Rel(heatingRepository, deviceDb, "HeatingSystems table", "SQL")
Rel(commandRepository, deviceDb, "Commands table", "SQL")

' Event streaming
Rel(deviceManager, messageQueue, "События устройств", "Kafka")
Rel(sensorService, messageQueue, "События датчиков", "Kafka")
Rel(heatingService, messageQueue, "События отопления", "Kafka")
Rel(commandService, messageQueue, "События команд", "Kafka")
Rel(telemetryCollector, messageQueue, "События телеметрии", "Kafka")

' Real-time коммуникация
Rel(deviceManager, realTimeQueue, "Состояния устройств", "Redis Pub/Sub")
Rel(sensorService, realTimeQueue, "Данные датчиков", "Redis Pub/Sub")
Rel(heatingService, realTimeQueue, "Состояния отопления", "Redis Pub/Sub")

' Связи с физическими устройствами
Rel(sensorService, sensors, "Чтение данных", "API устройств")
Rel(heatingService, heatingDevices, "Управление", "API устройств")
Rel(telemetryCollector, sensors, "Сбор телеметрии", "API устройств")
Rel(telemetryCollector, heatingDevices, "Сбор телеметрии", "API устройств")

@enduml
```

### Управление устройством

```puml

@startuml
!theme plain
title Последовательность управления устройством (TO BE)

actor Пользователь as User
autonumber
participant "Web App" as WebApp
participant "API Gateway" as Gateway
participant "Command Controller" as CmdCtrl
participant "Command Service" as CmdService
participant "Device Manager" as DeviceManager
participant "Heating Service" as HeatingService
participant "Command Repository" as CmdRepo
participant "Device Repository" as DeviceRepo
participant "Device Database" as DeviceDB
participant "Message Queue" as MQ
participant "Real-time Queue" as RTQ
participant "Heating Device" as Device

== Отправка команды управления устройством ==

User -> WebApp: Изменить температуру\n(Установить 22°C)
WebApp -> Gateway: POST /devices/123/commands\nAuthorization: Bearer <token>\n{type: "SET_TEMPERATURE", value: 22}
Gateway -> Gateway: Валидация токена
Gateway -> CmdCtrl: Execute command\nUserId: 456, DeviceId: 123

CmdCtrl -> DeviceRepo: findById(deviceId: 123)
DeviceRepo -> DeviceDB: SELECT * FROM devices\nWHERE id = 123 AND user_id = 456
DeviceDB --> DeviceRepo: Device data or null
DeviceRepo --> CmdCtrl: Device object or null

alt Устройство не найдено или не принадлежит пользователю
    CmdCtrl --> Gateway: 404 Not Found\n{error: "Device not found"}
    Gateway --> WebApp: HTTP 404
    WebApp --> User: Ошибка: устройство не найдено
else Устройство найдено
    CmdCtrl -> CmdService: executeCommand(deviceId, commandType, value, userId)

    CmdService -> CmdRepo: createCommand(deviceId, commandType, value, userId, status: "PENDING")
    CmdRepo -> DeviceDB: INSERT INTO commands\n(device_id, type, value, user_id, status, created_at)
    DeviceDB --> CmdRepo: Command created
    CmdRepo --> CmdService: Command object with ID

    CmdService -> MQ: publishEvent("command.created", {commandId, deviceId, type, value})
    MQ --> CmdService: Event published

    CmdService -> DeviceManager: sendCommand(deviceId, commandType, value)
    DeviceManager -> HeatingService: setTemperature(deviceId: 123, temperature: 22)

    HeatingService -> Device: HTTP POST /api/control\n{action: "SET_TEMP", value: 22}

    alt Устройство недоступно или ошибка
        Device --> HeatingService: HTTP 500 or Timeout
        HeatingService --> DeviceManager: {success: false, error: "Device unreachable"}
        DeviceManager --> CmdService: Command failed

        CmdService -> CmdRepo: updateCommandStatus(commandId, status: "FAILED", error: "Device unreachable")
        CmdRepo -> DeviceDB: UPDATE commands\nSET status = 'FAILED', error = 'Device unreachable'
        DeviceDB --> CmdRepo: Updated
        CmdRepo --> CmdService: Success

        CmdService -> MQ: publishEvent("command.failed", {commandId, deviceId, error})
        MQ --> CmdService: Event published

        CmdService -> RTQ: publish("device.123.status", {status: "OFFLINE", error: "Unreachable"})
        RTQ --> CmdService: Published

        CmdService --> CmdCtrl: {success: false, error: "Device unreachable"}
        CmdCtrl --> Gateway: 422 Unprocessable Entity\n{error: "Device unreachable"}
        Gateway --> WebApp: HTTP 422
        WebApp --> User: Ошибка: устройство недоступно

    else Команда выполнена успешно
        Device --> HeatingService: HTTP 200\n{status: "OK", currentTemp: 22}
        HeatingService --> DeviceManager: {success: true, currentTemp: 22}
        DeviceManager --> CmdService: Command executed successfully

        CmdService -> CmdRepo: updateCommandStatus(commandId, status: "COMPLETED", response: {currentTemp: 22})
        CmdRepo -> DeviceDB: UPDATE commands\nSET status = 'COMPLETED', response = '{"currentTemp": 22}'
        DeviceDB --> CmdRepo: Updated
        CmdRepo --> CmdService: Success

        CmdService -> DeviceRepo: updateDeviceState(deviceId, {temperature: 22, lastSeen: NOW()})
        DeviceRepo -> DeviceDB: UPDATE devices\nSET state = '{"temperature": 22}', last_seen = NOW()
        DeviceDB --> DeviceRepo: Updated
        DeviceRepo --> CmdService: Success

        CmdService -> MQ: publishEvent("command.completed", {commandId, deviceId, result})
        MQ --> CmdService: Event published

        CmdService -> RTQ: publish("device.123.status", {temperature: 22, status: "ONLINE"})
        RTQ --> CmdService: Published

        CmdService --> CmdCtrl: {success: true, result: {temperature: 22}}
        CmdCtrl --> Gateway: 200 OK\n{commandId, status: "COMPLETED", result: {temperature: 22}}
        Gateway --> WebApp: HTTP 200 + Result
        WebApp --> User: Температура установлена: 22°C
    end
end

== Real-time обновление статуса ==

note over RTQ, WebApp: WebSocket соединение активно
RTQ -> WebApp: WebSocket message\n{type: "DEVICE_UPDATE", deviceId: 123, temperature: 22}
WebApp -> WebApp: Обновить UI с новым состоянием
WebApp --> User: Отображение текущей температуры: 22°C

== Мониторинг выполнения команды ==

User -> WebApp: Проверить статус команды
WebApp -> Gateway: GET /commands/789\nAuthorization: Bearer <token>
Gateway -> CmdCtrl: Get command status

CmdCtrl -> CmdRepo: findById(commandId: 789, userId: 456)
CmdRepo -> DeviceDB: SELECT * FROM commands\nWHERE id = 789 AND user_id = 456
DeviceDB --> CmdRepo: Command data
CmdRepo --> CmdCtrl: Command object

CmdCtrl --> Gateway: 200 OK\n{commandId: 789, status: "COMPLETED", createdAt, completedAt}
Gateway --> WebApp: HTTP 200 + Command details
WebApp --> User: Статус команды: Выполнена

@enduml
```

## 3. Scenario Service

**Домен:** Автоматизация и пользовательские сценарии

**Ответственность**: Правила автоматизации, триггеры, планировщик

```puml

@startuml
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Component.puml

title Scenario Service - Components Diagram (TO BE)

Person(user, "Пользователь", "Владелец умного дома")
Container(webApp, "Web Application", "SPA", "Веб-интерфейс")
Container(apiGateway, "API Gateway", "Nginx", "Единая точка входа")

Container_Boundary(scenarioService, "Scenario Service") {
    Component(scenarioController, "Scenario Controller", "Go/Node.js", "REST API для управления сценариями, CRUD операции")
    Component(ruleController, "Rule Controller", "Go/Node.js", "REST API для управления правилами автоматизации")
    Component(automationController, "Automation Controller", "Go/Node.js", "REST API для управления автоматическими действиями")
    Component(schedulerController, "Scheduler Controller", "Go/Node.js", "REST API для управления расписаниями выполнения")

    Component(scenarioManager, "Scenario Manager", "Go/Node.js", "Основная логика управления пользовательскими сценариями")
    Component(ruleEngine, "Rule Engine", "Go/Node.js", "Обработка правил, условий и логических операций")
    Component(automationService, "Automation Service", "Go/Node.js", "Выполнение автоматических действий по сценариям")
    Component(eventProcessor, "Event Processor", "Go/Node.js", "Обработка событий от устройств и триггеров")
    Component(schedulerService, "Scheduler Service", "Go/Node.js", "Управление расписаниями и временными триггерами")
    Component(validationService, "Validation Service", "Go/Node.js", "Валидация сценариев и правил")

    Component(scenarioRepository, "Scenario Repository", "Go/Node.js", "Доступ к данным сценариев в БД")
    Component(ruleRepository, "Rule Repository", "Go/Node.js", "Доступ к данным правил и условий")
    Component(automationRepository, "Automation Repository", "Go/Node.js", "История выполнения автоматических действий")
    Component(scheduleRepository, "Schedule Repository", "Go/Node.js", "Доступ к данным расписаний и триггеров")
}

ContainerDb(scenarioDb, "Scenario Database", "PostgreSQL", "Сценарии автоматизации, правила")
Container(messageQueue, "Message Queue", "Apache Kafka", "Event streaming")
Container(realTimeQueue, "Real-time Queue", "Redis Pub/Sub", "Real-time уведомления")
Container_Ext(deviceService, "Device Service", "Go/Node.js", "Управление устройствами")
Container_Ext(notificationService, "Notification Service", "Go/Node.js", "Уведомления")

' Внешние связи
Rel(user, webApp, "Использует", "HTTPS")
Rel(webApp, apiGateway, "API calls", "HTTPS")
Rel(apiGateway, scenarioController, "Управление сценариями", "HTTP/gRPC")
Rel(apiGateway, ruleController, "Управление правилами", "HTTP/gRPC")
Rel(apiGateway, automationController, "Управление автоматизацией", "HTTP/gRPC")
Rel(apiGateway, schedulerController, "Управление расписаниями", "HTTP/gRPC")

' Внутренние связи контроллеров
Rel(scenarioController, scenarioManager, "Основные операции", "")
Rel(scenarioController, validationService, "Валидация сценариев", "")
Rel(scenarioController, scenarioRepository, "CRUD операции", "")

Rel(ruleController, ruleEngine, "Обработка правил", "")
Rel(ruleController, validationService, "Валидация правил", "")
Rel(ruleController, ruleRepository, "CRUD операции", "")

Rel(automationController, automationService, "Управление автоматизацией", "")
Rel(automationController, automationRepository, "История выполнения", "")

Rel(schedulerController, schedulerService, "Управление расписаниями", "")
Rel(schedulerController, scheduleRepository, "CRUD операции", "")

' Связи сервисов
Rel(scenarioManager, ruleEngine, "Применение правил", "")
Rel(scenarioManager, automationService, "Запуск автоматизации", "")
Rel(scenarioManager, eventProcessor, "Обработка событий", "")

Rel(ruleEngine, eventProcessor, "Оценка условий", "")
Rel(automationService, eventProcessor, "Триггеры действий", "")
Rel(schedulerService, eventProcessor, "Временные события", "")

Rel(eventProcessor, scenarioManager, "Активация сценариев", "")

' Связи с базой данных
Rel(scenarioRepository, scenarioDb, "Scenarios table", "SQL")
Rel(ruleRepository, scenarioDb, "Rules table", "SQL")
Rel(automationRepository, scenarioDb, "Automations table", "SQL")
Rel(scheduleRepository, scenarioDb, "Schedules table", "SQL")

' Event streaming
Rel(scenarioManager, messageQueue, "События сценариев", "Kafka")
Rel(ruleEngine, messageQueue, "События правил", "Kafka")
Rel(automationService, messageQueue, "События автоматизации", "Kafka")
Rel(eventProcessor, messageQueue, "События устройств", "Kafka")
Rel(schedulerService, messageQueue, "События расписаний", "Kafka")

' Real-time коммуникация
Rel(scenarioManager, realTimeQueue, "Статусы сценариев", "Redis Pub/Sub")
Rel(automationService, realTimeQueue, "Статусы автоматизации", "Redis Pub/Sub")

' Связи с другими сервисами
Rel(automationService, deviceService, "Команды устройствам", "gRPC")
Rel(automationService, notificationService, "Отправка уведомлений", "gRPC")

@enduml

```

### Выполнение сценария

```puml
@startuml
!theme plain
title Последовательность выполнения сценария автоматизации (TO BE)

autonumber
participant "Message Queue" as MQ
participant "Scenario Service" as ScenarioService
participant "Scenario Controller" as ScenarioCtrl
participant "Rule Engine" as RuleEngine
participant "Action Executor" as ActionExecutor
participant "Scenario Repository" as ScenarioRepo
participant "Execution Repository" as ExecutionRepo
participant "Scenario Database" as ScenarioDB
participant "Device Service" as DeviceService
participant "Notification Service" as NotificationService
participant "Telemetry Service" as TelemetryService
participant "Real-time Queue" as RTQ
participant "Web App" as WebApp
actor Пользователь as User

== Trigger сценария от телеметрии ==

MQ -> ScenarioService: Event: "scenario.triggered"\n{scenarioId: 456, trigger: "LOW_TEMP", deviceId: 123, value: 18.5}
ScenarioService -> ScenarioCtrl: processScenarioTrigger(scenarioId, triggerData)

ScenarioCtrl -> ScenarioRepo: findById(scenarioId: 456)
ScenarioRepo -> ScenarioDB: SELECT * FROM scenarios\nWHERE id = 456 AND active = true
ScenarioDB --> ScenarioRepo: Scenario data or null
ScenarioRepo --> ScenarioCtrl: Scenario object or null

alt Сценарий не найден или неактивен
    ScenarioCtrl --> ScenarioService: {success: false, error: "Scenario not found or inactive"}
    ScenarioService -> MQ: publishEvent("scenario.failed", {scenarioId, error: "Not found"})
    MQ --> ScenarioService: Event published
else Сценарий найден и активен
    ScenarioCtrl -> ExecutionRepo: createExecution(scenarioId, trigger, status: "STARTED")
    ExecutionRepo -> ScenarioDB: INSERT INTO scenario_executions\n(scenario_id, trigger_data, status, started_at)
    ScenarioDB --> ExecutionRepo: Execution created
    ExecutionRepo --> ScenarioCtrl: Execution object with ID

    ScenarioCtrl -> RuleEngine: evaluateConditions(scenario, triggerData)

    == Вычисление условий сценария ==

    RuleEngine -> RuleEngine: Парсинг условий сценария:\n- IF temperature < 20°C\n- AND time between 18:00-06:00\n- AND heating_mode = "auto"

    RuleEngine -> TelemetryService: getCurrentTemperature(deviceId: 123)
    TelemetryService --> RuleEngine: {temperature: 18.5, timestamp}

    RuleEngine -> RuleEngine: Проверка времени:\n- Текущее время: 22:30\n- В диапазоне 18:00-06:00: ✓

    RuleEngine -> DeviceService: getDeviceState(heatingDeviceId: 456)
    DeviceService --> RuleEngine: {mode: "auto", status: "ON"}

    RuleEngine -> RuleEngine: Вычисление результата:\n- temperature (18.5) < 20: ✓\n- time in range: ✓\n- mode = "auto": ✓\n→ Все условия выполнены

    RuleEngine --> ScenarioCtrl: {conditionsMet: true, details: {...}}

    alt Условия не выполнены
        ScenarioCtrl -> ExecutionRepo: updateExecution(executionId, status: "SKIPPED", reason: "Conditions not met")
        ExecutionRepo -> ScenarioDB: UPDATE scenario_executions\nSET status = 'SKIPPED', completed_at = NOW()
        ScenarioDB --> ExecutionRepo: Updated
        ExecutionRepo --> ScenarioCtrl: Success

        ScenarioCtrl --> ScenarioService: {success: true, result: "SKIPPED"}
        ScenarioService -> MQ: publishEvent("scenario.skipped", {scenarioId, executionId, reason})
        MQ --> ScenarioService: Event published

    else Условия выполнены - выполняем действия
        ScenarioCtrl -> ActionExecutor: executeActions(scenario.actions, context)

        == Выполнение действий сценария ==

        ActionExecutor -> ActionExecutor: Парсинг действий:\n1. Увеличить температуру на радиаторе на 2°C\n2. Отправить уведомление пользователю\n3. Включить дополнительный обогреватель

        loop Для каждого действия
            alt Действие: Увеличить температуру радиатора
                ActionExecutor -> DeviceService: sendCommand(deviceId: 456, action: "INCREASE_TEMP", value: 2)
                DeviceService -> DeviceService: Выполнение команды управления устройством
                DeviceService --> ActionExecutor: {success: true, newTemperature: 22}

            else Действие: Отправить уведомление
                ActionExecutor -> NotificationService: sendNotification(userId, type: "SCENARIO_EXECUTED", message: "Включено дополнительное отопление")
                NotificationService --> ActionExecutor: {success: true, notificationId}

            else Действие: Включить обогреватель
                ActionExecutor -> DeviceService: sendCommand(deviceId: 789, action: "TURN_ON")
                DeviceService --> ActionExecutor: {success: true, deviceStatus: "ON"}
            end
        end

        ActionExecutor -> ActionExecutor: Сбор результатов выполнения:\n- Радиатор: успешно (22°C)\n- Уведомление: отправлено\n- Обогреватель: включен

        ActionExecutor --> ScenarioCtrl: {success: true, executedActions: 3, results: [...]}

        ScenarioCtrl -> ExecutionRepo: updateExecution(executionId, status: "COMPLETED", results: executionResults)
        ExecutionRepo -> ScenarioDB: UPDATE scenario_executions\nSET status = 'COMPLETED', results = '...', completed_at = NOW()
        ScenarioDB --> ExecutionRepo: Updated
        ExecutionRepo --> ScenarioCtrl: Success

        ScenarioCtrl -> ScenarioRepo: updateScenarioStats(scenarioId, lastExecuted: NOW(), executionCount++)
        ScenarioRepo -> ScenarioDB: UPDATE scenarios\nSET last_executed = NOW(), execution_count = execution_count + 1
        ScenarioDB --> ScenarioRepo: Updated
        ScenarioRepo --> ScenarioCtrl: Success

        ScenarioCtrl --> ScenarioService: {success: true, result: "COMPLETED", executedActions: 3}
    end
end

== Публикация результатов ==

ScenarioService -> MQ: publishEvent("scenario.completed", {scenarioId, executionId, results})
MQ --> ScenarioService: Event published

ScenarioService -> RTQ: publish("scenario.456.status", {status: "COMPLETED", executedAt, actions: 3})
RTQ --> ScenarioService: Published

== Real-time обновление пользовательского интерфейса ==

note over RTQ, WebApp: WebSocket соединение активно
RTQ -> WebApp: WebSocket message\n{type: "SCENARIO_EXECUTED", scenarioId: 456, name: "Подогрев при низкой температуре"}
WebApp -> WebApp: Обновить страницу сценариев\nс отметкой о выполнении
WebApp --> User: Уведомление: "Сценарий выполнен: Включено дополнительное отопление"

== Ручное выполнение сценария пользователем ==

User -> WebApp: Нажать "Выполнить сценарий"
WebApp -> ScenarioService: POST /scenarios/456/execute\nAuthorization: Bearer <token>
ScenarioService -> ScenarioCtrl: executeScenario(scenarioId: 456, triggeredBy: "USER", userId: 123)

ScenarioCtrl -> ScenarioRepo: findById(scenarioId: 456)
ScenarioRepo -> ScenarioDB: SELECT * FROM scenarios\nWHERE id = 456 AND user_id = 123
ScenarioDB --> ScenarioRepo: Scenario data
ScenarioRepo --> ScenarioCtrl: Scenario object

ScenarioCtrl -> ExecutionRepo: createExecution(scenarioId, trigger: "MANUAL", userId: 123)
ExecutionRepo -> ScenarioDB: INSERT INTO scenario_executions\n(scenario_id, triggered_by, user_id, status, started_at)
ScenarioDB --> ExecutionRepo: Manual execution created
ExecutionRepo --> ScenarioCtrl: Execution object

note right: Ручное выполнение пропускает проверку условий
ScenarioCtrl -> ActionExecutor: executeActions(scenario.actions, manualContext)

ActionExecutor -> ActionExecutor: Выполнение всех действий\nбез проверки условий
ActionExecutor --> ScenarioCtrl: {success: true, executedActions: 3}

ScenarioCtrl -> ExecutionRepo: updateExecution(executionId, status: "COMPLETED")
ExecutionRepo -> ScenarioDB: UPDATE scenario_executions\nSET status = 'COMPLETED', completed_at = NOW()
ScenarioDB --> ExecutionRepo: Updated

ScenarioCtrl --> ScenarioService: {success: true, executionId, message: "Scenario executed manually"}
ScenarioService --> WebApp: HTTP 200 OK\n{executionId, status: "COMPLETED"}
WebApp --> User: "Сценарий выполнен успешно"

== Мониторинг выполнения сценариев ==

User -> WebApp: Открыть историю выполнения
WebApp -> ScenarioService: GET /scenarios/456/executions?limit=10\nAuthorization: Bearer <token>
ScenarioService -> ExecutionRepo: findByScenarioId(scenarioId: 456, userId: 123, limit: 10)
ExecutionRepo -> ScenarioDB: SELECT * FROM scenario_executions\nWHERE scenario_id = 456 AND user_id = 123\nORDER BY started_at DESC LIMIT 10
ScenarioDB --> ExecutionRepo: Execution history
ExecutionRepo --> ScenarioService: Executions array

ScenarioService --> WebApp: HTTP 200 OK\n{executions: [...]}
WebApp --> User: Отображение истории:\n- 22:30 - Выполнен автоматически\n- 21:15 - Выполнен вручную\n- 20:45 - Пропущен (условия не выполнены)

@enduml
```

## 4. Telemetry Service

**Домен** Сбор и обработка телеметрии

**Ответственность:** Метрики устройств, аналитика

```puml

@startuml
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Component.puml

title Telemetry Service - Components Diagram (TO BE)

Person(user, "Пользователь", "Владелец умного дома")
Container(webApp, "Web Application", "SPA", "Веб-интерфейс")
Container(apiGateway, "API Gateway", "Nginx", "Единая точка входа")

Container_Boundary(telemetryService, "Telemetry Service") {
    Component(telemetryController, "Telemetry Controller", "Go/Node.js", "REST API для получения телеметрических данных")
    Component(analyticsController, "Analytics Controller", "Go/Node.js", "REST API для аналитики и отчетов")
    Component(alertsController, "Alerts Controller", "Go/Node.js", "REST API для управления алертами по телеметрии")

    Component(dataCollector, "Data Collector", "Go/Node.js", "Сбор телеметрических данных от устройств")
    Component(dataProcessor, "Data Processor", "Go/Node.js", "Обработка и валидация входящих данных")
    Component(aggregationService, "Aggregation Service", "Go/Node.js", "Агрегация данных по временным интервалам")
    Component(analyticsEngine, "Analytics Engine", "Go/Node.js", "Анализ трендов, вычисление метрик")
    Component(alertsEngine, "Alerts Engine", "Go/Node.js", "Обработка алертов и пороговых значений")

    Component(telemetryRepository, "Telemetry Repository", "Go/Node.js", "Доступ к телеметрическим данным в БД")
    Component(metricsRepository, "Metrics Repository", "Go/Node.js", "Доступ к агрегированным метрикам")
    Component(alertsRepository, "Alerts Repository", "Go/Node.js", "Доступ к настройкам и истории алертов")
}

ContainerDb(telemetryDb, "Telemetry Database", "PostgreSQL", "Телеметрические данные, метрики, алерты")
Container(messageQueue, "Message Queue", "Apache Kafka", "Event streaming")
Container(realTimeQueue, "Real-time Queue", "Redis Pub/Sub", "Real-time уведомления")
Container(timeSeriesDb, "Time Series Database", "InfluxDB", "Высокопроизводительное хранение временных рядов")
Container_Ext(deviceService, "Device Service", "Go/Node.js", "Источник телеметрии устройств")
Container_Ext(notificationService, "Notification Service", "Go/Node.js", "Уведомления об алертах")

' Внешние связи
Rel(user, webApp, "Использует", "HTTPS")
Rel(webApp, apiGateway, "API calls", "HTTPS")
Rel(apiGateway, telemetryController, "Получение телеметрии", "HTTP/gRPC")
Rel(apiGateway, analyticsController, "Аналитика и отчеты", "HTTP/gRPC")
Rel(apiGateway, alertsController, "Управление алертами", "HTTP/gRPC")

' Внутренние связи контроллеров
Rel(telemetryController, dataCollector, "Получение данных", "")
Rel(telemetryController, telemetryRepository, "Чтение телеметрии", "")

Rel(analyticsController, analyticsEngine, "Аналитика", "")
Rel(analyticsController, metricsRepository, "Получение метрик", "")

Rel(alertsController, alertsEngine, "Управление алертами", "")
Rel(alertsController, alertsRepository, "CRUD операции", "")

' Связи сервисов
Rel(dataCollector, dataProcessor, "Обработка данных", "")
Rel(dataProcessor, aggregationService, "Агрегация", "")
Rel(dataProcessor, alertsEngine, "Проверка пороговых значений", "")

Rel(aggregationService, analyticsEngine, "Анализ агрегированных данных", "")
Rel(analyticsEngine, alertsEngine, "Аналитические алерты", "")

' Связи с базами данных
Rel(telemetryRepository, telemetryDb, "Telemetry table", "SQL")
Rel(telemetryRepository, timeSeriesDb, "Raw telemetry data", "InfluxQL")
Rel(metricsRepository, telemetryDb, "Metrics table", "SQL")
Rel(alertsRepository, telemetryDb, "Alerts table", "SQL")

' Event streaming
Rel(dataCollector, messageQueue, "События сбора данных", "Kafka")
Rel(dataProcessor, messageQueue, "События обработки", "Kafka")
Rel(aggregationService, messageQueue, "События агрегации", "Kafka")
Rel(alertsEngine, messageQueue, "События алертов", "Kafka")
Rel(messageQueue, dataCollector, "События от устройств", "Kafka")

' Real-time коммуникация
Rel(dataProcessor, realTimeQueue, "Real-time данные", "Redis Pub/Sub")
Rel(alertsEngine, realTimeQueue, "Real-time алерты", "Redis Pub/Sub")
Rel(analyticsEngine, realTimeQueue, "Real-time аналитика", "Redis Pub/Sub")

' Связи с другими сервисами
Rel(deviceService, dataCollector, "Телеметрия устройств", "gRPC/HTTP")
Rel(alertsEngine, notificationService, "Отправка алертов", "gRPC")

@enduml

```

### Обработчка телеметрии

```puml

@startuml
!theme plain
title Последовательность обработки телеметрии (TO BE)

autonumber
participant "Temperature Sensor" as Sensor
participant "Telemetry Collector" as Collector
participant "Telemetry Controller" as TelemetryCtrl
participant "Telemetry Processor" as TelemetryProc
participant "Analytics Service" as Analytics
participant "Alert Service" as AlertService
participant "Telemetry Repository" as TelemetryRepo
participant "Device Repository" as DeviceRepo
participant "Telemetry Database" as TelemetryDB
participant "Device Database" as DeviceDB
participant "Message Queue" as MQ
participant "Real-time Queue" as RTQ
participant "Cache Service" as Cache
participant "Web App" as WebApp
participant "Scenario Service" as ScenarioService

== Сбор телеметрии от датчика ==

Sensor -> Collector: HTTP POST /telemetry\n{deviceId: 123, temperature: 18.5, humidity: 65, timestamp}
Collector -> Collector: Валидация данных\nи формата сообщения

alt Невалидные данные
    Collector --> Sensor: HTTP 400 Bad Request\n{error: "Invalid data format"}
else Данные валидны
    Collector -> DeviceRepo: findById(deviceId: 123)
    DeviceRepo -> DeviceDB: SELECT * FROM devices WHERE id = 123
    DeviceDB --> DeviceRepo: Device data or null
    DeviceRepo --> Collector: Device object or null

    alt Устройство не найдено
        Collector --> Sensor: HTTP 404 Not Found\n{error: "Device not found"}
    else Устройство найдено
        Collector -> TelemetryCtrl: processTelemetry(deviceId, telemetryData)

        TelemetryCtrl -> TelemetryProc: processRawTelemetry(deviceId, data)

        == Обработка и обогащение данных ==

        TelemetryProc -> TelemetryProc: Обогащение данных\n(добавление метаданных, калибровка)
        TelemetryProc -> Analytics: calculateTrends(deviceId, currentTemp: 18.5)
        Analytics -> TelemetryRepo: getRecentData(deviceId, period: "1h")
        TelemetryRepo -> TelemetryDB: SELECT * FROM telemetry\nWHERE device_id = 123 AND timestamp > NOW() - INTERVAL '1 hour'
        TelemetryDB --> TelemetryRepo: Historical data
        TelemetryRepo --> Analytics: Recent telemetry array

        Analytics -> Analytics: Расчет тенденций:\n- Средняя температура за час: 19.2°C\n- Тренд: понижение на 0.7°C
        Analytics --> TelemetryProc: {avgTemp: 19.2, trend: "decreasing", rate: -0.7}

        == Сохранение данных ==

        TelemetryProc -> TelemetryRepo: saveTelemetry(enrichedData)
        TelemetryRepo -> TelemetryDB: INSERT INTO telemetry\n(device_id, temperature, humidity, trends, timestamp)
        TelemetryDB --> TelemetryRepo: Record saved
        TelemetryRepo --> TelemetryProc: Success

        TelemetryProc -> DeviceRepo: updateDeviceStatus(deviceId, lastSeen: NOW(), status: "ONLINE")
        DeviceRepo -> DeviceDB: UPDATE devices\nSET last_seen = NOW(), status = 'ONLINE'
        DeviceDB --> DeviceRepo: Updated
        DeviceRepo --> TelemetryProc: Success

        == Кэширование для быстрого доступа ==

        TelemetryProc -> Cache: setCurrent(deviceId, currentState: {temp: 18.5, humidity: 65})
        Cache --> TelemetryProc: Cached

        == Публикация событий ==

        TelemetryProc -> MQ: publishEvent("telemetry.received", {deviceId, data, trends})
        MQ --> TelemetryProc: Event published

        TelemetryProc -> RTQ: publish("device.123.telemetry", {temperature: 18.5, humidity: 65, trends})
        RTQ --> TelemetryProc: Published

        == Проверка алертов ==

        TelemetryProc -> AlertService: checkAlerts(deviceId, currentTemp: 18.5)
        AlertService -> AlertService: Проверка правил:\n- Минимальная температура: 20°C\n- Текущая: 18.5°C → ALERT!

        alt Найдены алерты
            AlertService -> MQ: publishEvent("alert.triggered", {deviceId, type: "LOW_TEMPERATURE", value: 18.5, threshold: 20})
            MQ --> AlertService: Alert event published

            AlertService -> RTQ: publish("alerts.123", {type: "LOW_TEMPERATURE", message: "Температура ниже нормы"})
            RTQ --> AlertService: Alert published

        else Алертов нет
            note right: Продолжаем без алертов
        end

        AlertService --> TelemetryProc: Alert check completed

        == Trigger сценариев автоматизации ==

        TelemetryProc -> ScenarioService: checkTriggers(deviceId, telemetryData)
        ScenarioService -> ScenarioService: Поиск активных сценариев\nс триггером по температуре

        alt Найдены сценарии для выполнения
            ScenarioService -> MQ: publishEvent("scenario.triggered", {scenarioId: 456, trigger: "LOW_TEMP", deviceId})
            MQ --> ScenarioService: Scenario event published
            note right: Сценарий выполнится асинхронно
        else Сценариев нет
            note right: Продолжаем без сценариев
        end

        ScenarioService --> TelemetryProc: Scenario check completed

        TelemetryProc --> TelemetryCtrl: Processing completed successfully
        TelemetryCtrl --> Collector: Success
        Collector --> Sensor: HTTP 200 OK\n{status: "processed"}
    end
end

== Real-time обновление пользовательского интерфейса ==

note over RTQ, WebApp: WebSocket соединение активно
RTQ -> WebApp: WebSocket message\n{type: "TELEMETRY_UPDATE", deviceId: 123, temperature: 18.5, humidity: 65}
WebApp -> WebApp: Обновить dashboard\nс новыми данными

alt Есть алерт
    RTQ -> WebApp: WebSocket message\n{type: "ALERT", deviceId: 123, message: "Низкая температура"}
    WebApp -> WebApp: Показать уведомление\nо критической температуре
end

== Batch обработка для аналитики ==

note over Analytics: Выполняется каждые 5 минут
Analytics -> TelemetryRepo: getRecentBatch(period: "5m")
TelemetryRepo -> TelemetryDB: SELECT * FROM telemetry\nWHERE timestamp > NOW() - INTERVAL '5 minutes'
TelemetryDB --> TelemetryRepo: Batch telemetry data
TelemetryRepo --> Analytics: Telemetry batch

Analytics -> Analytics: Агрегация данных:\n- Средние значения по комнатам\n- Обнаружение аномалий\n- Прогнозирование потребления
Analytics -> TelemetryRepo: saveAggregatedData(aggregatedStats)
TelemetryRepo -> TelemetryDB: INSERT INTO telemetry_aggregated\n(period, room_id, avg_temp, anomalies)
TelemetryDB --> TelemetryRepo: Aggregated data saved

Analytics -> MQ: publishEvent("analytics.completed", {period, insights})
MQ --> Analytics: Analytics event published

@enduml

```

## 5. Notification Service

**Домен:** Уведомления и коммуникации

**Ответственности:** Push-уведомления, email, SMS, WebSocket

```puml

@startuml
!include https://raw.githubusercontent.com/plantuml-stdlib/C4-PlantUML/master/C4_Component.puml

title Notification Service - Components Diagram (TO BE)

Person(user, "Пользователь", "Владелец умного дома")
Container(webApp, "Web Application", "SPA", "Веб-интерфейс")
Container(apiGateway, "API Gateway", "Nginx", "Единая точка входа")

Container_Boundary(notificationService, "Notification Service") {
    Component(notificationController, "Notification Controller", "Go/Node.js", "REST API для управления уведомлениями, получение истории")
    Component(emailController, "Email Controller", "Go/Node.js", "REST API для email уведомлений и рассылок")
    Component(smsController, "SMS Controller", "Go/Node.js", "REST API для SMS алертов и уведомлений")

    Component(notificationManager, "Notification Manager", "Go/Node.js", "Основная логика управления уведомлениями")
    Component(emailService, "Email Service", "Go/Node.js", "Отправка email через SMTP провайдеров")
    Component(smsService, "SMS Service", "Go/Node.js", "Отправка SMS через SMS провайдеров")
    Component(templateService, "Template Service", "Go/Node.js", "Управление шаблонами уведомлений")
    Component(subscriptionService, "Subscription Service", "Go/Node.js", "Управление подписками пользователей на уведомления")
    Component(eventProcessor, "Event Processor", "Go/Node.js", "Обработка событий из Kafka для триггеров уведомлений")

    Component(notificationRepository, "Notification Repository", "Go/Node.js", "Доступ к истории отправленных уведомлений")
    Component(templateRepository, "Template Repository", "Go/Node.js", "Доступ к шаблонам уведомлений в БД")
    Component(subscriptionRepository, "Subscription Repository", "Go/Node.js", "Доступ к подпискам пользователей")
    Component(deviceTokenRepository, "Device Token Repository", "Go/Node.js", "Доступ к токенам устройств для push")
}

ContainerDb(notificationDb, "Notification Database", "PostgreSQL", "История уведомлений, шаблоны, подписки")
Container(messageQueue, "Message Queue", "Apache Kafka", "Event streaming")
Container(realTimeQueue, "Real-time Queue", "Redis Pub/Sub", "Real-time уведомления")
System_Ext(emailProvider, "Email Provider", "SendGrid / AWS SES")
System_Ext(smsProvider, "SMS Provider", "Twilio / AWS SNS")
Container_Ext(authService, "Auth Service", "Go/Node.js", "Пользовательские данные")
Container_Ext(telemetryService, "Telemetry Service", "Go/Node.js", "События алертов")

' Внешние связи
Rel(user, webApp, "Использует", "HTTPS")
Rel(webApp, apiGateway, "API calls", "HTTPS")
Rel(apiGateway, notificationController, "Управление уведомлениями", "HTTP/gRPC")
Rel(apiGateway, emailController, "Email уведомления", "HTTP/gRPC")
Rel(apiGateway, smsController, "SMS уведомления", "HTTP/gRPC")

' Внутренние связи контроллеров
Rel(notificationController, notificationManager, "Основные операции", "")
Rel(notificationController, notificationRepository, "История уведомлений", "")
Rel(notificationController, subscriptionService, "Управление подписками", "")

Rel(emailController, emailService, "Отправка email", "")
Rel(emailController, templateService, "Шаблоны email", "")

Rel(smsController, smsService, "Отправка SMS", "")
Rel(smsController, templateService, "Шаблоны SMS", "")

' Связи сервисов
Rel(notificationManager, emailService, "Email уведомления", "")
Rel(notificationManager, smsService, "SMS уведомления", "")
Rel(notificationManager, templateService, "Получение шаблонов", "")
Rel(notificationManager, subscriptionService, "Проверка подписок", "")

Rel(eventProcessor, notificationManager, "Триггер уведомлений", "")

Rel(templateService, templateRepository, "CRUD шаблонов", "")
Rel(subscriptionService, subscriptionRepository, "CRUD подписок", "")

' Связи с базой данных
Rel(notificationRepository, notificationDb, "Notifications table", "SQL")
Rel(templateRepository, notificationDb, "Templates table", "SQL")
Rel(subscriptionRepository, notificationDb, "Subscriptions table", "SQL")

' Event streaming
Rel(messageQueue, eventProcessor, "События от других сервисов", "Kafka")
Rel(notificationManager, messageQueue, "События уведомлений", "Kafka")
Rel(emailService, messageQueue, "События email", "Kafka")
Rel(smsService, messageQueue, "События SMS", "Kafka")

' Real-time коммуникация
Rel(notificationManager, realTimeQueue, "Real-time уведомления", "Redis Pub/Sub")
Rel(emailService, realTimeQueue, "Email статусы", "Redis Pub/Sub")
Rel(smsService, realTimeQueue, "SMS статусы", "Redis Pub/Sub")

' Связи с внешними провайдерами
Rel(emailService, emailProvider, "Отправка email", "SMTP/API")
Rel(smsService, smsProvider, "Отправка SMS", "HTTPS API")

' Связи с другими сервисами
Rel(authService, subscriptionService, "Данные пользователей", "gRPC")
Rel(telemetryService, eventProcessor, "События алертов", "Kafka")

@enduml

```
