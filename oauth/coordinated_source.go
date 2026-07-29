package oauth

type sourceCoordination struct {
	userID      string
	coordinator RefreshCoordinator
	loader      CredentialLoader
}

func withSourceCoordination(coordination *sourceCoordination) SourceOption {
	return func(options *sourceOptions) error {
		if coordination == nil || coordination.userID == "" ||
			isNilCoordinatorValue(coordination.coordinator) || coordination.loader == nil {
			return ErrInvalidOption
		}
		options.coordination = coordination
		return nil
	}
}
