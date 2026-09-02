package storage

type SQLCommand struct {
	Name string
	SQL  string
}

var createFoldersTable = SQLCommand{
	Name: "folders",
	SQL: `
		CREATE TABLE IF NOT EXISTS folders (
		    id TEXT PRIMARY KEY,
		    name TEXT NOT NULL,
		    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`,
}

var createConnectionsTable = SQLCommand{
	Name: "connections",
	SQL: `
		CREATE TABLE IF NOT EXISTS connections (
		    id TEXT PRIMARY KEY,
		    label TEXT NOT NULL,

		    folder_id TEXT,

		    last_used DATETIME,
		    is_active INTEGER NOT NULL DEFAULT 0,
		    is_pinned INTEGER NOT NULL DEFAULT 0,

		    host TEXT NOT NULL,
		    port TEXT NOT NULL,
		    user TEXT NOT NULL,

		    auth_method TEXT NOT NULL,
		    key_path TEXT,

		    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

		    FOREIGN KEY (folder_id) REFERENCES folders(id)
		        ON DELETE SET NULL
		);
	`,
}

var createTagsTable = SQLCommand{
	Name: "tags",
	SQL: `
		CREATE TABLE IF NOT EXISTS tags (
		    id TEXT PRIMARY KEY,
		    name TEXT NOT NULL UNIQUE,
		    color TEXT,
		    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`,
}

var createConnectionsTagsRelationTable = SQLCommand{
	Name: "connections_tags",
	SQL: `
		CREATE TABLE IF NOT EXISTS connection_tags (
		    connection_id TEXT NOT NULL,
		    tag_id TEXT NOT NULL,

		    PRIMARY KEY (connection_id, tag_id),

		    FOREIGN KEY (connection_id) REFERENCES connections(id)
		        ON DELETE CASCADE,

		    FOREIGN KEY (tag_id) REFERENCES tags(id)
		        ON DELETE CASCADE
		);
	`,
}

var createConnectionsLabelIndex = SQLCommand{
	Name: "idx_connections_label",
	SQL: `
		CREATE UNIQUE INDEX IF NOT EXISTS idx_connections_label ON connections(LOWER(label));
	`,
}

var createConnectionsFolderIndex = SQLCommand{
	Name: "idx_connections_folder_id",
	SQL: `
		CREATE INDEX IF NOT EXISTS idx_connections_folder_id ON connections(folder_id);
	`,
}

var createConnectionsLastUsedIndex = SQLCommand{
	Name: "idx_connections_last_used",
	SQL: `
		CREATE INDEX IF NOT EXISTS idx_connections_last_used ON connections(last_used);
	`,
}

var createConnectionTagsConnectionIdIndex = SQLCommand{
	Name: "idx_connection_tags_connection_id",
	SQL: `
		CREATE INDEX IF NOT EXISTS idx_connection_tags_connection_id ON connection_tags(connection_id);
	`,
}

var createConnectionTagsTagIdIndex = SQLCommand{
	Name: "idx_connection_tags_tag_id",
	SQL: `
		CREATE INDEX IF NOT EXISTS idx_connection_tags_tag_id ON connection_tags(tag_id);
	`,
}

var stmts = []SQLCommand{
	createFoldersTable,
	createConnectionsTable,
	createTagsTable,
	createConnectionsTagsRelationTable,

	createConnectionsLabelIndex,
	createConnectionsFolderIndex,
	createConnectionsLastUsedIndex,
	createConnectionTagsConnectionIdIndex,
	createConnectionTagsTagIdIndex,
}
