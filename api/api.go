//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=/app/mock/mock_$GOPACKAGE/$GOFILE
package api

type Api interface {
	Do(body interface{}) error
	GetResult() *[]byte
}
