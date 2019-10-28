package orm

type ormErr string

func (e ormErr) Error() string {
	return string(e)
}

const (
	errInvalidTable                 ormErr = "InvalidTable"
	errInvalidScanHolder            ormErr = "InvalidScanHolder"
	errScanHolderMustBeValidPointer ormErr = "ScanHolderMustBeValidPointer"
	errSelectQueryNeedDataHolder    ormErr = "SelectQueryNeedDataHolder"
)
