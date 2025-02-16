package timeprovider

import "time"

var (
	Moscow = time.FixedZone("Moscow", 8)
)

type TimeProvider interface {
	Now() time.Time
}

type TimerProviderImpl struct {
	location *time.Location
}

func NewProvider(location *time.Location) *TimerProviderImpl {
	return &TimerProviderImpl{
		location: location,
	}
}

func (tp *TimerProviderImpl) Now() time.Time {
	return time.Now().In(tp.location)
}

type ConstTimeProvider struct {
	time.Time
}

func NewConstProvider(time time.Time) *ConstTimeProvider {
	return &ConstTimeProvider{time}
}

func (tp *ConstTimeProvider) Now() time.Time {
	return tp.Time
}
