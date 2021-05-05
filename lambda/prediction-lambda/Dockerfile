FROM public.ecr.aws/lambda/python:3.8

COPY requirements.txt .
RUN pip install --trusted-host pypi.org --trusted-host pypi.python.org --trusted-host=files.pythonhosted.org -r requirements.txt

COPY . .

CMD ["app.handler"]
