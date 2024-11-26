// Клиент для API Платформы Талант v2
//
// Документация к API: http://talent.kruzhok.org/v2/docs/
package talentv2

//go:generate go run github.com/shagohead/gotools@v0.1.3 ogen -target ./ -package talentv2 -clean openapi.yaml

const (
	ProductionURL = "https://talent.kruzhok.org/v2"
	StageURL      = "https://talent.test.kruzhok.org/v2"
	InternalURL   = "http://t2-api:8000/v2"
)

func New(url string, sec SecuritySource, opts ...ClientOption) (*Client, error) {
	if url == "" {
		url = InternalURL
	}
	return NewClient(url, sec, opts...)
}
