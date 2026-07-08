CREATE TABLE task_groups (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL
);
INSERT OR IGNORE INTO task_groups(name)
VALUES ('default');

INSERT OR IGNORE INTO settings(key,value)
VALUES(
    'selected_group',
    (SELECT id FROM task_groups WHERE name='default')
);


CREATE TABLE tasks_new (
    id INTEGER PRIMARY KEY,
    group_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    done BOOLEAN NOT NULL DEFAULT FALSE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY(group_id)
        REFERENCES task_groups(id)
);

INSERT INTO tasks_new(
    id,
    group_id,
    title,
    done,
    created_at
)
SELECT
    id,
    (SELECT id FROM task_groups WHERE name='default'),
    title,
    done,
    created_at
FROM tasks;

DROP TABLE tasks;

ALTER TABLE tasks_new
RENAME TO tasks;

CREATE INDEX idx_tasks_group_id
ON tasks(group_id);