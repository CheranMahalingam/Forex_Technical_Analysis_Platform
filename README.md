# Forex Analytics Platform
A platform providing foreign exchange rates, global market news & forecasting tools in real-time

<p align="center">
  <img src="images/home_page.png">
</p>

<h2>Exchange Rates</h2>
<p align="center">
  <img src="images/exchange_rate.png">
</p>

<h2>LSTM Generated Forecasts</h2>
<p align="center">
  <img src="images/forecast.png">
</p>

<h2>Global Market News</h2>
<p align="center">
  <img src="images/news.png">
</p>

# AWS Deployment Instructions
Login to the public AWS Elastic Container Registry with
```
aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws
```

cd into the directory containing lambda functions with
```
cd lambda
```

Start docker locally and run
```
sam build
```

To deploy the lambda functions run
```
sam deploy -g
```

# Platform Architecture
<p align="center">
  <img src="images/architecture_diagram.png">
</p>

# Future Features
- Develop backtrading system
- Add economic calendar
- Setup user authentication
- Add support for more currency pairs
- Use Elasticsearch to filter news
