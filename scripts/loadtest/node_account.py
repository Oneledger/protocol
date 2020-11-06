from sdkcom import oltest, fullnode_dev, fullnode_prod, addValidatorWalletAccounts

if __name__ == "__main__":
    if oltest == "1":
        account = addValidatorWalletAccounts(fullnode_dev)
    else:
        account = addValidatorWalletAccounts(fullnode_prod)
    print account
