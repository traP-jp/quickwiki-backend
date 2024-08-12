# データベーススキーマ

## テーブル: wikis

| カラム名       | データ型          | 制約                                      |
|----------------|-------------------|-------------------------------------------|
| id             | INT(11)           | NOT NULL PRIMARY KEY AUTO_INCREMENT       |
| name           | TEXT              | NOT NULL                                  |
| type           | TEXT              | NOT NULL                                  |
| created_at     | TIMESTAMP         | NOT NULL DEFAULT CURRENT_TIMESTAMP        |
| updated_at     | TIMESTAMP         | NOT NULL DEFAULT CURRENT_TIMESTAMP        |
| owner_traq_id  | INTEGER           | NOT NULL                                  |
| content        | TEXT              | NOT NULL                                  |

## テーブル: lectures

| カラム名       | データ型          | 制約                                      |
|----------------|-------------------|-------------------------------------------|
| id             | INT(11)           | NOT NULL PRIMARY KEY AUTO_INCREMENT       |
| title          | TEXT              | NOT NULL                                  |
| description    | TEXT              | NOT NULL                                  |
| content        | TEXT              | NOT NULL                                  |
| folder_id      | INT(11)           |                                           |

## テーブル: tags

| カラム名       | データ型          | 制約                                      |
|----------------|-------------------|-------------------------------------------|
| id             | INT(11)           | NOT NULL PRIMARY KEY AUTO_INCREMENT       |
| name           | TEXT              | NOT NULL                                  |

## テーブル: tags_in_wiki

| カラム名       | データ型          | 制約                                      |
|----------------|-------------------|-------------------------------------------|
| wiki_id        |                   | PRIMARY KEY                               |
| tag_id         |                   | PRIMARY KEY                               |
|                |                   | FOREIGN KEY (wiki_id) REFERENCES wikis(id)|
|                |                   | FOREIGN KEY (tag_id) REFERENCES tags(id)  |

## テーブル: folders

| カラム名       | データ型          | 制約                                      |
|----------------|-------------------|-------------------------------------------|
| id             | INT(11)           | NOT NULL PRIMARY KEY AUTO_INCREMENT       |
| name           | TEXT              | NOT NULL                                  |
| parent_id      | INT(11)           |                                           |
|                |                   | UNIQUE KEY (name, parent_id)              |