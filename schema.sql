CREATE TABLE wikis (
    id INT(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    owner_traq_id INTEGER NOT NULL,
    content TEXT NOT NULL
);
CREATE TABLE lectures (
    id INT(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    folder_id INT(11)
    folder_path TEXT NOT NULL,
);
CREATE TABLE tags (
    id INT(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name TEXT NOT NULL
);
CREATE TABLE tags_in_wiki (
    PRIMARY KEY (wiki_id, tag_id),
    FOREIGN KEY (wiki_id)
        REFERENCES wikis(id),
    FOREIGN KEY (tag_id)
        REFERENCES tags(id)
);
CREATE TABLE folders (
    id INT(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name TEXT NOT NULL,
    parent_id INT(11),  -- 0 if root
    UNIQUE KEY (name, parent_id)
);