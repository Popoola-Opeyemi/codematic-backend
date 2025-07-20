package utils

type PaymentProviders string

const (
	ProviderPaystack    PaymentProviders = "paystack"
	ProviderFlutterWave PaymentProviders = "flutterwave"
)

func (p PaymentProviders) String() string {
	return string(p)
}
