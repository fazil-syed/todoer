-- MIGRATION 1
CREATE TABLE IF NOT EXISTS task_groups (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
);

-- INITIAL TASK TABLE


CREATE TABLE IF NOT EXISTS tasks (

    id INTEGER PRIMARY KEY AUTOINCREMENT,

    title TEXT NOT NULL,

    done BOOLEAN NOT NULL DEFAULT FALSE,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP

);



-- DEFAULT GROUP and Settings

INSERT OR IGNORE INTO task_groups (name) VALUES ('default');

INSERT OR IGNORE INTO settings (key,value) VALUES('selected_group',(SELECT id from task_groups WHERE name = 'default'));


-- MIGRATION 2

CREATE TABLE IF NOT EXISTS tasks_new (

    id INTEGER PRIMARY KEY AUTOINCREMENT,

    group_id INTEGER NOT NULL,

    title TEXT NOT NULL,

    done BOOLEAN NOT NULL DEFAULT FALSE,

    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (group_id) REFERENCES task_groups(id)

);


-- Copy existing tasks into the  default group

INSERT INTO tasks_new(id,group_id,title,done,created_at)
SELECT 
    id,
    (SELECT id FROM task_groups WHERE name = 'default'),
    title,
    done,
    created_at
FROM tasks;

DROP TABLE tasks;

ALTER TABLE tasks_new RENAME TO tasks;

CREATE INDEX idx_tasks_group_id ON tasks(group_id);

