from sdkcom import fullnode, addValidatorAccounts

if __name__ == "__main__":
    account = addValidatorAccounts(fullnode)
    if len(account) > 0:
        print account
