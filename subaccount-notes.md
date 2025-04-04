Seems its:

- From account (typically an email address)
- From type (funding / trading / margin / etc)
- To account
- To type

Need to use an API key with permissions for remitting/from account being used.
(Or maybe it's always the master account?)

So we need to add another dimension to the config:

```yaml
offchain:
  exchanges:
    okx:
      api_key: "env:OKX_API_KEY"
      secret_key: "env:OKX_API_SECRET"
      passphrase: "env:OKX_API_PASSPHRASE"
      subaccounts:
        - id: "test1@example.com"
          api_key: "env:OKX_API_KEY_TEST1"
          secret_key: "env:OKX_API_SECRET_TEST1"
          passphrase: "env:OKX_API_PASSPHRASE_TEST1"
        - id: "test2@example.com"
          api_key: "env:OKX_API_KEY_TEST2"
          secret_key: "env:OKX_API_SECRET_TEST2"
          passphrase: "env:OKX_API_PASSPHRASE_TEST2"
    bybit:
      api_key: "env:BYBIT_API_KEY"
      secret_key: "env:BYBIT_API_SECRET"
    binance:
      api_key: "env:BINANCE_API_KEY"
      secret_key: "env:BINANCE_API_SECRET"
    binanceus:
      api_key: "env:BINANCEUS_API_KEY"
      secret_key: "env:BINANCEUS_API_SECRET"
```

To support, we need special endpoints we can host:

- GET /v1/platforms/:platform/accounts-types
