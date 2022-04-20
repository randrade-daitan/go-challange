DROP TABLE IF EXISTS task;
CREATE TABLE task (
  id         INT AUTO_INCREMENT NOT NULL,
  name       VARCHAR(255) NOT NULL,
  completed  BOOL NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
);

INSERT INTO task
  (name, completed)
VALUES
  ('A Tour of Go', '1'),
  ('How to Write Go Code', '1'),
  ('Writing Web Applications', '1'),
  ('Accessing a relational database', '1'),
  ('DB Setup', '0'),
  ('CRUD of Tasks', '0'),
  ('REST API of Tasks', '0'),
  ('ORM layer', '0'),
  ('gRPC API', '0'),
  ('Practical GO', '0');
