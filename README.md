# tournaments-game-creator

# Возможности сервиса:
- CRUD-операции над играми (игры как сессии)
- CRUD-операции над результатами (результаты эти игр, есть возможность указать нескольких победителей)

```mermaid
---
config:
  layout: dagre
---
erDiagram
    games ||--|{ results : ""
    game_types ||--|{ games : ""
    
    games {
        UUID game_id PK
        DATETIME game_start
        UUID game_type_id FK
    }
    
    game_types {
        UUID game_type_id PK
        TEXT platform_name
    }
    
    results {
        UUID result_id PK
        UUID game_id FK
        UUID winner_id FK
        TEXT comment
    }

```



_____________

# !!!
## В ходе разработки было принято решение изменить некоторые id сущностей с INT на UUID
## Миграции сделаны с помощью Liquibase, репозиторий с базой данных скрыт
