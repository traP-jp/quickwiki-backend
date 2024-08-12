CREATE TABLE wikis (
    id INT(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    owner_traq_id INTEGER NOT NULL,
    content TEXT NOT NULL, 
    tags TEXT[]
);
CREATE TABLE lectures (
    id INT(11) NOT NULL AUTO_INCREMENT,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    content TEXT NOT NULL,
    parent_id INT(11)
    KEY parent_id_idx (parent_id)
        FOREIGN KEY (parent_id)
        REFERENCES lectures(id)
);