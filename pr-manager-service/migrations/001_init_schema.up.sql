-- pull-request-manager-service

CREATE TABLE teams (
    team_name TEXT PRIMARY KEY NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users (
    user_id TEXT PRIMARY KEY NOT NULL,
    username TEXT NOT NULL,
    is_active BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE statuses (
    id INT PRIMARY KEY NOT NULL,
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE memberships (
    user_id TEXT NOT NULL REFERENCES users(user_id),
    team_name TEXT NOT NULL REFERENCES teams(team_name),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (user_id, team_name)
);

CREATE TABLE pull_requests (
    pull_request_id TEXT PRIMARY KEY NOT NULL,
    pull_request_name TEXT NOT NULL,
    author_id TEXT NOT NULL REFERENCES users(user_id),
    status_id INT NOT NULL DEFAULT 1 REFERENCES statuses(id),
    need_more_reviewers BOOLEAN NOT NULL DEFAULT true,
    mergedAt TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
    
CREATE TABLE reviewer_assignments (
    user_id TEXT NOT NULL REFERENCES users(user_id),
    pull_request_id TEXT NOT NULL REFERENCES pull_requests(pull_request_id),
    slot INT NOT NULL CHECK (slot in (1, 2)),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (pull_request_id, slot),
    UNIQUE (pull_request_id, user_id)
);

INSERT INTO statuses (id, name) VALUES (1, 'OPEN'), (2, 'MERGED');
CREATE INDEX team_name_mem_idx ON memberships (team_name);
CREATE INDEX user_id_mem_idx ON memberships (user_id);
CREATE INDEX pr_id_idx ON reviewer_assignments (pull_request_id);
CREATE INDEX usr_id_idx ON reviewer_assignments (user_id);
