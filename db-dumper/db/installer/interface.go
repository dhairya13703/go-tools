package installer

type UtilityInstaller interface {
    Install() error
    IsInstalled() bool
    GetName() string
}
