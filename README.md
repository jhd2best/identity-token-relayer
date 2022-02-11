# identity-token-relayer

the backend for mapping harmony identity token(base on HRC721) from ethereum

#### clone and build

```shell
# get the latest code
git clone https://github.com/harmony-one/identity-token-relayer
cd identity-token-relayer

# build binary
make
```

this project support deployed by Ansible, but you need prepare the `Google firebase service-account`
and `harmony wallet private key` in advance.

#### get firebase service-account

register a Google Firebase account [here](https://firebase.google.com), and create an empty `Cloud Firestore`. after
that you can get the `service-account` by
this [guide](https://firebase.google.com/docs/admin/setup?authuser=1#set-up-project-and-service-account).

replace the `firebase-service-account.json` file under `ansible/playbooks/files/<netwokr>` folder

#### get harmony wallet private key

download harmony cli wallet and create a new wallet. you can
see [runbook](https://docs.harmony.one/home/network/wallets/harmony-cli/create-import-wallet) for more detail.

replace the `harmony-testnet.key`(fill with private key of the wallet) file under `ansible/playbooks/files/<netwokr>` folder

#### contract deploy
deploy the `OwnershipValidator.sol` in this [repo](https://github.com/harmony-one/contract-libs/tree/main/contracts) on
Harmony. do not forget init NFT owner data and transfer the owner of contract to your Harmony wallet which private key
set in `config.yaml`

#### update configuration
update the fields as follows in `config.yaml`:
- update `OwnershipValidatorAddress` field to the address of `OwnershipValidator` deploy in last step
- update `DisableSentry` to `false` and `SentryDSN` to your Sentry project DSN if you want to gather error by Sentry
- update `RpcEndpoints` to your own ethereum RPC endpoints(such as `infura` etc.)

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