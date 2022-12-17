.DEFAULT_GOAL := help


## bootstrap
bootstrap:
	cdk bootstrap

## deploy
deploy:
	cdk deploy "*" --require-approval never --no-rollback
	rm -f main

## hotswap
hotswap:
	cdk deploy "*" --require-approval never--hotswap
	rm -f main

## synth
synth:
	cdk synth "*"
	rm -f main

## diff
diff:
	cdk diff "*"
	rm -f main

## destroy
_destroy:
	cdk destroy "*"
	rm -f main


## help
help:
	@make2help $(MAKEFILE_LIST)


.PHONY: help
.SILENT:
