
# load our xc wallet
 export PRIVATE_KEY=...

# lookup our exchange deposit addresses
OKX_DEPOSIT_ADDRESS=$(oc --exchange okx deposit --symbol SOL --network Solana)
BYBIT_DEPOSIT_ADDRESS=$(oc --exchange bybit deposit --symbol SOL --network SOL)

# our self-custody wallet
ADDRESS=$(xc --chain SOL address)

# Send to OKX
xc --chain SOL transfer $OKX_DEPOSIT_ADDRESS 0.03 -v

# transfer to trading account
oc --exchange okx transfer -vv --from core/funding --to core/trading --amount 0.03 --symbol SOL

# (?? make profit)

# transfer back to funding account
oc --exchange okx transfer -vv --from core/trading --to core/funding --amount 0.03 --symbol SOL

# withdraw back to xc
oc --exchange okx withdraw -vv --to $ADDRESS --amount 0.06 --symbol SOL --network Solana

# Send to Bybit
xc --chain SOL transfer $BYBIT_DEPOSIT_ADDRESS 0.06 -v

# transfer to trading account
oc --exchange bybit transfer -vv --from core/funding --to core/trading --amount 0.06 --symbol SOL

# (?? make profit)

# transfer back to funding account
oc --exchange bybit transfer -vv --from core/trading --to core/funding --amount 0.09 --symbol SOL

# withdraw back to xc
oc --exchange bybit withdraw -vv --to $ADDRESS --amount 0.09 --symbol SOL --network SOL
