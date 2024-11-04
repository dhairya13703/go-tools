package db

type BackupService interface {
    Backup(hostname, username, password, port, database, output string) error
}
