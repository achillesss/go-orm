package orm

type ormErr string

func (e ormErr) Error() string {
	return string(e)
}

const (
	ErrInvalidTable                 ormErr = "InvalidTable"
	ErrInvalidScanHolder            ormErr = "InvalidScanHolder"
	ErrScanHolderMustBeValidPointer ormErr = "ScanHolderMustBeValidPointer"
	ErrSelectQueryNeedDataHolder    ormErr = "SelectQueryNeedDataHolder"
	ErrNotSupportType               ormErr = "NotSupportType"
	ErrNotFound                     ormErr = "NotFound"
	ErrInvalidQuery                 ormErr = "InvalidQuery"
)
