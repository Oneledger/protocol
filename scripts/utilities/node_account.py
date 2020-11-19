from sdkcom import fullnode, addValidatorWalletAccounts

if __name__ == "__main__":
    account = addValidatorWalletAccounts(fullnode)
    print account
