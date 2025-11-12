-- pull-request-manager-service

CREATE TABLE teams (
    id UUID PRIMARY KEY NOT NULL,
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE users (
    id UUID PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
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
    user_id UUID NOT NULL REFERENCES users(id),
    team_id UUID NOT NULL REFERENCES teams(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (user_id, team_id)
);

CREATE TABLE pull_requests (
    id UUID PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    author_id UUID NOT NULL REFERENCES users(id),
    status_id INT NOT NULL DEFAULT 1 REFERENCES statuses(id),
    need_more_reviewers BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE reviewer_assignments (
    user_id UUID NOT NULL REFERENCES users(id),
    pull_request_id UUID NOT NULL REFERENCES pull_requests(id),
    slot INT NOT NULL CHECK (slot in (1, 2)),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (pull_request_id, slot),
    UNIQUE (pull_request_id, user_id)
);

INSERT INTO statuses (id, name) VALUEs (1, 'OPEN'), (2, 'MERGED');
CREATE INDEX team_id_mem_idx ON memberships (team_id);
CREATE INDEX user_id_mem_idx ON memberships (user_id);
CREATE INDEX pr_id_idx ON reviewer_assignments (pull_request_id);
CREATE INDEX usr_id_idx ON reviewer_assignments (user_id);
