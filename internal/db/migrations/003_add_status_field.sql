CREATE TABLE tasks_new (
    id INTEGER PRIMARY KEY,
    group_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    status TEXT NOT NULL
        CHECK (status in ('TODO','IN_PROGRESS','DONE'))
        DEFAULT "TODO",
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY(group_id)
        REFERENCES task_groups(id)
);


INSERT INTO tasks_new(
    id,
    group_id,
    title,
    status,
    created_at
)
SELECT
    id,
    group_id,
    title,
    CASE
        WHEN done = 1 THEN 'DONE'
        ELSE 'TODO'
    END    
    ,
    created_at
FROM tasks;



DROP TABLE tasks;

ALTER TABLE tasks_new
RENAME TO tasks;

CREATE INDEX idx_tasks_group_id
ON tasks(group_id);