#
FROM public.ecr.aws/lambda/python:3.9
COPY requirements.txt  .
RUN  pip3 install -r requirements.txt --target "${LAMBDA_TASK_ROOT}"
COPY index.py ${LAMBDA_TASK_ROOT}
CMD [ "index.handler" ]
