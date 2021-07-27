# Terraform Provider AWS LexBotV2

In the collection, you will use the AWS LexBotV2 provider as a bridge between Terraform and the AWS LexBotV2 API. 

## Build provider

Run the following command to build the provider

```shell
$ go build -o terraform-provider-awslexbotv2
```

## Test sample configuration

First, build and install the provider.

```shell
$ make install
```

Then, navigate to the `examples` directory. 

```shell
$ cd examples
```

Run the following command to initialize the workspace and apply the sample configuration.

```shell
$ terraform init && terraform apply
```
