package cloud

type ECloudProvider int64

const (
	GCP   = 0
	AWS   = 1
	Azure = 2
)

// TODO: Please make sure that all cloud providers respect the below interface and return the `CloudSession` struct, we can use a map to store open services and retrieve them ( A lot of work xd )

type ICloudSession interface {
	OpenSession(region string) (*Session, error)
}

type Session struct {
	CloudProvider ECloudProvider
	Session       *ICloudSession
}
