CREATE TABLE wikis (
    id INT(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    owner_traq_id CHAR(36) NOT NULL,
    content TEXT NOT NULL
);
CREATE TABLE messages (
    id INT(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    wiki_id INT(11) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_traq_id CHAR(36) NOT NULL,
    message_traq_id CHAR(36) NOT NULL,
    channel_id CHAR(36) NOT NULL,
    content TEXT NOT NULL,
    FOREIGN KEY (wiki_id) REFERENCES wikis(id)
);
CREATE TABLE messageStamps (
    id INT(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    message_id INT(11) NOT NULL,-- = messages.id
    stamp_traq_id CHAR(36) NOT NULL,
    count INT(11) NOT NULL,
    FOREIGN KEY (message_id) REFERENCES messages(id)
);
CREATE TABLE citedMessages (
    id INT(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    parent_message_id INT(11) NOT NULL,-- = messages.id
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_traq_id CHAR(36) NOT NULL,
    message_traq_id CHAR(36) NOT NULL,
    channel_id CHAR(36) NOT NULL,
    content TEXT NOT NULL,
    FOREIGN KEY (parent_message_id) REFERENCES messages(id)
);
CREATE TABLE memos (
    id INT(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    wiki_id INT(11) NOT NULL,
    title TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    owner_traq_id CHAR(36) NOT NULL,
    content TEXT NOT NULL,
    FOREIGN KEY (wiki_id) REFERENCES wikis(id)
);
CREATE TABLE tags (
    id INT(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    wiki_id INT(11) NOT NULL,
    name TEXT NOT NULL,
    tag_score FLOAT8 NOT NULL,
    FOREIGN KEY (wiki_id) REFERENCES wikis(id)
);
CREATE TABLE folders (
    id INT(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name TEXT NOT NULL,
    parent_id INT(11),-- 0 if root
    UNIQUE KEY (name, parent_id)
);
CREATE TABLE lectures (
    id INT(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    folder_id INT(11),
    folder_path TEXT NOT NULL,
    FOREIGN KEY (folder_id) REFERENCES folders(id)
);
CREATE TABLE favorites (
    id INT(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    user_traq_id CHAR(36) NOT NULL,
    wiki_id INT(11) NOT NULL,
    FOREIGN KEY (wiki_id) REFERENCES wikis(id)
);