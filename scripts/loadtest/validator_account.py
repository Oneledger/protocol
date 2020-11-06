from sdkcom import oltest, fullnode_dev, fullnode_prod, addValidatorAccounts

if __name__ == "__main__":
    if oltest == "1":
        account = addValidatorAccounts(fullnode_dev)
    else:
        account = addValidatorAccounts(fullnode_prod)
    if len(account) > 0:
        print account
