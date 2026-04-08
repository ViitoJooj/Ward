package services

type ProxyService struct{}

func NewProxyService() *ProxyService {
	return &ProxyService{}
}

func (s *ProxyService) MethodGet(url string) {

}
