import finnhub
import pandas as pd
from dotenv import load_dotenv
import os
load_dotenv()

finnhub_client = finnhub.Client(api_key=os.getenv("PROJECT_API_KEY"))

res = finnhub_client.forex_candles('OANDA:EUR_USD', '60', 1609200000, 1609357882)
print(finnhub_client.forex_rates(base='USD'))
df = pd.DataFrame(res)
print(df.describe())
