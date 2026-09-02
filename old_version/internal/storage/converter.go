package storage

func ConnectionToUpsertConnectionDto(conn Connection) UpsertConnectionDto {
	return UpsertConnectionDto{
		Label:      conn.Label,
		FolderID:   conn.FolderID,
		Host:       conn.Host,
		Port:       conn.Port,
		User:       conn.User,
		AuthMethod: conn.AuthMethod,
		KeyPath:    conn.KeyPath,
	}
}
