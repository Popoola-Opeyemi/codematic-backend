package provider

var registry = map[string]Service{}

func RegisterProvider(code string, p Service) {
	registry[code] = p
}

func GetProvider(code string) (Service, bool) {
	p, ok := registry[code]
	return p, ok
}
