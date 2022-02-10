# identity-token-relayer

the backend for mapping harmony identity token(base on HRC721) from ethereum

#### clone and build

```shell
# get the latest code
git clone https://github.com/jhd2best/identity-token-relayer
cd identity-token-relayer

# build binary
go build -ldflags "-s -w" -o relayer identity-token-relayer
```

this project support deployed by Ansible, but you need prepare the `Google firebase service-account`
and `harmony wallet private key` in advance.

#### get firebase service-account

register a Google Firebase account [here](https://firebase.google.com), and create an empty `Cloud Firestore`. after
that you can get the `service-account` by
this [guide](https://firebase.google.com/docs/admin/setup?authuser=1#set-up-project-and-service-account).

#### get harmony wallet private key

download harmony cli wallet and create a new wallet. you can
see [runbook](https://docs.harmony.one/home/network/wallets/harmony-cli/create-import-wallet) for more detail.

#### pre-deploy

deploy the `OwnershipValidator.sol` in this [repo](https://github.com/harmony-one/contract-libs/tree/main/contracts) on
Harmony and update the `OwnershipValidatorAddress` field in `config.yaml`. do not forget init NFT owner data and transfer
the owner of contract to your Harmony wallet which private key set in `config.yaml`

place your `firebase` and `harmony` related files under `ansible/playbooks/files/<netwokr>` folder, and replace
the `firebase-service-account.json` and `harmony-testnet.key`.

if you want change the name of files above, do not forget update the `config.yaml` too.

and then, change the ip and login username of your server in `ansible/hosts` file.

#### install service

```shell
ansible-playbook playbooks/install.yaml -i hosts -e "inventory=relayer-backend binary_path=<binary_path> network=<network>"
```

replace the `<binary_path>` to your path of relayer binary(default is `../relayer`) and `<network>` to the network you
want to map on harmony(default is `testnet`).

#### update service

```shell
ansible-playbook playbooks/update.yaml -i hosts -e "inventory=relayer-backend binary_path=<binary_path> network=<network>"
```

replace the variation in command too

#### restart service

```shell
ansible-playbook playbooks/restart.yaml -i hosts -e "inventory=relayer-backend network=<network>"
```

#### stop service

```shell
ansible-playbook playbooks/stop.yaml -i hosts -e "inventory=relayer-backend"
```

> NOTE: all commands above need run under `ansible` folder.